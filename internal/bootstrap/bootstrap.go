package bootstrap

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/jevvonn/sea-catering-be/config"
	"github.com/jevvonn/sea-catering-be/internal/infra/postgresql"
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

	// For migrating the database by command
	CommandHandler(db)

	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "Welcome to Sea Catering API",
		})
	})

	addr := fmt.Sprintf("localhost:%s", conf.AppPort)
	if conf.AppEnv == "production" {
		addr = fmt.Sprintf("0.0.0.0:%s", conf.AppPort)
	}

	return app.Listen(addr)
}
