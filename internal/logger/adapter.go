package logger

import (
	"github.com/nats-io/nats-server/v2/server"
	"go.uber.org/zap"
)

type ZapToNatsLogger struct {
	logger *zap.SugaredLogger
}

func (a *ZapToNatsLogger) Debugf(format string, v ...interface{}) {
	a.logger.Debugf(format, v...)
}

func (a *ZapToNatsLogger) Errorf(format string, v ...interface{}) {
	a.logger.Errorf(format, v...)
}

func (a *ZapToNatsLogger) Fatalf(format string, v ...interface{}) {
	a.logger.Fatalf(format, v...)
}

func (a *ZapToNatsLogger) Noticef(format string, v ...interface{}) {
	a.logger.Warnf(format, v...)
}

func (a *ZapToNatsLogger) Tracef(format string, v ...interface{}) {
	a.logger.Debugf(format, v...)
}

func (a *ZapToNatsLogger) Warnf(format string, v ...interface{}) {
	a.logger.Warnf(format, v...)
}

func NewZapToNatsLogger(logger *zap.SugaredLogger) server.Logger {
	return &ZapToNatsLogger{
		logger: logger,
	}
}
