package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestValidatePagination_ValidParameters(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/test?page=1&page_size=10&sort=name&order=asc", nil)

	ValidatePagination()(c)

	assert.False(t, c.IsAborted())
}

func TestValidatePagination_InvalidPage(t *testing.T) {
	gin.SetMode(gin.TestMode)

	testCases := []struct {
		name  string
		query string
	}{
		{"negative page", "/test?page=-1"},
		{"zero page", "/test?page=0"},
		{"non-numeric page", "/test?page=abc"},
		{"page too large", "/test?page=10001"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request, _ = http.NewRequest("GET", tc.query, nil)

			ValidatePagination()(c)

			assert.True(t, c.IsAborted())
			assert.Equal(t, http.StatusBadRequest, w.Code)
		})
	}
}

func TestValidatePagination_InvalidPageSize(t *testing.T) {
	gin.SetMode(gin.TestMode)

	testCases := []struct {
		name  string
		query string
	}{
		{"negative page_size", "/test?page_size=-1"},
		{"zero page_size", "/test?page_size=0"},
		{"non-numeric page_size", "/test?page_size=abc"},
		{"page_size too large", "/test?page_size=101"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request, _ = http.NewRequest("GET", tc.query, nil)

			ValidatePagination()(c)

			assert.True(t, c.IsAborted())
			assert.Equal(t, http.StatusBadRequest, w.Code)
		})
	}
}

func TestValidatePagination_InvalidSortField(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/test?sort=invalid_field", nil)

	ValidatePagination()(c)

	assert.True(t, c.IsAborted())
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Invalid sort field")
}

func TestValidatePagination_InvalidOrder(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/test?order=invalid", nil)

	ValidatePagination()(c)

	assert.True(t, c.IsAborted())
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Order must be 'asc' or 'desc'")
}

func TestSecurityHeaders(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/test", nil)

	SecurityHeaders()(c)

	assert.Equal(t, "nosniff", w.Header().Get("X-Content-Type-Options"))
	assert.Equal(t, "DENY", w.Header().Get("X-Frame-Options"))
	assert.Equal(t, "1; mode=block", w.Header().Get("X-XSS-Protection"))
	assert.Equal(t, "max-age=31536000; includeSubDomains", w.Header().Get("Strict-Transport-Security"))
	assert.Equal(t, "default-src 'self'", w.Header().Get("Content-Security-Policy"))
}
