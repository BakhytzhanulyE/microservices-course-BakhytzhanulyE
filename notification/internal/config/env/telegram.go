package env

import "github.com/caarlos0/env/v11"

type telegramEnvConfig struct {
	// Пустой токен выключает отправку: уведомления уходят в лог.
	Token  string `env:"TELEGRAM_BOT_TOKEN" envDefault:""`
	ChatID string `env:"TELEGRAM_CHAT_ID"   envDefault:""`
}

type telegramConfig struct {
	raw telegramEnvConfig
}

// NewTelegramConfig читает настройки Telegram-бота.
func NewTelegramConfig() (*telegramConfig, error) {
	var raw telegramEnvConfig
	if err := env.Parse(&raw); err != nil {
		return nil, err
	}

	return &telegramConfig{raw: raw}, nil
}

// Token — токен бота.
func (cfg *telegramConfig) Token() string { return cfg.raw.Token }

// ChatID — идентификатор чата, куда пишем.
func (cfg *telegramConfig) ChatID() string { return cfg.raw.ChatID }
