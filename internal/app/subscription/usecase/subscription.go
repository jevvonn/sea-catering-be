package usecase

import (
	"errors"
	"strings"
	"time"

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
	GetSpecific(ctx *fiber.Ctx) (dto.GetSubscriptionResponse, error)
	CreateSubscription(ctx *fiber.Ctx, req dto.CreateSubscriptionRequest) error
	UpdateSubscription(ctx *fiber.Ctx, req dto.UpdateSubscriptionRequest) error
	GetSubscriptionsReport(ctx *fiber.Ctx) (dto.GetSubscriptionReportResponse, error)
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
		allergies := []string{}
		if sub.Allergies != "" {
			allergies = strings.Split(sub.Allergies, ",")
		}
		isPaused := false
		if sub.PauseStartDate != nil && sub.PauseEndDate != nil {
			if sub.PauseStartDate.Before(time.Now()) && sub.PauseEndDate.After(time.Now()) {
				isPaused = true
			}
		}

		response = append(response, dto.GetSubscriptionResponse{
			ID:     sub.ID,
			UserID: sub.UserID,
			User: dto.GetUserResponse{
				ID:    sub.User.ID,
				Name:  sub.User.Name,
				Email: sub.User.Email,
			},
			PlanId:         sub.PlanId,
			Plans:          sub.Plans,
			Name:           sub.Name,
			PhoneNumber:    sub.PhoneNumber,
			Mealtypes:      strings.Split(sub.Mealtypes, ","),
			DeliveryDays:   strings.Split(sub.DeliveryDays, ","),
			Allergies:      allergies,
			TotalPrice:     sub.TotalPrice,
			Status:         sub.Status,
			PauseStartDate: sub.PauseStartDate,
			PauseEndDate:   sub.PauseEndDate,
			CreatedAt:      sub.CreatedAt,
			UpdatedAt:      sub.UpdatedAt,
			IsPaused:       isPaused,
		})
	}

	return response, nil
}

