package router

import (
	"pessoas-api/internal/infrastructure/http/handler"
	"pessoas-api/internal/infrastructure/http/middleware"

	"github.com/gin-gonic/gin"
)

// SetupRouter configures and returns the HTTP router with all routes.
// It accepts a PersonHandler to handle person-related endpoints.
func SetupRouter(personHandler *handler.PersonHandler) *gin.Engine {
	router := gin.New()

	router.Use(gin.Recovery())

	router.Use(middleware.LoggerMiddleware())

	api := router.Group("/api")
	{
		v1 := api.Group("/v1")
		{
			persons := v1.Group("/persons")
			{
				persons.POST("", personHandler.CreatePerson)
				persons.GET("", personHandler.ListPersons)
				persons.GET("/cpf/:cpf", personHandler.FindPersonByCPF)
			}
		}
	}

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})

	return router
}
