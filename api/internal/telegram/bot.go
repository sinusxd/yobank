package telegram

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
)

var BotInstance *tgbotapi.BotAPI

func StartBot(token string, webAppURL string) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Panic(err)
	}
	BotInstance = bot
	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil {
			button := tgbotapi.InlineKeyboardButton{
				Text: "Открыть мини-приложение",
				URL:  &webAppURL,
			}

			keyboard := tgbotapi.NewInlineKeyboardMarkup(
				tgbotapi.NewInlineKeyboardRow(button),
			)

			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Привет! Нажми кнопку ниже, чтобы открыть мини-приложение:")
			msg.ReplyMarkup = keyboard

			if _, err := bot.Send(msg); err != nil {
				log.Printf("Ошибка при отправке: %v", err)
			}
		}
	}
}

func NotifyTransfer(tgID int64, senderUsername string, amount int64, currency string, senderFromTg bool) {
	if BotInstance == nil {
		log.Println("Bot not initialized")
		return
	}
	text := ""
	if senderFromTg {
		text = fmt.Sprintf("💸 Вы получили перевод от @%s на сумму %.2f %s", senderUsername, float64(amount)/100, currency)
	} else {
		text = fmt.Sprintf("💸 Вы получили перевод от %s на сумму %.2f %s", senderUsername, float64(amount)/100, currency)
	}

	msg := tgbotapi.NewMessage(tgID, text)
	if _, err := BotInstance.Send(msg); err != nil {
		log.Printf("Не удалось отправить уведомление пользователю %d: %v", tgID, err)
	}
}
