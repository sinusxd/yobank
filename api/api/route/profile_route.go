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

func NewProfileRouter(env *bootstrap.Env, timeout time.Duration, db *gorm.DB, group *gin.RouterGroup) {
	ur := repository.NewUserRepository(db)
	pc := &controller.ProfileController{
		ProfileUsecase: service.NewProfileUsecase(ur, timeout),
	}
	group.GET("/profile", pc.Fetch)
}
