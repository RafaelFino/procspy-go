package handlers

import (
	"procspy/internal/procspy/config"
	"procspy/internal/procspy/service"
	"procspy/internal/procspy/storage"
	"testing"
)

// TestNewMatch testa criação de handler de match
func TestNewMatch(t *testing.T) {
	conn := storage.NewDbConnection(":memory:")
	defer conn.Close()

	matchService := service.NewMatch(conn)
	cfg := &config.Server{UserTarges: map[string]string{"user1": "url"}}
	usersService := service.NewUsers(cfg)

	handler := NewMatch(matchService, usersService)
	if handler == nil {
		t.Fatal("NewMatch retornou nil")
	}
}

// TestMatch_InsertMatch testa inserção de match via HTTP
func TestMatch_InsertMatch(t *testing.T) {
	conn := storage.NewDbConnection(":memory:")
	defer conn.Close()

	matchService := service.NewMatch(conn)
	cfg := &config.Server{UserTarges: map[string]string{"user1": "url"}}
	usersService := service.NewUsers(cfg)
	handler := NewMatch(matchService, usersService)

	gin := setupTestRouter()
	gin.POST("/match/:user", handler.InsertMatch)

	body := `{"user":"user1","name":"games","pattern":"steam","match":"steam.exe","elapsed":10.5}`
	req := makeTestRequest("POST", "/match/user1", body)
	w := executeRequest(gin, req)

	if w.Code != 201 {
		t.Errorf("Status = %d, esperado 201", w.Code)
	}
}

// TestMatch_InsertMatch_InvalidUser testa inserção com usuário inválido
func TestMatch_InsertMatch_InvalidUser(t *testing.T) {
	conn := storage.NewDbConnection(":memory:")
	defer conn.Close()

	matchService := service.NewMatch(conn)
	cfg := &config.Server{UserTarges: map[string]string{"user1": "url"}}
	usersService := service.NewUsers(cfg)
	handler := NewMatch(matchService, usersService)

	gin := setupTestRouter()
	gin.POST("/match/:user", handler.InsertMatch)

	body := `{"user":"invalid","name":"games","pattern":"steam","match":"steam.exe","elapsed":10.5}`
	req := makeTestRequest("POST", "/match/invalid", body)
	w := executeRequest(gin, req)

	if w.Code != 401 {
		t.Errorf("Status = %d, esperado 401", w.Code)
	}
}

// TestMatch_InsertMatch_InvalidJSON testa inserção com JSON inválido
func TestMatch_InsertMatch_InvalidJSON(t *testing.T) {
	conn := storage.NewDbConnection(":memory:")
	defer conn.Close()

	matchService := service.NewMatch(conn)
	cfg := &config.Server{UserTarges: map[string]string{"user1": "url"}}
	usersService := service.NewUsers(cfg)
	handler := NewMatch(matchService, usersService)

	gin := setupTestRouter()
	gin.POST("/match/:user", handler.InsertMatch)

	body := `{invalid json}`
	req := makeTestRequest("POST", "/match/user1", body)
	w := executeRequest(gin, req)

	if w.Code != 400 {
		t.Errorf("Status = %d, esperado 400", w.Code)
	}
}
