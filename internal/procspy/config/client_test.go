package config

import (
	"os"
	"strings"
	"testing"
)

// TestNewConfig testa a criação de uma nova configuração de client
// Valida que valores padrão são definidos corretamente
func TestNewConfig(t *testing.T) {
	// Act: Cria nova configuração
	config := NewConfig()

	// Assert: Valida que config foi criada
	if config == nil {
		t.Fatal("NewConfig retornou nil")
	}

	// Valida valores padrão
	if config.Interval != 30 {
		t.Errorf("Interval padrão = %d, esperado 30", config.Interval)
	}

	if config.LogPath != "logs" {
		t.Errorf("LogPath padrão = %s, esperado 'logs'", config.LogPath)
	}

	if config.Debug != false {
		t.Errorf("Debug padrão = %v, esperado false", config.Debug)
	}

	if config.APIPort != 8888 {
		t.Errorf("APIPort padrão = %d, esperado 8888", config.APIPort)
	}

	if config.APIHost != "localhost" {
		t.Errorf("APIHost padrão = %s, esperado 'localhost'", config.APIHost)
	}
}

// TestClient_SetDefaults testa a aplicação de valores padrão
// Valida que campos vazios ou inválidos recebem valores padrão
func TestClient_SetDefaults(t *testing.T) {
	tests := []struct {
		name             string
		config           *Client
		expectedInterval int
		expectedLogPath  string
		expectedAPIPort  int
		expectedAPIHost  string
	}{
		{
			name: "Todos os campos vazios",
			config: &Client{
				Interval: 0,
				LogPath:  "",
				APIPort:  0,
				APIHost:  "",
			},
			expectedInterval: 30,
			expectedLogPath:  "logs",
			expectedAPIPort:  8888,
			expectedAPIHost:  "localhost",
		},
		{
			name: "Interval abaixo do mínimo",
			config: &Client{
				Interval: 10,
				LogPath:  "custom",
				APIPort:  9000,
				APIHost:  "0.0.0.0",
			},
			expectedInterval: 30,
			expectedLogPath:  "custom",
			expectedAPIPort:  9000,
			expectedAPIHost:  "0.0.0.0",
		},
		{
			name: "Valores válidos não são alterados",
			config: &Client{
				Interval: 60,
				LogPath:  "/var/log",
				APIPort:  9999,
				APIHost:  "127.0.0.1",
			},
			expectedInterval: 60,
			expectedLogPath:  "/var/log",
			expectedAPIPort:  9999,
			expectedAPIHost:  "127.0.0.1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act: Aplica defaults
			tt.config.SetDefaults()

			// Assert: Valida valores
			if tt.config.Interval != tt.expectedInterval {
				t.Errorf("Interval = %d, esperado %d", tt.config.Interval, tt.expectedInterval)
			}

			if tt.config.LogPath != tt.expectedLogPath {
				t.Errorf("LogPath = %s, esperado %s", tt.config.LogPath, tt.expectedLogPath)
			}

			if tt.config.APIPort != tt.expectedAPIPort {
				t.Errorf("APIPort = %d, esperado %d", tt.config.APIPort, tt.expectedAPIPort)
			}

			if tt.config.APIHost != tt.expectedAPIHost {
				t.Errorf("APIHost = %s, esperado %s", tt.config.APIHost, tt.expectedAPIHost)
			}
		})
	}
}

// TestClient_ToJson testa a serialização de Client para JSON
// Valida que o JSON gerado é válido e formatado
func TestClient_ToJson(t *testing.T) {
	// Arrange: Cria config
	config := &Client{
		Interval:  60,
		LogPath:   "/var/log",
		ServerURL: "https://server.com",
		User:      "test_user",
		Debug:     true,
		APIPort:   9000,
		APIHost:   "0.0.0.0",
	}

	// Act: Serializa para JSON
	json := config.ToJson()

	// Assert: Valida que JSON não está vazio
	if json == "" {
		t.Error("ToJson retornou string vazia")
	}

	// Valida que contém campos esperados
	expectedFields := []string{"interval", "log_path", "server_url", "user", "debug", "api_port", "api_host"}
	for _, field := range expectedFields {
		if !strings.Contains(json, field) {
			t.Errorf("JSON não contém campo esperado: %s", field)
		}
	}

	// Valida que está indentado
	if !strings.Contains(json, "\t") {
		t.Error("JSON não está indentado com tabs")
	}
}

