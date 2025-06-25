package repository

import (
	"github.com/google/uuid"
	"github.com/jevvonn/sea-catering-be/internal/domain/dto"
	"github.com/jevvonn/sea-catering-be/internal/domain/entity"
	"gorm.io/gorm"
)

type TestimonialPostgreSQLItf interface {
	CreateTestimonial(req entity.Testimonial) error
	GetTestimonials(testimonialQuery dto.GetTestimonialQuery) ([]entity.Testimonial, error)
	GetSpecificTestimonial(req entity.Testimonial) (entity.Testimonial, error)
	DeleteTestimonial(req entity.Testimonial) error
}

type TestimonialPostgreSQL struct {
	db *gorm.DB
}

func NewTestimonialPostgreSQL(db *gorm.DB) TestimonialPostgreSQLItf {
	return &TestimonialPostgreSQL{db}
}

func (r *TestimonialPostgreSQL) CreateTestimonial(req entity.Testimonial) error {
	req.ID = uuid.New()
	err := r.db.Create(&req).Error
	if err != nil {
		return err
	}

	return nil
}

func (r *TestimonialPostgreSQL) GetTestimonials(testimonialQuery dto.GetTestimonialQuery) ([]entity.Testimonial, error) {
	var testimonials []entity.Testimonial
	query := r.db.Model(&entity.Testimonial{})

	query = query.Limit(testimonialQuery.Limit)
	query = query.Offset((testimonialQuery.Page - 1) * testimonialQuery.Limit)

	err := query.Find(&testimonials).Error
	if err != nil {
		return []entity.Testimonial{}, err
	}

	return testimonials, nil
}

func (r *TestimonialPostgreSQL) GetSpecificTestimonial(req entity.Testimonial) (entity.Testimonial, error) {
	var testimonial entity.Testimonial
	err := r.db.Where(&req).First(&testimonial).Error

	if err != nil {
		return entity.Testimonial{}, err
	}

	return testimonial, nil
}

func (r *TestimonialPostgreSQL) DeleteTestimonial(req entity.Testimonial) error {
	err := r.db.Delete(&req).Error

	if err != nil {
		return err
	}

	return nil
}
