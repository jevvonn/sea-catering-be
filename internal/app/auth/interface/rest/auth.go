package rest

import (
	"github.com/jevvonn/sea-catering-be/internal/app/auth/usecase"
	"github.com/jevvonn/sea-catering-be/internal/domain/dto"
	"github.com/jevvonn/sea-catering-be/internal/infra/validator"
	"github.com/jevvonn/sea-catering-be/internal/middleware"
	"github.com/jevvonn/sea-catering-be/internal/models"

	"github.com/gofiber/fiber/v2"
)

type AuthHandler struct {
	authUsecase usecase.AuthUsecaseItf
	validator   validator.ValidationService
}

func NewAuthHandler(
	router fiber.Router,
	authUsecase usecase.AuthUsecaseItf,
	validator validator.ValidationService,
) {
	handler := AuthHandler{authUsecase, validator}

	router.Post("/auth/login", handler.Login)
	router.Post("/auth/register", handler.Register)

	router.Get("/auth/session", middleware.Authenticated, handler.Session)
}

// @Tags         Auth
// @Summary      Login as User
// @Accept       json
// @Produce      json
// @Param        request  body  dto.LoginRequest  true  "Request body"
// @Router       /auth/login [post]
// @Success      200  {object}  dto.LoginResponse
// @Failure      400  {object}  models.JSONResponseModel
func (h *AuthHandler) Login(ctx *fiber.Ctx) error {
	var req dto.LoginRequest
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

	res, err := h.authUsecase.Login(ctx, req)
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
			Message: "User Logged In Successfully",
			Data:    res,
		},
	)
}

// @Tags         Auth
// @Summary      Register as User
// @Accept       json
// @Produce      json
// @Param        request  body  dto.RegisterRequest  true  "Request body"
// @Router       /auth/register [post]
// @Success      200  {object}  models.JSONResponseModel
// @Failure      400  {object}  models.JSONResponseModel
func (h *AuthHandler) Register(ctx *fiber.Ctx) error {
	var req dto.RegisterRequest
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

	err = h.authUsecase.Register(ctx, req)
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
			Message: "User Registered Successfully",
		},
	)
}

// @Tags         Auth
// @Summary      Get User Session
// @Produce      json
// @Router       /auth/session [get]
// @Security     BearerAuth
// @Success      200  {object}  dto.SessionResponse
// @Failure      401  {object}  models.JSONResponseModel
// @Failure      400  {object}  models.JSONResponseModel
func (h *AuthHandler) Session(ctx *fiber.Ctx) error {
	res, err := h.authUsecase.Session(ctx)
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
			Message: "Session Data",
			Data:    res,
		},
	)
}
