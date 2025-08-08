package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	
	api := router.Group("/api")
	{
		algorithms := api.Group("/algorithms")
		{
			algorithms.GET("", GetAlgorithms)
			algorithms.GET("/category/:category", GetAlgorithmsByCategory)
			algorithms.GET("/info/:id", GetAlgorithmInfo)
		}
	}
	
	return router
}

func TestGetAlgorithms(t *testing.T) {
	router := setupTestRouter()

	req, _ := http.NewRequest("GET", "/api/algorithms", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.True(t, response["success"].(bool))
	assert.NotNil(t, response["data"])
	assert.NotNil(t, response["count"])
}

func TestGetAlgorithmsByCategory(t *testing.T) {
	router := setupTestRouter()

	tests := []struct {
		category   string
		expectCode int
	}{
		{"sorting", http.StatusOK},
		{"searching", http.StatusOK},
		{"graph", http.StatusOK},
		{"invalid", http.StatusOK}, // 应该返回空数组，不是错误
	}

	for _, tt := range tests {
		t.Run(tt.category, func(t *testing.T) {
			req, _ := http.NewRequest("GET", "/api/algorithms/category/"+tt.category, nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectCode, w.Code)

			if w.Code == http.StatusOK {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.True(t, response["success"].(bool))
				assert.Equal(t, tt.category, response["category"])
			}
		})
	}
}

func TestGetAlgorithmInfo(t *testing.T) {
	router := setupTestRouter()

	tests := []struct {
		name       string
		algorithmID string
		expectCode int
	}{
		{"Valid algorithm", "bubble_sort", http.StatusOK},
		{"Invalid algorithm", "invalid_algorithm", http.StatusNotFound},
		{"Empty ID", "", http.StatusNotFound},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest("GET", "/api/algorithms/info/"+tt.algorithmID, nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectCode, w.Code)

			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)

			if w.Code == http.StatusOK {
				assert.True(t, response["success"].(bool))
				assert.NotNil(t, response["data"])
			} else {
				assert.False(t, response["success"].(bool))
				assert.NotNil(t, response["error"])
			}
		})
	}
}
