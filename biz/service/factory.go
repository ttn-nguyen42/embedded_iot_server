package service

import "sync"

var once sync.Once

var (
	roomService      *RoomService
	roomEventService *RoomEventService
)

func Init() {
	once.Do(func() {
		roomService = NewRoomService()
		roomEventService = NewRoomEventsService()
	})
}

func GetRoomService() *RoomService {
	return roomService
}

func GetRoomEventService() *RoomEventService {
	return roomEventService
}
