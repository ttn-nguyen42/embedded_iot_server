package eventsapi

import (
	"context"
	"sync"
)

var once sync.Once

var roomEventsHandler *RoomEventsHandler

func Init(ctx context.Context) {
	once.Do(func() {
		roomEventsHandler = NewRoomEventsPool()
	})
}

func GetRoomEventsHandler() *RoomEventsHandler {
	return roomEventsHandler
}
