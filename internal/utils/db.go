package utils

import (
	"fmt"
	"github.com/epiehl93/h24-notifier/config"
	"github.com/epiehl93/h24-notifier/pkg/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func SetupDB() (*gorm.DB, error) {
	var err error

	dsn := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=%s TimeZone=%s",
		config.C.Database.User,
		config.C.Database.Password,
		config.C.Database.DBName,
		config.C.Database.Host,
		config.C.Database.Port,
		config.C.Database.Sslmode,
		config.C.Database.Timezone)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
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
