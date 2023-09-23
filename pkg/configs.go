package pkg

import (
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"github.com/sirupsen/logrus"
)

const envPrefix = "APP"

type AppConfigs struct {
	ServerPort string `envconfig:"PORT"`
	// Безопасные логи (скрыть PAN и прочее)
	SafeLogs     bool   `envconfig:"SAFE_LOGS"`
	DbDSN        string `envconfig:"DB_DSN"`
	VaultAddress string `envconfig:"VAULT_ADDRESS"`
	VaultToken   string `envconfig:"VAULT_TOKEN"`
	VaultCapath  string `envconfig:"VAULT_CAPATH"`
}

func NewConfigs() *AppConfigs {
	return &AppConfigs{
		ServerPort: ":8090",
		SafeLogs:   true,
	}
}

func NewConfigsByEnv() *AppConfigs {
	if err := godotenv.Load(); err != nil {
		logrus.Print(".env file not found")
	}

	cfg := NewConfigs()
	if err := envconfig.Process(envPrefix, cfg); err != nil {
		logrus.Print("Error parse ENV")
	}

	return cfg
}
