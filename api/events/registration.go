package eventsapi

import (
	"labs/htmx-blog/internal/logger"

	"github.com/eclipse/paho.golang/autopaho"
	"github.com/eclipse/paho.golang/paho"
	"go.uber.org/zap"
)

func Register(cm *autopaho.ConnectionManager, connack *paho.Connack) {

}

func GlobalErrorHandler(err error) {
	logger := logger.Logger()

	logger.Error("MQTT Client", zap.Error(err))
}

func DisconnectHandler(d *paho.Disconnect) {
	logger := logger.Logger()

	logger.Error("MQTT Server Disconnect", zap.String("reason", d.Properties.ReasonString))
}
