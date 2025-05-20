package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"yobank/domain"
)

type UserController struct {
	UserService domain.UserService
}

func (uc *UserController) GetMe(c *gin.Context) {
	id, ok := c.Get("x-user-id")
	if !ok {
		c.JSON(http.StatusUnauthorized, domain.ErrorResponse{Message: "Unauthorized"})
		return
	}

	userID, err := strconv.ParseUint(id.(string), 10, 64)
	if err != nil {
		c.JSON(http.StatusUnauthorized, domain.ErrorResponse{Message: "Invalid user ID"})
		return
	}

	user, err := uc.UserService.GetUserInfoByID(c.Request.Context(), uint(userID))
	if err != nil {
		c.JSON(http.StatusNotFound, domain.ErrorResponse{Message: "Пользователь не найден"})
		return
	}

	c.JSON(http.StatusOK, user)
}
