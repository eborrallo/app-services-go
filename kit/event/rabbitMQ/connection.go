package rabbitMQ

import (
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/rabbitmq/amqp091-go"
)

const (
	// redeliveryCountHeader is the AMQP header key tracking how many times a
	// message has been redelivered. It is always stored on the wire as int32.
	redeliveryCountHeader = "redelivery_count"

	connectMaxRetries = 5
	connectBaseDelay  = 500 * time.Millisecond
)

type Settings struct {
	Secure   bool
	Hostname string
	Port     int
}
type ConnectionSettings struct {
	Username   string
	Password   string
	Vhost      string
	Connection Settings
}

type Connection struct {
	settings ConnectionSettings
	channel  *amqp091.Channel
	con      *amqp091.Connection
	// chMu serializes every operation performed on the single shared AMQP
	// channel. amqp091-go is NOT safe for concurrent use of one channel, so
	// Publish/Consume/Ack/Nack/Exchange/Queue calls must not overlap.
	chMu sync.Mutex
}

// NewConnection dials the broker and opens a channel. It returns an error on
// failure instead of terminating the process.
func NewConnection(settings ConnectionSettings) (*Connection, error) {
	connection := &Connection{settings: settings}
	if err := connection.connect(); err != nil {
		return nil, err
	}
	return connection, nil
}

func (c *Connection) connect() error {
	conn, err := c.amqpConnect()
	if err != nil {
		return err
	}
	c.con = conn

	channel, err := c.amqpChannel()
	if err != nil {
		_ = c.con.Close()
		c.con = nil
		return err
	}
	c.channel = channel
	return nil
}

func (c *Connection) amqpConnect() (*amqp091.Connection, error) {
	var protocol = "amqp"
	if c.settings.Connection.Secure {
		protocol = "amqps"
	}

	rabbitURI := fmt.Sprintf("%s://%s:%s@%s:%d/%s", protocol, c.settings.Username,
		c.settings.Password, c.settings.Connection.Hostname, c.settings.Connection.Port, c.settings.Vhost)

	var lastErr error
	for attempt := 1; attempt <= connectMaxRetries; attempt++ {
		conn, err := amqp091.Dial(rabbitURI)
		if err == nil {
			return conn, nil
		}
		lastErr = err
		log.Printf("Can't connect to AMQP (attempt %d/%d): %s", attempt, connectMaxRetries, err)
		if attempt < connectMaxRetries {
			time.Sleep(connectBaseDelay * time.Duration(attempt))
		}
	}
	return nil, fmt.Errorf("can't connect to AMQP after %d attempts: %w", connectMaxRetries, lastErr)
}

func (c *Connection) amqpChannel() (*amqp091.Channel, error) {
	var lastErr error
	for attempt := 1; attempt <= connectMaxRetries; attempt++ {
		channel, err := c.con.Channel()
		if err == nil {
			return channel, nil
		}
		lastErr = err
		log.Printf("Can't create an AMQP channel (attempt %d/%d): %s", attempt, connectMaxRetries, err)
		if attempt < connectMaxRetries {
			time.Sleep(connectBaseDelay * time.Duration(attempt))
		}
	}
	return nil, fmt.Errorf("can't create an AMQP channel after %d attempts: %w", connectMaxRetries, lastErr)
}

// Close closes the channel and then the underlying connection.
func (c *Connection) Close() error {
	c.chMu.Lock()
	defer c.chMu.Unlock()

	var errs []error
	if c.channel != nil {
		if err := c.channel.Close(); err != nil {
			errs = append(errs, fmt.Errorf("closing channel: %w", err))
		}
		c.channel = nil
	}
	if c.con != nil {
		if err := c.con.Close(); err != nil {
			errs = append(errs, fmt.Errorf("closing connection: %w", err))
		}
		c.con = nil
	}
	return errors.Join(errs...)
}

func (c *Connection) Publish(exchange, routingKey string, msg amqp091.Publishing) error {
	c.chMu.Lock()
	defer c.chMu.Unlock()
	return c.channel.Publish(exchange, routingKey, false, false, msg)
}
func (c *Connection) Exchange(name string) error {
	c.chMu.Lock()
	defer c.chMu.Unlock()
	return c.channel.ExchangeDeclare(name, "topic", true, false, false, false, nil)
}

