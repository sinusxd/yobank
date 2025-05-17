package telegram

import (
	"context"
	initdata "github.com/telegram-mini-apps/init-data-golang"
)

type contextKey string

const _initDataKey contextKey = "init-data"

// Вспомогательная функция: кладём initData в context
func WithInitData(ctx context.Context, initData initdata.InitData) context.Context {
	return context.WithValue(ctx, _initDataKey, initData)
}

// Вспомогательная функция: извлекаем initData из context
func CtxInitData(ctx context.Context) (initdata.InitData, bool) {
	initData, ok := ctx.Value(_initDataKey).(initdata.InitData)
	return initData, ok
}
