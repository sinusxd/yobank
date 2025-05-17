package route

import (
	"time"
	"yobank/repository"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"yobank/api/controller"
	"yobank/bootstrap"
	"yobank/service"
)

func NewRefreshTokenRouter(env *bootstrap.Env, timeout time.Duration, db *gorm.DB, group *gin.RouterGroup) {
	ur := repository.NewUserRepository(db)
	rtc := &controller.RefreshTokenController{
		RefreshTokenUsecase: service.NewRefreshTokenUsecase(ur, timeout),
		Env:                 env,
	}
	group.POST("/refresh", rtc.RefreshToken)
}
