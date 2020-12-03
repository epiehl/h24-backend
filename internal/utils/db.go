package utils

import (
	"fmt"
	"github.com/epiehl93/h24-notifier/pkg/models"
	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"moul.io/zapgorm2"
)

func SetupDB() (*gorm.DB, error) {
	var err error

	dsn := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=%s TimeZone=%s",
		viper.GetString("database.user"),
		viper.GetString("database.password"),
		viper.GetString("database.dbname"),
		viper.GetString("database.host"),
		viper.GetString("database.port"),
		viper.GetString("database.sslmode"),
		viper.GetString("database.timezone"))

	logger := zapgorm2.New(LLogger)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{Logger: logger})
	if err != nil {
		return nil, err
	}

	return db, nil
}

func MigrateTables(db *gorm.DB) error {
	if err := db.AutoMigrate(&models.Item{}); err != nil {
		return err
	}
	if err := db.AutoMigrate(&models.Wishlist{}); err != nil {
		return err
	}
	if err := db.AutoMigrate(&models.Cycle{}); err != nil {
		return err
	}
	return nil
}
