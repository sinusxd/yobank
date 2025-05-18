package route

import (
	"time"
	"yobank/api/middleware"

	"github.com/gin-gonic/gin"
	"yobank/bootstrap"
)

func Setup(app bootstrap.Application, timeout time.Duration, gin *gin.Engine) {
	container := app.Container

	publicRouter := gin.Group("/api/v1")
	// All Public APIs
	NewEmailLoginRouter(container.Repos.User, container.Services.Login, container.Services.EmailCode, container.Services.Wallet, app.Env, publicRouter)
	NewTelegramLoginRouter(container.Services.User, container.Services.Login, app.Env, publicRouter)
	NewRateRouter(container.Services.Rate, publicRouter)

	protectedRouter := gin.Group("/api/v1")
	protectedRouter.Use(middleware.JwtAuthMiddleware(app.Env.AccessTokenSecret))
	NewWalletRouter(container.Services.Wallet, protectedRouter)
}
