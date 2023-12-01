package main

import (
	"context"
	privateapi "labs/htmx-blog/api/private"
	publicapi "labs/htmx-blog/api/public"
	"labs/htmx-blog/internal/app"
	"labs/htmx-blog/internal/configs"
	custdb "labs/htmx-blog/internal/db"
	"labs/htmx-blog/internal/events"
	custhttp "labs/htmx-blog/internal/http"
	custmqtt "labs/htmx-blog/internal/mqtt"
	"time"

	"go.uber.org/zap"
)

func main() {
	app.Run(
		time.Second*10,
		func(configs *configs.Configs, logger *zap.Logger) []app.Optioner {
			return []app.Optioner{
				app.WithHttpServer(custhttp.New(
					custhttp.WithGlobalConfigs(&configs.Public),
					custhttp.WithErrorHandler(publicapi.GlobalErrorHandler()),
					custhttp.WithRegistration(publicapi.ServiceRegistration()),
					custhttp.WithMiddleware(publicapi.Middlewares(&configs.Public)...),
				)),
				app.WithHttpServer(custhttp.New(
					custhttp.WithGlobalConfigs(&configs.Private),
					custhttp.WithErrorHandler(privateapi.GlobalErrorHandler()),
					custhttp.WithRegistration(privateapi.ServiceRegistration()),
					custhttp.WithMiddleware(privateapi.Middlewares(&configs.Private)...),
				)),
				app.WithNatsServer(events.New(
					events.WithGlobalConfigs(&configs.EventStore),
					events.WithZapLogger(logger.Sugar()),
				)),
				app.WithMqttServer(custmqtt.New(
					custmqtt.WithGlobalConfigs(&configs.MqttStore),
					custmqtt.WithZapLogger(logger),
				)),
				app.WithFactoryHook(func() error {
					custdb.Init(
						context.Background(),
						custdb.WithGlobalConfigs(&configs.Sqlite),
					)
					return nil
				}),
				app.WithShutdownHook(func(ctx context.Context) {
					custdb.Stop(ctx)
				}),
			}
		},
	)
}
