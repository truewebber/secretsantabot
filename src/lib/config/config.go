package config

import (
	"os"

	"github.com/mgutz/logxi/v1"
	"github.com/spf13/viper"
)

type (
	Config struct {
		*viper.Viper
	}

	DbConnectionConfig struct {
		Host     string
		Db       string
		User     string
		Password string
	}
)

const (
	DevEnvironment     = "dev"
	DefaultEnvironment = DevEnvironment
	DefaultFormat      = "yaml"

	envName        = "ENV"
	configPathName = "CONFIG_PATH"
)

var (
	ENV         = ""
	configPath  = ""
	configCache = make(map[string]*Config)
	envCache    = make(map[string]string)
)

func init() {
	ENV = os.Getenv(envName)
	configPath = os.Getenv(configPathName)
}

func Get() *Config {
	if len(ENV) == 0 {
		ENV = DefaultEnvironment
	}

	if _, ok := configCache[ENV]; !ok {
		configCache[ENV] = NewConfig(ENV)
	}

	return configCache[ENV]
}

func GetEnvParam(envParam string) string {
	if _, ok := envCache[envParam]; !ok {
		envCache[envParam] = os.Getenv(envParam)
	}

	return envCache[envParam]
}

func IsDev() bool {
	return ENV == DevEnvironment
}

func NewConfig(envName string) *Config {
	v := viper.New()

	log.Debug("Load environment", "env", envName)
	v.SetConfigName(envName)

	if len(configPath) != 0 {
		v.AddConfigPath(configPath)
	}

	v.AddConfigPath(".")
	v.SetConfigType(DefaultFormat)

	err := v.ReadInConfig()
	if err != nil {
		log.Fatal("Error load config", "error", err.Error())
	}

	return &Config{v}
}
