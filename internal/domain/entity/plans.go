package entity

import "time"

type Plans struct {
	ID       string  `gorm:"primaryKey;type:varchar(10)" json:"id,omitempty"`
	Name     string  `gorm:"type:varchar(255);not null" json:"name,omitempty"`
	Slogan   string  `gorm:"type:varchar(255);not null" json:"slogan,omitempty"`
	Price    float64 `gorm:"type:float;not null" json:"price,omitempty"`
	Features string  `gorm:"type:text;not null" json:"features,omitempty"`

	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}
