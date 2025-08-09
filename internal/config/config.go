package config

import (
	"fmt"
	"time"

	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
)

// Config - структура для парсинга файла конфигурации
type Config struct {
	Server   ServerConfig   `envPrefix:"SERVER_"` // объект конфигурации сервера
	Database DatabaseConfig `envPrefix:"DB_"`     // объект конфигурации базы данных
}

// ServerConfig - структура для конфигурации сервера
type ServerConfig struct {
	Host    string        `env:"HOST" envDefault:""`      // хост сервера
	Port    int           `env:"PORT" envDefault:"8080"`  // порт сервера
	Timeout time.Duration `env:"TIMEOUT" envDefault:"5s"` // таймаут сервера
}

// DatabaseConfig - структура для конфигурации базы данных
type DatabaseConfig struct {
	Host     string `env:"HOST" envDefault:"db"`                       // хост базы данных
	User     string `env:"USER" envDefault:"postgres"`                 // пользователь базы данных
	Password string `env:"PASSWORD" envDefault:"postgres"`             // пароль базы данных
	Name     string `env:"NAME" envDefault:"subscription_aggregation"` // название базы данных
	Port     int    `env:"PORT" envDefault:"5432"`                     // порт базы данных
}

// Init получает данные из переменных окружения и возвращает объект Config
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
