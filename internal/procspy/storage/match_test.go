package storage

import (
	"procspy/internal/procspy/domain"
	"testing"
)

// TestNewMatch testa criação de storage de match
func TestNewMatch(t *testing.T) {
	conn := NewDbConnection(":memory:")
	defer conn.Close()

	match := NewMatch(conn)
	if match == nil {
		t.Fatal("NewMatch retornou nil")
	}
}

// TestMatch_InsertMatch testa inserção de match
func TestMatch_InsertMatch(t *testing.T) {
	conn := NewDbConnection(":memory:")
	defer conn.Close()

	storage := NewMatch(conn)
	match := domain.NewMatch("user1", "games", "steam", "steam.exe", 10.5)

	err := storage.InsertMatch(match)
	if err != nil {
		t.Errorf("InsertMatch() erro = %v", err)
	}
}

// TestMatch_GetMatches testa busca de matches
func TestMatch_GetMatches(t *testing.T) {
	conn := NewDbConnection(":memory:")
	defer conn.Close()

	storage := NewMatch(conn)

	// Insere alguns matches
	storage.InsertMatch(domain.NewMatch("user1", "games", "steam", "steam.exe", 100.0))
	storage.InsertMatch(domain.NewMatch("user1", "games", "steam", "steam.exe", 50.0))
	storage.InsertMatch(domain.NewMatch("user1", "browsers", "chrome", "chrome.exe", 200.0))

	// Busca matches
	matches, err := storage.GetMatches("user1")
	if err != nil {
		t.Fatalf("GetMatches() erro = %v", err)
	}

	if len(matches) != 2 {
		t.Errorf("Esperado 2 matches (games e browsers), obteve %d", len(matches))
	}

	// Valida que games tem soma correta (150 = 100 + 50)
	if games, ok := matches["games"]; ok {
		if games < 149.0 || games > 151.0 {
			t.Errorf("games elapsed = %.2f, esperado ~150.00", games)
		}
	} else {
		t.Error("games não encontrado nos matches")
	}
}

// TestMatch_GetMatchesInfo testa busca de informações detalhadas
func TestMatch_GetMatchesInfo(t *testing.T) {
	conn := NewDbConnection(":memory:")
	defer conn.Close()

	storage := NewMatch(conn)

	// Insere matches
	storage.InsertMatch(domain.NewMatch("user1", "games", "steam", "steam.exe", 100.0))
	storage.InsertMatch(domain.NewMatch("user1", "games", "steam", "steam.exe", 50.0))

	// Busca info
	info, err := storage.GetMatchesInfo("user1")
	if err != nil {
		t.Fatalf("GetMatchesInfo() erro = %v", err)
	}

	if len(info) == 0 {
		t.Error("GetMatchesInfo retornou vazio")
	}

	// Valida informações
	if gamesInfo, ok := info["games"]; ok {
		if gamesInfo.Elapsed != 150.0 {
			t.Errorf("elapsed = %.2f, esperado 150.00", gamesInfo.Elapsed)
		}
		if gamesInfo.Ocurrences != 2 {
			t.Errorf("ocurrences = %d, esperado 2", gamesInfo.Ocurrences)
		}
	}
}

// TestMatch_Close testa fechamento
func TestMatch_Close(t *testing.T) {
	conn := NewDbConnection(":memory:")
	storage := NewMatch(conn)

	err := storage.Close()
	if err != nil {
		t.Errorf("Close() erro = %v", err)
	}
}

// TestMatch_InsertMatch_NilConnection testa inserção com conexão nil
func TestMatch_InsertMatch_NilConnection(t *testing.T) {
	storage := &Match{conn: nil}
	match := domain.NewMatch("user1", "games", "steam", "steam.exe", 10.5)

	err := storage.InsertMatch(match)
	if err == nil {
		t.Error("InsertMatch() deveria retornar erro com conexão nil")
	}
}

// TestMatch_GetMatches_EmptyResult testa busca sem resultados
func TestMatch_GetMatches_EmptyResult(t *testing.T) {
	conn := NewDbConnection(":memory:")
	defer conn.Close()

	storage := NewMatch(conn)

	// Busca matches de usuário que não existe
	matches, err := storage.GetMatches("user_inexistente")
	if err != nil {
		t.Fatalf("GetMatches() erro = %v", err)
	}

	if len(matches) != 0 {
		t.Errorf("Esperado 0 matches, obteve %d", len(matches))
	}
}

// TestMatch_GetMatchesInfo_EmptyResult testa busca de info sem resultados
func TestMatch_GetMatchesInfo_EmptyResult(t *testing.T) {
	conn := NewDbConnection(":memory:")
	defer conn.Close()

	storage := NewMatch(conn)

	// Busca info de usuário que não existe
	info, err := storage.GetMatchesInfo("user_inexistente")
	if err != nil {
		t.Fatalf("GetMatchesInfo() erro = %v", err)
	}

	if len(info) != 0 {
		t.Errorf("Esperado 0 infos, obteve %d", len(info))
	}
}

// TestMatch_Close_AlreadyClosed testa fechar conexão já fechada
func TestMatch_Close_AlreadyClosed(t *testing.T) {
	storage := &Match{conn: nil}

	err := storage.Close()
	if err != nil {
		t.Errorf("Close() com conn nil deveria retornar nil, obteve erro = %v", err)
	}
}

// TestMatch_InsertMatch_MultipleUsers testa inserção de múltiplos usuários
func TestMatch_InsertMatch_MultipleUsers(t *testing.T) {
	conn := NewDbConnection(":memory:")
	defer conn.Close()

	storage := NewMatch(conn)

	// Insere matches de diferentes usuários
	storage.InsertMatch(domain.NewMatch("user1", "games", "steam", "steam.exe", 100.0))
	storage.InsertMatch(domain.NewMatch("user2", "games", "steam", "steam.exe", 200.0))

	// Busca matches de user1
	matches1, err := storage.GetMatches("user1")
	if err != nil {
		t.Fatalf("GetMatches(user1) erro = %v", err)
	}

	// Busca matches de user2
	matches2, err := storage.GetMatches("user2")
	if err != nil {
		t.Fatalf("GetMatches(user2) erro = %v", err)
	}

	// Valida que cada usuário tem seus próprios matches
	if games1, ok := matches1["games"]; ok {
		if games1 != 100.0 {
			t.Errorf("user1 games elapsed = %.2f, esperado 100.00", games1)
		}
	}

	if games2, ok := matches2["games"]; ok {
		if games2 != 200.0 {
			t.Errorf("user2 games elapsed = %.2f, esperado 200.00", games2)
		}
	}
}