type QueueParameters struct {
	exchange           string
	name               string
	routingKeys        []string
	deadLetterExchange string
	deadLetterQueue    string
	messageTtl         int
}

func (c *Connection) Queue(p QueueParameters) error {
	const durable = true
	const exclusive = false
	const autoDelete = false

	c.chMu.Lock()
	defer c.chMu.Unlock()

	args := c.queueArguments(p.deadLetterExchange, p.deadLetterQueue, p.messageTtl)
	_, err := c.channel.QueueDeclare(p.name, durable, autoDelete, exclusive, false, args)
	if err != nil {
		return err
	}

	for _, routingKey := range p.routingKeys {
		err = c.channel.QueueBind(p.name, routingKey, p.exchange, false, nil)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *Connection) queueArguments(deadLetterExchange string, deadLetterQueue string, messageTtl int) map[string]interface{} {
	args := make(map[string]interface{})

	if deadLetterExchange != "" {
		args["x-dead-letter-exchange"] = deadLetterExchange
	}
	if deadLetterQueue != "" {
		args["x-dead-letter-routing-key"] = deadLetterQueue
	}
	if messageTtl > 0 {
		args["x-message-ttl"] = messageTtl
	}

	return args
}

func (c *Connection) Consume(queueName string) (<-chan amqp091.Delivery, error) {
	c.chMu.Lock()
	defer c.chMu.Unlock()
	return c.channel.Consume(queueName, queueName+"_consumer", false, false, false, false, nil)
}
func (c *Connection) Ack(message amqp091.Delivery) error {
	c.chMu.Lock()
	defer c.chMu.Unlock()
	return message.Ack(false)
}

func (c *Connection) Nack(message amqp091.Delivery, requeue bool) error {
	c.chMu.Lock()
	defer c.chMu.Unlock()
	return message.Nack(false, requeue)
}

func (c *Connection) Retry(message ConsumeMessage, queue string, exchange string) error {
	exch := Retry(exchange)
	options := c.messageOptions(message)
	return c.Publish(exch, queue, options)
}

func (c *Connection) DeadLetter(message ConsumeMessage, queue string, exchange string) error {
	exch := DeadLetter(exchange)
	options := c.messageOptions(message)
	return c.Publish(exch, queue, options)
}

func (c *Connection) messageOptions(message ConsumeMessage) amqp091.Publishing {
	messageId := message.MessageId
	contentType := message.ContentType
	contentEncoding := message.ContentEncoding
	priority := message.Priority

	options := amqp091.Publishing{
		MessageId:       messageId,
		Headers:         c.incrementRedeliveryCount(message),
		ContentType:     contentType,
		ContentEncoding: contentEncoding,
		Priority:        priority,
		Body:            message.Body,
	}

	return options
}

func (c *Connection) incrementRedeliveryCount(message ConsumeMessage) amqp091.Table {
	headers := message.Headers
	if headers == nil {
		headers = amqp091.Table{}
	}

	count, ok := redeliveryCount(headers)
	if !ok {
		// First delivery, or an unexpected/missing header type: (re)start the
		// count explicitly rather than silently skipping the increment.
		headers[redeliveryCountHeader] = int32(1)
		return headers
	}
	headers[redeliveryCountHeader] = count + 1

	return headers
}

// redeliveryCount reads the redelivery count from an AMQP table, tolerating the
// various integer types a table may carry, and reports whether a usable value
// was present. The value is always written back as int32 elsewhere.
func redeliveryCount(headers amqp091.Table) (int32, bool) {
	if headers == nil {
		return 0, false
	}
	raw, found := headers[redeliveryCountHeader]
	if !found {
		return 0, false
	}
	switch v := raw.(type) {
	case int32:
		return v, true
	case int:
		return int32(v), true
	case int64:
		return int32(v), true
	case int16:
		return int32(v), true
	case int8:
		return int32(v), true
	default:
		return 0, false
	}
}
