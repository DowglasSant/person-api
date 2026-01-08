package middleware

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func ValidatePagination() gin.HandlerFunc {
	return func(c *gin.Context) {
		if pageStr := c.Query("page"); pageStr != "" {
			page, err := strconv.Atoi(pageStr)
			if err != nil || page < 1 {
				c.JSON(http.StatusBadRequest, gin.H{
					"error":   "invalid_parameter",
					"message": "Page must be a positive integer",
				})
				c.Abort()
				return
			}
			if page > 10000 {
				c.JSON(http.StatusBadRequest, gin.H{
					"error":   "invalid_parameter",
					"message": "Page number too large (max: 10000)",
				})
				c.Abort()
				return
			}
		}

		if pageSizeStr := c.Query("page_size"); pageSizeStr != "" {
			pageSize, err := strconv.Atoi(pageSizeStr)
			if err != nil || pageSize < 1 {
				c.JSON(http.StatusBadRequest, gin.H{
					"error":   "invalid_parameter",
					"message": "Page size must be a positive integer",
				})
				c.Abort()
				return
			}
			if pageSize > 100 {
				c.JSON(http.StatusBadRequest, gin.H{
					"error":   "invalid_parameter",
					"message": "Page size too large (max: 100)",
				})
				c.Abort()
				return
			}
		}

		if sort := c.Query("sort"); sort != "" {
			allowedFields := map[string]bool{
				"id":         true,
				"name":       true,
				"cpf":        true,
				"email":      true,
				"created_at": true,
				"updated_at": true,
			}
			if !allowedFields[sort] {
				c.JSON(http.StatusBadRequest, gin.H{
					"error":   "invalid_parameter",
					"message": "Invalid sort field. Allowed: id, name, cpf, email, created_at, updated_at",
				})
				c.Abort()
				return
			}
		}

		if order := c.Query("order"); order != "" {
			if order != "asc" && order != "desc" {
				c.JSON(http.StatusBadRequest, gin.H{
					"error":   "invalid_parameter",
					"message": "Order must be 'asc' or 'desc'",
				})
				c.Abort()
				return
			}
		}

		c.Next()
	}
}

// SecurityHeaders adds security-related HTTP headers
func SecurityHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("X-Content-Type-Options", "nosniff")
		c.Writer.Header().Set("X-Frame-Options", "DENY")
		c.Writer.Header().Set("X-XSS-Protection", "1; mode=block")
		c.Writer.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		c.Writer.Header().Set("Content-Security-Policy", "default-src 'self'")
		c.Next()
	}
}
