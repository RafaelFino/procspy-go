package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

// TestNewHealthcheck testa criação de handler de healthcheck
func TestNewHealthcheck(t *testing.T) {
	handler := NewHealthcheck()
	if handler == nil {
		t.Fatal("NewHealthcheck retornou nil")
	}

	if handler.startTime.IsZero() {
		t.Error("startTime não foi inicializado")
	}
}

// TestHealthcheck_GetStatus testa endpoint de healthcheck
func TestHealthcheck_GetStatus(t *testing.T) {
	// Arrange: Cria handler e router
	handler := NewHealthcheck()
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/healthcheck", handler.GetStatus)

	// Act: Faz request
	req := httptest.NewRequest("GET", "/healthcheck", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assert: Valida response
	if w.Code != http.StatusOK {
		t.Errorf("Status = %d, esperado %d", w.Code, http.StatusOK)
	}

	// Valida que response contém JSON
	if w.Header().Get("Content-Type") != "application/json; charset=utf-8" {
		t.Error("Content-Type deveria ser application/json")
	}
}
