package app

import (
	"context"
	"errors"
	"labs/htmx-blog/internal/configs"
	"labs/htmx-blog/internal/events"
	custhttp "labs/htmx-blog/internal/http"
	"labs/htmx-blog/internal/logger"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func Run(shutdownTimeout time.Duration, registration RegistrationFunc) {
	ctx := context.Background()

	log.Print("Run: initializing configurations")
	configs.Init(ctx)

	globalConfigs := configs.Get()

	loggerConfigs := globalConfigs.Logger

	options := registration(globalConfigs)

	opts := Options{}
	for _, optioner := range options {
		optioner(&opts)
	}

	log.Print("Run: initializing loggers")
	logger.Init(ctx, logger.WithGlobalConfigs(&loggerConfigs))

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	for _, s := range opts.httpServers {
		s := s
		go func() {
			log.Printf("Run: start HTTP server name = %s", s.Name())
			if err := s.Start(); err != nil && !errors.Is(err, http.ErrServerClosed) {
				log.Printf("Run: start HTTP server err = %s", err)
			}
		}()
	}

	for _, s := range opts.natsServers {
		s := s
		go func() {
			log.Printf("Run: start embedded NATS server name = %s", s.Name())
			if err := s.Start(); err != nil {
				log.Printf("Run: start embedded NATS server err = %s", err)
			}
		}()
	}

	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	for _, s := range opts.httpServers {
		s := s
		log.Printf("Run: stop HTTP server name = %s", s.Name())
		if err := s.Stop(ctx); err != nil {
			log.Fatal(err)
		}
	}

	for _, s := range opts.natsServers {
		s := s
		log.Printf("Run: stop NATS embedded server name = %s", s.Name())
		if err := s.Stop(ctx); err != nil {
			log.Fatal(err)
		}
	}

	log.Print("Run: shutdown complete")
}

type RegistrationFunc func(configs *configs.Configs) []Optioner

type Options struct {
	httpServers []*custhttp.HttpServer
	natsServers []*events.EmbeddedNats
}

type Optioner func(opts *Options)

func WithHttpServer(server *custhttp.HttpServer) Optioner {
	return func(opts *Options) {
		if server != nil {
			opts.httpServers = append(opts.httpServers, server)
		}
	}
}

func WithNatsServer(server *events.EmbeddedNats) Optioner {
	return func(opts *Options) {
		if server != nil {
			opts.natsServers = append(opts.natsServers, server)
		}
	}
}
