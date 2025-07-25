package logger

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

var Log *zap.Logger

func InitLogger() {
	var err error
	Log, err = zap.NewProduction()

	if err != nil {
		panic(err)
	}
}

func WarnError(c *gin.Context, message string, err error) {
	Log.Warn(message,
		zap.Error(err),
		zap.String("path", c.FullPath()),
		zap.String("method", c.Request.Method),
	)
}

func WarnFields(c *gin.Context, message string, fields ...zap.Field) {
	base := []zap.Field{
		zap.String("path", c.FullPath()),
		zap.String("method", c.Request.Method),
	}

	Log.Warn(message, append(base, fields...)...)
}
