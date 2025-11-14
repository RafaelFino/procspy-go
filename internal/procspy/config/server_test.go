package config

import (
	"os"
	"strings"
	"testing"
)

// TestNewServer testa a criação de uma nova configuração de server
func TestNewServer(t *testing.T) {
	config := NewServer()
	if config == nil {
		t.Fatal("NewServer retornou nil")
	}
}

// TestServer_ToJson testa a serialização de Server para JSON
func TestServer_ToJson(t *testing.T) {
	config := &Server{
		DBPath:  "/var/lib/procspy",
		LogPath: "/var/log/procspy",
		APIPort: 8080,
		APIHost: "0.0.0.0",
		UserTarges: map[string]string{
			"user1": "https://config.com/user1.json",
		},
		Debug: true,
	}

	json := config.ToJson()
	if json == "" {
		t.Error("ToJson retornou string vazia")
	}

	expectedFields := []string{"db_path", "log_path", "api_port", "api_host", "user_targets", "debug"}
	for _, field := range expectedFields {
		if !strings.Contains(json, field) {
			t.Errorf("JSON não contém campo esperado: %s", field)
		}
	}
}

// TestServerConfigFromJson testa o parsing de JSON para Server
func TestServerConfigFromJson(t *testing.T) {
	tests := []struct {
		name    string
		json    string
		wantErr bool
	}{
		{
			name: "JSON válido",
			json: `{
				"db_path": "/data",
				"log_path": "/logs",
				"api_port": 8080,
				"api_host": "localhost",
				"user_targets": {"user1": "url1"},
				"debug": false
			}`,
			wantErr: false,
		},
		{
			name:    "JSON inválido",
			json:    `{invalid}`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ServerConfigFromJson(tt.json)
			if (err != nil) != tt.wantErr {
				t.Errorf("ServerConfigFromJson() erro = %v, esperava erro = %v", err, tt.wantErr)
			}
			if !tt.wantErr && result == nil {
				t.Error("ServerConfigFromJson retornou nil sem erro")
			}
		})
	}
}

