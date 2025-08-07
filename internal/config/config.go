package config

import (
	"fmt"
	"github.com/caarlos0/env/v6"
	"time"

	"github.com/joho/godotenv"
)

// Config - структура для парсинга файла конфигурации
type Config struct {
	Server   ServerConfig   `envPrefix:"SERVER_"`
	Database DatabaseConfig `envPrefix:"DB_"`
}

// ServerConfig - структура для конфигурации сервера
type ServerConfig struct {
	Host    string        `env:"HOST" envDefault:""`
	Port    int           `env:"PORT" envDefault:"8080"`
	Timeout time.Duration `env:"TIMEOUT" envDefault:"5s"`
}

// DatabaseConfig - структура для конфигурации базы данных
type DatabaseConfig struct {
	Host     string `env:"HOST" envDefault:"db"`
	User     string `env:"USER" envDefault:"postgres"`
	Password string `env:"PASSWORD" envDefault:"postgres"`
	Name     string `env:"NAME" envDefault:"subscription_aggregation"`
	Port     int    `env:"PORT" envDefault:"5432"`
}

// Init пполучает данные из переменных окружения и возвращает объект Config
func Init() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, fmt.Errorf("error loading env variables: %w\n", err)
	}

	conf := Config{}
	err = env.Parse(&conf)
	if err != nil {
		return nil, err
	}

	return &conf, nil
}
