package bootstrap

import (
	"app-services-go/cmd/api/bootstrap/providers"
	"app-services-go/configs"
	"app-services-go/kit/cache"
	"app-services-go/kit/command"
	"app-services-go/kit/event"
	"app-services-go/kit/event/rabbitMQ"
	"app-services-go/kit/query"
	"app-services-go/kit/storage"

	"github.com/redis/go-redis/v9"
)

func Container() (*command.CommandBus, *query.QueryBus, redis.UniversalClient, error) {

	dbConf, err := configs.GetDatabaseConfig()
	if err != nil {
		return nil, nil, nil, err
	}
	db, err := storage.MySqlConnection(dbConf)
	if err != nil {
		return nil, nil, nil, err
	}

	redisConf, err := configs.GetRedisConfig()
	if err != nil {
		return nil, nil, nil, err
	}
	redisConnection := cache.RedisConnection(redisConf)

	rabbitMQConf, err := configs.GetRabbitMQConfig()
	if err != nil {
		return nil, nil, nil, err
	}
	rabbitConnection, err := rabbitMQ.NewConnection(rabbitMQ.ConnectionSettings{
		Username: rabbitMQConf.RabbitMQUser,
		Password: rabbitMQConf.RabbitMQPassword,
		Vhost:    rabbitMQConf.RabbitMQVhost,
		Connection: rabbitMQ.Settings{
			Secure:   rabbitMQConf.RabbitMQSecure,
			Hostname: rabbitMQConf.RabbitMQHostname,
			Port:     rabbitMQConf.RabbitMQPort,
		},
	})
	if err != nil {
		return nil, nil, nil, err
	}
	rabbitFormatter := *rabbitMQ.NewQueueFormatter("domain_event")

	var (
		commandBus = command.NewCommandBus()
		eventBus   = rabbitMQ.NewEventBus(rabbitConnection, rabbitMQConf.RabbitMQExchange, rabbitFormatter, rabbitMQConf.RabbitMQMaxRetries)
		queryBus   = query.NewQueryBus()
	)
	subscribers := []event.Subscriber{}

	if err := providers.AuthContainer(db, commandBus, queryBus, eventBus, &subscribers, redisConnection); err != nil {
		return nil, nil, nil, err
	}
	providers.CoursesContainer(db, commandBus, queryBus, eventBus, &subscribers)

	subscribersCopy := make([]event.Subscriber, len(subscribers))
	copy(subscribersCopy, subscribers)

	configurator := rabbitMQ.NewConfigurator(rabbitConnection, rabbitFormatter, rabbitMQConf.RabbitMQRetryTtl)
	err = configurator.Configure(rabbitMQConf.RabbitMQExchange, subscribers)

	if err != nil {
		return nil, nil, nil, err
	}
	for _, subscriber := range subscribers {
		eventBus.Subscribe(subscriber)
	}
	return commandBus, queryBus, redisConnection, nil
}
