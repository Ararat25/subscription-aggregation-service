package database

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Dbinstance - структура для базы данных
type Dbinstance struct {
	Db *gorm.DB
}

// DB - глобальная переменная для работы с базой данных
var DB Dbinstance

// ConnectDB производит соединение с бд
func ConnectDB(dbHost string, dbUser string, dbPassword string, dbName string, dbPort int) error {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=Europe/Moscow",
		dbHost,
		dbUser,
		dbPassword,
		dbName,
		dbPort,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Error),
	})
	if err != nil {
		return err
	}

	db.Logger = logger.Default.LogMode(logger.Error)

	DB = Dbinstance{
		Db: db,
	}

	return nil
}
