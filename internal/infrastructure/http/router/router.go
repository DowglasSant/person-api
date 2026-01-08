package router

import (
	"pessoas-api/internal/infrastructure/http/handler"
	"pessoas-api/internal/infrastructure/http/middleware"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetupRouter(personHandler *handler.PersonHandler, authHandler *handler.AuthHandler) *gin.Engine {
	router := gin.New()

	router.Use(gin.Recovery())

	router.Use(middleware.SecurityHeaders())

	router.Use(middleware.CORS())

	rateLimiter := middleware.NewRateLimiter(60)
	router.Use(rateLimiter.RateLimit())

	router.Use(middleware.LoggerMiddleware())

	// Public routes
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	api := router.Group("/api")
	{
		v1 := api.Group("/v1")
		{
			// Public authentication routes (no JWT required)
			auth := v1.Group("/auth")
			{
				auth.POST("/register", authHandler.Register)
				auth.POST("/login", authHandler.Login)
			}

			// Protected routes (JWT required)
			protected := v1.Group("")
			protected.Use(middleware.JWTAuth())
			{
				persons := protected.Group("/persons")
				persons.Use(middleware.ValidatePagination())
				{
					persons.POST("", personHandler.CreatePerson)
					persons.GET("", personHandler.ListPersons)
					persons.GET("/cpf/:cpf", personHandler.FindPersonByCPF)
				}
			}
		}
	}

	return router
}
