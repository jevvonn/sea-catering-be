package entity

import (
	"time"

	"github.com/google/uuid"
)

type Subscription struct {
	ID uuid.UUID `gorm:"primaryKey" json:"id,omitempty"`

	UserID uuid.UUID `gorm:"type:uuid;not null;" json:"user_id,omitempty"`
	User   User      `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"user,omitempty"`

	PlanId string `gorm:"type:varchar(10);not null;" json:"plan_id,omitempty"`
	Plans  Plans  `gorm:"foreignKey:PlanId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"plan,omitempty"`

	Name        string `gorm:"type:varchar(255);not null" json:"name,omitempty"`
	PhoneNumber string `gorm:"type:varchar(255);not null" json:"phone_number,omitempty"`

	Mealtypes    string `gorm:"type:text;not null" json:"mealtype,omitempty"`
	DeliveryDays string `gorm:"type:text;not null" json:"delivery_days,omitempty"`
	Allergies    string `gorm:"type:text;not null" json:"allergies,omitempty"`

	TotalPrice float64 `gorm:"type:decimal(10,2);not null" json:"total_price,omitempty"`

	Status         string     `gorm:"type:varchar(50);not null;default:'ACTIVE'" json:"status,omitempty"`
	PauseStartDate *time.Time `gorm:"type:timestamp" json:"pause_start_date,omitempty"`
	PauseEndDate   *time.Time `gorm:"type:timestamp" json:"pause_end_date,omitempty"`

	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}
