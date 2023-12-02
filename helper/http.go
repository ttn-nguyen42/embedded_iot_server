package helper

import "labs/htmx-blog/biz/models"

func GetPageAndLimit(common *models.ListCommon) (page uint64, limit uint64) {
	page = 0
	if common.Page >= 0 {
		page = common.Page
	}
	limit = 10
	if common.Limit >= 10 {
		limit = common.Limit
	}
	return
}
