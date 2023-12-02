package helper

import (
	custerror "labs/htmx-blog/internal/error"
	"labs/htmx-blog/internal/logger"

	"go.uber.org/zap"
)

var commonEventMessage string = "events handler error"

func EventHandlerErrorHandler(err error) {
	custError, ok := err.(*custerror.CustomError)
	if ok {
		logger.SInfo(commonEventMessage,
			zap.Error(err),
			zap.Uint32("type", custError.Code))
	} else {
		logger.SInfo(commonEventMessage,
			zap.Error(err))
	}
}
