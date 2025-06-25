package dto

type UpdatePlansRequest struct {
	Name     string  `json:"name,omitempty" validate:"required"`
	Slogan   string  `json:"slogan,omitempty" validate:"required"`
	Price    float64 `json:"price,omitempty" validate:"required,numeric,min=0"`
	Features string  `json:"features,omitempty" validate:"required"`
}