func (u *SubscriptionUsecase) GetSpecific(ctx *fiber.Ctx) (dto.GetSubscriptionResponse, error) {
	userId := ctx.Locals("userId").(string)
	role := ctx.Locals("role").(string)
	param := ctx.Params("id")

	subscriptionId, err := uuid.Parse(param)
	if err != nil {
		return dto.GetSubscriptionResponse{}, errors.New("invalid subscription ID format")
	}

	subscription := entity.Subscription{
		ID: subscriptionId,
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

	allergies := []string{}
	if result.Allergies != "" {
		allergies = strings.Split(result.Allergies, ",")
	}
	isPaused := false
	if result.PauseStartDate != nil && result.PauseEndDate != nil {
		if result.PauseStartDate.Before(time.Now()) && result.PauseEndDate.After(time.Now()) {
			isPaused = true
		}
	}

	response := dto.GetSubscriptionResponse{
		ID:     result.ID,
		UserID: result.UserID,
		User: dto.GetUserResponse{
			ID:    result.User.ID,
			Name:  result.User.Name,
			Email: result.User.Email,
		},
		PlanId:         result.PlanId,
		Plans:          result.Plans,
		Name:           result.Name,
		PhoneNumber:    result.PhoneNumber,
		Mealtypes:      strings.Split(result.Mealtypes, ","),
		DeliveryDays:   strings.Split(result.DeliveryDays, ","),
		Allergies:      allergies,
		TotalPrice:     result.TotalPrice,
		Status:         result.Status,
		PauseStartDate: result.PauseStartDate,
		PauseEndDate:   result.PauseEndDate,
		CreatedAt:      result.CreatedAt,
		UpdatedAt:      result.UpdatedAt,
		IsPaused:       isPaused,
	}

	return response, nil
}

func (u *SubscriptionUsecase) CreateSubscription(ctx *fiber.Ctx, req dto.CreateSubscriptionRequest) error {
	userId := ctx.Locals("userId").(string)

	plans, err := u.plansRepo.GetSpecificPlans(entity.Plans{
		ID: req.PlanId,
	})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("plan not found")
		} else {
			return err
		}
	}

	checkedSub, err := u.subRepo.GetSpecific(entity.Subscription{
		UserID: uuid.MustParse(userId),
		PlanId: req.PlanId,
		Status: constant.SubscriptionStatusActive,
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

	totalPrice := plans.Price * float64(mealTypeLength) * float64(deliveryDaysLength) * constant.SubscriptionTAX

	subscription := entity.Subscription{
		ID:           uuid.New(),
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
	param := ctx.Params("id")

	subscriptionId, err := uuid.Parse(param)
	if err != nil {
		return errors.New("invalid subscription ID format")
	}

	subscription, err := u.subRepo.GetSpecific(entity.Subscription{
		ID: subscriptionId,
	})
	if err != nil {
		return err
	}

	if subscription.UserID != uuid.MustParse(userId) && role != constant.RoleAdmin {
		return errors.New("unauthorized access to subscription")
	}

	if subscription.Status == constant.SubscriptionStatusCancelled {
		return errors.New("subscription is already cancelled")
	}

	pauseStartDate := new(time.Time)
	pauseEndDate := new(time.Time)
	parseLayout := "02-01-2006"

	if req.PauseStartDate != "" {
		if req.PauseEndDate == "" {
			return errors.New("pause end date is required when pause start date is provided")
		}

		parsedPauseStartDate, err := time.Parse(parseLayout, req.PauseStartDate)
		if err != nil {
			return errors.New("invalid pause start date format")
		}

		parsedPauseEndDate, err := time.Parse(parseLayout, req.PauseEndDate)
		if err != nil {
			return errors.New("invalid pause end date format")
		}

		now := time.Now()
		startOfToday := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

		if parsedPauseStartDate.Add(1 * time.Hour).Before(startOfToday) {
			return errors.New("pause start date cannot be in the past")
		}

		if parsedPauseStartDate.After(parsedPauseEndDate) {
			return errors.New("pause start date cannot be after pause end date")
		}

		pauseStartDate = &parsedPauseStartDate
		pauseEndDate = &parsedPauseEndDate
	} else {
		pauseStartDate = nil
		pauseEndDate = nil
	}

	if req.Status == constant.SubscriptionStatusActive {
		pauseStartDate = nil
		pauseEndDate = nil
	}

	subUpdate := entity.Subscription{
		ID:             subscription.ID,
		Name:           req.Name,
		PhoneNumber:    req.PhoneNumber,
		Status:         req.Status,
		PauseStartDate: pauseStartDate,
		PauseEndDate:   pauseEndDate,
	}

	if err := u.subRepo.UpdateSubscription(subUpdate); err != nil {
		return err
	}

	return nil
}

func (u *SubscriptionUsecase) GetSubscriptionsReport(ctx *fiber.Ctx) (dto.GetSubscriptionReportResponse, error) {
	queryStartDate := ctx.Query("start_date")
	queryEndDate := ctx.Query("end_date")

	startDate := new(time.Time)
	endDate := new(time.Time)
	parseLayout := "02-01-2006"

	if queryStartDate != "" && queryEndDate != "" {
		parsedStarDate, err := time.Parse(parseLayout, queryStartDate)
		if err != nil {
			return dto.GetSubscriptionReportResponse{}, errors.New("invalid start date format, expected dd-mm-yyyy")
		}
		parsedEndDate, err := time.Parse(parseLayout, queryEndDate)
		if err != nil {
			return dto.GetSubscriptionReportResponse{}, errors.New("invalid end date format, expected dd-mm-yyyy")
		}

		if parsedStarDate.After(parsedEndDate) {
			return dto.GetSubscriptionReportResponse{}, errors.New("start date cannot be after end date")
		}

		startDate = &parsedStarDate
		endDate = &parsedEndDate
	} else {
		now := time.Now()
		daysEnd := now.Add(24 * time.Hour * 30)
		startDate = &now
		endDate = &daysEnd
	}

	// If no date range is provided, use the current date as both start and end date
	getActiveSubscriptions, err := u.subRepo.GetActiveSubscriptions(nil, nil)
	if err != nil {
		return dto.GetSubscriptionReportResponse{}, err
	}

	// Get active subscriptions within the specified date range
	getActiveSubscriptionsByDate, err := u.subRepo.GetActiveSubscriptions(startDate, endDate)
	if err != nil {
		return dto.GetSubscriptionReportResponse{}, err
	}

	allActiveSubscriptions := len(getActiveSubscriptions)
	activeSubscriptionsByDate := len(getActiveSubscriptionsByDate)

	totalRevenue := 0.0
	for _, sub := range getActiveSubscriptions {
		totalRevenue += sub.TotalPrice
	}

	totalRevenueByDate := 0.0
	for _, sub := range getActiveSubscriptionsByDate {
		totalRevenueByDate += sub.TotalPrice
	}

	return dto.GetSubscriptionReportResponse{
		ActiveSubscriptionsByDate: activeSubscriptionsByDate,
		TotalRevenue:              totalRevenue,
		TotalActiveSubscriptions:  allActiveSubscriptions,
		TotalRevenueByDate:        totalRevenueByDate,
	}, nil
}
