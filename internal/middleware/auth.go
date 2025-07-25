package middleware

import (
	"net/http"
	"strings"

	"userapi/internal/handler"
	"userapi/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func JWTMiddleware(secret []byte, redis *service.RedisService) gin.HandlerFunc {
	return func(c *gin.Context) {

		authHeader := c.GetHeader("Authorization")

		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			handler.JSONErrorMsg(c, http.StatusUnauthorized, "missing or malformed token")

			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

		token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
			return secret, nil
		})

		if err != nil || !token.Valid {
			handler.JSONErrorMsg(c, http.StatusUnauthorized, "invalid or expired token")

			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)

		if !ok {
			handler.JSONErrorMsg(c, http.StatusUnauthorized, "invalid token claims")

			return
		}

		if jti, ok := claims["jti"].(string); ok {
			isBlacklisted, err := redis.IsBlacklisted(c.Request.Context(), jti)

			if err != nil {
				handler.JSONErrorMsg(c, http.StatusInternalServerError, "internal redis error")

				return
			}

			if isBlacklisted {
				handler.JSONErrorMsg(c, http.StatusUnauthorized, "token revoked")

				return
			}

			c.Set("jti", jti)
		} else {
			handler.JSONErrorMsg(c, http.StatusUnauthorized, "token missing jti")

			return
		}

		if userID, ok := claims["user_id"].(string); ok {
			c.Set("user_id", userID)
		}

		if login, ok := claims["login"].(string); ok {
			c.Set("login", login)
		}

		if role, ok := claims["role"].(bool); ok {
			c.Set("role", role)
		}

		c.Next()
	}
}

func RequireAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		raw, exists := c.Get("role")
		if !exists {
			handler.JSONErrorMsg(c, http.StatusForbidden, "admin only")
			return
		}

		isAdmin, ok := raw.(bool)
		if !ok || !isAdmin {
			handler.JSONErrorMsg(c, http.StatusForbidden, "admin only")
			return
		}

		c.Next()
	}
}
