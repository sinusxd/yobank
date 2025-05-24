package telegram

import (
	"context"
	"encoding/json"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"yobank/domain"
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

func NotifyTransfer(tgID int64, senderUsername string, amount int64, currency string, senderFromTg bool) error {
	if BotInstance == nil {
		log.Println("Bot not initialized")
		return fmt.Errorf("BotInstance is nil")
	}
	text := ""
	if senderFromTg {
		text = fmt.Sprintf("💸 Вы получили перевод от @%s на сумму %.2f %s", senderUsername, float64(amount)/100, currency)
	} else {
		text = fmt.Sprintf("💸 Вы получили перевод от %s на сумму %.2f %s", senderUsername, float64(amount)/100, currency)
	}

	msg := tgbotapi.NewMessage(tgID, text)
	_, err := BotInstance.Send(msg)
	return err
}

func HandleTransferNotificationMessage(ctx context.Context, message []byte) error {
	var event domain.TransferNotificationEvent
	if err := json.Unmarshal(message, &event); err != nil {
		log.Printf("Ошибка при разборе JSON: %v", err)
		return err
	}

	return NotifyTransfer(
		*event.ReceiverTgID,
		event.SenderUsername,
		event.Amount,
		event.Currency,
		true,
	)

}

func NotifyTopUp(tgID int64, amount int64, currency string) {
	if BotInstance == nil {
		log.Println("BotInstance не инициализирован")
		return
	}

	text := fmt.Sprintf("💰 Ваш кошелёк был пополнен на %.2f %s", float64(amount)/100, currency)

	msg := tgbotapi.NewMessage(tgID, text)
	if _, err := BotInstance.Send(msg); err != nil {
		log.Printf("Ошибка отправки уведомления о пополнении: %v", err)
	}
}
