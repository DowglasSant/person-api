package main

import (
	"log"

	personService "pessoas-api/internal/domain/person/service"
	"pessoas-api/internal/infrastructure/database"
	"pessoas-api/internal/infrastructure/http/handler"
	"pessoas-api/internal/infrastructure/http/router"
	personPersistence "pessoas-api/internal/infrastructure/persistence/person"
)

func main() {
	config := database.LoadConfig()

	db, err := database.NewPostgresConnection(config)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	personRepo := personPersistence.NewPersonRepository(db)

	personService := personService.NewPersonService(personRepo)

	personHandler := handler.NewPersonHandler(personService)

	r := router.SetupRouter(personHandler)

	log.Println("Starting server on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
