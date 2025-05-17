package route

import (
	"github.com/gin-gonic/gin"
	"yobank/api/controller"
	"yobank/bootstrap"
	"yobank/domain"
)

func NewEmailLoginRouter(loginService domain.LoginService, codeService domain.EmailCodeService, userRepo domain.UserRepository, env *bootstrap.Env, group *gin.RouterGroup) {
	lc := &controller.EmailLoginController{
		LoginService:   loginService,
		CodeService:    codeService,
		UserRepository: userRepo,
		Env:            env,
	}
	group.POST("/auth/email/request-code", lc.RequestCode)
	group.POST("/auth/email/verify-code", lc.VerifyCode)
}
