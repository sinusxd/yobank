package route

import (
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"yobank/api/controller"
	"yobank/bootstrap"
	"yobank/postgres"
	"yobank/service"
)

func NewTaskRouter(env *bootstrap.Env, timeout time.Duration, db *gorm.DB, group *gin.RouterGroup) {
	tr := postgres.NewTaskRepository(db)
	tc := &controller.TaskController{
		TaskUsecase: service.NewTaskUsecase(tr, timeout),
	}
	group.POST("/task", tc.Create)
	group.GET("/task", tc.Fetch)
}
