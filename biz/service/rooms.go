package service

import (
	"context"
	"labs/htmx-blog/biz/models"
	"labs/htmx-blog/helper"
	"time"

	custdb "labs/htmx-blog/internal/db"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
)

type RoomService struct {
	db *custdb.LayeredDb
}

func NewRoomService() *RoomService {
	return &RoomService{
		db: custdb.Layered(),
	}
}

func (s *RoomService) GetRooms(ctx context.Context, req *models.GetRoomsRequest) (*models.GetRoomsResponse, error) {
	var m []models.Room
	query := sq.Select("*").From("rooms")

	page, limit := helper.GetPageAndLimit(&req.ListCommon)
	query = query.Limit(limit).Offset(page * limit)

	if err := s.db.Select(ctx, query, &m); err != nil {
		return nil, err
	}

	return &models.GetRoomsResponse{
		Rooms: m,
	}, nil
}

func (s *RoomService) GetRoom(ctx context.Context, req *models.GetRoomRequest) (*models.GetRoomResponse, error) {
	var m models.Room

	query := sq.Select("*").From("rooms").Where("id = ?", req.Id)

	if err := s.db.Select(ctx, query, &m); err != nil {
		return nil, err
	}

	return &models.GetRoomResponse{
		Room: &m,
	}, nil
}

func (s *RoomService) AddRoom(ctx context.Context, req *models.CreateRoomRequest) (*models.CreateRoomResponse, error) {
	var m models.Room
	m = models.Room{
		Name:   req.Name,
		Status: req.Status,
	}

	id, _ := uuid.NewUUID()

	query := sq.Insert("rooms").Columns(
		"id",
		"created_at",
		"updated_at",
		"name",
		"status",
	).Values(
		id,
		time.Now().Format(time.RFC3339),
		time.Now().Format(time.RFC3339),
		m.Name,
		m.Status,
	)

	if err := s.db.Insert(ctx, query); err != nil {
		return nil, err
	}

	return &models.CreateRoomResponse{
		Id: id.ID(),
	}, nil
}

func (s *RoomService) UpdateRoom(ctx context.Context, req *models.UpdateRoomRequest) (*models.UpdateRoomResponse, error) {
	query := sq.Update("rooms").Where("id = ?", req.Id)

	if req.Name != "" {
		query.Set("name", req.Name)
	}

	if req.Status != "" {
		query.Set("status", req.Status)
	}

	if err := s.db.Update(ctx, query); err != nil {
		return nil, err
	}

	return &models.UpdateRoomResponse{
		Id: req.Id,
	}, nil
}
