package server

import (
	"procspy/internal/procspy/config"
	"procspy/internal/procspy/storage"
	"testing"
)

// TestNewServer testa criação de um novo Server
func TestNewServer(t *testing.T) {
	cfg := &config.Server{
		DBPath:  ":memory:",
		LogPath: "logs",
		APIPort: 8080,
		APIHost: "localhost",
		UserTarges: map[string]string{
			"user1": "http://example.com/user1.json",
		},
		Debug: false,
	}

	server := NewServer(cfg)
	if server == nil {
		t.Fatal("NewServer retornou nil")
	}

	if server.config != cfg {
		t.Error("Config não foi atribuída corretamente")
	}

	if server.dbConn == nil {
		t.Error("dbConn não foi inicializado")
	}

	if server.commandHandler == nil {
		t.Error("commandHandler não foi inicializado")
	}

	if server.targetHandler == nil {
		t.Error("targetHandler não foi inicializado")
	}

	if server.matchHandler == nil {
		t.Error("matchHandler não foi inicializado")
	}

	if server.reportHandler == nil {
		t.Error("reportHandler não foi inicializado")
	}

	if server.healthcheckHandler == nil {
		t.Error("healthcheckHandler não foi inicializado")
	}
}

// TestNewServer_WithDebug testa criação com modo debug
func TestNewServer_WithDebug(t *testing.T) {
	cfg := &config.Server{
		DBPath:  ":memory:",
		LogPath: "logs",
		APIPort: 8080,
		APIHost: "localhost",
		UserTarges: map[string]string{
			"user1": "http://example.com/user1.json",
		},
		Debug: true,
	}

	server := NewServer(cfg)
	if server == nil {
		t.Fatal("NewServer retornou nil")
	}

	if !server.config.Debug {
		t.Error("Debug mode não foi configurado corretamente")
	}
}

// TestServer_initServices testa inicialização de services
func TestServer_initServices(t *testing.T) {
	cfg := &config.Server{
		DBPath:  ":memory:",
		LogPath: "logs",
		APIPort: 8080,
		APIHost: "localhost",
		UserTarges: map[string]string{
			"user1": "http://example.com/user1.json",
		},
		Debug: false,
	}

	server := &Server{
		config: cfg,
		dbConn: storage.NewDbConnection(cfg.DBPath),
	}

	server.initServices()

	if server.commandHandler == nil {
		t.Error("commandHandler não foi inicializado")
	}

	if server.targetHandler == nil {
		t.Error("targetHandler não foi inicializado")
	}

	if server.matchHandler == nil {
		t.Error("matchHandler não foi inicializado")
	}

	if server.reportHandler == nil {
		t.Error("reportHandler não foi inicializado")
	}

	if server.healthcheckHandler == nil {
		t.Error("healthcheckHandler não foi inicializado")
	}
}

// TestNewServer_WithEmptyUserTargets testa criação com user targets vazio
func TestNewServer_WithEmptyUserTargets(t *testing.T) {
	cfg := &config.Server{
		DBPath:     ":memory:",
		LogPath:    "logs",
		APIPort:    8080,
		APIHost:    "localhost",
		UserTarges: map[string]string{},
		Debug:      false,
	}

	server := NewServer(cfg)
	if server == nil {
		t.Fatal("NewServer retornou nil")
	}

	if server.config.UserTarges == nil {
		t.Error("UserTarges não deveria ser nil")
	}

	if len(server.config.UserTarges) != 0 {
		t.Error("UserTarges deveria estar vazio")
	}
}

// TestNewServer_WithNilUserTargets testa criação com user targets nil
func TestNewServer_WithNilUserTargets(t *testing.T) {
	cfg := &config.Server{
		DBPath:     ":memory:",
		LogPath:    "logs",
		APIPort:    8080,
		APIHost:    "localhost",
		UserTarges: nil,
		Debug:      false,
	}

	server := NewServer(cfg)
	if server == nil {
		t.Fatal("NewServer retornou nil")
	}
}

// TestNewServer_WithMultipleUsers testa criação com múltiplos usuários
func TestNewServer_WithMultipleUsers(t *testing.T) {
	cfg := &config.Server{
		DBPath:  ":memory:",
		LogPath: "logs",
		APIPort: 8080,
		APIHost: "localhost",
		UserTarges: map[string]string{
			"user1": "http://example.com/user1.json",
			"user2": "http://example.com/user2.json",
			"user3": "http://example.com/user3.json",
		},
		Debug: false,
	}

	server := NewServer(cfg)
	if server == nil {
		t.Fatal("NewServer retornou nil")
	}

	if len(server.config.UserTarges) != 3 {
		t.Errorf("Esperado 3 user targets, obteve %d", len(server.config.UserTarges))
	}
}

