package rest

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jevvonn/sea-catering-be/internal/app/subscription/usecase"
	"github.com/jevvonn/sea-catering-be/internal/constant"
	"github.com/jevvonn/sea-catering-be/internal/domain/dto"
	"github.com/jevvonn/sea-catering-be/internal/infra/validator"
	"github.com/jevvonn/sea-catering-be/internal/middleware"
	"github.com/jevvonn/sea-catering-be/internal/models"
)

type SubscriptionHandler struct {
	subUsecase usecase.SubscriptionUsecaseItf
	validator  validator.ValidationService
}

func NewSubscriptionHandler(
	router fiber.Router,
	subUsecase usecase.SubscriptionUsecaseItf,
	validator validator.ValidationService,
) {
	handler := SubscriptionHandler{subUsecase, validator}

	router.Get("/subscriptions", middleware.Authenticated, handler.GetSubscriptions)
	router.Get("/subscriptions/report", middleware.Authenticated, middleware.RequireRoles(constant.RoleAdmin), handler.GetSubscriptionsReport)

	router.Get("/subscriptions/:id", middleware.Authenticated, handler.GetSpecific)
	router.Post("/subscriptions", middleware.Authenticated, handler.CreateSubscription)
	router.Put("/subscriptions/:id", middleware.Authenticated, handler.UpdateSubscription)
}

// @Tags         Subscription
// @Summary      Get All My Subscriptions
// @Accept       json
// @Produce      json
// @Router       /subscriptions [get]
// @Security     BearerAuth
// @Success      200  {object}  models.JSONResponseModel{data=[]dto.GetSubscriptionResponse}
// @Failure      400  {object}  models.JSONResponseModel
func (h *SubscriptionHandler) GetSubscriptions(ctx *fiber.Ctx) error {
	subscriptions, err := h.subUsecase.GetSubscriptions(ctx)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(
			models.JSONResponseModel{
				Message: "Failed to retrieve subscriptions",
				Errors:  err.Error(),
			},
		)
	}

	return ctx.Status(fiber.StatusOK).JSON(
		models.JSONResponseModel{
			Message: "Subscriptions retrieved successfully",
			Data:    subscriptions,
		},
	)
}

// @Tags         Subscription
// @Summary      Get Spesicfic Subscription
// @Accept       json
// @Produce      json
// @Param        subscriptionId path string true "Subscription ID"
// @Router       /subscriptions/{subscriptionId} [get]
// @Security     BearerAuth
// @Success      200  {object}  models.JSONResponseModel{data=dto.GetSubscriptionResponse}
// @Failure      400  {object}  models.JSONResponseModel
func (h *SubscriptionHandler) GetSpecific(ctx *fiber.Ctx) error {
	subscriptionId := ctx.Params("id")
	if subscriptionId == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(
			models.JSONResponseModel{
				Message: "Subscription ID is required",
			},
		)
	}

	subscription, err := h.subUsecase.GetSpecific(ctx)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(
			models.JSONResponseModel{
				Message: "Failed to retrieve subscription",
				Errors:  err.Error(),
			},
		)
	}

	return ctx.Status(fiber.StatusOK).JSON(
		models.JSONResponseModel{
			Message: "Subscription retrieved successfully",
			Data:    subscription,
		},
	)
}

// @Tags         Subscription
// @Summary      Create Subscription
// @Accept       json
// @Produce      json
// @Param        request body dto.CreateSubscriptionRequest true "Request body"
// @Router       /subscriptions [post]
// @Security     BearerAuth
// @Success      200  {object}  models.JSONResponseModel
// @Failure      400  {object}  models.JSONResponseModel
func (h *SubscriptionHandler) CreateSubscription(ctx *fiber.Ctx) error {
	var req dto.CreateSubscriptionRequest
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(
			models.JSONResponseModel{
				Message: "Invalid request body",
				Errors:  err.Error(),
			},
		)
	}

	if err := h.validator.Validate(req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(
			err.(*validator.ValidationError),
		)
	}

	if err := h.subUsecase.CreateSubscription(ctx, req); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(
			models.JSONResponseModel{
				Message: "Failed to create subscription",
				Errors:  err.Error(),
			},
		)
	}

	return ctx.Status(fiber.StatusCreated).JSON(
		models.JSONResponseModel{
			Message: "Subscription created successfully",
		},
	)
}

// @Tags         Subscription
// @Summary      Update Subscription
// @Accept       json
// @Produce      json
// @Param        subscriptionId path string true "Subscription ID"
// @Param        request body dto.UpdateSubscriptionRequest true "Request body"
// @Router       /subscriptions/{subscriptionId} [put]
// @Security     BearerAuth
// @Success      200  {object}  models.JSONResponseModel
// @Failure      400  {object}  models.JSONResponseModel
func (h *SubscriptionHandler) UpdateSubscription(ctx *fiber.Ctx) error {
	subscriptionId := ctx.Params("id")
	if subscriptionId == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(
			models.JSONResponseModel{
				Message: "Subscription ID is required",
			},
		)
	}

	var req dto.UpdateSubscriptionRequest
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(
			models.JSONResponseModel{
				Message: "Invalid request body",
				Errors:  err.Error(),
			},
		)
	}

	if err := h.validator.Validate(req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(
			err.(*validator.ValidationError),
		)
	}

	if err := h.subUsecase.UpdateSubscription(ctx, req); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(
			models.JSONResponseModel{
				Message: "Failed to update subscription",
				Errors:  err.Error(),
			},
		)
	}

	return ctx.Status(fiber.StatusOK).JSON(
		models.JSONResponseModel{
			Message: "Subscription updated successfully",
		},
	)
}

// @Tags         Subscription
// @Summary      Get Report Subscription
// @Accept       json
// @Produce      json
// @Param        start_date query string false "e.g 29-06-2025"
// @Param        end_date query string false "e.g 30-06-2025"
// @Router       /subscriptions/report [get]
// @Security     BearerAuth
// @Success      200  {object}  models.JSONResponseModel{data=dto.GetSubscriptionReportResponse}
// @Failure      400  {object}  models.JSONResponseModel
func (h *SubscriptionHandler) GetSubscriptionsReport(ctx *fiber.Ctx) error {
	report, err := h.subUsecase.GetSubscriptionsReport(ctx)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(
			models.JSONResponseModel{
				Message: "Failed to retrieve subscription report",
				Errors:  err.Error(),
			},
		)
	}

	return ctx.Status(fiber.StatusOK).JSON(
		models.JSONResponseModel{
			Message: "Subscription report retrieved successfully",
			Data:    report,
		},
	)
}
