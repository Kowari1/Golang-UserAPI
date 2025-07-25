package middleware

import (
	"net/http"
	"userapi/internal/logger"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func ErrorRecovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if rec := recover(); rec != nil {
				logger.Log.Error("panic recovered", zap.Any("panic", rec))
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
			}
		}()
		c.Next()
	}
}
