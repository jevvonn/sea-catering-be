package main

import (
	"fmt"

	"github.com/jevvonn/sea-catering-be/internal/bootstrap"
)

// @title SEA Catering API
// @version 1.0
// @description This is an API for SEA Catering application.
// @contact.name Jevon Mozart
// @contact.email jmcb1602@gmail.com
// @BasePath /
func main() {
	if err := bootstrap.Start(); err != nil {
		panic(fmt.Sprintf("Failed to start application: %v", err))
	}
}
