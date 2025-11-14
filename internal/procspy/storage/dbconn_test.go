package storage

import (
	"testing"
)

// TestNewDbConnection testa a criação de uma nova conexão de banco
func TestNewDbConnection(t *testing.T) {
	conn := NewDbConnection(":memory:")
	if conn == nil {
		t.Fatal("NewDbConnection retornou nil")
	}
	if conn.path != ":memory:" {
		t.Errorf("path = %s, esperado ':memory:'", conn.path)
	}
}

// TestDbConnection_GetConn testa obtenção de conexão
func TestDbConnection_GetConn(t *testing.T) {
	conn := NewDbConnection(":memory:")
	db, err := conn.GetConn()
	if err != nil {
		t.Fatalf("GetConn() erro = %v", err)
	}
	if db == nil {
		t.Fatal("GetConn retornou nil")
	}
	defer conn.Close()
}

// TestDbConnection_Close testa fechamento de conexão
func TestDbConnection_Close(t *testing.T) {
	conn := NewDbConnection(":memory:")
	_, err := conn.GetConn()
	if err != nil {
		t.Fatalf("GetConn() erro = %v", err)
	}

	err = conn.Close()
	if err != nil {
		t.Errorf("Close() erro = %v", err)
	}

	// Testar fechar novamente (deveria ser seguro)
	err = conn.Close()
	if err != nil {
		t.Errorf("Close() segunda vez erro = %v", err)
	}
}

// TestDbConnection_Exec testa execução de query
func TestDbConnection_Exec(t *testing.T) {
	conn := NewDbConnection(":memory:")
	defer conn.Close()

	// Cria tabela de teste
	err := conn.Exec("CREATE TABLE test (id INTEGER PRIMARY KEY, name TEXT)")
	if err != nil {
		t.Fatalf("Exec() erro ao criar tabela = %v", err)
	}

	// Insere dados
	err = conn.Exec("INSERT INTO test (name) VALUES (?)", "test_name")
	if err != nil {
		t.Errorf("Exec() erro ao inserir = %v", err)
	}
}

// TestDbConnection_makeDBPath testa geração de caminho do banco
func TestDbConnection_makeDBPath(t *testing.T) {
	conn := NewDbConnection("/data")
	path := conn.makeDBPath()
	expected := "/data/procspy.db"
	if path != expected {
		t.Errorf("makeDBPath() = %s, esperado %s", path, expected)
	}
}

// TestDbConnection_Exec_InvalidQuery testa execução de query inválida
func TestDbConnection_Exec_InvalidQuery(t *testing.T) {
	conn := NewDbConnection(":memory:")
	defer conn.Close()

	// Tenta executar query inválida
	err := conn.Exec("INVALID SQL QUERY")
	if err == nil {
		t.Error("Exec() deveria retornar erro para query inválida")
	}
}

// TestDbConnection_Exec_WithoutConnection testa execução sem conexão prévia
func TestDbConnection_Exec_WithoutConnection(t *testing.T) {
	conn := NewDbConnection(":memory:")
	defer conn.Close()

	// Executa query sem chamar GetConn() primeiro
	err := conn.Exec("CREATE TABLE test2 (id INTEGER)")
	if err != nil {
		t.Errorf("Exec() deveria criar conexão automaticamente, erro = %v", err)
	}
}

// TestDbConnection_GetConn_Reuse testa reutilização de conexão existente
func TestDbConnection_GetConn_Reuse(t *testing.T) {
	conn := NewDbConnection(":memory:")
	defer conn.Close()

	// Primeira chamada
	db1, err := conn.GetConn()
	if err != nil {
		t.Fatalf("GetConn() primeira chamada erro = %v", err)
	}

	// Segunda chamada deve retornar a mesma conexão
	db2, err := conn.GetConn()
	if err != nil {
		t.Fatalf("GetConn() segunda chamada erro = %v", err)
	}

	if db1 != db2 {
		t.Error("GetConn() deveria reutilizar a mesma conexão")
	}
}

// TestDbConnection_Exec_MultipleOperations testa múltiplas operações
func TestDbConnection_Exec_MultipleOperations(t *testing.T) {
	conn := NewDbConnection(":memory:")
	defer conn.Close()

	// Cria tabela
	err := conn.Exec("CREATE TABLE test3 (id INTEGER PRIMARY KEY, value TEXT)")
	if err != nil {
		t.Fatalf("Exec() erro ao criar tabela = %v", err)
	}

	// Insere múltiplos registros
	for i := 0; i < 5; i++ {
		err = conn.Exec("INSERT INTO test3 (value) VALUES (?)", "value")
		if err != nil {
			t.Errorf("Exec() erro ao inserir registro %d = %v", i, err)
		}
	}

	// Update
	err = conn.Exec("UPDATE test3 SET value = ? WHERE id = ?", "updated", 1)
	if err != nil {
		t.Errorf("Exec() erro ao atualizar = %v", err)
	}

	// Delete
	err = conn.Exec("DELETE FROM test3 WHERE id = ?", 2)
	if err != nil {
		t.Errorf("Exec() erro ao deletar = %v", err)
	}
}
