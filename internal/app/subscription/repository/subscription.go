package repository

import (
	"github.com/google/uuid"
	"github.com/jevvonn/sea-catering-be/internal/domain/entity"
	"gorm.io/gorm"
)

type SubscriptionPostgreSQLItf interface {
	GetSubscriptions(cond entity.Subscription) ([]entity.Subscription, error)
	GetSpecific(subscription entity.Subscription) (entity.Subscription, error)
	CreateSubscription(subscription entity.Subscription) error
	UpdateSubscription(subscription entity.Subscription) error
}

type SubscriptionPostgreSQL struct {
	db *gorm.DB
}

func NewSubscriptionPostgreSQL(db *gorm.DB) SubscriptionPostgreSQLItf {
	return &SubscriptionPostgreSQL{db}
}

func (r *SubscriptionPostgreSQL) GetSubscriptions(cond entity.Subscription) ([]entity.Subscription, error) {
	var subscriptions []entity.Subscription
	if err := r.db.Preload("Plans").Preload("User").Where(cond).Find(&subscriptions).Error; err != nil {
		return nil, err
	}

	return subscriptions, nil
}

func (r *SubscriptionPostgreSQL) GetSpecific(subscription entity.Subscription) (entity.Subscription, error) {
	var result entity.Subscription

	if err := r.db.Preload("Plans").Preload("User").First(&result, &subscription).Error; err != nil {
		return entity.Subscription{}, err
	}

	return result, nil
}

func (r *SubscriptionPostgreSQL) CreateSubscription(subscription entity.Subscription) error {
	if err := r.db.Create(&subscription).Error; err != nil {
		return err
	}

	return nil
}

func (r *SubscriptionPostgreSQL) UpdateSubscription(subscription entity.Subscription) error {
	if subscription.ID == uuid.Nil {
		return gorm.ErrRecordNotFound
	}

	data := map[string]any{}

	if subscription.Name != "" {
		data["name"] = subscription.Name
	}

	if subscription.PhoneNumber != "" {
		data["phone_number"] = subscription.PhoneNumber
	}
	if subscription.Status != "" {
		data["status"] = subscription.Status
	}
	if subscription.PauseStartDate != nil {
		data["pause_start_date"] = subscription.PauseStartDate
	}
	if subscription.PauseEndDate != nil {
		data["pause_end_date"] = subscription.PauseEndDate
	}

	if err := r.db.Model(entity.Subscription{}).Where("id = ?", subscription.ID).Updates(&data).Error; err != nil {
		return err
	}

	return nil
}
