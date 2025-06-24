package bootstrap

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/jevvonn/sea-catering-be/config"
	userRepo "github.com/jevvonn/sea-catering-be/internal/app/user/repository"

	authHandler "github.com/jevvonn/sea-catering-be/internal/app/auth/interface/rest"
	authUsecase "github.com/jevvonn/sea-catering-be/internal/app/auth/usecase"

	"github.com/jevvonn/sea-catering-be/internal/infra/postgresql"
	"github.com/jevvonn/sea-catering-be/internal/infra/validator"

	"github.com/gofiber/swagger"
	_ "github.com/jevvonn/sea-catering-be/docs"
)

const idleTimeout = 5 * time.Second

func Start() error {
	app := fiber.New(fiber.Config{
		IdleTimeout: idleTimeout,
	})
	conf := config.New()

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable TimeZone=Asia/Jakarta",
		conf.DbHost,
		conf.DbPort,
		conf.DbUser,
		conf.DbPassword,
		conf.DbName,
	)

	db, err := postgresql.New(dsn)
	if err != nil {
		panic(err)
	}

	validator := validator.NewValidator()

	// For migrating the database by command
	CommandHandler(db)

	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "Welcome to Sea Catering API",
		})
	})

	app.Get("/docs/*", swagger.HandlerDefault)

	apiRouter := app.Group("/api")

	userRepo := userRepo.NewUserPostgreSQL(db)

	authUsecase := authUsecase.NewAuthUsecase(userRepo)

	authHandler.NewAuthHandler(apiRouter, authUsecase, validator)

	addr := fmt.Sprintf("localhost:%s", conf.AppPort)
	if conf.AppEnv == "production" {
		addr = fmt.Sprintf("0.0.0.0:%s", conf.AppPort)
	}

	return app.Listen(addr)
}
