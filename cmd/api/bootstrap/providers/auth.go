package providers

import (
	"app-services-go/configs"
	"app-services-go/internal/auth/application/creating"
	"app-services-go/internal/auth/application/sending_email_validation"
	"app-services-go/internal/auth/infrastructure/storage/mysql"
	"app-services-go/kit/command"
	"app-services-go/kit/email/smtp"
	"app-services-go/kit/event"
	"app-services-go/kit/event/rabbitMQ"
	"app-services-go/kit/query"
	"database/sql"
)

func AuthContainer(db *sql.DB, commandBus *command.CommandBus, queryBus *query.QueryBus, eventBus *rabbitMQ.EventBus, subscribers *[]event.Subscriber) {
	dbConf, _ := configs.GetDatabaseConfig()

	userRepository := mysql.NewUserRepository(db, dbConf)
	creatingUserService := creating.NewUserService(userRepository, eventBus)
	emailSender := smtp.NewSender()
	sendVerificationEmailService := sending_email_validation.NewEmailValidatorSenderService(emailSender)

	createUserCommandHandler := creating.NewUserCommandHandler(creatingUserService)
	commandBus.Register(creating.AuthCommandType, createUserCommandHandler)

	*subscribers = append(*subscribers, creating.NewSendEmailVerificationOnUserCreated(sendVerificationEmailService))

}
