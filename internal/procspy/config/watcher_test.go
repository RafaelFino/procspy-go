package config

import (
	"os"
	"strings"
	"testing"
)

// TestNewWatcher testa a criação de uma nova configuração de watcher
func TestNewWatcher(t *testing.T) {
	config := NewWatcher()
	if config == nil {
		t.Fatal("NewWatcher retornou nil")
	}

	// Valida valores padrão
	if config.Interval != 10 {
		t.Errorf("Interval padrão = %d, esperado 10", config.Interval)
	}
	if config.LogPath != "logs" {
		t.Errorf("LogPath padrão = %s, esperado 'logs'", config.LogPath)
	}
	if config.ProcspyURL != "http://localhost:8888" {
		t.Errorf("ProcspyURL padrão = %s, esperado 'http://localhost:8888'", config.ProcspyURL)
	}
	if config.StartCmd != "" {
		t.Errorf("StartCmd padrão = %s, esperado ''", config.StartCmd)
	}
}

// TestWatcher_SetDefaults testa a aplicação de valores padrão
func TestWatcher_SetDefaults(t *testing.T) {
	tests := []struct {
		name               string
		config             *Watcher
		expectedInterval   int
		expectedLogPath    string
		expectedProcspyURL string
	}{
		{
			name: "Todos os campos vazios",
			config: &Watcher{
				Interval:   0,
				LogPath:    "",
				ProcspyURL: "",
			},
			expectedInterval:   10,
			expectedLogPath:    "logs",
			expectedProcspyURL: "http://localhost:8888",
		},
		{
			name: "Interval abaixo do mínimo",
			config: &Watcher{
				Interval:   5,
				LogPath:    "custom",
				ProcspyURL: "http://custom:9999",
			},
			expectedInterval:   10,
			expectedLogPath:    "custom",
			expectedProcspyURL: "http://custom:9999",
		},
		{
			name: "Valores válidos não são alterados",
			config: &Watcher{
				Interval:   30,
				LogPath:    "/var/log",
				ProcspyURL: "http://server:8080",
			},
			expectedInterval:   30,
			expectedLogPath:    "/var/log",
			expectedProcspyURL: "http://server:8080",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.config.SetDefaults()

			if tt.config.Interval != tt.expectedInterval {
				t.Errorf("Interval = %d, esperado %d", tt.config.Interval, tt.expectedInterval)
			}
			if tt.config.LogPath != tt.expectedLogPath {
				t.Errorf("LogPath = %s, esperado %s", tt.config.LogPath, tt.expectedLogPath)
			}
			if tt.config.ProcspyURL != tt.expectedProcspyURL {
				t.Errorf("ProcspyURL = %s, esperado %s", tt.config.ProcspyURL, tt.expectedProcspyURL)
			}
		})
	}
}

// TestWatcher_ToJson testa a serialização de Watcher para JSON
func TestWatcher_ToJson(t *testing.T) {
	config := &Watcher{
		Interval:   20,
		LogPath:    "/var/log",
		ProcspyURL: "http://localhost:8888",
		StartCmd:   "systemctl restart procspy-client",
	}

	json := config.ToJson()
	if json == "" {
		t.Error("ToJson retornou string vazia")
	}

	expectedFields := []string{"interval", "log_path", "procspy_url", "start_cmd"}
	for _, field := range expectedFields {
		if !strings.Contains(json, field) {
			t.Errorf("JSON não contém campo esperado: %s", field)
		}
	}
}

