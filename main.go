package main

import (
	"context"
	eventsapi "labs/htmx-blog/api/events"
	privateapi "labs/htmx-blog/api/private"
	publicapi "labs/htmx-blog/api/public"
	"labs/htmx-blog/biz/models"
	"labs/htmx-blog/internal/app"
	"labs/htmx-blog/internal/cache"
	"labs/htmx-blog/internal/configs"
	custdb "labs/htmx-blog/internal/db"
	"labs/htmx-blog/internal/events"
	custhttp "labs/htmx-blog/internal/http"
	"labs/htmx-blog/internal/logger"
	custmqtt "labs/htmx-blog/internal/mqtt"
	"time"

	"go.uber.org/zap"
)

func main() {
	app.Run(
		time.Second*10,
		func(configs *configs.Configs, zl *zap.Logger) []app.Optioner {
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
					events.WithZapLogger(zl.Sugar()),
				)),
				app.WithMqttServer(custmqtt.New(
					custmqtt.WithGlobalConfigs(&configs.MqttStore),
					custmqtt.WithZapLogger(zl),
				)),
				app.WithFactoryHook(func() error {
					ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
					defer cancel()

					custdb.Init(
						context.Background(),
						custdb.WithGlobalConfigs(&configs.Sqlite),
					)

					custdb.Migrate(&models.Room{})

					cache.Init()
					custdb.LayeredInit()

					eventsapi.Init(ctx)

					custmqtt.InitClient(
						context.Background(),
						custmqtt.WithClientGlobalConfigs(&configs.MqttStore),
						custmqtt.WithOnReconnection(eventsapi.Register),
						custmqtt.WithOnConnectError(func(err error) {
							logger.Error("MQTT Connection failed", zap.Error(err))
						}),
						custmqtt.WithClientError(eventsapi.ClientErrorHandler),
						custmqtt.WithOnServerDisconnect(eventsapi.DisconnectHandler),
						custmqtt.WithHandlerRegister(eventsapi.RouterHandler()),
					)
					return nil
				}),
				app.WithShutdownHook(func(ctx context.Context) {
					custdb.Stop(ctx)
					custmqtt.StopClient(ctx)
					logger.Close()
				}),
			}
		},
	)
}
