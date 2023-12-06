package service

import (
	"context"
	"errors"
	"labs/htmx-blog/biz/models"
	custdb "labs/htmx-blog/internal/db"
	custerror "labs/htmx-blog/internal/error"
	"labs/htmx-blog/internal/logger"
	"time"

	"go.uber.org/zap"
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

type UpdateRoomEventTracker struct {
	req       *models.RoomEventUpdateRequest
	timestamp time.Time
	id        uint32
	ctx       context.Context
}

func (s *RoomEventService) UpdateRoomStatus(ctx context.Context, req *models.RoomEventUpdateRequest) error {
	logger.SInfo("UpdateRoomStatus", logger.Json("req", req))

	tracker := &UpdateRoomEventTracker{
		req: req,
		ctx: ctx,
	}

	if err := s.validateUpdateEvent(tracker); err != nil {
		return err
	}

	if err := s.updateRoomStatus(tracker); err != nil {
		return err
	}

	return nil
}

func (s *RoomEventService) validateUpdateEvent(tracker *UpdateRoomEventTracker) error {
	req := tracker.req

	switch req.Status {
	case models.RoomStatus_Empty:
	case models.RoomStatus_Occupied:
	default:
		logger.SWarn(
			"UpdateRoomStatus: status invalid",
			zap.String("status", req.Status),
		)
		return custerror.FormatInvalidArgument("UpdateRoomStatus: status invalid")
	}

	timestamp, err := time.Parse(time.RFC3339, req.Timestamp)
	if err != nil {
		logger.SWarn(
			"UpdateRoomStatus: timestamp format is not RFC3339",
			zap.String("timestamp", req.Timestamp),
		)
		return custerror.FormatInvalidArgument("UpdateRoomStatus: timestamp is not RFC3339")
	}

	tracker.timestamp = timestamp
	tracker.id = req.Id
	return nil
}

func (s *RoomEventService) updateRoomStatus(tracker *UpdateRoomEventTracker) error {
	req := tracker.req
	rs := s.roomService

	resp, err := rs.UpdateRoom(tracker.ctx, &models.UpdateRoomRequest{
		Timestamp: tracker.timestamp,
		Id:        tracker.id,
		Status:    req.Status,
	})

	if err != nil {
		if errors.Is(err, custerror.ErrorNotFound) {
			logger.SWarn("updateRoomStatus: room not found")
			return err
		}
		logger.SInfo("updateRoomStatus: UpdateRoom", zap.Error(err))
		return custerror.FormatInternalError("updateRoomStatus: UpdateRoom failed")
	}

	logger.SInfo("updateRoomStatus: success",
		zap.Uint32("id", resp.Id),
		zap.String("status", req.Status))

	return nil
}
