package repository

import "github.com/epiehl93/h24-notifier/pkg/models"

type CycleRepository interface {
	Create(cycle *models.Cycle) error
	GetLastSuccessfulCycle(cycleType models.CycleType) (*models.Cycle, error)
}
