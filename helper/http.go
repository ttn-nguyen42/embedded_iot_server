package helper

import (
	"github.com/gofiber/fiber/v2"
	"labs/htmx-blog/biz/models"
)

func GetPageAndLimit(common *models.ListCommon) (page uint64, limit uint64) {
	page = 0
	if common.Page > 0 {
		page = common.Page
	}
	limit = 10
	if common.Limit > 0 {
		limit = common.Limit
	}
	return
}

func GetPageAndLimitFromCtx(ctx *fiber.Ctx) (page uint64, limit uint64) {
	pageQuery := ctx.QueryInt("page")
	limitQuery := ctx.QueryInt("limit")

	if pageQuery >= 0 {
		page = uint64(pageQuery)
	}

	if limitQuery >= 0 {
		limit = uint64(limitQuery)
	}

	return page, limit
}
