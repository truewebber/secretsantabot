package main

import (
	"fmt"
	"syscall"

	"github.com/Netflix/go-env"

	"github.com/truewebber/secretsantabot/domain/log"
	"github.com/truewebber/secretsantabot/domain/signal"
	"github.com/truewebber/secretsantabot/port/telegram"
	"github.com/truewebber/secretsantabot/service"
)

func main() {
	logger := log.NewZapWrapper()

	mustRun(logger)

	if err := logger.Close(); err != nil {
		panic(err)
	}
}

func mustRun(logger log.Logger) {
	cfg := mustLoadConfig()

	logger.Info("Config inited")

	appConfig := service.NewConfig(service.ChatService{PostgresURI: cfg.PostgresURI})
	application := service.MustNewApplication(appConfig)
	bot := telegram.MustNewTelegramBot(cfg.TelegramToken, application, logger)

	logger.Info("Starting application")

	ctx := signal.ContextClosableOnSignals(syscall.SIGINT, syscall.SIGTERM)

	bot.Run(ctx)

	logger.Info("Application stopped")
}

type config struct {
	TelegramToken string `env:"TELEGRAM_TOKEN,required=true"`
	PostgresURI   string `env:"POSTGRES_URI,required=true"`
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
