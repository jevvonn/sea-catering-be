package rest

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jevvonn/sea-catering-be/internal/app/testimonial/usecase"
	"github.com/jevvonn/sea-catering-be/internal/domain/dto"
	"github.com/jevvonn/sea-catering-be/internal/infra/validator"
	"github.com/jevvonn/sea-catering-be/internal/models"
)

type TestimonialHandler struct {
	testimonialUsacase usecase.TestimonialUsecaseItf
	validator          validator.ValidationService
}

func NewTestimonialHandler(
	router fiber.Router,
	testimonialUsacase usecase.TestimonialUsecaseItf,
	validator validator.ValidationService,
) {
	handler := TestimonialHandler{testimonialUsacase, validator}

	router.Get("/testimonials", handler.GetTestimonials)
	router.Post("/testimonials", handler.CreateTestimonial)
}

// @Tags         Testimonial
// @Summary      Get Testimonial
// @Accept       json
// @Produce      json
// @Param        limit query int false "Limit"
// @Param        page query int false "Page"
// @Router       /testimonials [get]
// @Success      200  {object}  models.JSONResponseModel
// @Failure      400  {object}  models.JSONResponseModel
func (h *TestimonialHandler) GetTestimonials(ctx *fiber.Ctx) error {
	var req dto.GetTestimonialQuery
	err := ctx.QueryParser(&req)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(
			models.JSONResponseModel{
				Message: "Invalid Request",
				Errors:  err.Error(),
			},
		)
	}

	err = h.validator.Validate(req)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(
			err.(*validator.ValidationError),
		)
	}

	response, err := h.testimonialUsacase.GetTestimonials(ctx, req)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(
			models.JSONResponseModel{
				Message: "Invalid Request",
				Errors:  err.Error(),
			},
		)
	}

	return ctx.Status(fiber.StatusCreated).JSON(
		models.JSONResponseModel{
			Message: "Testimonials Found",
			Data:    response,
		},
	)
}

// @Tags         Testimonial
// @Summary      Create a new Testimonial
// @Accept       json
// @Produce      json
// @Param        request  body  dto.TestimonialRequest  true  "Request body"
// @Router       /testimonials [post]
// @Success      200  {object}  models.JSONResponseModel
// @Failure      400  {object}  models.JSONResponseModel
func (h *TestimonialHandler) CreateTestimonial(ctx *fiber.Ctx) error {
	var req dto.TestimonialRequest
	err := ctx.BodyParser(&req)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(
			models.JSONResponseModel{
				Message: "Invalid Request",
				Errors:  err.Error(),
			},
		)
	}

	err = h.validator.Validate(req)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(
			err.(*validator.ValidationError),
		)
	}

	err = h.testimonialUsacase.CreateTestimonial(req)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(
			models.JSONResponseModel{
				Message: "Invalid Request",
				Errors:  err.Error(),
			},
		)
	}

	return ctx.Status(fiber.StatusCreated).JSON(
		models.JSONResponseModel{
			Message: "Testimonials Created Successfully",
		},
	)
}