// TestWatcherConfigFromJson testa o parsing de JSON para Watcher
func TestWatcherConfigFromJson(t *testing.T) {
	tests := []struct {
		name    string
		json    string
		wantErr bool
	}{
		{
			name: "JSON válido",
			json: `{
				"interval": 15,
				"log_path": "/logs",
				"procspy_url": "http://localhost:8888",
				"start_cmd": "systemctl restart procspy"
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
			result, err := WatcherConfigFromJson(tt.json)
			if (err != nil) != tt.wantErr {
				t.Errorf("WatcherConfigFromJson() erro = %v, esperava erro = %v", err, tt.wantErr)
			}
			if !tt.wantErr && result == nil {
				t.Error("WatcherConfigFromJson retornou nil sem erro")
			}
		})
	}
}

// TestWatcherConfigFromFile testa leitura de arquivo
func TestWatcherConfigFromFile(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "watcher_config_*.json")
	if err != nil {
		t.Fatalf("Erro ao criar arquivo temporário: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	validJSON := `{"interval": 20, "procspy_url": "http://localhost:8888"}`
	tmpFile.Write([]byte(validJSON))
	tmpFile.Close()

	config, err := WatcherConfigFromFile(tmpFile.Name())
	if err != nil {
		t.Fatalf("WatcherConfigFromFile() erro = %v", err)
	}
	if config == nil {
		t.Fatal("WatcherConfigFromFile retornou nil")
	}
}

// TestWatcherConfigFromFile_NotFound testa arquivo inexistente
func TestWatcherConfigFromFile_NotFound(t *testing.T) {
	_, err := WatcherConfigFromFile("/nonexistent/file.json")
	if err == nil {
		t.Error("Deveria retornar erro para arquivo inexistente")
	}
}

// TestWatcherConfigFromFile_InvalidJSON testa arquivo com JSON inválido
func TestWatcherConfigFromFile_InvalidJSON(t *testing.T) {
	// Arrange: Cria arquivo temporário com JSON inválido
	tmpFile, err := os.CreateTemp("", "watcher_config_invalid_*.json")
	if err != nil {
		t.Fatalf("Erro ao criar arquivo temporário: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	invalidJSON := `{invalid json content}`
	tmpFile.Write([]byte(invalidJSON))
	tmpFile.Close()

	// Act: Tenta ler configuração
	config, err := WatcherConfigFromFile(tmpFile.Name())

	// Assert: Valida erro
	if err == nil {
		t.Error("WatcherConfigFromFile deveria retornar erro para JSON inválido")
	}
	if config != nil {
		t.Error("WatcherConfigFromFile deveria retornar nil para JSON inválido")
	}
}

// TestWatcher_ToJson_EmptyConfig testa ToJson com config vazia
func TestWatcher_ToJson_EmptyConfig(t *testing.T) {
	// Arrange: Cria config vazia
	config := &Watcher{}

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

// TestWatcherConfigFromJson_EmptyString testa parsing de string vazia
func TestWatcherConfigFromJson_EmptyString(t *testing.T) {
	// Act: Tenta fazer parsing de string vazia
	config, err := WatcherConfigFromJson("")

	// Assert: Valida erro
	if err == nil {
		t.Error("WatcherConfigFromJson deveria retornar erro para string vazia")
	}
	if config != nil {
		t.Error("WatcherConfigFromJson deveria retornar nil para string vazia")
	}
}

// TestWatcherConfigFromJson_MinimalJSON testa parsing de JSON mínimo
func TestWatcherConfigFromJson_MinimalJSON(t *testing.T) {
	// Act: Faz parsing de JSON mínimo
	config, err := WatcherConfigFromJson("{}")

	// Assert: Valida sucesso e defaults aplicados
	if err != nil {
		t.Errorf("WatcherConfigFromJson retornou erro inesperado: %v", err)
	}
	if config == nil {
		t.Fatal("WatcherConfigFromJson retornou nil")
	}

	// Valida que defaults foram aplicados
	if config.Interval != 10 {
		t.Errorf("Interval = %d, esperado 10 (default)", config.Interval)
	}
	if config.LogPath != "logs" {
		t.Errorf("LogPath = %s, esperado 'logs' (default)", config.LogPath)
	}
	if config.ProcspyURL != "http://localhost:8888" {
		t.Errorf("ProcspyURL = %s, esperado 'http://localhost:8888' (default)", config.ProcspyURL)
	}
}

// TestWatcherConfigFromJson_NullJSON testa parsing de JSON null
func TestWatcherConfigFromJson_NullJSON(t *testing.T) {
	// Act: Tenta fazer parsing de null
	config, err := WatcherConfigFromJson("null")

	// Assert: Valida que não houve erro (null é JSON válido)
	if err != nil {
		t.Errorf("WatcherConfigFromJson retornou erro inesperado: %v", err)
	}
	if config == nil {
		t.Error("WatcherConfigFromJson retornou nil para JSON null")
	}
}

// TestWatcher_Serialization_RoundTrip testa serialização e desserialização
func TestWatcher_Serialization_RoundTrip(t *testing.T) {
	// Arrange: Cria config original
	original := &Watcher{
		Interval:   25,
		LogPath:    "/custom/logs",
		ProcspyURL: "http://192.168.1.100:9999",
		StartCmd:   "/usr/bin/start-procspy.sh",
	}

	// Act: Serializa e desserializa
	json := original.ToJson()
	restored, err := WatcherConfigFromJson(json)

	// Assert: Valida sucesso
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
	if restored.ProcspyURL != original.ProcspyURL {
		t.Errorf("ProcspyURL não preservado: %s != %s", restored.ProcspyURL, original.ProcspyURL)
	}
	if restored.StartCmd != original.StartCmd {
		t.Errorf("StartCmd não preservado: %s != %s", restored.StartCmd, original.StartCmd)
	}
}

// TestWatcher_SetDefaults_EdgeCases testa SetDefaults com casos extremos
func TestWatcher_SetDefaults_EdgeCases(t *testing.T) {
	tests := []struct {
		name             string
		interval         int
		expectedInterval int
	}{
		{"Interval negativo", -5, 10},
		{"Interval zero", 0, 10},
		{"Interval 1", 1, 10},
		{"Interval 9", 9, 10},
		{"Interval 10", 10, 10},
		{"Interval 11", 11, 11},
		{"Interval muito grande", 999999, 999999},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			config := &Watcher{Interval: tt.interval}

			// Act
			config.SetDefaults()

			// Assert
			if config.Interval != tt.expectedInterval {
				t.Errorf("Interval = %d, esperado %d", config.Interval, tt.expectedInterval)
			}
		})
	}
}

// TestWatcher_SetDefaults_PartialConfig testa SetDefaults com config parcial
func TestWatcher_SetDefaults_PartialConfig(t *testing.T) {
	// Arrange: Cria config com alguns campos definidos
	config := &Watcher{
		Interval:   0,                  // Deve receber default
		LogPath:    "custom/path",      // Não deve mudar
		ProcspyURL: "",                 // Deve receber default
		StartCmd:   "/custom/start.sh", // Não deve mudar
	}

	// Act: Aplica defaults
	config.SetDefaults()

	// Assert: Valida que apenas campos vazios foram alterados
	if config.Interval != 10 {
		t.Errorf("Interval = %d, esperado 10", config.Interval)
	}
	if config.LogPath != "custom/path" {
		t.Errorf("LogPath = %s, esperado 'custom/path'", config.LogPath)
	}
	if config.ProcspyURL != "http://localhost:8888" {
		t.Errorf("ProcspyURL = %s, esperado 'http://localhost:8888'", config.ProcspyURL)
	}
	if config.StartCmd != "/custom/start.sh" {
		t.Errorf("StartCmd = %s, esperado '/custom/start.sh'", config.StartCmd)
	}
}

// TestWatcherConfigFromJson_WithIntervalBelowMinimum testa parsing com interval abaixo do mínimo
func TestWatcherConfigFromJson_WithIntervalBelowMinimum(t *testing.T) {
	// Arrange: JSON com interval abaixo do mínimo
	jsonStr := `{
		"interval": 5,
		"log_path": "/logs",
		"procspy_url": "http://localhost:8888"
	}`

	// Act: Faz parsing
	config, err := WatcherConfigFromJson(jsonStr)

	// Assert: Valida sucesso
	if err != nil {
		t.Errorf("WatcherConfigFromJson retornou erro inesperado: %v", err)
	}
	if config == nil {
		t.Fatal("WatcherConfigFromJson retornou nil")
	}

	// Valida que interval foi corrigido para o mínimo
	if config.Interval != 10 {
		t.Errorf("Interval = %d, esperado 10 (mínimo)", config.Interval)
	}
}

// TestWatcher_ToJson_WithSpecialCharacters testa ToJson com caracteres especiais
func TestWatcher_ToJson_WithSpecialCharacters(t *testing.T) {
	// Arrange: Cria config com caracteres especiais
	config := &Watcher{
		Interval:   15,
		LogPath:    "/path/with spaces/and-special_chars",
		ProcspyURL: "http://localhost:8888/api/v1",
		StartCmd:   "bash -c 'systemctl restart procspy'",
	}

	// Act: Serializa para JSON
	json := config.ToJson()

	// Assert: Valida que JSON foi gerado
	if json == "" {
		t.Error("ToJson retornou string vazia")
	}

	// Valida que caracteres especiais estão presentes
	if !strings.Contains(json, "spaces") {
		t.Error("Caracteres especiais não foram preservados")
	}
}

// TestWatcher_ToJson_EdgeCases testa ToJson com valores extremos
func TestWatcher_ToJson_EdgeCases(t *testing.T) {
	tests := []struct {
		name   string
		config *Watcher
	}{
		{
			name: "Valores negativos",
			config: &Watcher{
				Interval: -1,
			},
		},
		{
			name: "Valores zero",
			config: &Watcher{
				Interval: 0,
			},
		},
		{
			name: "Strings vazias",
			config: &Watcher{
				LogPath:    "",
				ProcspyURL: "",
				StartCmd:   "",
			},
		},
		{
			name: "Valores muito grandes",
			config: &Watcher{
				Interval: 999999999,
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

// TestWatcherConfigFromJson_CompleteConfig testa parsing de config completa
func TestWatcherConfigFromJson_CompleteConfig(t *testing.T) {
	// Arrange: JSON completo
	jsonStr := `{
		"interval": 30,
		"log_path": "/var/log/procspy",
		"procspy_url": "http://192.168.1.100:8888",
		"start_cmd": "systemctl restart procspy-client.service"
	}`

	// Act: Faz parsing
	config, err := WatcherConfigFromJson(jsonStr)

	// Assert: Valida sucesso
	if err != nil {
		t.Fatalf("WatcherConfigFromJson retornou erro: %v", err)
	}
	if config == nil {
		t.Fatal("WatcherConfigFromJson retornou nil")
	}

	// Valida todos os campos
	if config.Interval != 30 {
		t.Errorf("Interval = %d, esperado 30", config.Interval)
	}
	if config.LogPath != "/var/log/procspy" {
		t.Errorf("LogPath = %s, esperado '/var/log/procspy'", config.LogPath)
	}
	if config.ProcspyURL != "http://192.168.1.100:8888" {
		t.Errorf("ProcspyURL = %s, esperado 'http://192.168.1.100:8888'", config.ProcspyURL)
	}
	if config.StartCmd != "systemctl restart procspy-client.service" {
		t.Errorf("StartCmd = %s, esperado 'systemctl restart procspy-client.service'", config.StartCmd)
	}
}
