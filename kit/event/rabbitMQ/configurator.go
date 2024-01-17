package rabbitMQ

import (
	"app-services-go/kit/event"
)

type Configurator struct {
	connection         *Connection
	queueNameFormatter QueueFormatter
	messageRetryTtl    int
}

func NewConfigurator(connection *Connection, queueNameFormatter QueueFormatter, messageRetryTtl int) *Configurator {
	return &Configurator{
		connection:         connection,
		queueNameFormatter: queueNameFormatter,
		messageRetryTtl:    messageRetryTtl,
	}
}

func (b *Configurator) Configure(exchange string, subscribers []event.Subscriber) error {
	retryExchange := Retry(exchange)
	deadLetterExchange := DeadLetter(exchange)

	err := b.connection.Exchange(exchange)
	if err != nil {
		return err
	}

	err = b.connection.Exchange(retryExchange)
	if err != nil {
		return err
	}

	err = b.connection.Exchange(deadLetterExchange)
	if err != nil {
		return err
	}

	for _, subscriber := range subscribers {
		err := b.addQueue(subscriber, exchange)
		if err != nil {
			return err
		}
	}
	return nil
}

func (b *Configurator) addQueue(subscriber event.Subscriber, exchange string) error {
	retryExchange := Retry(exchange)
	deadLetterExchange := DeadLetter(exchange)
	routingKeys := b.routingKeysFor(subscriber)

	queue := b.queueNameFormatter.Format(subscriber)
	deadLetterQueue := b.queueNameFormatter.formatDeadLetter(subscriber)
	retryQueue := b.queueNameFormatter.formatRetry(subscriber)

	err := b.connection.Queue(QueueParameters{
		exchange:    exchange,
		name:        queue,
		routingKeys: routingKeys,
	})
	if err != nil {
		return err
	}
	err = b.connection.Queue(QueueParameters{
		exchange:           retryExchange,
		name:               retryQueue,
		routingKeys:        []string{queue},
		messageTtl:         b.messageRetryTtl,
		deadLetterQueue:    queue,
		deadLetterExchange: exchange,
	})
	if err != nil {
		return err
	}
	err = b.connection.Queue(QueueParameters{
		exchange:    deadLetterExchange,
		name:        deadLetterQueue,
		routingKeys: []string{queue},
	})
	if err != nil {
		return err
	}
	return nil
}

func (b *Configurator) routingKeysFor(subscriber event.Subscriber) []string {
	var routingKeys []string
	events := subscriber.SubscribedTo()
	routingKeys = append(routingKeys, string(events.Type()))

	queue := b.queueNameFormatter.Format(subscriber)
	routingKeys = append(routingKeys, queue)
	return routingKeys

}