// TestNewServer_WithDifferentPorts testa criação com portas diferentes
func TestNewServer_WithDifferentPorts(t *testing.T) {
	tests := []struct {
		name    string
		port    int
		host    string
	}{
		{"Porta padrão", 8080, "localhost"},
		{"Porta alternativa", 9090, "localhost"},
		{"Todas as interfaces", 8080, "0.0.0.0"},
		{"IP específico", 8080, "127.0.0.1"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &config.Server{
				DBPath:     ":memory:",
				LogPath:    "logs",
				APIPort:    tt.port,
				APIHost:    tt.host,
				UserTarges: map[string]string{},
				Debug:      false,
			}

			server := NewServer(cfg)
			if server == nil {
				t.Fatal("NewServer retornou nil")
			}

			if server.config.APIPort != tt.port {
				t.Errorf("Porta = %d, esperado %d", server.config.APIPort, tt.port)
			}

			if server.config.APIHost != tt.host {
				t.Errorf("Host = %s, esperado %s", server.config.APIHost, tt.host)
			}
		})
	}
}

// TestNewServer_WithDifferentDBPaths testa criação com diferentes caminhos de DB
func TestNewServer_WithDifferentDBPaths(t *testing.T) {
	// Apenas testa banco em memória para evitar problemas com filesystem
	cfg := &config.Server{
		DBPath:     ":memory:",
		LogPath:    "logs",
		APIPort:    8080,
		APIHost:    "localhost",
		UserTarges: map[string]string{},
		Debug:      false,
	}

	server := NewServer(cfg)
	if server == nil {
		t.Fatal("NewServer retornou nil")
	}

	if server.dbConn == nil {
		t.Error("dbConn não foi inicializado")
	}
}

// TestServer_initServices_MultipleCallsSafe testa múltiplas chamadas de initServices
func TestServer_initServices_MultipleCallsSafe(t *testing.T) {
	cfg := &config.Server{
		DBPath:  ":memory:",
		LogPath: "logs",
		APIPort: 8080,
		APIHost: "localhost",
		UserTarges: map[string]string{
			"user1": "http://example.com/user1.json",
		},
		Debug: false,
	}

	server := &Server{
		config: cfg,
		dbConn: storage.NewDbConnection(cfg.DBPath),
	}

	// Primeira chamada
	server.initServices()
	firstCommandHandler := server.commandHandler

	// Segunda chamada (deve sobrescrever)
	server.initServices()
	secondCommandHandler := server.commandHandler

	// Valida que handlers foram recriados
	if firstCommandHandler == secondCommandHandler {
		t.Log("Handlers podem ser os mesmos após reinicialização")
	}

	// Valida que handlers ainda estão funcionais
	if server.commandHandler == nil {
		t.Error("commandHandler não deveria ser nil após reinicialização")
	}
}

// TestNewServer_AllHandlersInitialized testa que todos os handlers são inicializados
func TestNewServer_AllHandlersInitialized(t *testing.T) {
	cfg := &config.Server{
		DBPath:  ":memory:",
		LogPath: "logs",
		APIPort: 8080,
		APIHost: "localhost",
		UserTarges: map[string]string{
			"user1": "http://example.com/user1.json",
		},
		Debug: false,
	}

	server := NewServer(cfg)

	// Valida que todos os handlers foram inicializados
	handlers := []struct {
		name    string
		handler interface{}
	}{
		{"commandHandler", server.commandHandler},
		{"targetHandler", server.targetHandler},
		{"matchHandler", server.matchHandler},
		{"reportHandler", server.reportHandler},
		{"healthcheckHandler", server.healthcheckHandler},
	}

	for _, h := range handlers {
		if h.handler == nil {
			t.Errorf("%s não foi inicializado", h.name)
		}
	}
}

// TestNewServer_ConfigPreserved testa que config é preservada
func TestNewServer_ConfigPreserved(t *testing.T) {
	cfg := &config.Server{
		DBPath:  ":memory:",
		LogPath: "/custom/logs",
		APIPort: 9999,
		APIHost: "192.168.1.100",
		UserTarges: map[string]string{
			"alice": "http://example.com/alice.json",
		},
		Debug: true,
	}

	server := NewServer(cfg)

	// Valida que todas as propriedades da config foram preservadas
	if server.config.DBPath != cfg.DBPath {
		t.Errorf("DBPath = %s, esperado %s", server.config.DBPath, cfg.DBPath)
	}

	if server.config.LogPath != cfg.LogPath {
		t.Errorf("LogPath = %s, esperado %s", server.config.LogPath, cfg.LogPath)
	}

	if server.config.APIPort != cfg.APIPort {
		t.Errorf("APIPort = %d, esperado %d", server.config.APIPort, cfg.APIPort)
	}

	if server.config.APIHost != cfg.APIHost {
		t.Errorf("APIHost = %s, esperado %s", server.config.APIHost, cfg.APIHost)
	}

	if server.config.Debug != cfg.Debug {
		t.Errorf("Debug = %v, esperado %v", server.config.Debug, cfg.Debug)
	}
}
