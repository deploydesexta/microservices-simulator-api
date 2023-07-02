package config

import (
	"github.com/kelseyhightower/envconfig"
	"github.com/stackus/dotenv"
	"os"
)

type (
	SecurityConfig struct {
		Id          string `envconfig:"SECURITY_ID" required:"true"`
		PublicKey   string `envconfig:"SECURITY_PUBLIC_KEY" required:"true"`
		PrivateKey  string `envconfig:"SECURITY_PRIVATE_KEY" required:"true"`
		Salts       int    `envconfig:"SECURITY_SALTS" required:"true"`
		TokenMaxAge int    `envconfig:"SECURITY_TOKEN_MAXAGE_DAYS" required:"true"`
		TokenDomain string `envconfig:"SECURITY_TOKEN_DOMAIN" required:"true"`
		TokenName   string `envconfig:"SECURITY_TOKEN_NAME" required:"true"`
	}

	RedisConfig struct {
		Host     string `envconfig:"REDIS_HOST" required:"true"`
		Password string `envconfig:"REDIS_PASS" required:"true"`
		Db       string `envconfig:"REDIS_DB" required:"true"`
	}

	Config struct {
		Environment string `envconfig:"ENVIRONMENT" default:"development"`
		LogLevel    string `envconfig:"LOG_LEVEL" default:"INFO"`
		Security    SecurityConfig
		Redis       RedisConfig
	}
)

func InitConfig() (cfg Config, err error) {
	if err = dotenv.Load(dotenv.EnvironmentFiles(getEnvironment())); err != nil {
		return
	}
	err = envconfig.Process("", &cfg)
	return
}

func getEnvironment() string {
	environment := os.Getenv("ENVIRONMENT")
	if environment == "" {
		environment = "development"
	}
	return environment
}
