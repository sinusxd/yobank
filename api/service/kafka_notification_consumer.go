package service

import (
	"context"
	"encoding/json"
	"github.com/IBM/sarama"
	"log"
	"os"
	"os/signal"
	"time"
	"yobank/domain"
	"yobank/internal/mailer"
	"yobank/internal/telegram"
)

type NotificationConsumer struct {
	mailer *mailer.GoMailer
}

func (NotificationConsumer) Setup(_ sarama.ConsumerGroupSession) error   { return nil }
func (NotificationConsumer) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }

func (c NotificationConsumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		var event domain.TransferNotificationEvent
		if err := json.Unmarshal(msg.Value, &event); err != nil {
			log.Printf("Ошибка разбора сообщения: %v", err)
			session.MarkMessage(msg, "")
			continue
		}

		success := false
		for i := 0; i < 3; i++ {
			if err := c.tryNotify(event); err != nil {
				log.Printf("Попытка %d уведомления неудачна: %v", i+1, err)
				time.Sleep(2 * time.Second)
			} else {
				success = true
				break
			}
		}

		if !success {
			log.Printf("Не удалось отправить уведомление после 3 попыток: %+v", event)
		}

		session.MarkMessage(msg, "")
	}
	return nil
}

func (c NotificationConsumer) tryNotify(event domain.TransferNotificationEvent) error {
	if event.UseTelegram {
		return telegram.NotifyTransfer(*event.ReceiverTgID, event.SenderUsername, event.Amount, event.Currency, event.UseTelegram)
	} else {
		return c.mailer.SendTransferNotification(event.ReceiverEmail, event.SenderUsername, event.Amount, event.Currency)
	}
}

func StartKafkaConsumer(mailer *mailer.GoMailer, brokers []string, topic string, groupID string) error {
	config := sarama.NewConfig()
	config.Version = sarama.V2_1_0_0
	config.Consumer.Offsets.Initial = sarama.OffsetOldest

	consumerGroup, err := sarama.NewConsumerGroup(brokers, groupID, config)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		sigchan := make(chan os.Signal, 1)
		signal.Notify(sigchan, os.Interrupt)
		<-sigchan
		cancel()
	}()

	handler := NotificationConsumer{
		mailer: mailer,
	}

	for {
		if err := consumerGroup.Consume(ctx, []string{topic}, handler); err != nil {
			log.Printf("Consumer error: %v", err)
			break
		}

		if ctx.Err() != nil {
			break
		}
	}

	return consumerGroup.Close()
}
