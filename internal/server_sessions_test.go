package internal

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"

	"github.com/serg666/card_storage/models"
	"github.com/serg666/card_storage/pkg"
)

func TestApiSessionCreate_Fail(t *testing.T) {
	r, _, _ := getSessionsRouter()

	req, _ := http.NewRequest("POST", "/v1/sessions", nil)

	testHTTPResponse(t, r, req, func(w *httptest.ResponseRecorder) bool {
		return w.Code == http.StatusBadRequest
	})
}

func TestApiSessionCreate_Complete(t *testing.T) {
	r, mRepo, kRepo := getSessionsRouter()
	mRepo.On("Create", mock.Anything).Return(&models.Session{}, nil)
	kRepo.On("Query", mock.Anything).Return([]*models.Key{}, nil)

	jsonBody, _ := json.Marshal(map[string]string{
		"key":  "session-unique-key",
		"body": `{"key1":"data1","key2":"data2"}`,
	})
	req, _ := http.NewRequest("POST", "/v1/sessions", bytes.NewBuffer(jsonBody))

	testHTTPResponse(t, r, req, func(w *httptest.ResponseRecorder) bool {
		return w.Code == http.StatusCreated
	})
}

func TestApiSessionGet_NotFound(t *testing.T) {
	r, mRepo, _ := getSessionsRouter()
	mRepo.On("Get", mock.Anything).Return(&models.Session{}, errors.New("Some error"))

	req, _ := http.NewRequest("GET", "/v1/sessions/unknown-key", nil)

	testHTTPResponse(t, r, req, func(w *httptest.ResponseRecorder) bool {
		return w.Code == http.StatusNotFound
	})
}

func TestApiSessionGet_Complete(t *testing.T) {
	r, mRepo, _ := getSessionsRouter()
	mRepo.On("Get", mock.Anything).Return(&models.Session{}, nil)

	req, _ := http.NewRequest("GET", "/v1/sessions/session-key", nil)

	testHTTPResponse(t, r, req, func(w *httptest.ResponseRecorder) bool {
		return w.Code == http.StatusOK
	})
}

func getSessionsRouter() (*gin.Engine, *mockSessionRepository, *mockKeyRepository) {
	mCardRepo := NewMockCardRepository()
	mSessionRepo := NewMockSessionRepository()
	mKeyRepo := NewMockKeyRepository()
	s := NewServer(pkg.NewConfigs(), mCardRepo, mSessionRepo, mKeyRepo)
	s.registerRoutes()
	return s.r, mSessionRepo, mKeyRepo
}
