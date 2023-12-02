package eventsapi

import (
	"labs/htmx-blog/biz/models"
	"labs/htmx-blog/helper"
	custcon "labs/htmx-blog/internal/concurrent"
	custerror "labs/htmx-blog/internal/error"
	"labs/htmx-blog/internal/logger"

	"github.com/bytedance/sonic"
	"github.com/eclipse/paho.golang/paho"
	"github.com/panjf2000/ants/v2"
	"go.uber.org/zap"
)

type RoomEventsHandler struct {
	pool *ants.Pool
}

func NewRoomEventsPool() *RoomEventsHandler {
	return &RoomEventsHandler{
		pool: custcon.New(100),
	}
}

func (h *RoomEventsHandler) Handle(p *paho.Publish) error {
	topicParts := helper.ParseTopic(p.Topic)

	msg := p.Payload

	var msgModel models.Room
	if err := sonic.Unmarshal(msg, &msgModel); err != nil {
		return custerror.FormatInvalidArgument("RoomEventsHandler.Unmarshal: err = %s", err)
	}

	logger.SInfo("message received",
		zap.String("where", "eventsapi.RoomEventsHandler"),
		zap.String("topic", topicParts[0]),
		zap.String("clientId", topicParts[1]),
		logger.Json("message", msgModel),
	)

	return nil
}
