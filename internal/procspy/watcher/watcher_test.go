package watcher

import (
	"net/http"
	"net/http/httptest"
	"procspy/internal/procspy/config"
	"strings"
	"testing"
)

// TestNewWatcher testa criação de um novo Watcher
func TestNewWatcher(t *testing.T) {
	cfg := &config.Watcher{
		Interval:   10,
		LogPath:    "logs",
		ProcspyURL: "http://localhost:8888",
		StartCmd:   "systemctl restart procspy-client",
	}

	watcher := NewWatcher(cfg)
	if watcher == nil {
		t.Fatal("NewWatcher retornou nil")
	}

	if watcher.config != cfg {
		t.Error("Config não foi atribuída corretamente")
	}

	if watcher.enabled {
		t.Error("Watcher deveria iniciar desabilitado")
	}
}

// TestWatcher_httpGet testa requisições HTTP GET
func TestWatcher_httpGet(t *testing.T) {
	cfg := &config.Watcher{
		Interval:   10,
		ProcspyURL: "http://localhost:8888",
	}

	watcher := NewWatcher(cfg)

	t.Run("GET com sucesso", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"status":"ok"}`))
		}))
		defer server.Close()

		body, status, err := watcher.httpGet(server.URL)
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

		body, status, err := watcher.httpGet(server.URL)
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

	t.Run("GET com URL inválida", func(t *testing.T) {
		_, status, err := watcher.httpGet("http://invalid-url-that-does-not-exist:99999")
		if err == nil {
			t.Error("Esperado erro com URL inválida")
		}
		if status != http.StatusInternalServerError {
			t.Errorf("Status esperado 500, obtido %d", status)
		}
	})
}

// TestWatcher_check testa verificação de health check
func TestWatcher_check(t *testing.T) {
	t.Run("Check com procspy up", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"status":"healthy"}`))
		}))
		defer server.Close()

		cfg := &config.Watcher{
			Interval:   10,
			ProcspyURL: server.URL,
			StartCmd:   "",
		}

		watcher := NewWatcher(cfg)
		watcher.check()
		// Não deve executar comando de start
	})

	t.Run("Check com procspy down e sem comando", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		}))
		defer server.Close()

		cfg := &config.Watcher{
			Interval:   10,
			ProcspyURL: server.URL,
			StartCmd:   "",
		}

		watcher := NewWatcher(cfg)
		watcher.check()
		// Não deve executar comando pois StartCmd está vazio
	})

	t.Run("Check com procspy down e com comando", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusServiceUnavailable)
		}))
		defer server.Close()

		cfg := &config.Watcher{
			Interval:   10,
			ProcspyURL: server.URL,
			StartCmd:   "comando_invalido_xyz",
		}

		watcher := NewWatcher(cfg)
		watcher.check()
		// Deve tentar executar comando mas falhar
	})

	t.Run("Check com erro de conexão", func(t *testing.T) {
		cfg := &config.Watcher{
			Interval:   10,
			ProcspyURL: "http://invalid-url:99999",
			StartCmd:   "",
		}

		watcher := NewWatcher(cfg)
		watcher.check()
		// Deve tratar erro de conexão
	})
}

// TestWatcher_Stop testa parada do watcher
func TestWatcher_Stop(t *testing.T) {
	cfg := &config.Watcher{
		Interval:   10,
		ProcspyURL: "http://localhost:8888",
	}

	watcher := NewWatcher(cfg)
	watcher.enabled = true

	watcher.Stop()

	if watcher.enabled {
		t.Error("Watcher deveria estar desabilitado após Stop")
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
