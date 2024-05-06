package rabbitMQ

import (
	"app-services-go/kit/event"
	"context"
	"encoding/json"
	"log"
	"os"

	"github.com/rabbitmq/amqp091-go"
)

// EventBus is an rabbitMq implementation of the event.Bus.
type EventBus struct {
	connection         *Connection
	exchange           string
	queueNameFormatter QueueFormatter
	maxRetries         int
}
type ConsumeMessage amqp091.Delivery

// NewEventBus initializes a new EventBus.
func NewEventBus(connection *Connection, exchange string, queueNameFormatter QueueFormatter, maxRetries int) *EventBus {
	return &EventBus{
		connection:         connection,
		exchange:           exchange,
		queueNameFormatter: queueNameFormatter,
		maxRetries:         maxRetries,
	}
}

type optionsEvent struct {
	MessageId       string
	ContentType     string
	ContentEncoding string
}

func (b *EventBus) Publish(ctx context.Context, events []event.Event) error {

	for _, evt := range events {
		var routingKey = evt.Type()
		var content, _ = json.Marshal(evt)
		var options = optionsEvent{
			MessageId:       evt.ID(),
			ContentType:     "application/json",
			ContentEncoding: "utf-8",
		}

		var err = b.connection.Publish(b.exchange, string(routingKey),
			amqp091.Publishing{
				MessageId:   options.MessageId,
				ContentType: options.ContentType,
				Body:        content,
			})
		if err != nil {
			return err
		}
	}

	return nil
}

// Subscribe implements the event.Bus interface.
func (b *EventBus) Subscribe(subscriber event.Subscriber) {
	queueName := b.queueNameFormatter.Format(subscriber)
	messageChannel, err := b.connection.Consume(queueName)
	handleError(err, "Could not register consumer")

	go func() {
		log.Printf("Consumer ready, Name %s_consumer PID: %d", queueName, os.Getpid())
		for d := range messageChannel {
			//log.Printf("Received a message: %s", d.Body)
			queueName := b.queueNameFormatter.Format(subscriber)
			eventDeserializer := Configure(subscriber)

			consumerFactory := NewConsumerFactory(eventDeserializer, b.connection, b.maxRetries)
			consumer := consumerFactory.Build(subscriber, b.exchange, queueName)

			consumer.OnMessage(ConsumeMessage(d))
		}
	}()

}
