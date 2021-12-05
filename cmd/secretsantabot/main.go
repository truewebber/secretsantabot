package main

import (
	"fmt"
	"syscall"

	"github.com/Netflix/go-env"

	"github.com/truewebber/secretsantabot/internal/log"
	"github.com/truewebber/secretsantabot/internal/port"
	"github.com/truewebber/secretsantabot/internal/service"
	"github.com/truewebber/secretsantabot/internal/signal"
)

func main() {
	logger := log.NewZapWrapper()

	if err := run(logger); err != nil {
		panic(err)
	}

	if err := logger.Close(); err != nil {
		panic(err)
	}
}

func run(logger log.Logger) error {
	cfg := mustLoadConfig()

	appConfig := service.NewConfig(logger)
	application := service.MustNewApplication(appConfig)
	bot := port.MustNewTelegramBot(cfg.TelegramToken, application, logger)

	ctx := signal.ContextClosableOnSignals(syscall.SIGINT, syscall.SIGTERM)

	if err := bot.Run(ctx); err != nil {
		return fmt.Errorf("bot run: %w", err)
	}

	return nil
}

type config struct {
	TelegramToken string `env:"TELEGRAM_TOKEN,required=true"`
	PostgresURI   string `env:"DATABASE_URL,required=true"`
}

func mustLoadConfig() *config {
	cfg, err := loadConfig()
	if err != nil {
		panic(err)
	}

	return cfg
}

func loadConfig() (*config, error) {
	var c config

	if _, err := env.UnmarshalFromEnviron(&c); err != nil {
		return nil, fmt.Errorf("config unmarshal: %w", err)
	}

	return &c, nil
}
