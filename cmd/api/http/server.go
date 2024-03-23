package http

import (
	"app-services-go/configs"
	"context"
	"fmt"
	"log"
	"net/http"
	"time"
    "os"
	"os/signal"

	"github.com/gin-gonic/gin"
)

type Server struct {
	httpAddr string
	Engine   *gin.Engine

	readTimeout time.Duration
	writeTimeout time.Duration
	shutdownTimeout time.Duration
}

func New(ctx context.Context, config configs.ServerConfig) (context.Context, Server) {
	gin.SetMode(config.Mode)

	srv := Server{
		Engine:   gin.New(),
		httpAddr: fmt.Sprintf("%s:%d", config.Host, config.Port),

	    readTimeout:  5 * time.Second,
		writeTimeout: 10 * time.Second,

		shutdownTimeout: config.ShutdownTimeout,
	}

	return serverContext(ctx), srv
}



func (s *Server) Run(ctx context.Context) error {
	log.Println("Server running on", s.httpAddr)

	srv := &http.Server{
		Addr:    s.httpAddr,
        ReadTimeout:  s.readTimeout,
        WriteTimeout: s.writeTimeout,
		Handler: s.Engine,
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
