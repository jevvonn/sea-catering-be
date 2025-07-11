package rest

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jevvonn/sea-catering-be/internal/app/plans/usecase"
	"github.com/jevvonn/sea-catering-be/internal/constant"
	"github.com/jevvonn/sea-catering-be/internal/domain/dto"
	"github.com/jevvonn/sea-catering-be/internal/infra/validator"
	"github.com/jevvonn/sea-catering-be/internal/middleware"
	"github.com/jevvonn/sea-catering-be/internal/models"
)

type PlansHandler struct {
	plansUsecase usecase.PlansUsecaseItf
	validator    validator.ValidationService
}

func NewPlansHandler(
	router fiber.Router,
	plansUsecase usecase.PlansUsecaseItf,
	validator validator.ValidationService,
) {
	handler := PlansHandler{plansUsecase, validator}

	router.Get("/plans", handler.GetPlans)
	router.Put("/plans/:id", middleware.Authenticated, middleware.RequireRoles(constant.RoleAdmin), handler.UpdatePlan)
}

// @Tags         Plans
// @Summary      Get All Meal Plans
// @Accept       json
// @Produce      json
// @Router       /plans [get]
// @Success      200  {object}  models.JSONResponseModel{data=[]entity.Plans}
// @Failure      400  {object}  models.JSONResponseModel
func (h *PlansHandler) GetPlans(ctx *fiber.Ctx) error {
	plans, err := h.plansUsecase.GetPlans()
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(
			models.JSONResponseModel{
				Message: "Invalid Request",
				Errors:  err.Error(),
			},
		)
	}

	return ctx.Status(fiber.StatusOK).JSON(
		models.JSONResponseModel{
			Message: "Plans retrieved successfully",
			Data:    plans,
		},
	)
}

// @Tags         Plans
// @Summary      Update a Testimonial
// @Accept       json
// @Produce      json
// @Param        plansId path string true "Plans ID"
// @Param        request  body  dto.UpdatePlansRequest  true  "Request body"
// @Router       /plans/{plansId} [put]
// @Security     BearerAuth
// @Success      200  {object}  models.JSONResponseModel
// @Failure      400  {object}  models.JSONResponseModel
func (h *PlansHandler) UpdatePlan(ctx *fiber.Ctx) error {
	var req dto.UpdatePlansRequest
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

	if err := h.plansUsecase.UpdatePlan(ctx, req); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(
			models.JSONResponseModel{
				Message: "Failed to update plan",
				Errors:  err.Error(),
			},
		)
	}

	return ctx.Status(fiber.StatusOK).JSON(
		models.JSONResponseModel{
			Message: "Plan updated successfully",
		},
	)
}
