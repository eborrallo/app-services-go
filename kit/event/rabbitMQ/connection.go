package rabbitMQ

import (
	"fmt"
	"log"

	"github.com/rabbitmq/amqp091-go"
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

func handleError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}

}

type Connection struct {
	settings ConnectionSettings
	channel  *amqp091.Channel
	con      *amqp091.Connection
}

func NewConnection(settings ConnectionSettings) *Connection {
	connection := &Connection{settings: settings}
	connection.connect()
	return connection
}
func (c *Connection) connect() {
	c.con = c.amqpConnect()
	c.channel = c.amqpChannel()
}

func (c *Connection) amqpConnect() *amqp091.Connection {
	var protocol = "amqp"
	if c.settings.Connection.Secure {
		protocol = "amqps"
	}

	rabbitURI := fmt.Sprintf("%s://%s:%s@%s:%d/%s", protocol, c.settings.Username,
		c.settings.Password, c.settings.Connection.Hostname, c.settings.Connection.Port, c.settings.Vhost)
	conn, err := amqp091.Dial(rabbitURI)
	handleError(err, "Can't connect to AMQP")
	//defer conn.Close()
	return conn
}
func (c *Connection) amqpChannel() *amqp091.Channel {

	amqpChannel, err := c.con.Channel()
	handleError(err, "Can't create a amqpChannel")

	//defer amqpChannel.Close()
	return amqpChannel
}

func (c *Connection) Publish(exchange, routingKey string, msg amqp091.Publishing) error {
	return c.channel.Publish(exchange, routingKey, false, false, msg)
}
func (c *Connection) Exchange(name string) error {
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
	return c.channel.Consume(queueName, queueName+"_consumer", false, false, false, false, nil)
}
func (c *Connection) Ack(message amqp091.Delivery) error {
	return message.Ack(false)
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

func (c *Connection) incrementRedeliveryCount(message ConsumeMessage) map[string]interface{} {
	headers := message.Headers
	if headers == nil {
		headers = make(map[string]interface{})
	}
	if c.hasBeenRedelivered(message) {
		count, ok := headers["redelivery_count"].(int32)
		if ok {
			headers["redelivery_count"] = count + 1
		}
	} else {
		headers["redelivery_count"] = 1
	}

	return headers
}

func (c *Connection) hasBeenRedelivered(message ConsumeMessage) bool {
	headers := message.Headers
	if headers != nil {
		_, ok := headers["redelivery_count"].(int32)
		return ok
	}
	return false
}
