package database

import (
	"affiliates-backoffice-backend/pkg/config"
	"affiliates-backoffice-backend/pkg/log"
	"fmt"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
	"moul.io/zapgorm2"
)

var DB *gorm.DB

func Connect(logger log.LoggerI) {
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		config.Vars.DbHost,
		config.Vars.DbPort,
		config.Vars.DbUser,
		config.Vars.DbPassword,
		config.Vars.DbName)

	newLogger := zapgorm2.New(logger.Desugar())
	newLogger.LogLevel = gormlogger.Info
	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		logger.Fatal("Error connecting to database: " + err.Error())
	}

	sqlDB, err := DB.DB()
	if err != nil {
		logger.Fatal("Error connecting to database: " + err.Error())
	}

	sqlDB.SetMaxIdleConns(config.Vars.DbMaxIdleConns)
	sqlDB.SetMaxOpenConns(config.Vars.DbMaxOpenConns)
	sqlDB.SetConnMaxLifetime(time.Hour)

	logger.Info("connection opened to database")
}
