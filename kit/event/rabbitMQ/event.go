package rabbitMQ

import (
	"app-services-go/kit/event"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/rabbitmq/amqp091-go"
)

// Ensure *EventBus satisfies the event.Bus interface.
var _ event.Bus = (*EventBus)(nil)

// EventBus is an rabbitMq implementation of the event.Bus.
type EventBus struct {
	connection         *Connection
	exchange           string
	queueNameFormatter QueueFormatter
	maxRetries         int

	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup
}
type ConsumeMessage amqp091.Delivery

// NewEventBus initializes a new EventBus.
func NewEventBus(connection *Connection, exchange string, queueNameFormatter QueueFormatter, maxRetries int) *EventBus {
	ctx, cancel := context.WithCancel(context.Background())
	return &EventBus{
		connection:         connection,
		exchange:           exchange,
		queueNameFormatter: queueNameFormatter,
		maxRetries:         maxRetries,
		ctx:                ctx,
		cancel:             cancel,
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
		content, err := json.Marshal(evt)
		if err != nil {
			return err
		}
		var options = optionsEvent{
			MessageId:       evt.ID(),
			ContentType:     "application/json",
			ContentEncoding: "utf-8",
		}

		if err := b.connection.Publish(b.exchange, string(routingKey),
			amqp091.Publishing{
				MessageId:       options.MessageId,
				ContentType:     options.ContentType,
				ContentEncoding: options.ContentEncoding,
				Body:            content,
			}); err != nil {
			return err
		}
	}

	return nil
}

// Subscribe implements the event.Bus interface.
//
// The Bus interface is intentionally return-less (frozen), so a registration
// failure is logged here rather than aborting the process with log.Fatalf.
// The underlying registration returns the error for any caller able to use it.
func (b *EventBus) Subscribe(subscriber event.Subscriber) {
	if err := b.subscribe(subscriber); err != nil {
		log.Printf("Could not register consumer for %T: %s", subscriber, err)
	}
}

func (b *EventBus) subscribe(subscriber event.Subscriber) error {
	queueName := b.queueNameFormatter.Format(subscriber)
	messageChannel, err := b.connection.Consume(queueName)
	if err != nil {
		return fmt.Errorf("could not register consumer %s: %w", queueName, err)
	}

	// Build the deserializer/consumer once, before the delivery loop.
	eventDeserializer := Configure(subscriber)
	consumerFactory := NewConsumerFactory(eventDeserializer, b.connection, b.maxRetries)
	consumer := consumerFactory.Build(subscriber, b.exchange, queueName)

	b.wg.Add(1)
	go func() {
		defer b.wg.Done()
		log.Printf("Consumer ready, Name %s_consumer PID: %d", queueName, os.Getpid())
		for {
			if b.ctx.Err() != nil {
				return
			}
			select {
			case <-b.ctx.Done():
				return
			case d, ok := <-messageChannel:
				if !ok {
					return
				}
				consumer.OnMessage(ConsumeMessage(d))
			}
		}
	}()

	return nil
}

// Close signals all consumer goroutines to stop, waits for any in-flight
// message to finish processing (draining), and closes the connection.
func (b *EventBus) Close() error {
	b.cancel()
	b.wg.Wait()
	return b.connection.Close()
}
