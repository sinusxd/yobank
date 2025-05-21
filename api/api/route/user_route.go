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

	group.GET("/users/id/:id", uc.GetByID)
	group.GET("/users/email/:email", uc.GetByEmail)
	group.GET("/users/telegram/:telegramId", uc.GetByTelegramID)
	group.GET("/users/username/:username", uc.GetByUsername)
}
