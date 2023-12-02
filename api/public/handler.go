package publicapi

import (
	"labs/htmx-blog/biz/models"
	"labs/htmx-blog/biz/service"
	"labs/htmx-blog/helper"
	custerror "labs/htmx-blog/internal/error"
	"labs/htmx-blog/internal/logger"
	"time"

	"github.com/bytedance/sonic"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

func ApiRoomsHandler(ctx *fiber.Ctx) error {
	reqModel := &models.GetRoomsRequest{}
	reqModel.Page, reqModel.Limit = helper.GetPageAndLimitFromCtx(ctx)

	logger.SInfo("ApiRoomsHandler", zap.Any("req", reqModel))

	resp, err := service.GetRoomService().GetRooms(ctx.Context(), reqModel)
	if err != nil {
		return err
	}

	ctx.JSON(resp)
	return nil
}

func ApiAddRoomHandler(ctx *fiber.Ctx) error {
	reqModel := &models.CreateRoomRequest{}
	if err := sonic.Unmarshal(ctx.Request().Body(), &reqModel); err != nil {
		logger.SInfo("ApiAddRoomHandler: Unmarshal error", zap.Error(err))
		return custerror.ErrorInvalidArgument
	}

	if err := validateAddRoomRequest(reqModel); err != nil {
		logger.SInfo("ApiAddRoomHandler: validation error", zap.Error(err))
		return err
	}

	logger.SInfo("ApiAddRoomHandler", zap.Any("req", reqModel))

	resp, err := service.GetRoomService().AddRoom(ctx.Context(), reqModel)
	if err != nil {
		return err
	}

	ctx.JSON(resp)
	return nil
}

func validateAddRoomRequest(req *models.CreateRoomRequest) error {
	if req.Name == "" {
		return custerror.FormatInvalidArgument("name is empty")
	}
	return nil
}

func UiDashboard(ctx *fiber.Ctx) error {

	return ctx.Render("index", fiber.Map{
		"PartialTablePath":            "/partials/room_table",
		"PartialTableRefreshInterval": "5s",
	})
}

func UiPartialTable(ctx *fiber.Ctx) error {
	reqModel := &models.GetRoomsRequest{}
	reqModel.Page, reqModel.Limit = helper.GetPageAndLimitFromCtx(ctx)

	logger.SInfo("UiPartialTable", zap.Any("req", reqModel))

	resp, err := service.GetRoomService().GetRooms(ctx.Context(), reqModel)
	if err != nil {
		return err
	}

	modifiedRooms := []models.Room{}
	for _, r := range resp.Rooms {
		t, _ := time.Parse(time.RFC3339, r.UpdatedAt)
		modifiedUpdateTime := t.Format(time.RFC822)
		modifiedRoom := r
		modifiedRoom.UpdatedAt = modifiedUpdateTime
		modifiedRooms = append(modifiedRooms, modifiedRoom)
	}

	return ctx.Render("partials/room_table", fiber.Map{
		"Rooms": modifiedRooms,
	})
}
