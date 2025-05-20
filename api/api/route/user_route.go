package route

import (
	"github.com/gin-gonic/gin"
	"yobank/api/controller"
	"yobank/domain"
)

func NewUserRouter(userService domain.UserService, group *gin.RouterGroup) {
	uc := &controller.UserController{
		UserService: userService,
	}

	group.GET("/users/me", uc.GetMe)
}
