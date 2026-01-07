package router

import (
	"pessoas-api/internal/infrastructure/http/handler"
	"pessoas-api/internal/infrastructure/http/middleware"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetupRouter(personHandler *handler.PersonHandler) *gin.Engine {
	router := gin.New()

	router.Use(gin.Recovery())

	router.Use(middleware.LoggerMiddleware())

	// Swagger documentation
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

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
