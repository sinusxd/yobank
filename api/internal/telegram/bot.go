package telegram

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
)

func StartBot(token string, webAppURL string) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Panic(err)
	}

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
