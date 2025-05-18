package route

import (
	"yobank/api/controller"
	"yobank/bootstrap"
	"yobank/domain"

	"github.com/gin-gonic/gin"
)

func NewEmailLoginRouter(userRepo domain.UserRepository, loginService domain.LoginService, codeService domain.EmailCodeService, walletService domain.WalletService, env *bootstrap.Env, group *gin.RouterGroup) {
	lc := &controller.EmailLoginController{
		UserRepository: userRepo,
		LoginService:   loginService,
		CodeService:    codeService,
		WalletService:  walletService,
		Env:            env,
	}
	group.POST("/auth/email/request-code", lc.RequestCode)
	group.POST("/auth/email/verify-code", lc.VerifyCode)
}
