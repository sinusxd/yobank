package route

import (
	"time"

	"github.com/gin-gonic/gin"
	"yobank/bootstrap"
)

func Setup(app bootstrap.Application, timeout time.Duration, gin *gin.Engine) {
	container := bootstrap.BuildContainer(app)
	publicRouter := gin.Group("/api/v1")
	// All Public APIs
	NewEmailLoginRouter(container.Services.Login, container.Services.EmailCode, container.Repos.User, app.Env, publicRouter)
	NewTelegramLoginRouter(container.Repos.User, container.Services.Login, app.Env, publicRouter)
}
