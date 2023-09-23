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

func TestApiCardList_NotFound(t *testing.T) {
	r, mRepo, _ := getCardsRouter()
	mRepo.On("ListByPan", mock.Anything, mock.Anything).Return([]*models.Card{}, errors.New("Any error"))

	req, _ := http.NewRequest("GET", "/v1/cards?pan=1111222211112222", nil)

	testHTTPResponse(t, r, req, func(w *httptest.ResponseRecorder) bool {
		return w.Code == http.StatusNotFound
	})
}

func TestApiCardList_Complete(t *testing.T) {
	r, mRepo, _ := getCardsRouter()
	mRepo.On("ListByPan", mock.Anything, mock.Anything).Return([]*models.Card{}, nil)

	req, _ := http.NewRequest("GET", "/v1/cards?pan=1111222211112222", nil)

	testHTTPResponse(t, r, req, func(w *httptest.ResponseRecorder) bool {
		return w.Code == http.StatusOK
	})
}

func TestApiCardCreate_Fail(t *testing.T) {
	r, _, _ := getCardsRouter()

	req, _ := http.NewRequest("POST", "/v1/cards", nil)

	testHTTPResponse(t, r, req, func(w *httptest.ResponseRecorder) bool {
		return w.Code == http.StatusBadRequest
	})
}

func TestApiCardCreate_Complete(t *testing.T) {
	r, mRepo, kRepo := getCardsRouter()
	mRepo.On("Create", mock.Anything).Return(&models.Card{}, nil)
	kRepo.On("Query", mock.Anything).Return([]*models.Key{}, nil)

	jsonBody, _ := json.Marshal(map[string]string{
		"pan":      "1111 2222 3333 4444",
		"holder":   "HOLDER NAME",
		"exp_date": "26/01",
	})
	req, _ := http.NewRequest("POST", "/v1/cards", bytes.NewBuffer(jsonBody))

	testHTTPResponse(t, r, req, func(w *httptest.ResponseRecorder) bool {
		return w.Code == http.StatusCreated
	})
}

func TestApiCardDelete_Complete(t *testing.T) {
	r, mRepo, _ := getCardsRouter()
	mRepo.On("Delete", mock.Anything).Return(nil)

	req, _ := http.NewRequest("DELETE", "/v1/cards/1", nil)

	testHTTPResponse(t, r, req, func(w *httptest.ResponseRecorder) bool {
		return w.Code == http.StatusOK
	})
}

func testHTTPResponse(t *testing.T, r *gin.Engine, req *http.Request, f func(w *httptest.ResponseRecorder) bool) {
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if !f(w) {
		t.Fail()
	}
}

func getCardsRouter() (*gin.Engine, *mockCardRepository, *mockKeyRepository) {
	mCardRepo := NewMockCardRepository()
	mSessionRepo := NewMockSessionRepository()
	mKeyRepo := NewMockKeyRepository()
	s := NewServer(pkg.NewConfigs(), mCardRepo, mSessionRepo, mKeyRepo)
	s.registerRoutes()
	return s.r, mCardRepo, mKeyRepo
}
