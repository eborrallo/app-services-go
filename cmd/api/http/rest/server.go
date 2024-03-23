package rest

import (
	"app-services-go/cmd/api/http/rest/controllers/courses"
	"app-services-go/cmd/api/http/rest/controllers/health"
	"app-services-go/cmd/api/http/rest/middleware/cache"
	"app-services-go/cmd/api/http/rest/middleware/logging"
	"app-services-go/cmd/api/http/rest/middleware/recovery"
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
		g := api.Group("/course")
		{
			g.POST("/", courses.CreateHandler(commandBus))
			g.GET("/:id", courses.RetrieveHandler(queryBus))
		}
	}

}


