package env

import "github.com/caarlos0/env/v11"

type telegramEnvConfig struct {
	BotToken     string `env:"TELEGRAM_BOT_TOKEN,required"`
	SkipAPICheck bool   `env:"TELEGRAM_SKIP_API_CHECK"`
}

type TelegramConfig struct {
	raw telegramEnvConfig
}

func NewTelegramConfig() (*TelegramConfig, error) {
	var raw telegramEnvConfig
	if err := env.Parse(&raw); err != nil {
		return nil, err
	}

	return &TelegramConfig{raw: raw}, nil
}

func (cfg *TelegramConfig) BotToken() string {
	return cfg.raw.BotToken
}

func (cfg *TelegramConfig) SkipAPICheck() bool {
	return cfg.raw.SkipAPICheck
}
