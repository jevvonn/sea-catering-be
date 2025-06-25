package dto

import "github.com/google/uuid"

type GetUserResponse struct {
	ID    uuid.UUID `json:"id,omitempty"`
	Email string    `json:"email,omitempty"`
	Name  string    `json:"name,omitempty"`
}
