package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/jevvonn/sea-catering-be/internal/domain/entity"
)

type GetSubscriptionResponse struct {
	ID uuid.UUID `json:"id,omitempty"`

	UserID uuid.UUID   `json:"user_id,omitempty"`
	User   entity.User `json:"user,omitempty"`

	PlanId string       `json:"plan_id,omitempty"`
	Plans  entity.Plans `json:"plan,omitempty"`

	Name        string `json:"name,omitempty"`
	PhoneNumber string `json:"phone_number,omitempty"`

	Mealtypes    []string `json:"mealtype,omitempty"`
	DeliveryDays []string `json:"delivery_days,omitempty"`
	Allergies    []string `json:"allergies,omitempty"`

	TotalPrice float64 `json:"total_price,omitempty"`

	Status         string     `json:"status,omitempty"`
	PauseStartDate *time.Time `json:"pause_start_date,omitempty"`
	PauseEndDate   *time.Time `json:"pause_end_date,omitempty"`

	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}
