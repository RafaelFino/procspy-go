package handlers

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"procspy/internal/procspy/config"
	"procspy/internal/procspy/service"
	"testing"

	"github.com/gin-gonic/gin"
)

// TestValidateUser_ValidUser testa validação de usuário válido
func TestValidateUser_ValidUser(t *testing.T) {
	// Arrange: Cria service de users
	cfg := &config.Server{
		UserTarges: map[string]string{
			"user1": "http://example.com/user1.json",
		},
	}
	users := service.NewUsers(cfg)

	// Cria contexto gin com parâmetro user
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Params = gin.Params{{Key: "user", Value: "user1"}}

	// Act: Valida usuário
	user, err := ValidateUser(users, ctx)

	// Assert: Valida resultado
	if err != nil {
		t.Errorf("ValidateUser() erro = %v, esperado nil", err)
	}

	if user != "user1" {
		t.Errorf("user = %s, esperado 'user1'", user)
	}
}

// TestValidateUser_InvalidUser testa validação de usuário inválido
func TestValidateUser_InvalidUser(t *testing.T) {
	// Arrange: Cria service de users
	cfg := &config.Server{
		UserTarges: map[string]string{
			"user1": "http://example.com/user1.json",
		},
	}
	users := service.NewUsers(cfg)

	// Cria contexto gin com usuário inexistente
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Params = gin.Params{{Key: "user", Value: "nonexistent"}}

	// Act: Valida usuário
	_, err := ValidateUser(users, ctx)

	// Assert: Valida que retornou erro
	if err == nil {
		t.Error("ValidateUser() deveria retornar erro para usuário inexistente")
	}
}

// Helper functions for testing handlers

func setupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	return gin.New()
}

func makeTestRequest(method, url, body string) *http.Request {
	if body == "" {
		return httptest.NewRequest(method, url, nil)
	}
	return httptest.NewRequest(method, url, bytes.NewBufferString(body))
}

func executeRequest(router *gin.Engine, req *http.Request) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	return w
}
