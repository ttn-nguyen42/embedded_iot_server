package app

import (
	"context"
	"errors"
	"labs/htmx-blog/internal/configs"
	"labs/htmx-blog/internal/events"
	custhttp "labs/htmx-blog/internal/http"
	"labs/htmx-blog/internal/logger"
	custmqtt "labs/htmx-blog/internal/mqtt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"

	"go.uber.org/zap"
)

func Run(shutdownTimeout time.Duration, registration RegistrationFunc) {
	ctx := context.Background()
	configs.Init(ctx)

	globalConfigs := configs.Get()

	loggerConfigs := globalConfigs.Logger
	logger.Init(ctx, logger.WithGlobalConfigs(&loggerConfigs))

	options := registration(globalConfigs, logger.Logger())

	opts := Options{}
	for _, optioner := range options {
		optioner(&opts)
	}

	logger := zap.L().Sugar()

	logger.Infof("Run: configs = %s", globalConfigs.String())

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	for _, s := range opts.httpServers {
		s := s
		go func() {
			logger.Infof("Run: start HTTP server name = %s", s.Name())
			if err := s.Start(); err != nil && !errors.Is(err, http.ErrServerClosed) {
				logger.Infof("Run: start HTTP server err = %s", err)
			}
		}()
	}

	for _, s := range opts.natsServers {
		s := s
		go func() {
			logger.Infof("Run: start embedded NATS server name = %s", s.Name())
			if err := s.Start(); err != nil {
				logger.Infof("Run: start embedded NATS server err = %s", err)
			}
		}()
	}

	for _, s := range opts.mqttServers {
		s := s
		go func() {
			logger.Infof("Run: start embedded MQTT server name = %s", s.Name())
			if err := s.Start(); err != nil {
				logger.Infof("Run: start embedded MQTT server err = %s", err)
			}
		}()
	}

	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	var wg sync.WaitGroup

	wg.Add(3)

	go func() {
		defer wg.Done()
		for _, s := range opts.httpServers {
			s := s
			logger.Infof("Run: stop HTTP server name = %s", s.Name())
			if err := s.Stop(ctx); err != nil {
				log.Fatal(err)
			}
		}
	}()

	go func() {
		for _, s := range opts.natsServers {
			s := s
			logger.Infof("Run: stop NATS embedded server name = %s", s.Name())
			if err := s.Stop(ctx); err != nil {
				log.Fatal(err)
			}
		}
	}()

	go func() {
		for _, s := range opts.mqttServers {
			s := s
			logger.Infof("Run: stop MQTT embedded server name = %s", s.Name())
			if err := s.Stop(ctx); err != nil {
				log.Fatal(err)
			}
		}
	}()

	wg.Wait()

	zap.L().Sync()
	log.Print("Run: shutdown complete")
}

type RegistrationFunc func(configs *configs.Configs, logger *zap.Logger) []Optioner

type Options struct {
	httpServers []*custhttp.HttpServer
	natsServers []*events.EmbeddedNats
	mqttServers []*custmqtt.EmbeddedMqtt
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

func WithMqttServer(server *custmqtt.EmbeddedMqtt) Optioner {
	return func(opts *Options) {
		if server != nil {
			opts.mqttServers = append(opts.mqttServers, server)
		}
	}
}