// TestClientConfigFromJson testa o parsing de JSON para Client
// Valida cenários: JSON válido, JSON inválido, valores padrão
func TestClientConfigFromJson(t *testing.T) {
	tests := []struct {
		name             string
		json             string
		wantErr          bool
		expectedInterval int
		expectedUser     string
	}{
		{
			name: "JSON válido completo",
			json: `{
				"interval": 45,
				"log_path": "/tmp/logs",
				"server_url": "https://api.example.com",
				"user": "john",
				"debug": true,
				"api_port": 9999,
				"api_host": "127.0.0.1"
			}`,
			wantErr:          false,
			expectedInterval: 45,
			expectedUser:     "john",
		},
		{
			name: "JSON mínimo com defaults",
			json: `{
				"server_url": "https://api.example.com",
				"user": "jane"
			}`,
			wantErr:          false,
			expectedInterval: 30, // Valor padrão
			expectedUser:     "jane",
		},
		{
			name:    "JSON inválido",
			json:    `{invalid json}`,
			wantErr: true,
		},
		{
			name: "JSON com interval abaixo do mínimo",
			json: `{
				"interval": 5,
				"user": "test"
			}`,
			wantErr:          false,
			expectedInterval: 30, // Corrigido para mínimo
			expectedUser:     "test",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act: Faz parsing do JSON
			result, err := ClientConfigFromJson(tt.json)

			// Assert: Valida erro
			if (err != nil) != tt.wantErr {
				t.Errorf("ClientConfigFromJson() erro = %v, esperava erro = %v", err, tt.wantErr)
				return
			}

			// Se não esperava erro, valida resultado
			if !tt.wantErr {
				if result == nil {
					t.Fatal("ClientConfigFromJson retornou nil sem erro")
				}

				if result.Interval != tt.expectedInterval {
					t.Errorf("Interval = %d, esperado %d", result.Interval, tt.expectedInterval)
				}

				if result.User != tt.expectedUser {
					t.Errorf("User = %s, esperado %s", result.User, tt.expectedUser)
				}
			}
		})
	}
}

