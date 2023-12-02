package service

import (
	"context"
	custdb "labs/htmx-blog/internal/db"
)

type RoomEventService struct {
	db          *custdb.LayeredDb
	roomService *RoomService
}

func NewRoomEventsService() *RoomEventService {
	return &RoomEventService{
		db:          custdb.Layered(),
		roomService: GetRoomService(),
	}
}

func (s *RoomEventService) UpdateRoomStatus(ctx context.Context) error {
	return nil
}
