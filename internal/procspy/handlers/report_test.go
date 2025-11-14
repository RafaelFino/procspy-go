package handlers

import (
	"procspy/internal/procspy/config"
	"procspy/internal/procspy/service"
	"procspy/internal/procspy/storage"
	"testing"
	"time"
)

// TestNewReport testa criação de handler de report
func TestNewReport(t *testing.T) {
	cfg := &config.Server{UserTarges: map[string]string{"user1": "url"}}
	targetService := service.NewTarget(cfg)
	usersService := service.NewUsers(cfg)
	
	conn := storage.NewDbConnection(":memory:")
	defer conn.Close()
	matchService := service.NewMatch(conn)
	commandService := service.NewCommand(conn)

	handler := NewReport(targetService, usersService, matchService, commandService)
	if handler == nil {
		t.Fatal("NewReport retornou nil")
	}
}

// TestFormatInterval testa formatação de intervalo de tempo
func TestFormatInterval(t *testing.T) {
	tests := []struct {
		name     string
		seconds  float64
		scale    time.Duration
		expected string
	}{
		{"1 segundo", 1.0, time.Second, "1s"},
		{"60 segundos", 60.0, time.Second, "1m0s"},
		{"3600 segundos", 3600.0, time.Second, "1h0m0s"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FormatInterval(tt.seconds, tt.scale)
			if result != tt.expected {
				t.Errorf("FormatInterval() = %s, esperado %s", result, tt.expected)
			}
		})
	}
}

// TestReport_GetReport_InvalidUser testa busca com usuário inválido
func TestReport_GetReport_InvalidUser(t *testing.T) {
	cfg := &config.Server{UserTarges: map[string]string{"user1": "url"}}
	targetService := service.NewTarget(cfg)
	usersService := service.NewUsers(cfg)
	
	conn := storage.NewDbConnection(":memory:")
	defer conn.Close()
	matchService := service.NewMatch(conn)
	commandService := service.NewCommand(conn)
	handler := NewReport(targetService, usersService, matchService, commandService)

	gin := setupTestRouter()
	gin.GET("/report/:user", handler.GetReport)

	req := makeTestRequest("GET", "/report/invalid", "")
	w := executeRequest(gin, req)

	if w.Code != 401 {
		t.Errorf("Status = %d, esperado 401", w.Code)
	}
}
