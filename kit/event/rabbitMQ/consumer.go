package rabbitMQ

import (
	"app-services-go/kit/event"
	"log"
	"reflect"

	"github.com/rabbitmq/amqp091-go"
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
		// Poison message: it can never be deserialized, so route it to the
		// dead-letter queue instead of leaving it unacked forever.
		log.Printf("Error deserializing event: %s", err)
		if dlErr := c.deadLetter(message); dlErr != nil {
			log.Printf("Error dead-lettering undeserializable message: %s", dlErr)
			// Could not republish to the DLQ: requeue=false so the broker's
			// own dead-letter policy applies; do not ack (would lose it).
			if nackErr := c.connection.Nack(amqp091.Delivery(message), false); nackErr != nil {
				log.Printf("Error nacking message: %s", nackErr)
			}
			return
		}
		if ackErr := c.connection.Ack(amqp091.Delivery(message)); ackErr != nil {
			log.Printf("Error acknowledging dead-lettered message: %s", ackErr)
		}
		return
	}

	if err := c.subscriber.On(evt); err != nil {
		log.Printf("Error processing event %v: %v", reflect.TypeOf(evt), err)
		if hErr := c.handleError(message, reflect.TypeOf(evt)); hErr != nil {
			// The retry/dead-letter republish failed: do NOT ack, requeue so
			// the message is redelivered rather than lost.
			log.Printf("Error republishing failed event %v: %v", reflect.TypeOf(evt), hErr)
			if nackErr := c.connection.Nack(amqp091.Delivery(message), true); nackErr != nil {
				log.Printf("Error nacking message: %s", nackErr)
			}
			return
		}
		// Republish confirmed: only now is it safe to ack the original.
		if ackErr := c.connection.Ack(amqp091.Delivery(message)); ackErr != nil {
			log.Printf("Error acknowledging message: %s", ackErr)
		}
		return
	}

	if err := c.connection.Ack(amqp091.Delivery(message)); err != nil {
		log.Printf("Error acknowledging message: %s", err)
	}
}

func (c *Consumer) handleError(message ConsumeMessage, eventType reflect.Type) error {
	if c.hasBeenRedeliveredTooMuch(message) {
		log.Printf("Sending event %v to dead letter queue", eventType)
		return c.deadLetter(message)
	}

	count, ok := redeliveryCount(message.Headers)
	if !ok {
		count = 1
	} else {
		count++
	}
	log.Printf("Retrying event %v num %v", eventType, count)
	return c.retry(message)
}

func (c *Consumer) retry(message ConsumeMessage) error {
	return c.connection.Retry(message, c.queueName, c.exchange)
}

func (c *Consumer) deadLetter(message ConsumeMessage) error {
	return c.connection.DeadLetter(message, c.queueName, c.exchange)
}

func (c *Consumer) hasBeenRedeliveredTooMuch(message ConsumeMessage) bool {
	if count, ok := redeliveryCount(message.Headers); ok {
		return int(count) >= c.maxRetries
	}
	return false
}
