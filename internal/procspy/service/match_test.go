package service

import (
	"procspy/internal/procspy/domain"
	"procspy/internal/procspy/storage"
	"testing"
)

// TestNewMatch testa criação de service de match
func TestNewMatch(t *testing.T) {
	conn := storage.NewDbConnection(":memory:")
	defer conn.Close()

	service := NewMatch(conn)
	if service == nil {
		t.Fatal("NewMatch retornou nil")
	}
}

// TestMatch_InsertMatch testa inserção com validação de elapsed
func TestMatch_InsertMatch(t *testing.T) {
	conn := storage.NewDbConnection(":memory:")
	defer conn.Close()

	service := NewMatch(conn)
	match := domain.NewMatch("user1", "games", "steam", "steam.exe", 10.0)

	err := service.InsertMatch(match)
	if err != nil {
		t.Errorf("InsertMatch() erro = %v", err)
	}
}

// TestMatch_InsertMatch_MaxElapsed testa validação de elapsed máximo
func TestMatch_InsertMatch_MaxElapsed(t *testing.T) {
	conn := storage.NewDbConnection(":memory:")
	defer conn.Close()

	service := NewMatch(conn)
	match := domain.NewMatch("user1", "games", "steam", "steam.exe", 200.0) // Acima do máximo

	err := service.InsertMatch(match)
	if err != nil {
		t.Errorf("InsertMatch() erro = %v", err)
	}

	// Valida que elapsed foi limitado
	if match.Elapsed != MATCH_MAX_ELAPSED {
		t.Errorf("Elapsed = %.2f, esperado %.2f (máximo)", match.Elapsed, MATCH_MAX_ELAPSED)
	}
}

// TestMatch_GetMatches testa busca de matches
func TestMatch_GetMatches(t *testing.T) {
	conn := storage.NewDbConnection(":memory:")
	defer conn.Close()

	service := NewMatch(conn)
	service.InsertMatch(domain.NewMatch("user1", "games", "steam", "steam.exe", 100.0))

	matches, err := service.GetMatches("user1")
	if err != nil {
		t.Errorf("GetMatches() erro = %v", err)
	}

	if matches == nil {
		t.Error("GetMatches retornou nil")
	}
}

// TestMatch_GetMatchesInfo testa busca de informações detalhadas
func TestMatch_GetMatchesInfo(t *testing.T) {
	conn := storage.NewDbConnection(":memory:")
	defer conn.Close()

	service := NewMatch(conn)
	service.InsertMatch(domain.NewMatch("user1", "games", "steam", "steam.exe", 100.0))

	info, err := service.GetMatchesInfo("user1")
	if err != nil {
		t.Errorf("GetMatchesInfo() erro = %v", err)
	}

	if info == nil {
		t.Error("GetMatchesInfo retornou nil")
	}
}

// TestMatch_Close testa fechamento do service
func TestMatch_Close(t *testing.T) {
	conn := storage.NewDbConnection(":memory:")

	service := NewMatch(conn)
	err := service.Close()
	if err != nil {
		t.Errorf("Close() erro = %v", err)
	}
}

// TestMatch_GetMatches_EmptyResult testa busca sem resultados
func TestMatch_GetMatches_EmptyResult(t *testing.T) {
	conn := storage.NewDbConnection(":memory:")
	defer conn.Close()

	service := NewMatch(conn)

	matches, err := service.GetMatches("user_inexistente")
	if err != nil {
		t.Errorf("GetMatches() erro = %v", err)
	}

	if len(matches) != 0 {
		t.Errorf("Esperado 0 matches, obteve %d", len(matches))
	}
}
