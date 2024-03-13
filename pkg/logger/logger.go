package logger

import (
	"context"
	"geolocation-service/internal/model"
	"github.com/sirupsen/logrus"
)

func GetLogger(ctx context.Context) *logrus.Entry {
	logger, ok := ctx.Value(model.ContextLogger).(*logrus.Entry)
	if !ok {
		return logrus.NewEntry(logrus.New())
	}
	return logger
}
