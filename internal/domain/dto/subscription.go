package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/jevvonn/sea-catering-be/internal/domain/entity"
)

type CreateSubscriptionRequest struct {
	PlanId string `json:"plan_id,omitempty" validate:"required"`

	Name        string `json:"name,omitempty" validate:"required"`
	PhoneNumber string `json:"phone_number,omitempty" validate:"required"`

	Mealtypes    []string `json:"mealtype,omitempty" validate:"required,min=1,dive,oneof=Breakfast Lunch Dinner"`
	DeliveryDays []string `json:"delivery_days,omitempty" validate:"required,min=1,dive,oneof=Monday Tuesday Wednesday Thursday Friday Saturday Sunday"`
	Allergies    []string `json:"allergies,omitempty" validate:"required,dive"`
}

type UpdateSubscriptionRequest struct {
	Name        string `json:"name,omitempty"`
	PhoneNumber string `json:"phone_number,omitempty"`

	Status         string `json:"status,omitempty" validate:"omitempty,oneof=ACTIVE CANCELLED"`
	PauseStartDate string `json:"pause_start_date,omitempty" example:"27-06-2025"`
	PauseEndDate   string `json:"pause_end_date,omitempty" example:"30-06-2025"`
}

type GetSubscriptionResponse struct {
	ID uuid.UUID `json:"id,omitempty"`

	UserID uuid.UUID       `json:"user_id,omitempty"`
	User   GetUserResponse `json:"user,omitempty"`

	PlanId string       `json:"plan_id,omitempty"`
	Plans  entity.Plans `json:"plan,omitempty"`

	Name        string `json:"name"`
	PhoneNumber string `json:"phone_number"`

	Mealtypes    []string `json:"mealtype"`
	DeliveryDays []string `json:"delivery_days"`
	Allergies    []string `json:"allergies"`

	TotalPrice float64 `json:"total_price"`

	Status         string     `json:"status"`
	PauseStartDate *time.Time `json:"pause_start_date"`
	PauseEndDate   *time.Time `json:"pause_end_date"`

	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}
