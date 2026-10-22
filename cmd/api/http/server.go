package http

import (
	"app-services-go/configs"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

type Server struct {
	httpAddr string
	Engine   *gin.Engine

	readTimeout     time.Duration
	writeTimeout    time.Duration
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

	return ctx, srv
}

func (s *Server) Run(parent context.Context) error {
	log.Println("Server running on", s.httpAddr)

	// signal.NotifyContext cancels ctx on SIGINT or SIGTERM, so Docker/K8s
	// stop signals (SIGTERM) trigger the same graceful shutdown as Ctrl-C.
	ctx, stop := signal.NotifyContext(parent, os.Interrupt, syscall.SIGTERM)
	defer stop()

	srv := &http.Server{
		Addr:         s.httpAddr,
		ReadTimeout:  s.readTimeout,
		WriteTimeout: s.writeTimeout,
		Handler:      s.Engine,
	}

	// Surface the listen error through a channel instead of log.Fatal inside the
	// goroutine (which would call os.Exit and skip every deferred cleanup).
	serverErr := make(chan error, 1)
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			serverErr <- err
		}
	}()

	select {
	case err := <-serverErr:
		return err
	case <-ctx.Done():
		// A shutdown signal was received; fall through to graceful shutdown.
	}

	ctxShutDown, cancel := context.WithTimeout(context.Background(), s.shutdownTimeout)
	defer cancel()

	return srv.Shutdown(ctxShutDown)
}
