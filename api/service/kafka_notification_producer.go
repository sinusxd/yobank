package service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/IBM/sarama"
	"time"
	"yobank/domain"
)

type kafkaNotificationProducer struct {
	producer sarama.SyncProducer
	topic    string
}

func NewKafkaNotificationProducer(brokers []string, topic string) (domain.NotificationProducer, error) {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5
	config.Producer.Timeout = 5 * time.Second

	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		return nil, err
	}

	return &kafkaNotificationProducer{
		producer: producer,
		topic:    topic,
	}, nil
}

func (k *kafkaNotificationProducer) SendTransferNotificationEvent(ctx context.Context, event domain.TransferNotificationEvent) error {
	data, err := json.Marshal(event)
	if err != nil {
		fmt.Println("Ошибка при сериализации события:", err)
		return err
	}

	msg := &sarama.ProducerMessage{
		Topic: k.topic,
		Value: sarama.ByteEncoder(data),
	}

	_, _, err = k.producer.SendMessage(msg)
	fmt.Println("Ошибка при отправке сообщения:", err)
	return err
}
