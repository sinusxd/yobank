package route

import (
	"github.com/gin-gonic/gin"
	"yobank/api/controller"
	"yobank/bootstrap"
	"yobank/domain"
)

func NewTelegramLoginRouter(userRepo domain.UserRepository, loginService domain.LoginService, env *bootstrap.Env, group *gin.RouterGroup) {
	tc := &controller.TelegramLoginController{
		UserRepo:     userRepo,
		LoginService: loginService,
		Env:          env,
	}

	group.POST("/auth/telegram/login", tc.Login)
}
