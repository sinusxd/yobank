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

func (uc *UserController) GetByEmail(c *gin.Context) {
	email := c.Param("email")

	user, err := uc.UserService.GetUserInfoByEmail(c.Request.Context(), email)
	if err != nil {
		c.JSON(http.StatusNotFound, domain.ErrorResponse{Message: "Пользователь не найден"})
		return
	}

	c.JSON(http.StatusOK, user)
}

func (uc *UserController) GetByTelegramID(c *gin.Context) {
	tgIDStr := c.Param("telegramId")
	tgID, err := strconv.ParseInt(tgIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{Message: "Некорректный telegram ID"})
		return
	}

	user, err := uc.UserService.GetUserInfoByTelegramID(c.Request.Context(), tgID)
	if err != nil {
		c.JSON(http.StatusNotFound, domain.ErrorResponse{Message: "Пользователь не найден"})
		return
	}

	c.JSON(http.StatusOK, user)
}

func (uc *UserController) GetByUsername(c *gin.Context) {
	username := c.Param("username")

	user, err := uc.UserService.GetByUsername(c.Request.Context(), username)
	if err != nil {
		c.JSON(http.StatusNotFound, domain.ErrorResponse{Message: "Пользователь не найден"})
		return
	}

	c.JSON(http.StatusOK, user)
}

func (uc *UserController) GetByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{Message: "Некорректный ID"})
		return
	}

	user, err := uc.UserService.GetUserInfoByID(c.Request.Context(), uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, domain.ErrorResponse{Message: "Пользователь не найден"})
		return
	}

	c.JSON(http.StatusOK, user)
}
