package service

import (
	"context"
	"labs/htmx-blog/biz/models"
	"labs/htmx-blog/helper"
	"time"

	custdb "labs/htmx-blog/internal/db"
	"labs/htmx-blog/internal/logger"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"go.uber.org/zap"
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
	m := []models.Room{}
	query := sq.Select("*").From("rooms")

	page, limit := helper.GetPageAndLimit(&req.ListCommon)
	offset := page * limit
	query = query.Limit(limit).Offset(offset)

	logger.SDebug("RoomService.GetRooms",
		zap.Uint64("page", page),
		zap.Uint64("limit", limit),
		zap.Uint64("offset", offset))

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
	m := models.Room{
		Name:   req.Name,
		Status: models.RoomStatus_Empty,
	}

	id, _ := uuid.NewUUID()

	query := sq.Insert("rooms").Columns(
		"id",
		"created_at",
		"updated_at",
		"name",
		"status",
	).Values(
		id.ID(),
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
		query = query.Set("name", req.Name)
	}

	if req.Status != "" {
		query = query.Set("status", req.Status)
	}

	if !req.Timestamp.IsZero() {
		query = query.Set("updated_at", req.Timestamp.Format(time.RFC3339))
	}

	if err := s.db.Update(ctx, query); err != nil {
		return nil, err
	}

	return &models.UpdateRoomResponse{
		Id: req.Id,
	}, nil
}