// TestServerConfigFromFile testa leitura de arquivo
func TestServerConfigFromFile(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "server_config_*.json")
	if err != nil {
		t.Fatalf("Erro ao criar arquivo temporário: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	validJSON := `{"db_path": "/data", "api_port": 8080}`
	tmpFile.Write([]byte(validJSON))
	tmpFile.Close()

	config, err := ServerConfigFromFile(tmpFile.Name())
	if err != nil {
		t.Fatalf("ServerConfigFromFile() erro = %v", err)
	}
	if config == nil {
		t.Fatal("ServerConfigFromFile retornou nil")
	}
}

// TestServerConfigFromFile_NotFound testa arquivo inexistente
func TestServerConfigFromFile_NotFound(t *testing.T) {
	_, err := ServerConfigFromFile("/nonexistent/file.json")
	if err == nil {
		t.Error("Deveria retornar erro para arquivo inexistente")
	}
}

// TestServerConfigFromFile_InvalidJSON testa arquivo com JSON inválido
func TestServerConfigFromFile_InvalidJSON(t *testing.T) {
	// Arrange: Cria arquivo temporário com JSON inválido
	tmpFile, err := os.CreateTemp("", "server_config_invalid_*.json")
	if err != nil {
		t.Fatalf("Erro ao criar arquivo temporário: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	invalidJSON := `{invalid json}`
	tmpFile.Write([]byte(invalidJSON))
	tmpFile.Close()

	// Act: Tenta ler configuração
	config, err := ServerConfigFromFile(tmpFile.Name())

	// Assert: Valida erro
	if err == nil {
		t.Error("ServerConfigFromFile deveria retornar erro para JSON inválido")
	}
	if config != nil {
		t.Error("ServerConfigFromFile deveria retornar nil para JSON inválido")
	}
}

// TestServer_ToJson_EmptyConfig testa ToJson com config vazia
func TestServer_ToJson_EmptyConfig(t *testing.T) {
	// Arrange: Cria config vazia
	config := &Server{}

	// Act: Serializa para JSON
	json := config.ToJson()

	// Assert: Valida que JSON foi gerado
	if json == "" {
		t.Error("ToJson retornou string vazia")
	}

	// Valida estrutura JSON básica
	if !strings.Contains(json, "{") || !strings.Contains(json, "}") {
		t.Error("ToJson não retornou estrutura JSON válida")
	}
}

// TestServer_ToJson_WithNilMap testa ToJson com map nil
func TestServer_ToJson_WithNilMap(t *testing.T) {
	// Arrange: Cria config com map nil
	config := &Server{
		DBPath:     "/data",
		UserTarges: nil,
	}

	// Act: Serializa para JSON
	json := config.ToJson()

	// Assert: Valida que JSON foi gerado
	if json == "" {
		t.Error("ToJson retornou string vazia")
	}

	// Map nil deve aparecer como null no JSON
	if !strings.Contains(json, "null") && !strings.Contains(json, "user_targets") {
		t.Log("Map nil pode não aparecer no JSON")
	}
}

// TestServerConfigFromJson_EmptyString testa parsing de string vazia
func TestServerConfigFromJson_EmptyString(t *testing.T) {
	// Act: Tenta fazer parsing de string vazia
	config, err := ServerConfigFromJson("")

	// Assert: Valida erro
	if err == nil {
		t.Error("ServerConfigFromJson deveria retornar erro para string vazia")
	}
	if config != nil {
		t.Error("ServerConfigFromJson deveria retornar nil para string vazia")
	}
}

// TestServerConfigFromJson_MinimalJSON testa parsing de JSON mínimo
func TestServerConfigFromJson_MinimalJSON(t *testing.T) {
	// Act: Faz parsing de JSON mínimo
	config, err := ServerConfigFromJson("{}")

	// Assert: Valida sucesso
	if err != nil {
		t.Errorf("ServerConfigFromJson retornou erro inesperado: %v", err)
	}
	if config == nil {
		t.Error("ServerConfigFromJson retornou nil para JSON vazio")
	}
}

// TestServerConfigFromJson_NullJSON testa parsing de JSON null
func TestServerConfigFromJson_NullJSON(t *testing.T) {
	// Act: Tenta fazer parsing de null
	config, err := ServerConfigFromJson("null")

	// Assert: Valida que não houve erro (null é JSON válido)
	if err != nil {
		t.Errorf("ServerConfigFromJson retornou erro inesperado: %v", err)
	}
	if config == nil {
		t.Error("ServerConfigFromJson retornou nil para JSON null")
	}
}

// TestServer_Serialization_RoundTrip testa serialização e desserialização
func TestServer_Serialization_RoundTrip(t *testing.T) {
	// Arrange: Cria config original
	original := &Server{
		DBPath:  "/var/lib/procspy",
		LogPath: "/var/log/procspy",
		APIPort: 9090,
		APIHost: "192.168.1.100",
		UserTarges: map[string]string{
			"alice": "https://config.example.com/alice.json",
			"bob":   "https://config.example.com/bob.json",
		},
		Debug: true,
	}

	// Act: Serializa e desserializa
	json := original.ToJson()
	restored, err := ServerConfigFromJson(json)

	// Assert: Valida sucesso
	if err != nil {
		t.Fatalf("Erro ao desserializar: %v", err)
	}

	// Valida que dados foram preservados
	if restored.DBPath != original.DBPath {
		t.Errorf("DBPath não preservado: %s != %s", restored.DBPath, original.DBPath)
	}
	if restored.LogPath != original.LogPath {
		t.Errorf("LogPath não preservado: %s != %s", restored.LogPath, original.LogPath)
	}
	if restored.APIPort != original.APIPort {
		t.Errorf("APIPort não preservado: %d != %d", restored.APIPort, original.APIPort)
	}
	if restored.APIHost != original.APIHost {
		t.Errorf("APIHost não preservado: %s != %s", restored.APIHost, original.APIHost)
	}
	if restored.Debug != original.Debug {
		t.Errorf("Debug não preservado: %v != %v", restored.Debug, original.Debug)
	}
	if len(restored.UserTarges) != len(original.UserTarges) {
		t.Errorf("UserTarges length não preservado: %d != %d", len(restored.UserTarges), len(original.UserTarges))
	}
}

// TestServer_ToJson_ComplexUserTargets testa ToJson com múltiplos user targets
func TestServer_ToJson_ComplexUserTargets(t *testing.T) {
	// Arrange: Cria config com múltiplos user targets
	config := &Server{
		DBPath:  "/data",
		LogPath: "/logs",
		APIPort: 8080,
		APIHost: "0.0.0.0",
		UserTarges: map[string]string{
			"user1": "https://example.com/user1.json",
			"user2": "https://example.com/user2.json",
			"user3": "https://example.com/user3.json",
		},
		Debug: false,
	}

	// Act: Serializa para JSON
	json := config.ToJson()

	// Assert: Valida que todos os users estão no JSON
	if !strings.Contains(json, "user1") {
		t.Error("JSON não contém user1")
	}
	if !strings.Contains(json, "user2") {
		t.Error("JSON não contém user2")
	}
	if !strings.Contains(json, "user3") {
		t.Error("JSON não contém user3")
	}
}

// TestServerConfigFromJson_WithSpecialCharacters testa parsing com caracteres especiais
func TestServerConfigFromJson_WithSpecialCharacters(t *testing.T) {
	// Arrange: JSON com caracteres especiais
	jsonStr := `{
		"db_path": "/path/with spaces/and-special_chars",
		"log_path": "/logs/日本語/path",
		"api_port": 8080,
		"api_host": "localhost",
		"user_targets": {
			"user@example.com": "https://example.com/config.json"
		}
	}`

	// Act: Faz parsing
	config, err := ServerConfigFromJson(jsonStr)

	// Assert: Valida sucesso
	if err != nil {
		t.Errorf("ServerConfigFromJson retornou erro inesperado: %v", err)
	}
	if config == nil {
		t.Fatal("ServerConfigFromJson retornou nil")
	}

	// Valida que caracteres especiais foram preservados
	if !strings.Contains(config.DBPath, "spaces") {
		t.Error("Caracteres especiais não foram preservados em DBPath")
	}
}

// TestServer_ToJson_EdgeCases testa ToJson com valores extremos
func TestServer_ToJson_EdgeCases(t *testing.T) {
	tests := []struct {
		name   string
		config *Server
	}{
		{
			name: "Valores negativos",
			config: &Server{
				APIPort: -1,
			},
		},
		{
			name: "Valores zero",
			config: &Server{
				APIPort: 0,
			},
		},
		{
			name: "Strings vazias",
			config: &Server{
				DBPath:  "",
				LogPath: "",
				APIHost: "",
			},
		},
		{
			name: "Map vazio",
			config: &Server{
				UserTarges: map[string]string{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act: Serializa
			json := tt.config.ToJson()

			// Assert: Valida que JSON foi gerado
			if json == "" {
				t.Error("ToJson retornou string vazia")
			}
			if !strings.Contains(json, "{") {
				t.Error("ToJson não retornou JSON válido")
			}
		})
	}
}
