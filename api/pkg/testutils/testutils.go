package testutils

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"swipelearn-api/internal/models"
)

// TestLogger creates a logger for testing (discards output)
func TestLogger() *logrus.Logger {
	logger := logrus.New()
	logger.SetLevel(logrus.ErrorLevel) // Only show errors during tests
	return logger
}

// CreateTestUser creates a test user model
func CreateTestUser() *models.User {
	return &models.User{
		ID:           uuid.New(),
		Email:        "test@example.com",
		Name:         "Test User",
		PasswordHash: "$2a$10$example.hash",
	}
}

// CreateTestDeck creates a test deck model
func CreateTestDeck(userID uuid.UUID) *models.Deck {
	return &models.Deck{
		ID:          uuid.New(),
		UserID:      userID,
		Name:        "Test Deck",
		Description: "Test Description",
	}
}

// CreateTestFlashcard creates a test flashcard model
func CreateTestFlashcard(userID, deckID uuid.UUID) *models.Flashcard {
	return &models.Flashcard{
		ID:          uuid.New(),
		UserID:      userID,
		DeckID:      deckID,
		Front:       "Test Front",
		Back:        "Test Back",
		Difficulty:  2.5,
		Interval:    1,
		EaseFactor:  2.5,
		ReviewCount: 0,
	}
}

// GinTestRouter creates a Gin router for testing
func GinTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	return gin.New()
}

// MakeRequest creates an HTTP test request
func MakeRequest(t *testing.T, method, path string, body interface{}, headers map[string]string) *httptest.ResponseRecorder {
	router := GinTestRouter()

	var reqBody *bytes.Buffer
	if body != nil {
		jsonBody, err := json.Marshal(body)
		require.NoError(t, err)
		reqBody = bytes.NewBuffer(jsonBody)
	} else {
		reqBody = bytes.NewBuffer(nil)
	}

	req, err := http.NewRequest(method, path, reqBody)
	require.NoError(t, err)

	req.Header.Set("Content-Type", "application/json")
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	return w
}

// AssertJSONResponse asserts that the response is valid JSON and contains expected data
func AssertJSONResponse(t *testing.T, w *httptest.ResponseRecorder, expectedStatus int, expectedBody interface{}) {
	assert.Equal(t, expectedStatus, w.Code)

	if expectedBody != nil {
		var actualBody, expectedJSON interface{}

		err := json.Unmarshal(w.Body.Bytes(), &actualBody)
		require.NoError(t, err)

		switch v := expectedBody.(type) {
		case string:
			err = json.Unmarshal([]byte(v), &expectedJSON)
		default:
			expectedBytes, err := json.Marshal(v)
			require.NoError(t, err)
			err = json.Unmarshal(expectedBytes, &expectedJSON)
		}
		require.NoError(t, err)

		assert.Equal(t, expectedJSON, actualBody)
	}
}

// AssertErrorResponse asserts that the response contains an error with expected message
func AssertErrorResponse(t *testing.T, w *httptest.ResponseRecorder, expectedStatus int, expectedError string) {
	assert.Equal(t, expectedStatus, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Contains(t, response, "error")
	if expectedError != "" {
		assert.Equal(t, expectedError, response["error"])
	}
}

// RandomString generates a random string for testing
func RandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, length)
	for i := range result {
		result[i] = charset[uuid.New().ID()%uint32(len(charset))]
	}
	return string(result)
}

// RandomEmail generates a random email for testing
func RandomEmail() string {
	return RandomString(8) + "@example.com"
}

// CompareUUIDs safely compares UUID pointers and values
func CompareUUIDs(t *testing.T, expected, actual interface{}) {
	switch exp := expected.(type) {
	case uuid.UUID:
		switch act := actual.(type) {
		case uuid.UUID:
			assert.Equal(t, exp, act)
		case *uuid.UUID:
			require.NotNil(t, act)
			assert.Equal(t, exp, *act)
		default:
			t.Fatalf("unexpected actual type: %T", actual)
		}
	case *uuid.UUID:
		require.NotNil(t, exp)
		switch act := actual.(type) {
		case uuid.UUID:
			assert.Equal(t, *exp, act)
		case *uuid.UUID:
			require.NotNil(t, act)
			assert.Equal(t, *exp, *act)
		default:
			t.Fatalf("unexpected actual type: %T", actual)
		}
	default:
		t.Fatalf("unexpected expected type: %T", expected)
	}
}
