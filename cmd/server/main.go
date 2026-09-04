package main

import (
	"log"

	"github.com/techcollectiveassam/flood-response-api/internal/api"
	"github.com/techcollectiveassam/flood-response-api/internal/app"
)

func main() {
	application, err := app.New()
	if err != nil {
		log.Fatal(err)
	}
	defer application.Close()

	features := app.NewFeatures(application)

	router := api.NewRouter(features)

	log.Printf("server listening on :%s", application.Config.Port)

	if err := router.Run(":" + application.Config.Port); err != nil {
		log.Fatal(err)
	}
}
