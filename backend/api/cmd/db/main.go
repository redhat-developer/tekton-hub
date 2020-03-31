package main

import (
	"fmt"
	"log"
	"os"

	"github.com/jinzhu/gorm"
	"github.com/redhat-developer/tekton-hub/backend/api/pkg/app"
	"github.com/redhat-developer/tekton-hub/backend/api/pkg/models"

	// Blank
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

func main() {
	app, err := app.FromEnv("db")
	if err != nil {
		fmt.Fprintf(os.Stderr, "FATAL: failed to initialise: %s", err)
		os.Exit(1)
	}

	conn := app.Database().ConnectionString()

	db, err := gorm.Open("postgres", conn)

	defer db.Close()

	if err != nil {
		log.Fatalf("db connection failed: %s", err)
	}

	// Disable table name's pluralization globally
	db.SingularTable(true)

	// Create Tables
	err = models.CreateAndInitialiseTables(db)
	if err != nil {
		log.Println("Db initialisation failed", err)
	}

}
