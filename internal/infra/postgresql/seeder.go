package postgresql

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/jevvonn/sea-catering-be/internal/constant"
	"github.com/jevvonn/sea-catering-be/internal/domain/entity"
	utils "github.com/jevvonn/sea-catering-be/internal/lib"
	"gorm.io/gorm"
)

func Seed(db *gorm.DB) {
	hashedPassword, err := utils.HashPassword("password")
	if err != nil {
		panic(err)
	}

	adminAccount := entity.User{
		ID:       uuid.New(),
		Name:     "Admin",
		Email:    "admin@gmail.com",
		Password: hashedPassword,
		Role:     constant.RoleAdmin,
	}

	userAccount := entity.User{
		ID:       uuid.New(),
		Name:     "User",
		Email:    "user@gmail.com",
		Password: hashedPassword,
		Role:     constant.RoleUser,
	}

	dietPlan := entity.Plans{
		ID:       "diet",
		Name:     "Diet Plan",
		Slogan:   "Light & Nutritious",
		Price:    30000,
		Features: "300-400 calories per meal, High fiber content, Low fat recipes, Portion controlled, Fresh vegetables daily",
	}

	proteinPlan := entity.Plans{
		ID:       "protein",
		Name:     "Protein Plan",
		Slogan:   "Power & Performance",
		Price:    40000,
		Features: "25-35g protein per meal, Lean meat & fish, Post-workout friendly, Balanced macronutrients, Athletic performance focused",
	}

	royalPlan := entity.Plans{
		ID:       "royal",
		Name:     "Royal Plan",
		Slogan:   "Luxury & Elegance",
		Price:    60000,
		Features: "Premium ingredients, Chef-crafted recipes, Restaurant quality, Exclusive menu items",
	}

	err = db.Create(&dietPlan).Error
	if err != nil {
		panic(err)
	}

	err = db.Create(&proteinPlan).Error
	if err != nil {
		panic(err)
	}

	err = db.Create(&royalPlan).Error
	if err != nil {
		panic(err)
	}

	err = db.Create(&adminAccount).Error
	if err != nil {
		panic(err)
	}

	err = db.Create(&userAccount).Error
	if err != nil {
		panic(err)
	}

	fmt.Println("Database seeded successfully with initial data.")
}
