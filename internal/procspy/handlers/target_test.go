package handlers

import (
	"procspy/internal/procspy/config"
	"procspy/internal/procspy/service"
	"procspy/internal/procspy/storage"
	"testing"
)

// TestNewTarget testa criação de handler de target
func TestNewTarget(t *testing.T) {
	cfg := &config.Server{UserTarges: map[string]string{"user1": "url"}}
	targetService := service.NewTarget(cfg)
	usersService := service.NewUsers(cfg)
	
	conn := storage.NewDbConnection(":memory:")
	defer conn.Close()
	matchService := service.NewMatch(conn)

	handler := NewTarget(targetService, usersService, matchService)
	if handler == nil {
		t.Fatal("NewTarget retornou nil")
	}
}

// TestTarget_GetTargets_InvalidUser testa busca com usuário inválido
func TestTarget_GetTargets_InvalidUser(t *testing.T) {
	cfg := &config.Server{UserTarges: map[string]string{"user1": "url"}}
	targetService := service.NewTarget(cfg)
	usersService := service.NewUsers(cfg)
	
	conn := storage.NewDbConnection(":memory:")
	defer conn.Close()
	matchService := service.NewMatch(conn)
	handler := NewTarget(targetService, usersService, matchService)

	gin := setupTestRouter()
	gin.GET("/targets/:user", handler.GetTargets)

	req := makeTestRequest("GET", "/targets/invalid", "")
	w := executeRequest(gin, req)

	if w.Code != 401 {
		t.Errorf("Status = %d, esperado 401", w.Code)
	}
}
