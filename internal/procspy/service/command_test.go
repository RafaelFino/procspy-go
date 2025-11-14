package service

import (
	"procspy/internal/procspy/domain"
	"procspy/internal/procspy/storage"
	"testing"
)

// TestNewCommand testa criação de service de command
func TestNewCommand(t *testing.T) {
	conn := storage.NewDbConnection(":memory:")
	defer conn.Close()

	service := NewCommand(conn)
	if service == nil {
		t.Fatal("NewCommand retornou nil")
	}
}

// TestCommand_InsertCommand testa inserção de command
func TestCommand_InsertCommand(t *testing.T) {
	conn := storage.NewDbConnection(":memory:")
	defer conn.Close()

	service := NewCommand(conn)
	cmd := domain.NewCommand("user1", "games", "notify-send test", "success")

	err := service.InsertCommand(cmd)
	if err != nil {
		t.Errorf("InsertCommand() erro = %v", err)
	}
}

// TestCommand_GetCommands testa busca de commands
func TestCommand_GetCommands(t *testing.T) {
	conn := storage.NewDbConnection(":memory:")
	defer conn.Close()

	service := NewCommand(conn)
	service.InsertCommand(domain.NewCommand("user1", "games", "cmd1", "out1"))

	commands, err := service.GetCommands("user1")
	if err != nil {
		t.Errorf("GetCommands() erro = %v", err)
	}

	if commands == nil {
		t.Error("GetCommands retornou nil")
	}
}

// TestCommand_Close testa fechamento do service
func TestCommand_Close(t *testing.T) {
	conn := storage.NewDbConnection(":memory:")

	service := NewCommand(conn)
	err := service.Close()
	if err != nil {
		t.Errorf("Close() erro = %v", err)
	}
}

// TestCommand_GetCommands_EmptyResult testa busca sem resultados
func TestCommand_GetCommands_EmptyResult(t *testing.T) {
	conn := storage.NewDbConnection(":memory:")
	defer conn.Close()

	service := NewCommand(conn)

	commands, err := service.GetCommands("user_inexistente")
	if err != nil {
		t.Errorf("GetCommands() erro = %v", err)
	}

	if len(commands) != 0 {
		t.Errorf("Esperado 0 commands, obteve %d", len(commands))
	}
}

// TestCommand_InsertCommand_MultipleCommands testa inserção de múltiplos commands
func TestCommand_InsertCommand_MultipleCommands(t *testing.T) {
	conn := storage.NewDbConnection(":memory:")
	defer conn.Close()

	service := NewCommand(conn)

	// Insere múltiplos commands
	for i := 0; i < 5; i++ {
		cmd := domain.NewCommand("user1", "test", "cmd", "output")
		err := service.InsertCommand(cmd)
		if err != nil {
			t.Errorf("InsertCommand() erro = %v", err)
		}
	}

	// Busca e valida
	commands, err := service.GetCommands("user1")
	if err != nil {
		t.Errorf("GetCommands() erro = %v", err)
	}

	if len(commands) != 5 {
		t.Errorf("Esperado 5 commands, obteve %d", len(commands))
	}
}
