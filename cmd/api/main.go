package main

import (
	"log"

	operatorService "pessoas-api/internal/domain/operator/service"
	personService "pessoas-api/internal/domain/person/service"
	"pessoas-api/internal/infrastructure/database"
	"pessoas-api/internal/infrastructure/http/handler"
	"pessoas-api/internal/infrastructure/http/router"
	operatorPersistence "pessoas-api/internal/infrastructure/persistence/operator"
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

// @tag.name         Authentication
// @tag.description  Operator authentication and registration

// @tag.name         Persons
// @tag.description  CRUD operations for person management

func main() {
	config := database.LoadConfig()

	db, err := database.NewPostgresConnection(config)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Initialize repositories
	personRepo := personPersistence.NewPersonRepository(db)
	operatorRepo := operatorPersistence.NewOperatorRepository(db)

	// Initialize services
	personSvc := personService.NewPersonService(personRepo)
	authSvc := operatorService.NewAuthService(operatorRepo)

	// Initialize handlers
	personHandler := handler.NewPersonHandler(personSvc)
	authHandler := handler.NewAuthHandler(authSvc)

	// Setup router
	r := router.SetupRouter(personHandler, authHandler)

	log.Println("Starting server on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
