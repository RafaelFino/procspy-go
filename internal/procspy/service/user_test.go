package service

import (
	"procspy/internal/procspy/config"
	"testing"
)

// TestNewUsers testa criação de service de users
func TestNewUsers(t *testing.T) {
	cfg := &config.Server{
		UserTarges: map[string]string{
			"user1": "http://example.com/user1.json",
		},
	}

	service := NewUsers(cfg)
	if service == nil {
		t.Fatal("NewUsers retornou nil")
	}
}

// TestUsers_GetUsers testa busca de usuários
func TestUsers_GetUsers(t *testing.T) {
	cfg := &config.Server{
		UserTarges: map[string]string{
			"user1": "http://example.com/user1.json",
			"user2": "http://example.com/user2.json",
		},
	}

	service := NewUsers(cfg)
	users, err := service.GetUsers()

	if err != nil {
		t.Errorf("GetUsers() erro = %v", err)
	}

	if len(users) != 2 {
		t.Errorf("Esperado 2 usuários, obteve %d", len(users))
	}
}

// TestUsers_Exists testa verificação de existência de usuário
func TestUsers_Exists(t *testing.T) {
	cfg := &config.Server{
		UserTarges: map[string]string{
			"user1": "http://example.com/user1.json",
		},
	}

	service := NewUsers(cfg)

	if !service.Exists("user1") {
		t.Error("user1 deveria existir")
	}

	if service.Exists("user2") {
		t.Error("user2 não deveria existir")
	}
}

// TestUsers_GetUsers_Empty testa busca de usuários com config vazia
func TestUsers_GetUsers_Empty(t *testing.T) {
	cfg := &config.Server{
		UserTarges: map[string]string{},
	}

	service := NewUsers(cfg)
	users, err := service.GetUsers()

	if err != nil {
		t.Errorf("GetUsers() erro = %v", err)
	}

	if len(users) != 0 {
		t.Errorf("Esperado 0 usuários, obteve %d", len(users))
	}
}

// TestUsers_Exists_EmptyConfig testa verificação com config vazia
func TestUsers_Exists_EmptyConfig(t *testing.T) {
	cfg := &config.Server{
		UserTarges: map[string]string{},
	}

	service := NewUsers(cfg)

	if service.Exists("anyuser") {
		t.Error("Nenhum usuário deveria existir em config vazia")
	}
}
