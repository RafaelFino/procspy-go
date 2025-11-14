package handlers

import (
	"procspy/internal/procspy/config"
	"procspy/internal/procspy/service"
	"procspy/internal/procspy/storage"
	"testing"
)

// TestNewCommand testa criação de handler de command
func TestNewCommand(t *testing.T) {
	conn := storage.NewDbConnection(":memory:")
	defer conn.Close()

	commandService := service.NewCommand(conn)
	cfg := &config.Server{UserTarges: map[string]string{"user1": "url"}}
	usersService := service.NewUsers(cfg)

	handler := NewCommand(commandService, usersService)
	if handler == nil {
		t.Fatal("NewCommand retornou nil")
	}
}

// TestCommand_InsertCommand testa inserção de command via HTTP
func TestCommand_InsertCommand(t *testing.T) {
	conn := storage.NewDbConnection(":memory:")
	defer conn.Close()

	commandService := service.NewCommand(conn)
	cfg := &config.Server{UserTarges: map[string]string{"user1": "url"}}
	usersService := service.NewUsers(cfg)
	handler := NewCommand(commandService, usersService)

	gin := setupTestRouter()
	gin.POST("/command/:user", handler.InsertCommand)

	body := `{"user":"user1","name":"games","command_line":"notify-send test","return":"0","source":"client"}`
	req := makeTestRequest("POST", "/command/user1", body)
	w := executeRequest(gin, req)

	if w.Code != 201 {
		t.Errorf("Status = %d, esperado 201", w.Code)
	}
}

// TestCommand_InsertCommand_InvalidUser testa inserção com usuário inválido
func TestCommand_InsertCommand_InvalidUser(t *testing.T) {
	conn := storage.NewDbConnection(":memory:")
	defer conn.Close()

	commandService := service.NewCommand(conn)
	cfg := &config.Server{UserTarges: map[string]string{"user1": "url"}}
	usersService := service.NewUsers(cfg)
	handler := NewCommand(commandService, usersService)

	gin := setupTestRouter()
	gin.POST("/command/:user", handler.InsertCommand)

	body := `{"user":"invalid","name":"games","command_line":"test","return":"0"}`
	req := makeTestRequest("POST", "/command/invalid", body)
	w := executeRequest(gin, req)

	if w.Code != 401 {
		t.Errorf("Status = %d, esperado 401", w.Code)
	}
}

// TestCommand_InsertCommand_InvalidJSON testa inserção com JSON inválido
func TestCommand_InsertCommand_InvalidJSON(t *testing.T) {
	conn := storage.NewDbConnection(":memory:")
	defer conn.Close()

	commandService := service.NewCommand(conn)
	cfg := &config.Server{UserTarges: map[string]string{"user1": "url"}}
	usersService := service.NewUsers(cfg)
	handler := NewCommand(commandService, usersService)

	gin := setupTestRouter()
	gin.POST("/command/:user", handler.InsertCommand)

	body := `{invalid json}`
	req := makeTestRequest("POST", "/command/user1", body)
	w := executeRequest(gin, req)

	if w.Code != 400 {
		t.Errorf("Status = %d, esperado 400", w.Code)
	}
}
