package database

import (
	"fmt"

	"gorm.io/gorm/logger"

	"github.com/sarrooo/go-clean/internal/errcode"
	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func dsnBuilder() string {
	host := viper.GetString("DATABASE_HOST")
	port := viper.GetString("POSTGRES_PORT")
	user := viper.GetString("POSTGRES_USER")
	password := viper.GetString("POSTGRES_PASSWORD")
	dbName := viper.GetString("POSTGRES_DB")
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable TimeZone=Europe/Paris", host, port, user, password, dbName)
}

func NewGormClient() (*gorm.DB, error) {
	dsn := dsnBuilder()
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, fmt.Errorf("%w: %v", errcode.ErrConfigurationFailed, err)
	}

	err = migrate(db)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", errcode.ErrDatabaseMigrate, err)
	}

	if viper.GetBool("SEED_DB") {
		CreateTestingEntities(db)
	}

	return db, nil
}
