package client

import (
	"net/http"
	"net/http/httptest"
	"procspy/internal/procspy/config"
	"procspy/internal/procspy/domain"
	"strings"
	"testing"
	"time"
)

// TestNewSpy testa criação de um novo Spy
func TestNewSpy(t *testing.T) {
	cfg := &config.Client{
		Interval:  30,
		LogPath:   "logs",
		ServerURL: "http://localhost:8080",
		User:      "test_user",
		APIPort:   8888,
		APIHost:   "localhost",
	}

	spy := NewSpy(cfg)
	if spy == nil {
		t.Fatal("NewSpy retornou nil")
	}

	if spy.config != cfg {
		t.Error("Config não foi atribuída corretamente")
	}

	if spy.enabled {
		t.Error("Spy deveria iniciar desabilitado")
	}

	if spy.targets == nil {
		t.Error("Targets não foi inicializado")
	}

	if spy.commandBuf == nil {
		t.Error("commandBuf não foi inicializado")
	}

	if spy.matchBuf == nil {
		t.Error("matchBuf não foi inicializado")
	}
}

// TestSpy_IsEnabled testa verificação de estado
func TestSpy_IsEnabled(t *testing.T) {
	cfg := &config.Client{
		Interval:  30,
		ServerURL: "http://localhost:8080",
		User:      "test",
	}

	spy := NewSpy(cfg)

	if spy.IsEnabled() {
		t.Error("Spy deveria estar desabilitado inicialmente")
	}

	spy.enabled = true
	if !spy.IsEnabled() {
		t.Error("Spy deveria estar habilitado após enabled = true")
	}
}

// TestRoundFloat testa arredondamento de float
func TestRoundFloat(t *testing.T) {
	tests := []struct {
		name      string
		value     float64
		precision uint
		expected  float64
	}{
		{"2 casas decimais", 3.14159, 2, 3.14},
		{"1 casa decimal", 2.56, 1, 2.6},
		{"0 casas decimais", 5.7, 0, 6.0},
		{"Valor exato", 10.0, 2, 10.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := roundFloat(tt.value, tt.precision)
			if result != tt.expected {
				t.Errorf("roundFloat(%.5f, %d) = %.5f, esperado %.5f",
					tt.value, tt.precision, result, tt.expected)
			}
		})
	}
}

