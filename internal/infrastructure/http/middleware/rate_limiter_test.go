package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestRateLimiter_WithinLimit(t *testing.T) {
	gin.SetMode(gin.TestMode)

	rl := NewRateLimiter(10) // 10 requests per minute

	for i := 0; i < 10; i++ {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("GET", "/test", nil)
		c.Request.RemoteAddr = "192.168.1.1:1234"

		rl.RateLimit()(c)

		assert.False(t, c.IsAborted())
	}
}

func TestRateLimiter_ExceedsLimit(t *testing.T) {
	gin.SetMode(gin.TestMode)

	rl := NewRateLimiter(5) // 5 requests per minute

	// First 5 requests should succeed
	for i := 0; i < 5; i++ {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("GET", "/test", nil)
		c.Request.RemoteAddr = "192.168.1.1:1234"

		rl.RateLimit()(c)

		assert.False(t, c.IsAborted())
	}

	// 6th request should be rate limited
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/test", nil)
	c.Request.RemoteAddr = "192.168.1.1:1234"

	rl.RateLimit()(c)

	assert.True(t, c.IsAborted())
	assert.Equal(t, http.StatusTooManyRequests, w.Code)
	assert.Contains(t, w.Body.String(), "rate_limit_exceeded")
}

func TestRateLimiter_DifferentIPs(t *testing.T) {
	gin.SetMode(gin.TestMode)

	rl := NewRateLimiter(2) // 2 requests per minute

	// IP 1 - 2 requests (should both succeed)
	for i := 0; i < 2; i++ {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("GET", "/test", nil)
		c.Request.RemoteAddr = "192.168.1.1:1234"

		rl.RateLimit()(c)

		assert.False(t, c.IsAborted())
	}

	// IP 2 - should still be allowed (different IP)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/test", nil)
	c.Request.RemoteAddr = "192.168.1.2:1234"

	rl.RateLimit()(c)

	assert.False(t, c.IsAborted())
}

func TestRateLimiter_WindowReset(t *testing.T) {
	gin.SetMode(gin.TestMode)

	rl := &RateLimiter{
		visitors: make(map[string]*visitor),
		rate:     2,
		window:   100 * time.Millisecond, // Short window for testing
	}

	// First 2 requests
	for i := 0; i < 2; i++ {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("GET", "/test", nil)
		c.Request.RemoteAddr = "192.168.1.1:1234"

		rl.RateLimit()(c)

		assert.False(t, c.IsAborted())
	}

	// Wait for window to reset
	time.Sleep(150 * time.Millisecond)

	// Should be allowed again after window reset
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/test", nil)
	c.Request.RemoteAddr = "192.168.1.1:1234"

	rl.RateLimit()(c)

	assert.False(t, c.IsAborted())
}
