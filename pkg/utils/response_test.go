package utils

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func setupTestContext() (*gin.Context, *httptest.ResponseRecorder) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	return c, w
}

func TestSuccessResponse(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		message    string
		data       interface{}
	}{
		{
			name:       "Success with data",
			statusCode: http.StatusOK,
			message:    "Operation successful",
			data:       map[string]string{"key": "value"},
		},
		{
			name:       "Created with data",
			statusCode: http.StatusCreated,
			message:    "Resource created",
			data:       map[string]int{"id": 123},
		},
		{
			name:       "Success without data",
			statusCode: http.StatusOK,
			message:    "Success",
			data:       nil,
		},
		{
			name:       "Success with array data",
			statusCode: http.StatusOK,
			message:    "List retrieved",
			data:       []string{"item1", "item2", "item3"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, w := setupTestContext()

			SuccessResponse(c, tt.statusCode, tt.message, tt.data)

			if w.Code != tt.statusCode {
				t.Errorf("Expected status code %d, got %d", tt.statusCode, w.Code)
			}

			var response Response
			if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
				t.Fatalf("Failed to unmarshal response: %v", err)
			}

			if !response.Success {
				t.Error("Expected success to be true")
			}

			if response.Message != tt.message {
				t.Errorf("Expected message '%s', got '%s'", tt.message, response.Message)
			}

			if response.Error != "" {
				t.Errorf("Expected error to be empty, got '%s'", response.Error)
			}

			if tt.data != nil && response.Data == nil {
				t.Error("Expected data to be present")
			}
		})
	}
}

func TestErrorResponse(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		message    string
	}{
		{
			name:       "Bad request error",
			statusCode: http.StatusBadRequest,
			message:    "Invalid input",
		},
		{
			name:       "Unauthorized error",
			statusCode: http.StatusUnauthorized,
			message:    "Authentication required",
		},
		{
			name:       "Not found error",
			statusCode: http.StatusNotFound,
			message:    "Resource not found",
		},
		{
			name:       "Internal server error",
			statusCode: http.StatusInternalServerError,
			message:    "Internal server error",
		},
		{
			name:       "Forbidden error",
			statusCode: http.StatusForbidden,
			message:    "Access denied",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, w := setupTestContext()

			ErrorResponse(c, tt.statusCode, tt.message)

			if w.Code != tt.statusCode {
				t.Errorf("Expected status code %d, got %d", tt.statusCode, w.Code)
			}

			var response Response
			if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
				t.Fatalf("Failed to unmarshal response: %v", err)
			}

			if response.Success {
				t.Error("Expected success to be false")
			}

			if response.Error != tt.message {
				t.Errorf("Expected error '%s', got '%s'", tt.message, response.Error)
			}

			if response.Message != "" {
				t.Errorf("Expected message to be empty, got '%s'", response.Message)
			}

			if response.Data != nil {
				t.Errorf("Expected data to be nil, got %v", response.Data)
			}
		})
	}
}

func TestValidationErrorResponse(t *testing.T) {
	tests := []struct {
		name  string
		err   error
		want  string
	}{
		{
			name: "Simple error",
			err:  errors.New("validation failed"),
			want: "validation failed",
		},
		{
			name: "Empty error message",
			err:  errors.New(""),
			want: "",
		},
		{
			name: "Complex error message",
			err:  errors.New("field 'email' is required and must be a valid email"),
			want: "field 'email' is required and must be a valid email",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, w := setupTestContext()

			ValidationErrorResponse(c, tt.err)

			if w.Code != http.StatusBadRequest {
				t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, w.Code)
			}

			var response Response
			if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
				t.Fatalf("Failed to unmarshal response: %v", err)
			}

			if response.Success {
				t.Error("Expected success to be false")
			}

			if response.Error != tt.want {
				t.Errorf("Expected error '%s', got '%s'", tt.want, response.Error)
			}
		})
	}
}

func TestResponseStructure(t *testing.T) {
	t.Run("Success response structure", func(t *testing.T) {
		c, w := setupTestContext()
		testData := map[string]interface{}{
			"id":   1,
			"name": "test",
		}

		SuccessResponse(c, http.StatusOK, "Success", testData)

		var response map[string]interface{}
		if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
			t.Fatalf("Failed to unmarshal response: %v", err)
		}

		// Check required fields exist
		if _, exists := response["success"]; !exists {
			t.Error("Response missing 'success' field")
		}
		if _, exists := response["message"]; !exists {
			t.Error("Response missing 'message' field")
		}
		if _, exists := response["data"]; !exists {
			t.Error("Response missing 'data' field")
		}
	})

	t.Run("Error response structure", func(t *testing.T) {
		c, w := setupTestContext()

		ErrorResponse(c, http.StatusBadRequest, "Error message")

		var response map[string]interface{}
		if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
			t.Fatalf("Failed to unmarshal response: %v", err)
		}

		// Check required fields exist
		if _, exists := response["success"]; !exists {
			t.Error("Response missing 'success' field")
		}
		if _, exists := response["error"]; !exists {
			t.Error("Response missing 'error' field")
		}
	})
}

func TestContentType(t *testing.T) {
	tests := []struct {
		name string
		fn   func(*gin.Context)
	}{
		{
			name: "SuccessResponse",
			fn: func(c *gin.Context) {
				SuccessResponse(c, http.StatusOK, "test", nil)
			},
		},
		{
			name: "ErrorResponse",
			fn: func(c *gin.Context) {
				ErrorResponse(c, http.StatusBadRequest, "test")
			},
		},
		{
			name: "ValidationErrorResponse",
			fn: func(c *gin.Context) {
				ValidationErrorResponse(c, errors.New("test"))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, w := setupTestContext()
			tt.fn(c)

			contentType := w.Header().Get("Content-Type")
			if contentType != "application/json; charset=utf-8" {
				t.Errorf("Expected Content-Type 'application/json; charset=utf-8', got '%s'", contentType)
			}
		})
	}
}

func TestComplexDataTypes(t *testing.T) {
	type User struct {
		ID    int    `json:"id"`
		Name  string `json:"name"`
		Email string `json:"email"`
	}

	tests := []struct {
		name string
		data interface{}
	}{
		{
			name: "Struct data",
			data: User{ID: 1, Name: "John", Email: "john@example.com"},
		},
		{
			name: "Slice of structs",
			data: []User{
				{ID: 1, Name: "John", Email: "john@example.com"},
				{ID: 2, Name: "Jane", Email: "jane@example.com"},
			},
		},
		{
			name: "Nested map",
			data: map[string]interface{}{
				"user": map[string]interface{}{
					"id":   1,
					"name": "John",
				},
				"metadata": map[string]string{
					"created": "2023-01-01",
				},
			},
		},
		{
			name: "Mixed types",
			data: map[string]interface{}{
				"string":  "value",
				"number":  42,
				"boolean": true,
				"null":    nil,
				"array":   []int{1, 2, 3},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, w := setupTestContext()

			SuccessResponse(c, http.StatusOK, "Success", tt.data)

			if w.Code != http.StatusOK {
				t.Errorf("Expected status code 200, got %d", w.Code)
			}

			var response Response
			if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
				t.Fatalf("Failed to unmarshal response: %v", err)
			}

			if response.Data == nil {
				t.Error("Expected data to be present")
			}
		})
	}
}
