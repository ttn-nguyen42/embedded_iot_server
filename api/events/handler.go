package eventsapi

import (
	custcon "labs/htmx-blog/internal/concurrent"

	"github.com/panjf2000/ants/v2"
)

type RoomEventsHandler struct {
	pool *ants.Pool
}

func NewRoomEventsPool() *RoomEventsHandler {
	return &RoomEventsHandler{
		pool: custcon.New(100),
	}
}

func (h *RoomEventsHandler) Handle() error {
	return nil
}
