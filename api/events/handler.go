package eventsapi

import (
	"context"
	"labs/htmx-blog/biz/models"
	"labs/htmx-blog/biz/service"
	"labs/htmx-blog/helper"
	"labs/htmx-blog/internal/cache"
	custcon "labs/htmx-blog/internal/concurrent"
	custerror "labs/htmx-blog/internal/error"
	"labs/htmx-blog/internal/logger"
	"strconv"
	"time"

	"github.com/bytedance/sonic"
	"github.com/dgraph-io/ristretto"
	"github.com/eclipse/paho.golang/paho"
	"github.com/panjf2000/ants/v2"
	"go.uber.org/zap"
)

type RoomEventsHandler struct {
	pool  *ants.Pool
	cache *ristretto.Cache
}

func NewRoomEventsPool() *RoomEventsHandler {
	return &RoomEventsHandler{
		pool:  custcon.New(100),
		cache: cache.Cache(),
	}
}

func (h *RoomEventsHandler) Handle(p *paho.Publish) error {
	topicParts := helper.ParseTopic(p.Topic)

	parsedId, err := strconv.Atoi(topicParts[1])
	if err != nil {
		return custerror.FormatInvalidArgument("RoomEventsHandler.Handle: id not integer")
	}

	if parsedId < 0 {
		return custerror.FormatInvalidArgument("RoomEventsHandler.Handle: id invalid")
	}

	msg := p.Payload

	var msgModel models.RoomStatusChanged

	if err := sonic.Unmarshal(msg, &msgModel); err != nil {
		return custerror.FormatInvalidArgument("RoomEventsHandler.Unmarshal: err = %s", err)
	}

	logger.SInfo("message received",
		zap.String("where", "eventsapi.RoomEventsHandler"),
		zap.String("topic", topicParts[0]),
		zap.String("clientId", topicParts[1]),
		logger.Json("message", msgModel),
	)

	currentStatus, ok := h.cache.Get(parsedId)
	if ok {
		if currentStatus == msgModel.Status {
			logger.SInfo("read from cache",
				zap.String("status", msgModel.Status))
			return nil
		}
	}

	h.pool.Submit(func() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()

		err = service.GetRoomEventService().UpdateRoomStatus(ctx, &models.RoomEventUpdateRequest{
			Status:    msgModel.Status,
			Timestamp: msgModel.Timestamp,
			Id:        uint32(parsedId),
		})

		if err != nil {
			logger.SWarn("UpdateRoomStatus: failed",
				zap.Error(err))
			return
		}

		ok := h.cache.Set(parsedId, msgModel.Status, 0)
		if ok {
			logger.SWarn("UpdateRoomStatus: cache set success")
		} else {
			logger.SWarn("UpdateRoomStatus: cache set failed")
		}
	})

	return nil
}
