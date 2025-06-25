package usecase

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jevvonn/sea-catering-be/internal/app/plans/repository"
	"github.com/jevvonn/sea-catering-be/internal/domain/dto"
	"github.com/jevvonn/sea-catering-be/internal/domain/entity"
)

type PlansUsecaseItf interface {
	GetPlans() ([]entity.Plans, error)
	UpdatePlan(ctx *fiber.Ctx, plan dto.UpdatePlansRequest) error
}

type PlansUsecase struct {
	plansRepo repository.PlansPostgreSQLItf
}

func NewPlansUsecase(plansRepo repository.PlansPostgreSQLItf) PlansUsecaseItf {
	return &PlansUsecase{plansRepo}
}

func (u *PlansUsecase) GetPlans() ([]entity.Plans, error) {
	plans, err := u.plansRepo.GetPlans()
	if err != nil {
		return nil, err
	}

	return plans, nil
}

func (u *PlansUsecase) UpdatePlan(ctx *fiber.Ctx, plan dto.UpdatePlansRequest) error {
	planID := ctx.Params("id")

	updatedPlan := entity.Plans{
		ID:       planID,
		Name:     plan.Name,
		Slogan:   plan.Slogan,
		Price:    plan.Price,
		Features: plan.Features,
	}

	if err := u.plansRepo.UpdatePlan(updatedPlan); err != nil {
		return err
	}

	return nil
}
