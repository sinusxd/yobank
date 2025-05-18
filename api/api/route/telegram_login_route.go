package route

import (
	"yobank/api/controller"
	"yobank/bootstrap"
	"yobank/domain"

	"github.com/gin-gonic/gin"
)

func NewTelegramLoginRouter(userService domain.UserService, loginService domain.LoginService, env *bootstrap.Env, group *gin.RouterGroup) {
	tc := &controller.TelegramLoginController{
		UserService:  userService,
		LoginService: loginService,
		Env:          env,
	}

	group.POST("/auth/telegram/login", tc.Login)
}
