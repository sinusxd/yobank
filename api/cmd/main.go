package main

import (
	"time"
	"yobank/api/route"
	"yobank/internal/telegram"

	"github.com/gin-gonic/gin"
	"yobank/bootstrap"
)

func main() {
	app := bootstrap.App()
	env := app.Env

	defer app.CloseDBConnection()

	timeout := time.Duration(env.ContextTimeout) * time.Second
	gin := gin.Default()

	route.Setup(app, timeout, gin)

	go telegram.StartBot(env.TelegramBotToken, env.TelegramWebAppUrl)
	go app.Container.Services.Rate.StartScheduler()

	gin.Run(env.ServerAddress)
}
