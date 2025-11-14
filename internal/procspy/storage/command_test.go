package storage

import (
	"procspy/internal/procspy/domain"
	"testing"
)

// TestNewCommand testa criação de storage de command
func TestNewCommand(t *testing.T) {
	conn := NewDbConnection(":memory:")
	defer conn.Close()

	cmd := NewCommand(conn)
	if cmd == nil {
		t.Fatal("NewCommand retornou nil")
	}
}

// TestCommand_InsertCommand testa inserção de command
func TestCommand_InsertCommand(t *testing.T) {
	conn := NewDbConnection(":memory:")
	defer conn.Close()

	storage := NewCommand(conn)
	cmd := domain.NewCommand("user1", "games", "notify-send test", "success")

	err := storage.InsertCommand(cmd)
	if err != nil {
		t.Errorf("InsertCommand() erro = %v", err)
	}
}

// TestCommand_GetCommands testa busca de commands
func TestCommand_GetCommands(t *testing.T) {
	conn := NewDbConnection(":memory:")
	defer conn.Close()

	storage := NewCommand(conn)

	// Insere alguns commands
	cmd1 := domain.NewCommand("user1", "games", "cmd1", "out1")
	cmd2 := domain.NewCommand("user1", "browsers", "cmd2", "out2")

	storage.InsertCommand(cmd1)
	storage.InsertCommand(cmd2)

	// Busca commands
	commands, err := storage.GetCommands("user1")
	if err != nil {
		t.Fatalf("GetCommands() erro = %v", err)
	}

	if len(commands) == 0 {
		t.Error("GetCommands retornou vazio")
	}

	// Valida que pelo menos um command foi retornado
	if len(commands) < 2 {
		t.Errorf("Esperado pelo menos 2 commands, obteve %d", len(commands))
	}
}

// TestCommand_Close testa fechamento
func TestCommand_Close(t *testing.T) {
	conn := NewDbConnection(":memory:")
	storage := NewCommand(conn)

	err := storage.Close()
	if err != nil {
		t.Errorf("Close() erro = %v", err)
	}
}

// TestCommand_Close_AlreadyClosed testa fechar conexão já fechada
func TestCommand_Close_AlreadyClosed(t *testing.T) {
	storage := &Command{conn: nil}

	err := storage.Close()
	if err != nil {
		t.Errorf("Close() com conn nil deveria retornar nil, obteve erro = %v", err)
	}
}

// TestCommand_GetCommands_EmptyResult testa busca sem resultados
func TestCommand_GetCommands_EmptyResult(t *testing.T) {
	conn := NewDbConnection(":memory:")
	defer conn.Close()

	storage := NewCommand(conn)

	// Busca commands de usuário que não existe
	commands, err := storage.GetCommands("user_inexistente")
	if err != nil {
		t.Fatalf("GetCommands() erro = %v", err)
	}

	if len(commands) != 0 {
		t.Errorf("Esperado 0 commands, obteve %d", len(commands))
	}
}

// TestCommand_InsertCommand_MultipleUsers testa inserção de múltiplos usuários
func TestCommand_InsertCommand_MultipleUsers(t *testing.T) {
	conn := NewDbConnection(":memory:")
	defer conn.Close()

	storage := NewCommand(conn)

	// Insere commands de diferentes usuários
	cmd1 := domain.NewCommand("user1", "games", "cmd1", "out1")
	cmd2 := domain.NewCommand("user2", "browsers", "cmd2", "out2")

	storage.InsertCommand(cmd1)
	storage.InsertCommand(cmd2)

	// Busca commands de user1
	commands1, err := storage.GetCommands("user1")
	if err != nil {
		t.Fatalf("GetCommands(user1) erro = %v", err)
	}

	// Busca commands de user2
	commands2, err := storage.GetCommands("user2")
	if err != nil {
		t.Fatalf("GetCommands(user2) erro = %v", err)
	}

	// Valida que cada usuário tem seus próprios commands
	if len(commands1) != 1 {
		t.Errorf("user1 deveria ter 1 command, obteve %d", len(commands1))
	}

	if len(commands2) != 1 {
		t.Errorf("user2 deveria ter 1 command, obteve %d", len(commands2))
	}
}

// TestCommand_InsertCommand_WithAllFields testa inserção com todos os campos
func TestCommand_InsertCommand_WithAllFields(t *testing.T) {
	conn := NewDbConnection(":memory:")
	defer conn.Close()

	storage := NewCommand(conn)

	cmd := &domain.Command{
		User:        "user1",
		Name:        "test",
		CommandLine: "echo test",
		Return:      "0",
		Source:      "client",
		CommandLog:  "test output",
	}

	err := storage.InsertCommand(cmd)
	if err != nil {
		t.Errorf("InsertCommand() erro = %v", err)
	}

	// Busca e valida
	commands, err := storage.GetCommands("user1")
	if err != nil {
		t.Fatalf("GetCommands() erro = %v", err)
	}

	if len(commands) != 1 {
		t.Fatalf("Esperado 1 command, obteve %d", len(commands))
	}

	retrieved := commands[0]
	if retrieved.User != cmd.User {
		t.Errorf("User = %s, esperado %s", retrieved.User, cmd.User)
	}
	if retrieved.Name != cmd.Name {
		t.Errorf("Name = %s, esperado %s", retrieved.Name, cmd.Name)
	}
	if retrieved.CommandLine != cmd.CommandLine {
		t.Errorf("CommandLine = %s, esperado %s", retrieved.CommandLine, cmd.CommandLine)
	}
}
