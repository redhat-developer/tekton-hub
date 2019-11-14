package main

import (
	"log"

	"github.com/Pipelines-Marketplace/backend/pkg/models"
)

func main() {
	if err := models.StartConnection(); err != nil {
		log.Fatalln(err)
	}
}
