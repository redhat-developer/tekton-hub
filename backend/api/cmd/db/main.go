package main

import (
	"fmt"
	"os"

	"github.com/jinzhu/gorm"
	"github.com/redhat-developer/tekton-hub/backend/api/pkg/app"
	"github.com/redhat-developer/tekton-hub/backend/api/pkg/db/model"

	// Blank
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

func main() {
	app, err := app.BaseConfigFromEnv()
	if err != nil {
		fmt.Fprintf(os.Stderr, "FATAL: failed to initialise: %s", err)
		os.Exit(1)
	}

	log := app.Logger()
	defer log.Sync()

	conn := app.Database().ConnectionString()
	db, err := gorm.Open("postgres", conn)
	if err != nil {
		log.Fatalf("DB connection failed: %s", err)
	}

	// Create Tables
	if err = model.Migrate(db, log); err != nil {
		log.Fatal("DB initialisation failed", err)
	}

}
