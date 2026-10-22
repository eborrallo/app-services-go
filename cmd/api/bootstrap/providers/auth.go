package providers

import (
	"app-services-go/configs"
	"app-services-go/internal/auth/application/creating"
	"app-services-go/internal/auth/application/fetching"
	"app-services-go/internal/auth/application/forgot"
	"app-services-go/internal/auth/application/login"
	"app-services-go/internal/auth/application/message_generator"
	"app-services-go/internal/auth/application/remember"
	"app-services-go/internal/auth/application/sending_email_validation"
	"app-services-go/internal/auth/application/validating"
	"app-services-go/internal/auth/infrastructure/storage/mysql"
	authredis "app-services-go/internal/auth/infrastructure/storage/redis"
	"app-services-go/kit/blokchain"
	"app-services-go/kit/cache"
	"app-services-go/kit/command"
	"app-services-go/kit/email/smtp"
	"app-services-go/kit/event"
	"app-services-go/kit/event/rabbitMQ"
	"app-services-go/kit/query"
	"database/sql"
	"time"

	"github.com/redis/go-redis/v9"
)

// signMessageTTL is how long a generated web3 sign-in message stays valid in
// Redis before the user must request a new one.
const signMessageTTL = 5 * time.Minute

func AuthContainer(db *sql.DB, commandBus *command.CommandBus, queryBus *query.QueryBus, eventBus *rabbitMQ.EventBus, subscribers *[]event.Subscriber, redisConnection redis.UniversalClient) error {
	dbConf, err := configs.GetDatabaseConfig()
	if err != nil {
		return err
	}

	userRepository := mysql.NewUserRepository(db, dbConf)
	userMessageRepository := authredis.NewUserMessageRepository(cache.NewRedisCache(redisConnection, signMessageTTL))

	ethClient, err := blokchain.RpcConection()
	if err != nil {
		return err
	}
	signatureVerificator := blokchain.NewSignatureVerificator(ethClient)

	emailSender := smtp.NewSender()

	// Application services.
	creatingUserService := creating.NewUserService(userRepository, eventBus)
	loginUserService := login.NewLoginUserService(userRepository, eventBus)
	validateUserService := validating.NewValidateUserService(userRepository, eventBus)
	rememberPasswordService := remember.NewPasswordUserService(userRepository, eventBus)
	saveMessageService := message_generator.NewSaveMessageService(userMessageRepository, eventBus)
	forgotUserService := forgot.NewForgotUserService(emailSender)
	sendVerificationEmailService := sending_email_validation.NewEmailValidatorSenderService(emailSender)

	// Commands.
	commandBus.Register(creating.AuthCommandType, creating.NewUserCommandHandler(creatingUserService))
	commandBus.Register(creating.Web3AuthCommandType, creating.NewWeb3UserCommandHandler(creatingUserService))
	commandBus.Register(login.AuthCommandType, login.NewLoginCommandHandler(loginUserService))
	commandBus.Register(validating.AuthCommandType, validating.NewValidateUserCommandHandler(validateUserService))
	commandBus.Register(remember.AuthCommandType, remember.NewPasswordUserCommandHandler(rememberPasswordService))
	commandBus.Register(message_generator.AuthCommandType, message_generator.NewSaveMessageCommandHandler(saveMessageService))
	commandBus.Register(forgot.AuthCommandType, forgot.NewForgotCommandHandler(forgotUserService))

	// Queries.
	queryBus.Register(fetching.UserByEmailQueryType, fetching.NewUserByEmailQueryHandler(userRepository))
	queryBus.Register(fetching.UserByIdQueryType, fetching.NewUserByIdQueryHandler(userRepository))
	queryBus.Register(fetching.UserBySignatureQueryType, fetching.NewUserBySignatureQueryHandler(userRepository, userMessageRepository, signatureVerificator))

	// Subscribers.
	*subscribers = append(*subscribers, creating.NewSendEmailVerificationOnUserCreated(sendVerificationEmailService))

	return nil
}
