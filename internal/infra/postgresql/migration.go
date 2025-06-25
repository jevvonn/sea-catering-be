package postgresql

import (
	"fmt"

	"github.com/jevvonn/sea-catering-be/internal/domain/entity"
	"gorm.io/gorm"
)

func Migrate(db *gorm.DB, command string) {
	migrator := db.Migrator()
	tables := []any{
		&entity.User{},
		&entity.Testimonial{},
		&entity.Plans{},
		&entity.Subscription{},
	}

	var err error
	if command == "up" {
		err = migrator.AutoMigrate(tables...)
	}

	if command == "down" {
		err = migrator.DropTable(tables...)
	}

	if err != nil {
		panic(err)
	}

	fmt.Printf("Migration %s completed successfully\n", command)
}
