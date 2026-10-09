package config

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewDatabase(viper *viper.Viper) *gorm.DB {

	host := viper.GetString("DB_HOST")
	user := viper.GetString("DB_USER")
	password := viper.GetString("DB_PASSWORD")
	dbname := viper.GetString("DB_NAME")
	port := viper.GetString("DB_PORT")
	sslmode := viper.GetString("DB_SSLMODE")
	tz := viper.GetString("DB_TIMEZONE")

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s",
		host, user, password, dbname, port, sslmode, tz,
	)

	gormLogger := logger.Default.LogMode(logger.Info)
	if viper.GetString("APP_ENV") == "production" {
		gormLogger = logger.Default.LogMode(logger.Error)
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: gormLogger,
	})
	if err != nil {
		panic(fmt.Sprintf("[GORM] Failed to initialize PostgreSQL connection: %v", err))
	}

	sqlDB, err := db.DB()
	if err != nil {
		panic(fmt.Sprintf("[GORM] Failed to access underlying sql.DB: %v", err))
	}

	// Force a ping to verify connection at startup
	if err := sqlDB.Ping(); err != nil {
		panic(fmt.Sprintf("[GORM] Failed to ping database: %v", err))
	}

	maxIdle := viper.GetInt("DB_POOL_IDLE")
	if maxIdle == 0 {
		maxIdle = 10
	}
	maxOpen := viper.GetInt("DB_POOL_MAX")
	if maxOpen == 0 {
		maxOpen = 100
	}
	lifetime := viper.GetInt("DB_POOL_LIFETIME")
	if lifetime == 0 {
		lifetime = 300
	}

	sqlDB.SetMaxIdleConns(maxIdle)
	sqlDB.SetMaxOpenConns(maxOpen)
	sqlDB.SetConnMaxLifetime(time.Duration(lifetime) * time.Second)

	return db
}
