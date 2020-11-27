package test_utils

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func SetupTestDb() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open("file::memory?cache=shared"), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return db, nil
}
