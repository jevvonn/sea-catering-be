package usecase

import (
	"errors"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	plansRepo "github.com/jevvonn/sea-catering-be/internal/app/plans/repository"
	subRepo "github.com/jevvonn/sea-catering-be/internal/app/subscription/repository"
	"github.com/jevvonn/sea-catering-be/internal/constant"
	"github.com/jevvonn/sea-catering-be/internal/domain/dto"
	"github.com/jevvonn/sea-catering-be/internal/domain/entity"
	"gorm.io/gorm"
)

type SubscriptionUsecaseItf interface {
	GetSubscriptions(ctx *fiber.Ctx) ([]dto.GetSubscriptionResponse, error)
	GetSpecific(ctx *fiber.Ctx, subscriptionId string) (dto.GetSubscriptionResponse, error)
	CreateSubscription(ctx *fiber.Ctx, req dto.CreateSubscriptionRequest) error
	UpdateSubscription(ctx *fiber.Ctx, req dto.UpdateSubscriptionRequest) error
}

type SubscriptionUsecase struct {
	subRepo   subRepo.SubscriptionPostgreSQLItf
	plansRepo plansRepo.PlansPostgreSQLItf
}

func NewSubscriptionUsecase(
	subRepo subRepo.SubscriptionPostgreSQLItf,
	plansRepo plansRepo.PlansPostgreSQLItf,
) SubscriptionUsecaseItf {
	return &SubscriptionUsecase{subRepo, plansRepo}
}

func (u *SubscriptionUsecase) GetSubscriptions(ctx *fiber.Ctx) ([]dto.GetSubscriptionResponse, error) {
	userId := ctx.Locals("userId").(string)
	role := ctx.Locals("role").(string)

	condition := entity.Subscription{}
	if role != constant.RoleAdmin {
		condition.UserID = uuid.MustParse(userId)
	}

	subscriptions, err := u.subRepo.GetSubscriptions(condition)
	if err != nil {
		return nil, err
	}

	var response []dto.GetSubscriptionResponse
	for _, sub := range subscriptions {
		response = append(response, dto.GetSubscriptionResponse{
			ID:             sub.ID,
			UserID:         sub.UserID,
			User:           sub.User,
			PlanId:         sub.PlanId,
			Plans:          sub.Plans,
			Name:           sub.Name,
			PhoneNumber:    sub.PhoneNumber,
			Mealtypes:      strings.Split(sub.Mealtypes, ","),
			DeliveryDays:   strings.Split(sub.DeliveryDays, ","),
			Allergies:      strings.Split(sub.Allergies, ","),
			TotalPrice:     sub.TotalPrice,
			Status:         sub.Status,
			PauseStartDate: sub.PauseStartDate,
			PauseEndDate:   sub.PauseEndDate,
			CreatedAt:      sub.CreatedAt,
			UpdatedAt:      sub.UpdatedAt,
		})
	}

	return response, nil
}

func (u *SubscriptionUsecase) GetSpecific(ctx *fiber.Ctx, subscriptionId string) (dto.GetSubscriptionResponse, error) {
	userId := ctx.Locals("userId").(string)
	role := ctx.Locals("role").(string)

	subscription := entity.Subscription{
		ID: uuid.MustParse(subscriptionId),
	}

	if role != constant.RoleAdmin {
		subscription.UserID = uuid.MustParse(userId)
	}

	result, err := u.subRepo.GetSpecific(subscription)
	if err != nil {
		return dto.GetSubscriptionResponse{}, err
	}

	if result.UserID != uuid.MustParse(userId) && role != constant.RoleAdmin {
		return dto.GetSubscriptionResponse{}, errors.New("unauthorized access to subscription")
	}

	response := dto.GetSubscriptionResponse{
		ID:             result.ID,
		UserID:         result.UserID,
		User:           result.User,
		PlanId:         result.PlanId,
		Plans:          result.Plans,
		Name:           result.Name,
		PhoneNumber:    result.PhoneNumber,
		Mealtypes:      strings.Split(result.Mealtypes, ","),
		DeliveryDays:   strings.Split(result.DeliveryDays, ","),
		Allergies:      strings.Split(result.Allergies, ","),
		TotalPrice:     result.TotalPrice,
		Status:         result.Status,
		PauseStartDate: result.PauseStartDate,
		PauseEndDate:   result.PauseEndDate,
		CreatedAt:      result.CreatedAt,
		UpdatedAt:      result.UpdatedAt,
	}

	return response, nil
}

func (u *SubscriptionUsecase) CreateSubscription(ctx *fiber.Ctx, req dto.CreateSubscriptionRequest) error {
	userId := ctx.Locals("userId").(string)

	plans, err := u.plansRepo.GetSpecificPlans(entity.Plans{
		ID: req.PlanId,
	})
	if err != nil {
		return err
	}

	checkedSub, err := u.subRepo.GetSpecific(entity.Subscription{
		UserID: uuid.MustParse(userId),
		PlanId: req.PlanId,
	})
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	// Check if the user already has an active subscription for the same plan
	if checkedSub.ID != uuid.Nil {
		return errors.New("you already has an active subscription for this plan")
	}

	mealTypeLength := len(req.Mealtypes)
	deliveryDaysLength := len(req.DeliveryDays)
	allergiesLength := len(req.Allergies)

	totalPrice := plans.Price * float64(mealTypeLength) * float64(deliveryDaysLength) * float64(allergiesLength) * constant.SubscriptionTAX

	subscription := entity.Subscription{
		UserID:       uuid.MustParse(userId),
		PlanId:       req.PlanId,
		Name:         req.Name,
		PhoneNumber:  req.PhoneNumber,
		Mealtypes:    strings.Join(req.Mealtypes, ","),
		DeliveryDays: strings.Join(req.DeliveryDays, ","),
		Allergies:    strings.Join(req.Allergies, ","),
		Status:       constant.SubscriptionStatusActive,
		TotalPrice:   totalPrice,
	}

	if err := u.subRepo.CreateSubscription(subscription); err != nil {
		return err
	}

	return nil
}

func (u *SubscriptionUsecase) UpdateSubscription(ctx *fiber.Ctx, req dto.UpdateSubscriptionRequest) error {
	userId := ctx.Locals("userId").(string)
	role := ctx.Locals("role").(string)
	subscriptionId := ctx.Params("subscriptionId")

	subscription, err := u.subRepo.GetSpecific(entity.Subscription{
		ID: uuid.MustParse(subscriptionId),
	})
	if err != nil {
		return err
	}

	if subscription.UserID != uuid.MustParse(userId) && role != constant.RoleAdmin {
		return errors.New("unauthorized access to subscription")
	}

	subUpdate := entity.Subscription{
		ID:             subscription.ID,
		Name:           req.Name,
		PhoneNumber:    req.PhoneNumber,
		Status:         req.Status,
		PauseStartDate: &req.PauseStartDate,
		PauseEndDate:   &req.PauseEndDate,
	}

	if err := u.subRepo.UpdateSubscription(subUpdate); err != nil {
		return err
	}

	return nil
}
