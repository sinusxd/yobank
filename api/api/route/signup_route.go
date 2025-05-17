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

func NewSignupRouter(env *bootstrap.Env, timeout time.Duration, db *gorm.DB, group *gin.RouterGroup) {
	ur := repository.NewUserRepository(db)
	sc := &controller.SignupController{
		SignupUsecase: service.NewSignupUsecase(ur, timeout),
		Env:           env,
	}
	group.POST("/signup", sc.Signup)
}
