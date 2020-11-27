package adapter

import (
	"errors"
	"github.com/epiehl93/h24-notifier/pkg/models"
	"github.com/epiehl93/h24-notifier/pkg/repository"
	"gorm.io/gorm"
)

type CycleRepository struct {
	db *gorm.DB
}

func (a CycleRepository) GetLastSuccessfulCycle(cycleType models.CycleType) (*models.Cycle, error) {
	cycle := models.Cycle{}
	tx := a.db.Where("successful = ? and type = ?", true, cycleType).Last(&cycle)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return &cycle, nil
}

func (a CycleRepository) Create(cycle *models.Cycle) error {
	if cycle.Type == models.UnsetCycle {
		return errors.New("cycle type has to be set")
	}
	tx := a.db.Create(cycle)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}

func NewCycleRepository(db *gorm.DB) repository.CycleRepository {
	return &CycleRepository{db: db}
}
