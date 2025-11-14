package service

import (
	"procspy/internal/procspy/config"
	"testing"
)

// TestNewTarget testa criação de service de target
func TestNewTarget(t *testing.T) {
	cfg := &config.Server{
		UserTarges: map[string]string{
			"user1": "http://example.com/user1.json",
		},
	}

	service := NewTarget(cfg)
	if service == nil {
		t.Fatal("NewTarget retornou nil")
	}

	if service.urls == nil {
		t.Error("urls não foi inicializado")
	}
}

// TestTarget_GetTargets_NoUser testa busca de targets para usuário inexistente
func TestTarget_GetTargets_NoUser(t *testing.T) {
	cfg := &config.Server{
		UserTarges: map[string]string{},
	}

	service := NewTarget(cfg)
	targets, err := service.GetTargets("nonexistent")

	// Quando usuário não existe, retorna lista vazia sem erro
	if err != nil {
		t.Errorf("Não deveria retornar erro para usuário inexistente, got: %v", err)
	}

	if targets == nil {
		t.Error("Targets não deveria ser nil")
	}

	if len(targets.Targets) != 0 {
		t.Errorf("Targets deveria estar vazio, got %d targets", len(targets.Targets))
	}
}

// TestTarget_GetTargets_InvalidURL testa busca com URL inválida
func TestTarget_GetTargets_InvalidURL(t *testing.T) {
	cfg := &config.Server{
		UserTarges: map[string]string{
			"user1": "http://invalid-url-that-does-not-exist-12345.com/targets.json",
		},
	}

	service := NewTarget(cfg)
	_, err := service.GetTargets("user1")

	if err == nil {
		t.Error("GetTargets() deveria retornar erro para URL inválida")
	}
}

// TestTarget_getFromUrl_InvalidURL testa getFromUrl com URL inválida
func TestTarget_getFromUrl_InvalidURL(t *testing.T) {
	cfg := &config.Server{
		UserTarges: map[string]string{},
	}

	service := NewTarget(cfg)
	_, err := service.getFromUrl("http://invalid-url-that-does-not-exist-12345.com")

	if err == nil {
		t.Error("getFromUrl() deveria retornar erro para URL inválida")
	}
}
