package bootstrap

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	cors "github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/jevvonn/sea-catering-be/config"

	plansRepo "github.com/jevvonn/sea-catering-be/internal/app/plans/repository"
	subsRepo "github.com/jevvonn/sea-catering-be/internal/app/subscription/repository"
	testimonialRepo "github.com/jevvonn/sea-catering-be/internal/app/testimonial/repository"
	userRepo "github.com/jevvonn/sea-catering-be/internal/app/user/repository"

	authUsecase "github.com/jevvonn/sea-catering-be/internal/app/auth/usecase"
	plansUsecase "github.com/jevvonn/sea-catering-be/internal/app/plans/usecase"
	subsUsecase "github.com/jevvonn/sea-catering-be/internal/app/subscription/usecase"
	testimonialUsecase "github.com/jevvonn/sea-catering-be/internal/app/testimonial/usecase"

	authHandler "github.com/jevvonn/sea-catering-be/internal/app/auth/interface/rest"
	plansHandler "github.com/jevvonn/sea-catering-be/internal/app/plans/interface/rest"
	subsHandler "github.com/jevvonn/sea-catering-be/internal/app/subscription/interface/rest"
	testimonialHandler "github.com/jevvonn/sea-catering-be/internal/app/testimonial/interface/rest"

	"github.com/jevvonn/sea-catering-be/internal/infra/postgresql"
	"github.com/jevvonn/sea-catering-be/internal/infra/validator"

	"github.com/gofiber/fiber/v2/middleware/limiter"

	"github.com/gofiber/swagger"
	_ "github.com/jevvonn/sea-catering-be/docs"
)

const idleTimeout = 5 * time.Second

func Start() error {
	app := fiber.New(fiber.Config{
		IdleTimeout: idleTimeout,
	})

	app.Use(cors.New())

	app.Use(limiter.New(limiter.Config{
		Max:               100,
		Expiration:        30 * time.Second,
		LimiterMiddleware: limiter.SlidingWindow{},
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"error": "Too many requests, please try again later.",
			})
		},
	}))

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
			"docs":    "/docs",
		})
	})

	app.Get("/docs/*", swagger.HandlerDefault)

	apiRouter := app.Group("/api")

	userRepo := userRepo.NewUserPostgreSQL(db)
	testimonialRepo := testimonialRepo.NewTestimonialPostgreSQL(db)
	plansRepo := plansRepo.NewPlansPostgreSQL(db)
	subsRepo := subsRepo.NewSubscriptionPostgreSQL(db)

	authUsecase := authUsecase.NewAuthUsecase(userRepo)
	testimonialUsecase := testimonialUsecase.NewTestimonialUsecase(testimonialRepo)
	plansUsecase := plansUsecase.NewPlansUsecase(plansRepo)
	subsUsecase := subsUsecase.NewSubscriptionUsecase(subsRepo, plansRepo)

	authHandler.NewAuthHandler(apiRouter, authUsecase, validator)
	testimonialHandler.NewTestimonialHandler(apiRouter, testimonialUsecase, validator)
	plansHandler.NewPlansHandler(apiRouter, plansUsecase, validator)
	subsHandler.NewSubscriptionHandler(apiRouter, subsUsecase, validator)

	addr := fmt.Sprintf("localhost:%s", conf.AppPort)
	if conf.AppEnv == "production" {
		addr = fmt.Sprintf("0.0.0.0:%s", conf.AppPort)
	}

	return app.Listen(addr)
}
