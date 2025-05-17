package middleware

import (
	"github.com/gin-gonic/gin"
	initdata "github.com/telegram-mini-apps/init-data-golang"
	"strings"
	"time"
	"yobank/internal/telegram"
)

// Middleware: авторизация по Telegram Mini App initData
func TelegramAuthMiddleware(botToken string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authParts := strings.Split(c.GetHeader("Authorization"), " ")
		if len(authParts) != 2 || strings.ToLower(authParts[0]) != "tma" {
			c.AbortWithStatusJSON(401, gin.H{"message": "Unauthorized"})
			return
		}

		authData := authParts[1]

		// Проверка подписи (валидна 1 час)
		if err := initdata.Validate(authData, botToken, time.Hour); err != nil {
			c.AbortWithStatusJSON(401, gin.H{"message": "Invalid signature: " + err.Error()})
			return
		}

		// Парсим initData
		parsed, err := initdata.Parse(authData)
		if err != nil {
			c.AbortWithStatusJSON(400, gin.H{"message": "Parse error: " + err.Error()})
			return
		}

		c.Request = c.Request.WithContext(telegram.WithInitData(c.Request.Context(), parsed))
		c.Next()
	}
}
