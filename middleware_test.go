package turispro_user

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupRouter() *gin.Engine {
	r := gin.New()
	r.Use(InjectUser())
	r.GET("/test", func(c *gin.Context) {
		user := GetUser(c)
		if user == nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "user not found"})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"id":            user.ID,
			"email":         user.Email,
			"level":         user.Level,
			"tour_operator": user.TourOperator,
		})
	})
	return r
}

func TestInjectUser_Success(t *testing.T) {
	r := setupRouter()
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("X-User-ID", "123")
	req.Header.Set("X-User-Email", "test@email.cl")
	req.Header.Set("X-User-TourOperator", "test_tour_operator")
	req.Header.Set("X-User-Level", "1")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), `"id":"123"`)
	assert.Contains(t, w.Body.String(), `"email":"test@email.cl"`)
	assert.Contains(t, w.Body.String(), `"level":1`)
	assert.Contains(t, w.Body.String(), `"tour_operator":"test_tour_operator"`)
}

func TestInjectUser_MissingHeaders(t *testing.T) {
	r := setupRouter()
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("X-User-Level", "1")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
	assert.Contains(t, w.Body.String(), `"error":"missing user information"`)
}

func TestInjectUser_InvalidLevel(t *testing.T) {
	r := setupRouter()

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("X-User-Id", "123")
	req.Header.Set("X-User-Email", "asd@asd.cl")
	req.Header.Set("X-User-Level", "100")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
	assert.Contains(t, w.Body.String(), "invalid user level")
}

func TestInjectUser_InvalidLevelFormat(t *testing.T) {
	r := setupRouter()

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("X-User-Id", "123")
	req.Header.Set("X-User-Name", "Juan")
	req.Header.Set("X-User-Level", "abc") // no es número

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
	assert.Contains(t, w.Body.String(), "invalid or missing level header")
}

func TestGetUser_NoUserInContext(t *testing.T) {
	c, _ := gin.CreateTestContext(httptest.NewRecorder())

	user := GetUser(c)
	assert.Nil(t, user)
}
