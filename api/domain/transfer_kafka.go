package domain

import "context"

type TransferNotificationEvent struct {
	ReceiverID     int    `json:"ReceiverID"`
	ReceiverEmail  string `json:"ReceiverEmail"`
	ReceiverTgID   *int64 `json:"ReceiverTgID"`
	SenderUsername string `json:"SenderUsername"`
	Amount         int64  `json:"Amount"`
	Currency       string `json:"Currency"`
	UseTelegram    bool   `json:"UseTelegram"`
	RetryCount     int    `json:"RetryCount"`
}

type NotificationProducer interface {
	SendTransferNotificationEvent(ctx context.Context, event TransferNotificationEvent) error
}
