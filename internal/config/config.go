package config

import (
	"os"
	"time"

	_ "github.com/lib/pq" // postgres driver don`t delete
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

type ServerConfig struct {
	Name string `mapstructure:"name"`
	Port string `mapstructure:"port"`
}

type LoggerConfig struct {
	Level  string `mapstructure:"level"`
	Format string `mapstructure:"format"`
}

type DatabaseConfig struct {
	SQL struct {
		URL string `mapstructure:"url"`
	} `mapstructure:"sql"`
}

type OAuthConfig struct {
	Google struct {
		ClientID     string `mapstructure:"client_id"`
		ClientSecret string `mapstructure:"client_secret"`
		RedirectURL  string `mapstructure:"redirect_url"`
	}
}

type KafkaConfig struct {
	Brokers []string `mapstructure:"brokers"`
}

type JWTConfig struct {
	User struct {
		AccessToken struct {
			SecretKey     string        `mapstructure:"secret_key"`
			TokenLifetime time.Duration `mapstructure:"token_lifetime"`
		} `mapstructure:"access_token"`
		RefreshToken struct {
			SecretKey     string        `mapstructure:"secret_key"`
			EncryptionKey string        `mapstructure:"encryption_key"`
			TokenLifetime time.Duration `mapstructure:"token_lifetime"`
		} `mapstructure:"refresh_token"`
	} `mapstructure:"user"`
	Service struct {
		SecretKey string `mapstructure:"secret_key"`
	} `mapstructure:"service"`
}

type SwaggerConfig struct {
	Enabled bool   `mapstructure:"enabled"`
	URL     string `mapstructure:"url"`
	Port    string `mapstructure:"port"`
}

type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Logger   LoggerConfig   `mapstructure:"logger"`
	JWT      JWTConfig      `mapstructure:"jwt"`
	OAuth    OAuthConfig    `mapstructure:"oauth"`
	Kafka    KafkaConfig    `mapstructure:"kafka"`
	Database DatabaseConfig `mapstructure:"database"`
	Swagger  SwaggerConfig  `mapstructure:"swagger"`
}

func LoadConfig() (Config, error) {
	configPath := os.Getenv("KV_VIPER_FILE")
	if configPath == "" {
		return Config{}, errors.New("KV_VIPER_FILE env var is not set")
	}
	viper.SetConfigFile(configPath)

	if err := viper.ReadInConfig(); err != nil {
		return Config{}, errors.Errorf("error reading config file: %s", err)
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return Config{}, errors.Errorf("error unmarshalling config: %s", err)
	}

	return config, nil
}
