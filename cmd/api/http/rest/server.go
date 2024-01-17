package rest

import (
	"app-services-go/cmd/api/http/rest/controllers/courses"
	"app-services-go/cmd/api/http/rest/controllers/health"
	"app-services-go/cmd/api/http/rest/middleware/cache"
	"app-services-go/cmd/api/http/rest/middleware/logging"
	"app-services-go/cmd/api/http/rest/middleware/recovery"
	"app-services-go/configs"
	"app-services-go/kit/command"
	"app-services-go/kit/query"
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gin-gonic/gin"
)

type Server struct {
	httpAddr string
	engine   *gin.Engine

	shutdownTimeout time.Duration
	commandBus      command.Bus
	queryBus        query.Bus
	redisConnection redis.UniversalClient
}

func New(ctx context.Context, config configs.ServerConfig, commandBus command.Bus, queryBus query.Bus, redisConnection redis.UniversalClient) (context.Context, Server) {
	gin.SetMode(config.Mode)

	srv := Server{
		engine:   gin.New(),
		httpAddr: fmt.Sprintf("%s:%d", config.Host, config.Port),

		shutdownTimeout: config.ShutdownTimeout,

		commandBus:      commandBus,
		queryBus:        queryBus,
		redisConnection: redisConnection,
	}

	srv.registerRoutes()
	return serverContext(ctx), srv
}

func (s *Server) registerRoutes() {
	s.engine.Use(
		recovery.Middleware(),
		logging.Middleware(),
	)
	s.engine.GET("/health", cache.Middleware(s.redisConnection, 5*time.Second), health.CheckHandler())
	g := s.engine.Group("/course")
	{
		g.POST("/", courses.CreateHandler(s.commandBus))
		g.GET("/:id", courses.RetrieveHandler(s.queryBus))
	}
}

func (s *Server) Run(ctx context.Context) error {
	log.Println("Server running on", s.httpAddr)

	srv := &http.Server{
		Addr:    s.httpAddr,
		Handler: s.engine,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("rest shut down", err)
		}
	}()

	<-ctx.Done()
	ctxShutDown, cancel := context.WithTimeout(context.Background(), s.shutdownTimeout)
	defer cancel()

	return srv.Shutdown(ctxShutDown)
}

func serverContext(ctx context.Context) context.Context {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	ctx, cancel := context.WithCancel(ctx)
	go func() {
		<-c
		cancel()
	}()

	return ctx
}
