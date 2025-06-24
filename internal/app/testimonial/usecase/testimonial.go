package usecase

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jevvonn/sea-catering-be/internal/app/testimonial/repository"
	"github.com/jevvonn/sea-catering-be/internal/domain/dto"
	"github.com/jevvonn/sea-catering-be/internal/domain/entity"
)

type TestimonialUsecaseItf interface {
	GetTestimonials(ctx *fiber.Ctx, query dto.GetTestimonialQuery) (dto.GetTestimonialsResponse, error)
	CreateTestimonial(req dto.TestimonialRequest) error
}

type TestimonialUsecase struct {
	testimonialRepo repository.TestimonialPostgreSQLItf
}

func NewTestimonialUsecase(testimonialRepo repository.TestimonialPostgreSQLItf) TestimonialUsecaseItf {
	return &TestimonialUsecase{testimonialRepo}
}

func (u *TestimonialUsecase) GetTestimonials(ctx *fiber.Ctx, query dto.GetTestimonialQuery) (dto.GetTestimonialsResponse, error) {
	if query.Limit <= 0 {
		query.Limit = 10
	}

	if query.Page <= 0 {
		query.Page = 1
	}

	testimonials, err := u.testimonialRepo.GetTestimonials(query)
	if err != nil {
		return dto.GetTestimonialsResponse{}, err
	}

	return dto.GetTestimonialsResponse{
		Testimonials: testimonials,
		Total:        len(testimonials),
		Page:         query.Page,
		Limit:        query.Limit,
	}, nil
}

func (u *TestimonialUsecase) CreateTestimonial(req dto.TestimonialRequest) error {
	testimonial := entity.Testimonial{
		Name:    req.Name,
		Message: req.Message,
		Rating:  req.Rating,
	}

	err := u.testimonialRepo.CreateTestimonial(testimonial)
	if err != nil {
		return err
	}

	return nil
}
