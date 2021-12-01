package main

import (
	"fmt"

	"github.com/Netflix/go-env"

	"github.com/truewebber/secretsantabot/internal/log"
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
	return nil
}

type config struct {
	AppAddress     string `env:"APP_ADDRESS,required=true"`
	MetricsAddress string `env:"METRICS_ADDRESS,required=true"`

	ProjectsService struct {
		Address string `env:"PROJECTS_SERVICE_ADDRESS,required=true"`
		Token   string `env:"PROJECTS_SERVICE_TOKEN,required=true"`
	}
	RecipientService struct {
		Address string `env:"RECIPIENT_SERVICE_ADDRESS,required=true"`
	}
	NotificationSettingsService struct {
		Address string `env:"NOTIFICATION_SETTINGS_SERVICE_ADDRESS,required=true"`
	}
	AdPickerService struct {
		Address string `env:"AD_PICKER_SERVICE_ADDRESS,required=true"`
	}
	ArchiveService struct {
		Address string `env:"ARCHIVE_SERVICE_ADDRESS,required=true"`
	}
	LimitService struct {
		Host string `env:"LIMIT_SERVICE_HOST,required=true"`
	}
	Redis struct {
		HostWithPort string `env:"REDIS_HOST_WITH_PORT,required=true"`
	}
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
