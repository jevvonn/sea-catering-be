package repository

import (
	"github.com/jevvonn/sea-catering-be/internal/domain/entity"
	"gorm.io/gorm"
)

type PlansPostgreSQLItf interface {
	GetPlans() ([]entity.Plans, error)
	UpdatePlan(plan entity.Plans) error
}

type PlansPostgreSQL struct {
	db *gorm.DB
}

func NewPlansPostgreSQL(db *gorm.DB) PlansPostgreSQLItf {
	return &PlansPostgreSQL{db}
}

func (r *PlansPostgreSQL) GetPlans() ([]entity.Plans, error) {
	var plans []entity.Plans
	if err := r.db.Find(&plans).Error; err != nil {
		return nil, err
	}
	return plans, nil
}

func (r *PlansPostgreSQL) UpdatePlan(plan entity.Plans) error {
	if plan.ID == "" {
		return gorm.ErrRecordNotFound
	}

	if err := r.db.Updates(&plan).Error; err != nil {
		return err
	}

	return nil
}
