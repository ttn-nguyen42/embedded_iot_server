package publicapi

import (
	"labs/htmx-blog/biz/models"
	"labs/htmx-blog/biz/service"
	"labs/htmx-blog/helper"
	custerror "labs/htmx-blog/internal/error"
	"labs/htmx-blog/internal/logger"

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
