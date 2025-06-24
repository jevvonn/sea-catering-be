package dto

import "github.com/jevvonn/sea-catering-be/internal/domain/entity"

type TestimonialRequest struct {
	Name    string  `json:"name,omitempty" validate:"required"`
	Message string  `json:"message,omitempty" validate:"required"`
	Rating  float64 `json:"rating,omitempty" validate:"required,min=1,max=5,numeric"`
}

type GetTestimonialQuery struct {
	Limit int `query:"limit" validate:"omitempty,numeric,min=1,max=100"`
	Page  int `query:"page" validate:"omitempty,numeric,min=1"`
}

type GetTestimonialsResponse struct {
	Testimonials []entity.Testimonial `json:"testimonials"`
	Total        int                  `json:"total,omitempty"`
	Page         int                  `json:"page,omitempty"`
	Limit        int                  `json:"limit,omitempty"`
}
