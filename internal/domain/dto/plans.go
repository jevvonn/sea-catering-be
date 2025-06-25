package dto

type UpdatePlansRequest struct {
	Name     string  `json:"name,omitempty"`
	Slogan   string  `json:"slogan,omitempty"`
	Price    float64 `json:"price,omitempty" validate:"numeric,min=1"`
	Features string  `json:"features,omitempty"`
}
