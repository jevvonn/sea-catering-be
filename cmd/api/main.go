package main

import (
	"fmt"

	"github.com/jevvonn/sea-catering-be/internal/bootstrap"
)

func main() {
	if err := bootstrap.Start(); err != nil {
		panic(fmt.Sprintf("Failed to start application: %v", err))
	}
}