// TestSpy_httpGet testa requisições HTTP GET
func TestSpy_httpGet(t *testing.T) {
	cfg := &config.Client{
		Interval:  30,
		ServerURL: "http://localhost:8080",
		User:      "test",
		Debug:     false,
	}

	spy := NewSpy(cfg)

	t.Run("GET com sucesso", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"status":"ok"}`))
		}))
		defer server.Close()

		body, status, err := spy.httpGet(server.URL)
		if err != nil {
			t.Errorf("Erro inesperado: %v", err)
		}
		if status != http.StatusOK {
			t.Errorf("Status esperado 200, obtido %d", status)
		}
		if !strings.Contains(body, "ok") {
			t.Errorf("Body esperado conter 'ok', obtido %s", body)
		}
	})

	t.Run("GET com erro de servidor", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"error":"internal error"}`))
		}))
		defer server.Close()

		body, status, err := spy.httpGet(server.URL)
		if err != nil {
			t.Errorf("Erro inesperado: %v", err)
		}
		if status != http.StatusInternalServerError {
			t.Errorf("Status esperado 500, obtido %d", status)
		}
		if !strings.Contains(body, "error") {
			t.Errorf("Body esperado conter 'error', obtido %s", body)
		}
	})

	t.Run("GET com debug ativado", func(t *testing.T) {
		spy.config.Debug = true
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"debug":"true"}`))
		}))
		defer server.Close()

		body, status, err := spy.httpGet(server.URL)
		if err != nil {
			t.Errorf("Erro inesperado: %v", err)
		}
		if status != http.StatusOK {
			t.Errorf("Status esperado 200, obtido %d", status)
		}
		if !strings.Contains(body, "debug") {
			t.Errorf("Body esperado conter 'debug', obtido %s", body)
		}
		spy.config.Debug = false
	})
}

// TestSpy_httpPost testa requisições HTTP POST
func TestSpy_httpPost(t *testing.T) {
	cfg := &config.Client{
		Interval:  30,
		ServerURL: "http://localhost:8080",
		User:      "test",
		Debug:     false,
	}

	spy := NewSpy(cfg)

	t.Run("POST com sucesso", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodPost {
				t.Errorf("Método esperado POST, obtido %s", r.Method)
			}
			w.WriteHeader(http.StatusCreated)
			w.Write([]byte(`{"created":"true"}`))
		}))
		defer server.Close()

		body, status, err := spy.httpPost(server.URL, `{"test":"data"}`)
		if err != nil {
			t.Errorf("Erro inesperado: %v", err)
		}
		if status != http.StatusCreated {
			t.Errorf("Status esperado 201, obtido %d", status)
		}
		if !strings.Contains(body, "created") {
			t.Errorf("Body esperado conter 'created', obtido %s", body)
		}
	})

	t.Run("POST com debug ativado", func(t *testing.T) {
		spy.config.Debug = true
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusCreated)
			w.Write([]byte(`{"debug":"post"}`))
		}))
		defer server.Close()

		body, status, err := spy.httpPost(server.URL, `{"test":"data"}`)
		if err != nil {
			t.Errorf("Erro inesperado: %v", err)
		}
		if status != http.StatusCreated {
			t.Errorf("Status esperado 201, obtido %d", status)
		}
		if !strings.Contains(body, "debug") {
			t.Errorf("Body esperado conter 'debug', obtido %s", body)
		}
		spy.config.Debug = false
	})
}

// TestSpy_updateTargets testa atualização de targets
func TestSpy_updateTargets(t *testing.T) {
	t.Run("Atualização com sucesso", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"targets":[{"name":"test","pattern":".*test.*","limit":3600}]}`))
		}))
		defer server.Close()

		cfg := &config.Client{
			Interval:  30,
			ServerURL: server.URL,
			User:      "test",
		}

		spy := NewSpy(cfg)
		spy.updateTargets()

		if spy.targets == nil {
			t.Error("Targets não deveria ser nil")
		}
	})

	t.Run("Atualização com erro HTTP", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		}))
		defer server.Close()

		cfg := &config.Client{
			Interval:  30,
			ServerURL: server.URL,
			User:      "test",
		}

		spy := NewSpy(cfg)
		spy.updateTargets()

		if spy.targets == nil {
			t.Error("Targets não deveria ser nil mesmo com erro")
		}
	})

	t.Run("Atualização com JSON inválido", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`invalid json`))
		}))
		defer server.Close()

		cfg := &config.Client{
			Interval:  30,
			ServerURL: server.URL,
			User:      "test",
		}

		spy := NewSpy(cfg)
		spy.updateTargets()

		if spy.targets == nil {
			t.Error("Targets não deveria ser nil mesmo com JSON inválido")
		}
	})

	t.Run("Atualização com targets nil inicialmente", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"targets":[]}`))
		}))
		defer server.Close()

		cfg := &config.Client{
			Interval:  30,
			ServerURL: server.URL,
			User:      "test",
		}

		spy := NewSpy(cfg)
		spy.targets = nil
		spy.updateTargets()

		if spy.targets == nil {
			t.Error("Targets deveria ser inicializado")
		}
	})
}

// TestSpy_postMatch testa envio de matches
func TestSpy_postMatch(t *testing.T) {
	t.Run("POST match com sucesso", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusCreated)
			w.Write([]byte(`{"status":"created"}`))
		}))
		defer server.Close()

		cfg := &config.Client{
			Interval:  30,
			ServerURL: server.URL,
			User:      "test",
			Debug:     false,
		}

		spy := NewSpy(cfg)
		match := domain.NewMatch("test", "app", ".*app.*", "app.exe", 10.5)

		err := spy.postMatch(match)
		if err != nil {
			t.Errorf("Erro inesperado: %v", err)
		}
	})

	t.Run("POST match com debug", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusCreated)
			w.Write([]byte(`{"status":"created"}`))
		}))
		defer server.Close()

		cfg := &config.Client{
			Interval:  30,
			ServerURL: server.URL,
			User:      "test",
			Debug:     true,
		}

		spy := NewSpy(cfg)
		match := domain.NewMatch("test", "app", ".*app.*", "app.exe", 10.5)

		err := spy.postMatch(match)
		if err != nil {
			t.Errorf("Erro inesperado: %v", err)
		}
	})

	t.Run("POST match com erro HTTP", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		}))
		defer server.Close()

		cfg := &config.Client{
			Interval:  30,
			ServerURL: server.URL,
			User:      "test",
		}

		spy := NewSpy(cfg)
		match := domain.NewMatch("test", "app", ".*app.*", "app.exe", 10.5)

		err := spy.postMatch(match)
		if err == nil {
			t.Error("Esperado erro com status 500")
		}
	})

	t.Run("POST match nil", func(t *testing.T) {
		cfg := &config.Client{
			Interval:  30,
			ServerURL: "http://localhost:8080",
			User:      "test",
		}

		spy := NewSpy(cfg)
		err := spy.postMatch(nil)
		if err == nil {
			t.Error("Esperado erro com match nil")
		}
	})
}

// TestSpy_postCommand testa envio de comandos
func TestSpy_postCommand(t *testing.T) {
	t.Run("POST command com sucesso", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusCreated)
			w.Write([]byte(`{"status":"created"}`))
		}))
		defer server.Close()

		cfg := &config.Client{
			Interval:  30,
			ServerURL: server.URL,
			User:      "test",
		}

		spy := NewSpy(cfg)
		cmd := domain.NewCommand("test", "app", "shutdown", "executed")

		err := spy.postCommand(cmd)
		if err != nil {
			t.Errorf("Erro inesperado: %v", err)
		}
	})

	t.Run("POST command com erro HTTP", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusBadRequest)
		}))
		defer server.Close()

		cfg := &config.Client{
			Interval:  30,
			ServerURL: server.URL,
			User:      "test",
		}

		spy := NewSpy(cfg)
		cmd := domain.NewCommand("test", "app", "shutdown", "executed")

		err := spy.postCommand(cmd)
		if err == nil {
			t.Error("Esperado erro com status 400")
		}
	})

	t.Run("POST command nil", func(t *testing.T) {
		cfg := &config.Client{
			Interval:  30,
			ServerURL: "http://localhost:8080",
			User:      "test",
		}

		spy := NewSpy(cfg)
		err := spy.postCommand(nil)
		if err == nil {
			t.Error("Esperado erro com command nil")
		}
	})
}

// TestSpy_consumeBuffers testa consumo de buffers
func TestSpy_consumeBuffers(t *testing.T) {
	t.Run("Consumir buffers com sucesso", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusCreated)
			w.Write([]byte(`{"status":"created"}`))
		}))
		defer server.Close()

		cfg := &config.Client{
			Interval:  30,
			ServerURL: server.URL,
			User:      "test",
		}

		spy := NewSpy(cfg)
		match := domain.NewMatch("test", "app", ".*app.*", "app.exe", 10.5)
		cmd := domain.NewCommand("test", "app", "shutdown", "executed")

		spy.matchBuf <- match
		spy.commandBuf <- cmd

		spy.consumeBuffers()
		time.Sleep(100 * time.Millisecond)

		if len(spy.matchBuf) > 0 {
			t.Error("Match buffer deveria estar vazio")
		}
		if len(spy.commandBuf) > 0 {
			t.Error("Command buffer deveria estar vazio")
		}
	})

	t.Run("Consumir buffers com erro", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		}))
		defer server.Close()

		cfg := &config.Client{
			Interval:  30,
			ServerURL: server.URL,
			User:      "test",
		}

		spy := NewSpy(cfg)
		match := domain.NewMatch("test", "app", ".*app.*", "app.exe", 10.5)

		spy.matchBuf <- match
		spy.consumeBuffers()
		time.Sleep(100 * time.Millisecond)

		// Com erro, o match volta para o buffer (DLQ)
		if len(spy.matchBuf) == 0 {
			t.Error("Match deveria estar de volta no buffer após erro")
		}
	})

	t.Run("Consumir buffers nil", func(t *testing.T) {
		cfg := &config.Client{
			Interval:  30,
			ServerURL: "http://localhost:8080",
			User:      "test",
		}

		spy := NewSpy(cfg)
		spy.matchBuf = nil
		spy.commandBuf = nil

		spy.consumeBuffers()
		time.Sleep(50 * time.Millisecond)

		if spy.matchBuf == nil {
			t.Error("Match buffer deveria ser inicializado")
		}
		if spy.commandBuf == nil {
			t.Error("Command buffer deveria ser inicializado")
		}
	})
}

// TestSpy_kill testa terminação de processos
func TestSpy_kill(t *testing.T) {
	cfg := &config.Client{
		Interval:  30,
		ServerURL: "http://localhost:8080",
		User:      "test",
	}

	spy := NewSpy(cfg)

	t.Run("Kill com lista vazia", func(t *testing.T) {
		spy.kill("test", "pattern", []int{})
		// Não deve fazer nada, apenas não deve dar panic
	})

	t.Run("Kill com PID inválido", func(t *testing.T) {
		spy.kill("test", "pattern", []int{999999})
		// Deve tentar matar mas falhar gracefully
		time.Sleep(50 * time.Millisecond)
	})
}

// TestSpy_Stop testa parada do spy
func TestSpy_Stop(t *testing.T) {
	cfg := &config.Client{
		Interval:  30,
		ServerURL: "http://localhost:8080",
		User:      "test",
		APIPort:   0, // porta aleatória
		APIHost:   "localhost",
	}

	spy := NewSpy(cfg)
	spy.enabled = true

	spy.Stop()

	if spy.IsEnabled() {
		t.Error("Spy deveria estar desabilitado após Stop")
	}
}

// TestExecuteCommand testa execução de comandos
func TestExecuteCommand(t *testing.T) {
	t.Run("Comando inválido", func(t *testing.T) {
		_, err := executeCommand("comando_inexistente_xyz")
		if err == nil {
			t.Error("Esperado erro com comando inválido")
		}
	})
}

// TestSpy_run testa ciclo de scan
func TestSpy_run(t *testing.T) {
	t.Run("Run com targets vazios", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"targets":[]}`))
		}))
		defer server.Close()

		cfg := &config.Client{
			Interval:  30,
			ServerURL: server.URL,
			User:      "test",
		}

		spy := NewSpy(cfg)
		err := spy.run(time.Now().Add(-30 * time.Second))
		if err != nil {
			t.Errorf("Erro inesperado: %v", err)
		}
	})
}

// TestSpy_stopHttpServer testa parada do servidor HTTP
func TestSpy_stopHttpServer(t *testing.T) {
	cfg := &config.Client{
		Interval:  30,
		ServerURL: "http://localhost:8080",
		User:      "test",
	}

	spy := NewSpy(cfg)

	t.Run("Stop com servidor nil", func(t *testing.T) {
		spy.srv = nil
		spy.stopHttpServer()
		// Não deve dar panic
	})
}
