package main

import (
	"log"

	"github.com/florantos/orbital-command/internal/config"
)

func main() {
	_, err := config.Load()

	if err != nil {
		log.Fatalf("Failed to load config: %p", err)
	}

}
