package rest

import (
	"app-services-go/cmd/api/http/rest/controllers/auth_web2"
	"app-services-go/cmd/api/http/rest/controllers/auth_web3"
	"app-services-go/cmd/api/http/rest/controllers/courses"
	"app-services-go/cmd/api/http/rest/controllers/health"
	"app-services-go/cmd/api/http/rest/middleware/cache"
	"app-services-go/cmd/api/http/rest/middleware/logging"
	"app-services-go/cmd/api/http/rest/middleware/recovery"
	"app-services-go/internal/auth/application/message_generator"
	"app-services-go/kit/command"
	"app-services-go/kit/query"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(s *gin.Engine, commandBus command.Bus, queryBus query.Bus, redisConnection redis.UniversalClient) {
	s.Use(
		recovery.Middleware(),
		logging.Middleware(),
	)
	api := s.Group("/api")
	{
		api.GET("/health", cache.Middleware(redisConnection, 5*time.Second), health.CheckHandler())
		coursesApi := api.Group("/course")
		{
			coursesApi.POST("/", courses.CreateHandler(commandBus))
			coursesApi.GET("/:id", courses.RetrieveHandler(queryBus))
		}
		authApi := api.Group("/auth")
		{
			authApi.POST("/user", auth_web2.CreateHandler(commandBus))
			authApi.GET("/validate/:token", auth_web2.ValidateHandler(commandBus))
			authApi.POST("/login", auth_web2.LoginHandler(commandBus, queryBus))
			authApi.POST("/forgot", auth_web2.ForgotHandler(commandBus, queryBus))
			authApi.POST("/remember/:token", auth_web2.RememberHandler(commandBus))
			authApi.GET("/message/:address", auth_web3.MessageHandler(commandBus, message_generator.NewMessageGeneratorService()))
			authApi.POST("/validate/:address", auth_web3.ValidateMessageHandler(commandBus, queryBus))
		}
	}

}
