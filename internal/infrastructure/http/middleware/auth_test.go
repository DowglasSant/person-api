package middleware

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestJWTAuth_MissingAuthorizationHeader(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/test", nil)

	JWTAuth()(c)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "Authorization header is required")
}

func TestJWTAuth_InvalidHeaderFormat(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/test", nil)
	c.Request.Header.Set("Authorization", "InvalidFormat")

	JWTAuth()(c)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "Invalid authorization header format")
}

func TestJWTAuth_InvalidToken(t *testing.T) {
	gin.SetMode(gin.TestMode)
	os.Setenv("JWT_SECRET", "test-secret-key-minimum-32-characters-long")
	defer os.Unsetenv("JWT_SECRET")

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/test", nil)
	c.Request.Header.Set("Authorization", "Bearer invalid.token.here")

	JWTAuth()(c)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "Invalid or expired token")
}

func TestJWTAuth_ValidToken(t *testing.T) {
	gin.SetMode(gin.TestMode)
	os.Setenv("JWT_SECRET", "test-secret-key-minimum-32-characters-long")
	defer os.Unsetenv("JWT_SECRET")

	token, err := GenerateToken(123, "testuser")
	assert.NoError(t, err)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/test", nil)
	c.Request.Header.Set("Authorization", "Bearer "+token)

	var called bool
	handler := func(c *gin.Context) {
		called = true
		userID, exists := c.Get("user_id")
		assert.True(t, exists)
		assert.Equal(t, 123, userID)

		username, exists := c.Get("username")
		assert.True(t, exists)
		assert.Equal(t, "testuser", username)
	}

	JWTAuth()(c)
	if !c.IsAborted() {
		handler(c)
	}

	assert.True(t, called)
}

func TestGenerateToken(t *testing.T) {
	os.Setenv("JWT_SECRET", "test-secret-key-minimum-32-characters-long")
	defer os.Unsetenv("JWT_SECRET")

	token, err := GenerateToken(456, "johndoe")
	assert.NoError(t, err)
	assert.NotEmpty(t, token)
}

func TestGetJWTSecret_NotSet(t *testing.T) {
	os.Unsetenv("JWT_SECRET")

	assert.Panics(t, func() {
		getJWTSecret()
	})
}

func TestGetJWTSecret_TooShort(t *testing.T) {
	os.Setenv("JWT_SECRET", "short")
	defer os.Unsetenv("JWT_SECRET")

	assert.Panics(t, func() {
		getJWTSecret()
	})
}
