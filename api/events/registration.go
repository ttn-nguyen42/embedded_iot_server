package eventsapi

import (
	"context"
	"labs/htmx-blog/helper"
	"labs/htmx-blog/internal/logger"
	custmqtt "labs/htmx-blog/internal/mqtt"
	"time"

	"github.com/eclipse/paho.golang/autopaho"
	"github.com/eclipse/paho.golang/paho"
	"go.uber.org/zap"
)

func Register(cm *autopaho.ConnectionManager, connack *paho.Connack) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	subs := makeSubsciptions(ctx, cm, connack)
	if _, err := cm.Subscribe(ctx, &paho.Subscribe{
		Subscriptions: subs,
	}); err != nil {
		logger.SError("unable to make MQTT subscriptions",
			zap.String("where", "api.events.Register"),
			zap.Any("subs", subs),
		)
		return
	}

	logger.SInfo("MQTT subscriptions made success", zap.Any("subs", subs))
}

func makeSubsciptions(ctx context.Context, cm *autopaho.ConnectionManager, connack *paho.Connack) []paho.SubscribeOptions {
	return []paho.SubscribeOptions{
		{
			Topic: "room_events/#",
			QoS:   1, // at least once
		},
	}
}

func ClientErrorHandler(err error) {
	logger := logger.Logger()

	logger.Error("MQTT Client", zap.Error(err))
}

func DisconnectHandler(d *paho.Disconnect) {
	logger := logger.Logger()

	logger.Error("MQTT Server Disconnect", zap.String("reason", d.Properties.ReasonString))
}

func RouterHandler() custmqtt.RouterRegister {
	return func(router *paho.StandardRouter) {
		roomEventsHandler := GetRoomEventsHandler()
		router.RegisterHandler(
			"room_events/#",
			paho.MessageHandler(WrapForHandlers(roomEventsHandler.Handle)),
		)
	}
}

func WrapForHandlers(handler func(p *paho.Publish) error) func(p *paho.Publish) {
	return func(p *paho.Publish) {
		if err := handler(p); err != nil {
			helper.EventHandlerErrorHandler(err)
		}
	}
}
