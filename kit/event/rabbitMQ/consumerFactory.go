package rabbitMQ

import "app-services-go/kit/event"

type ConsumerFactory struct {
	deserializer *EventDeserializer
	connection   *Connection
	maxRetries   int
}

func NewConsumerFactory(
	deserializer *EventDeserializer,
	connection *Connection,
	maxRetries int) *ConsumerFactory {
	return &ConsumerFactory{
		deserializer: deserializer,
		connection:   connection,
		maxRetries:   maxRetries,
	}
}

func (c *ConsumerFactory) Build(
	subscriber event.Subscriber,
	exchange string,
	queueName string) *Consumer {
	return NewConsumer(
		subscriber,
		c.deserializer,
		c.connection,
		queueName,
		exchange,
		c.maxRetries,
	)
}
