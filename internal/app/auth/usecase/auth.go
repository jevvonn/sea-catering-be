package usecase

import (
	"errors"

	userRepo "github.com/jevvonn/sea-catering-be/internal/app/user/repository"
	"github.com/jevvonn/sea-catering-be/internal/domain/dto"
	"github.com/jevvonn/sea-catering-be/internal/domain/entity"
	"github.com/jevvonn/sea-catering-be/internal/infra/jwt"
	utils "github.com/jevvonn/sea-catering-be/internal/lib"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AuthUsecaseItf interface {
	Register(ctx *fiber.Ctx, req dto.RegisterRequest) error
	Login(ctx *fiber.Ctx, req dto.LoginRequest) (dto.LoginResponse, error)
	Session(ctx *fiber.Ctx) (dto.SessionResponse, error)
}

type AuthUsecase struct {
	userRepo userRepo.UserPostgreSQLItf
}

func NewAuthUsecase(
	userRepo userRepo.UserPostgreSQLItf,
) AuthUsecaseItf {
	return &AuthUsecase{userRepo}
}

func (u *AuthUsecase) Register(ctx *fiber.Ctx, req dto.RegisterRequest) error {
	// Check if username already exists
	user, err := u.userRepo.GetSpecificUser(entity.User{
		Email: req.Email,
	})

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	if user.ID != uuid.Nil {
		return errors.New("email already exists")
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return err
	}

	// Create user
	user = entity.User{
		ID:       uuid.New(),
		Name:     req.Name,
		Email:    req.Email,
		Password: hashedPassword,
	}
	err = u.userRepo.CreateUser(user)

	if err != nil {
		return err
	}

	return nil
}

func (u *AuthUsecase) Login(ctx *fiber.Ctx, req dto.LoginRequest) (dto.LoginResponse, error) {
	// Check if username exists
	user, err := u.userRepo.GetSpecificUser(entity.User{
		Email: req.Email,
	})

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return dto.LoginResponse{}, err
	}

	if user.ID == uuid.Nil {
		return dto.LoginResponse{}, errors.New("email or password is incorrect")
	}

	// Check password
	if !utils.VerifyPassword(req.Password, user.Password) {
		return dto.LoginResponse{}, errors.New("email or password is incorrect")
	}
	// Create Jwt token
	token, err := jwt.CreateAuthToken(user.ID.String(), user.Email, user.Role)

	if err != nil {
		return dto.LoginResponse{}, err
	}

	return dto.LoginResponse{
		UserId: user.ID.String(),
		Token:  token,
	}, nil
}

func (u *AuthUsecase) Session(ctx *fiber.Ctx) (dto.SessionResponse, error) {
	userId := ctx.Locals("userId").(string)

	uuidUser, err := uuid.Parse(userId)
	if err != nil {
		return dto.SessionResponse{}, err
	}

	user, err := u.userRepo.GetSpecificUser(entity.User{
		ID: uuidUser,
	})
	if err != nil {
		return dto.SessionResponse{}, err
	}

	return dto.SessionResponse{
		ID:    user.ID.String(),
		Name:  user.Name,
		Email: user.Email,
		Role:  user.Role,
	}, nil
}