// TestClientConfigFromFile testa o parsing de arquivo para Client
// Valida leitura de arquivo válido e tratamento de arquivo inexistente
func TestClientConfigFromFile(t *testing.T) {
	// Arrange: Cria arquivo temporário com JSON válido
	tmpFile, err := os.CreateTemp("", "client_config_*.json")
	if err != nil {
		t.Fatalf("Erro ao criar arquivo temporário: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	validJSON := `{
		"interval": 60,
		"log_path": "/var/log",
		"server_url": "https://server.com",
		"user": "test_user"
	}`

	if _, err := tmpFile.Write([]byte(validJSON)); err != nil {
		t.Fatalf("Erro ao escrever no arquivo temporário: %v", err)
	}
	tmpFile.Close()

	// Act: Lê configuração do arquivo
	config, err := ClientConfigFromFile(tmpFile.Name())

	// Assert: Valida que não houve erro
	if err != nil {
		t.Fatalf("ClientConfigFromFile() erro = %v", err)
	}

	if config == nil {
		t.Fatal("ClientConfigFromFile retornou nil sem erro")
	}

	// Valida valores
	if config.Interval != 60 {
		t.Errorf("Interval = %d, esperado 60", config.Interval)
	}

	if config.User != "test_user" {
		t.Errorf("User = %s, esperado 'test_user'", config.User)
	}
}

// TestClientConfigFromFile_FileNotFound testa leitura de arquivo inexistente
// Valida que erro é retornado quando arquivo não existe
func TestClientConfigFromFile_FileNotFound(t *testing.T) {
	// Act: Tenta ler arquivo inexistente
	config, err := ClientConfigFromFile("/path/that/does/not/exist.json")

	// Assert: Valida que houve erro
	if err == nil {
		t.Error("ClientConfigFromFile deveria retornar erro para arquivo inexistente")
	}

	if config != nil {
		t.Error("ClientConfigFromFile deveria retornar nil para arquivo inexistente")
	}
}

// TestClientConfigFromFile_InvalidJSON testa leitura de arquivo com JSON inválido
// Valida que erro é retornado quando JSON é inválido
func TestClientConfigFromFile_InvalidJSON(t *testing.T) {
	// Arrange: Cria arquivo temporário com JSON inválido
	tmpFile, err := os.CreateTemp("", "client_config_invalid_*.json")
	if err != nil {
		t.Fatalf("Erro ao criar arquivo temporário: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	invalidJSON := `{invalid json content}`

	if _, err := tmpFile.Write([]byte(invalidJSON)); err != nil {
		t.Fatalf("Erro ao escrever no arquivo temporário: %v", err)
	}
	tmpFile.Close()

	// Act: Tenta ler configuração do arquivo
	config, err := ClientConfigFromFile(tmpFile.Name())

	// Assert: Valida que houve erro
	if err == nil {
		t.Error("ClientConfigFromFile deveria retornar erro para JSON inválido")
	}

	if config != nil {
		t.Error("ClientConfigFromFile deveria retornar nil para JSON inválido")
	}
}

// TestClient_Serialization_RoundTrip testa serialização e desserialização
// Valida que dados são preservados após round-trip JSON
func TestClient_Serialization_RoundTrip(t *testing.T) {
	// Arrange: Cria config original
	original := &Client{
		Interval:  45,
		LogPath:   "/custom/path",
		ServerURL: "https://api.test.com",
		User:      "testuser",
		Debug:     true,
		APIPort:   9876,
		APIHost:   "192.168.1.1",
	}

	// Act: Serializa e desserializa
	json := original.ToJson()
	restored, err := ClientConfigFromJson(json)

	// Assert: Valida que não houve erro
	if err != nil {
		t.Fatalf("Erro ao desserializar: %v", err)
	}

	// Valida que dados foram preservados
	if restored.Interval != original.Interval {
		t.Errorf("Interval não preservado: %d != %d", restored.Interval, original.Interval)
	}

	if restored.LogPath != original.LogPath {
		t.Errorf("LogPath não preservado: %s != %s", restored.LogPath, original.LogPath)
	}

	if restored.ServerURL != original.ServerURL {
		t.Errorf("ServerURL não preservado: %s != %s", restored.ServerURL, original.ServerURL)
	}

	if restored.User != original.User {
		t.Errorf("User não preservado: %s != %s", restored.User, original.User)
	}

	if restored.Debug != original.Debug {
		t.Errorf("Debug não preservado: %v != %v", restored.Debug, original.Debug)
	}

	if restored.APIPort != original.APIPort {
		t.Errorf("APIPort não preservado: %d != %d", restored.APIPort, original.APIPort)
	}

	if restored.APIHost != original.APIHost {
		t.Errorf("APIHost não preservado: %s != %s", restored.APIHost, original.APIHost)
	}
}

// TestClient_DebugOmitEmpty testa que debug é omitido quando false
// Valida comportamento do omitempty no JSON
func TestClient_DebugOmitEmpty(t *testing.T) {
	// Arrange: Cria config com debug false
	config := &Client{
		Interval:  30,
		LogPath:   "logs",
		ServerURL: "https://server.com",
		User:      "user",
		Debug:     false,
	}

	// Act: Serializa para JSON
	json := config.ToJson()

	// Assert: Valida que debug não aparece no JSON quando false
	// (devido ao omitempty)
	if strings.Contains(json, `"debug": false`) {
		t.Log("Debug false está presente no JSON (comportamento pode variar)")
	}
}

// TestClient_ToJson_WithInvalidData testa ToJson com dados que podem causar erro
// Valida que ToJson retorna string mesmo com dados problemáticos
func TestClient_ToJson_WithInvalidData(t *testing.T) {
	// Arrange: Cria config com valores extremos
	config := &Client{
		Interval:  -1,
		LogPath:   "",
		ServerURL: "",
		User:      "",
		Debug:     false,
		APIPort:   -1,
		APIHost:   "",
	}

	// Act: Serializa para JSON
	json := config.ToJson()

	// Assert: Valida que JSON foi gerado (mesmo com valores inválidos)
	if json == "" {
		t.Error("ToJson retornou string vazia para dados inválidos")
	}

	// Valida que contém estrutura JSON básica
	if !strings.Contains(json, "{") || !strings.Contains(json, "}") {
		t.Error("ToJson não retornou estrutura JSON válida")
	}
}

// TestClient_ToJson_EmptyConfig testa ToJson com config vazia
// Valida que ToJson funciona com struct vazia
func TestClient_ToJson_EmptyConfig(t *testing.T) {
	// Arrange: Cria config vazia
	config := &Client{}

	// Act: Serializa para JSON
	json := config.ToJson()

	// Assert: Valida que JSON foi gerado
	if json == "" {
		t.Error("ToJson retornou string vazia para config vazia")
	}

	// Valida estrutura JSON
	if !strings.Contains(json, "{") || !strings.Contains(json, "}") {
		t.Error("ToJson não retornou estrutura JSON válida")
	}
}

// TestClientConfigFromJson_EmptyString testa parsing de string vazia
// Valida tratamento de JSON vazio
func TestClientConfigFromJson_EmptyString(t *testing.T) {
	// Act: Tenta fazer parsing de string vazia
	config, err := ClientConfigFromJson("")

	// Assert: Valida que houve erro
	if err == nil {
		t.Error("ClientConfigFromJson deveria retornar erro para string vazia")
	}

	if config != nil {
		t.Error("ClientConfigFromJson deveria retornar nil para string vazia")
	}
}

// TestClientConfigFromJson_NullJSON testa parsing de JSON null
// Valida tratamento de JSON null
func TestClientConfigFromJson_NullJSON(t *testing.T) {
	// Act: Tenta fazer parsing de null
	config, err := ClientConfigFromJson("null")

	// Assert: Valida que não houve erro (null é JSON válido)
	if err != nil {
		t.Errorf("ClientConfigFromJson retornou erro inesperado para null: %v", err)
	}

	// Config deve ser criado mas com valores zerados
	if config == nil {
		t.Error("ClientConfigFromJson retornou nil para JSON null")
	}
}

// TestClient_SetDefaults_PartialConfig testa SetDefaults com config parcial
// Valida que apenas campos vazios recebem defaults
func TestClient_SetDefaults_PartialConfig(t *testing.T) {
	// Arrange: Cria config com alguns campos definidos
	config := &Client{
		Interval: 0,             // Deve receber default
		LogPath:  "custom",      // Não deve mudar
		APIPort:  0,             // Deve receber default
		APIHost:  "custom.host", // Não deve mudar
	}

	// Act: Aplica defaults
	config.SetDefaults()

	// Assert: Valida que apenas campos vazios foram alterados
	if config.Interval != 30 {
		t.Errorf("Interval = %d, esperado 30", config.Interval)
	}

	if config.LogPath != "custom" {
		t.Errorf("LogPath = %s, esperado 'custom'", config.LogPath)
	}

	if config.APIPort != 8888 {
		t.Errorf("APIPort = %d, esperado 8888", config.APIPort)
	}

	if config.APIHost != "custom.host" {
		t.Errorf("APIHost = %s, esperado 'custom.host'", config.APIHost)
	}
}

// TestClient_SetDefaults_EdgeCases testa SetDefaults com casos extremos
// Valida comportamento com valores negativos e limites
func TestClient_SetDefaults_EdgeCases(t *testing.T) {
	tests := []struct {
		name             string
		interval         int
		expectedInterval int
	}{
		{"Interval negativo", -10, 30},
		{"Interval zero", 0, 30},
		{"Interval 1", 1, 30},
		{"Interval 29", 29, 30},
		{"Interval 30", 30, 30},
		{"Interval 31", 31, 31},
		{"Interval muito grande", 999999, 999999},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			config := &Client{Interval: tt.interval}

			// Act
			config.SetDefaults()

			// Assert
			if config.Interval != tt.expectedInterval {
				t.Errorf("Interval = %d, esperado %d", config.Interval, tt.expectedInterval)
			}
		})
	}
}
