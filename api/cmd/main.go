package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"time"
	"yobank/api/route"
	"yobank/bootstrap"
	"yobank/internal/telegram"
	"yobank/service"
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

	brokers := []string{"broker:9092"}
	topic := "transfer_notifications"
	groupID := "notification-workers"

	errChan := make(chan error, 1)

	go func() {
		if err := service.StartKafkaConsumer(app.Container.Services.Mailer, brokers, topic, groupID); err != nil {
			errChan <- err
		}
	}()

	go func() {
		if err := gin.Run(env.ServerAddress); err != nil {
			errChan <- err
		}
	}()

	if err := <-errChan; err != nil {
		log.Fatalf("Приложение аварийно завершено: %v", err)
	}
}
