package rabbitMQ

import (
	"app-services-go/kit/event"
	"errors"
	"fmt"
	"github.com/rabbitmq/amqp091-go"
	"log"
	"reflect"
)

type Consumer struct {
	subscriber   event.Subscriber
	deserializer *EventDeserializer
	connection   *Connection
	queueName    string
	exchange     string
	maxRetries   int
}

func NewConsumer(
	subscriber event.Subscriber,
	deserializer *EventDeserializer,
	connection *Connection,
	queueName string,
	exchange string,
	maxRetries int) *Consumer {
	return &Consumer{
		subscriber:   subscriber,
		deserializer: deserializer,
		connection:   connection,
		exchange:     exchange,
		queueName:    queueName,
		maxRetries:   maxRetries,
	}
}

func (c *Consumer) OnMessage(message ConsumeMessage) {
	evt, err := c.deserializer.Deserialize(message.Body)
	if err != nil {
		log.Printf("Error deserializing event %s", err)
		return
	}
	err = c.subscriber.On(evt)
	if err != nil {
		log.Printf("Error processing event %v: %v", reflect.TypeOf(evt), err)
		c.handleError(message, reflect.TypeOf(evt))
	}
	if err := c.connection.Ack(amqp091.Delivery(message)); err != nil {
		log.Printf("Error acknowledging message : %s", err)
	}

}

func (c *Consumer) handleError(message ConsumeMessage, eventType reflect.Type) {
	if c.hasBeenRedeliveredTooMuch(message) {
		log.Printf("Sending event %v to dead letter queue", eventType)
		c.deadLetter(message)
	} else {
		count, ok := message.Headers["redelivery_count"].(int32)
		if !ok {
			count = 1
		} else {
			count++
		}
		log.Printf("Retrying event %v num %v", eventType, count)
		c.retry(message)
	}
}

func (c *Consumer) retry(message ConsumeMessage) error {
	return c.connection.Retry(message, c.queueName, c.exchange)
}

func (c *Consumer) deadLetter(message ConsumeMessage) error {
	return c.connection.DeadLetter(message, c.queueName, c.exchange)
}
func (c *Consumer) hasBeenRedeliveredTooMuch(message ConsumeMessage) bool {
	if c.hasBeenRedelivered(message) {
		count, _ := message.Headers["redelivery_count"].(int32)
		return int(count) >= c.maxRetries
	}
	return false
}
func (c *Consumer) hasBeenRedelivered(message ConsumeMessage) bool {
	if _, found := message.Headers["redelivery_count"]; found {
		return true
	} else {
		return false
	}
}

func (c *Consumer) areStructPropertiesPopulated(eventPtr interface{}) error {
	val := reflect.ValueOf(eventPtr)
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		// Check if the field is zero-valued (unpopulated)
		if reflect.DeepEqual(field.Interface(), reflect.Zero(field.Type()).Interface()) {
			return errors.New(fmt.Sprintf("field %s is empty", val.Type().Field(i).Name))

		}
	}
	return nil
}
