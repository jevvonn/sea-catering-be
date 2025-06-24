package entity

import (
	"time"

	"github.com/google/uuid"
)

type Testimonial struct {
	ID      uuid.UUID `gorm:"primaryKey" json:"id,omitempty"`
	Name    string    `gorm:"type:varchar(255);not null" json:"name,omitempty"`
	Message string    `gorm:"type:text;not null;" json:"message,omitempty"`
	Rating  float64   `gorm:"type:float;not null;" json:"rating,omitempty"`

	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}
