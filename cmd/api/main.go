package main

import (
	"log"

	personService "pessoas-api/internal/domain/person/service"
	"pessoas-api/internal/infrastructure/database"
	"pessoas-api/internal/infrastructure/http/handler"
	"pessoas-api/internal/infrastructure/http/router"
	personPersistence "pessoas-api/internal/infrastructure/persistence/person"

	_ "pessoas-api/docs" // Swagger docs
)

// @title           Pessoas API
// @version         1.0
// @description     REST API for people management following hexagonal architecture principles
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.email  support@pessoasapi.com

// @license.name  MIT
// @license.url   https://opensource.org/licenses/MIT

// @host      localhost:8080
// @BasePath  /api/v1

// @schemes http https

// @tag.name         Health
// @tag.description  Endpoints for API health check

// @tag.name         Persons
// @tag.description  CRUD operations for person management

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
