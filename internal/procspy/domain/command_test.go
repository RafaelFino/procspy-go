package domain

import (
	"strings"
	"testing"
	"time"
)

// TestNewCommand testa a criação de um novo Command
// Valida que todos os campos são inicializados corretamente
func TestNewCommand(t *testing.T) {
	// Arrange: Define parâmetros
	user := "test_user"
	name := "games"
	commandLine := "notify-send 'Tempo esgotado'"
	commandReturn := "Success"

	// Act: Cria novo command
	result := NewCommand(user, name, commandLine, commandReturn)

	// Assert: Valida campos
	if result == nil {
		t.Fatal("NewCommand retornou nil")
	}

	if result.User != user {
		t.Errorf("User = %s, esperado %s", result.User, user)
	}

	if result.Name != name {
		t.Errorf("Name = %s, esperado %s", result.Name, name)
	}

	if result.CommandLine != commandLine {
		t.Errorf("CommandLine = %s, esperado %s", result.CommandLine, commandLine)
	}

	if result.Return != commandReturn {
		t.Errorf("Return = %s, esperado %s", result.Return, commandReturn)
	}

	// Valida que Source foi definido com valor padrão
	if result.Source != "procspy" {
		t.Errorf("Source = %s, esperado 'procspy'", result.Source)
	}

	// Valida que CommandLog foi inicializado vazio
	if result.CommandLog != "" {
		t.Errorf("CommandLog deveria estar vazio, obteve %s", result.CommandLog)
	}

	// Valida que CreatedAt foi definido
	if result.CreatedAt.IsZero() {
		t.Error("CreatedAt não foi inicializado")
	}

	// Valida que CreatedAt está próximo do tempo atual
	now := time.Now()
	diff := now.Sub(result.CreatedAt)
	if diff > time.Second {
		t.Errorf("CreatedAt está muito distante do tempo atual: %v", diff)
	}
}

// TestCommand_ToJson testa a serialização de Command para JSON
// Valida que o JSON gerado é válido e formatado com indentação
func TestCommand_ToJson(t *testing.T) {
	// Arrange: Cria command
	cmd := NewCommand("user1", "games", "echo test", "output")

	// Act: Serializa para JSON
	json := cmd.ToJson()

	// Assert: Valida que JSON não está vazio
	if json == "" {
		t.Error("ToJson retornou string vazia")
	}

	// Valida que contém campos esperados
	expectedFields := []string{"user", "name", "command_line", "command_return", "source"}
	for _, field := range expectedFields {
		if !strings.Contains(json, field) {
			t.Errorf("JSON não contém campo esperado: %s", field)
		}
	}

	// Valida que está indentado (contém espaços)
	if !strings.Contains(json, "  ") {
		t.Error("JSON não está indentado")
	}

	// Valida que contém valores esperados
	if !strings.Contains(json, "user1") {
		t.Error("JSON não contém user esperado")
	}

	if !strings.Contains(json, "games") {
		t.Error("JSON não contém name esperado")
	}
}

// TestCommand_ToLog testa a serialização compacta de Command
// Valida que o log é gerado sem indentação
func TestCommand_ToLog(t *testing.T) {
	// Arrange: Cria command
	cmd := NewCommand("user1", "browsers", "killall chrome", "killed")

	// Act: Serializa para log
	log := cmd.ToLog()

	// Assert: Valida que log não está vazio
	if log == "" {
		t.Error("ToLog retornou string vazia")
	}

	// Valida que contém dados esperados
	if !strings.Contains(log, "user1") || !strings.Contains(log, "browsers") {
		t.Error("ToLog não contém dados esperados")
	}

	// Valida que não contém indentação (formato compacto)
	if strings.Contains(log, "  ") {
		t.Error("ToLog não deveria conter indentação (deve ser compacto)")
	}
}

// TestCommandFromJson testa o parsing de JSON para Command
// Valida cenários: JSON válido, JSON inválido
func TestCommandFromJson(t *testing.T) {
	tests := []struct {
		name        string
		json        string
		wantErr     bool
		wantUser    string
		wantName    string
		wantCmdLine string
	}{
		{
			name: "JSON válido completo",
			json: `{
				"user": "test_user",
				"name": "games",
				"command_line": "notify-send test",
				"command_return": "success",
				"source": "limit"
			}`,
			wantErr:     false,
			wantUser:    "test_user",
			wantName:    "games",
			wantCmdLine: "notify-send test",
		},
		{
			name: "JSON válido mínimo",
			json: `{
				"user": "user1",
				"name": "browsers"
			}`,
			wantErr:     false,
			wantUser:    "user1",
			wantName:    "browsers",
			wantCmdLine: "",
		},
		{
			name:    "JSON inválido",
			json:    `{invalid json}`,
			wantErr: true,
		},
		{
			name:        "JSON vazio",
			json:        `{}`,
			wantErr:     false,
			wantUser:    "",
			wantName:    "",
			wantCmdLine: "",
		},
		{
			name:    "String vazia",
			json:    ``,
			wantErr: true,
		},
		{
			name: "JSON com campos extras",
			json: `{
				"user": "user2",
				"name": "media",
				"command_line": "pkill vlc",
				"extra_field": "ignored"
			}`,
			wantErr:     false,
			wantUser:    "user2",
			wantName:    "media",
			wantCmdLine: "pkill vlc",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act: Faz parsing do JSON
			result, err := CommandFromJson(tt.json)

			// Assert: Valida erro
			if (err != nil) != tt.wantErr {
				t.Errorf("CommandFromJson() erro = %v, esperava erro = %v", err, tt.wantErr)
				return
			}

			// Se não esperava erro, valida resultado
			if !tt.wantErr {
				if result == nil {
					t.Fatal("CommandFromJson retornou nil sem erro")
				}

				if result.User != tt.wantUser {
					t.Errorf("User = %s, esperado %s", result.User, tt.wantUser)
				}

				if result.Name != tt.wantName {
					t.Errorf("Name = %s, esperado %s", result.Name, tt.wantName)
				}

				if result.CommandLine != tt.wantCmdLine {
					t.Errorf("CommandLine = %s, esperado %s", result.CommandLine, tt.wantCmdLine)
				}
			}
		})
	}
}

// TestCommand_WithEmptyStrings testa Command com strings vazias
// Valida que campos vazios são tratados corretamente
func TestCommand_WithEmptyStrings(t *testing.T) {
	// Arrange & Act: Cria command com strings vazias
	cmd := NewCommand("", "", "", "")

	// Assert: Valida que command foi criado
	if cmd == nil {
		t.Fatal("NewCommand retornou nil")
	}

	// Valida que campos vazios são aceitos
	if cmd.User != "" {
		t.Errorf("User deveria estar vazio, obteve %s", cmd.User)
	}

	if cmd.Name != "" {
		t.Errorf("Name deveria estar vazio, obteve %s", cmd.Name)
	}

	// Valida que serialização funciona com strings vazias
	json := cmd.ToJson()
	if json == "" {
		t.Error("ToJson retornou string vazia")
	}

	log := cmd.ToLog()
	if log == "" {
		t.Error("ToLog retornou string vazia")
	}
}

// TestCommand_WithLongStrings testa Command com strings longas
// Valida que strings longas são tratadas corretamente
func TestCommand_WithLongStrings(t *testing.T) {
	// Arrange: Cria strings longas
	longUser := strings.Repeat("a", 1000)
	longName := strings.Repeat("b", 1000)
	longCommand := strings.Repeat("c", 5000)
	longReturn := strings.Repeat("d", 10000)

	// Act: Cria command com strings longas
	cmd := NewCommand(longUser, longName, longCommand, longReturn)

	// Assert: Valida que command foi criado
	if cmd == nil {
		t.Fatal("NewCommand retornou nil")
	}

	// Valida que strings foram armazenadas
	if len(cmd.User) != 1000 {
		t.Errorf("User length = %d, esperado 1000", len(cmd.User))
	}

	if len(cmd.CommandLine) != 5000 {
		t.Errorf("CommandLine length = %d, esperado 5000", len(cmd.CommandLine))
	}

	// Valida que serialização funciona
	json := cmd.ToJson()
	if json == "" {
		t.Error("ToJson retornou string vazia com strings longas")
	}
}

// TestCommand_Serialization_RoundTrip testa serialização e desserialização
// Valida que dados são preservados após round-trip JSON
func TestCommand_Serialization_RoundTrip(t *testing.T) {
	// Arrange: Cria command original
	original := NewCommand("user1", "games", "notify-send 'Test'", "Success")
	original.Source = "Limit"
	original.CommandLog = "Executed successfully"

	// Act: Serializa e desserializa
	json := original.ToJson()
	restored, err := CommandFromJson(json)

	// Assert: Valida que não houve erro
	if err != nil {
		t.Fatalf("Erro ao desserializar: %v", err)
	}

	// Valida que dados foram preservados
	if restored.User != original.User {
		t.Errorf("User não preservado: %s != %s", restored.User, original.User)
	}

	if restored.Name != original.Name {
		t.Errorf("Name não preservado: %s != %s", restored.Name, original.Name)
	}

	if restored.CommandLine != original.CommandLine {
		t.Errorf("CommandLine não preservado: %s != %s", restored.CommandLine, original.CommandLine)
	}

	if restored.Return != original.Return {
		t.Errorf("Return não preservado: %s != %s", restored.Return, original.Return)
	}

	if restored.Source != original.Source {
		t.Errorf("Source não preservado: %s != %s", restored.Source, original.Source)
	}

	if restored.CommandLog != original.CommandLog {
		t.Errorf("CommandLog não preservado: %s != %s", restored.CommandLog, original.CommandLog)
	}
}

// TestCommand_WithSpecialCharacters testa Command com caracteres especiais
// Valida que caracteres especiais são tratados corretamente no JSON
func TestCommand_WithSpecialCharacters(t *testing.T) {
	// Arrange: Cria command com caracteres especiais
	cmd := NewCommand(
		"user@domain.com",
		"test-name_123",
		"echo \"Hello World\" && ls -la",
		"Output:\n\tLine 1\n\tLine 2",
	)

	// Act: Serializa e desserializa
	json := cmd.ToJson()
	restored, err := CommandFromJson(json)

	// Assert: Valida que não houve erro
	if err != nil {
		t.Fatalf("Erro ao desserializar: %v", err)
	}

	// Valida que caracteres especiais foram preservados
	if restored.User != cmd.User {
		t.Errorf("User com @ não preservado: %s != %s", restored.User, cmd.User)
	}

	if restored.CommandLine != cmd.CommandLine {
		t.Errorf("CommandLine com aspas não preservado: %s != %s", restored.CommandLine, cmd.CommandLine)
	}

	if restored.Return != cmd.Return {
		t.Errorf("Return com newlines não preservado: %s != %s", restored.Return, cmd.Return)
	}
}

// TestCommand_DefaultSource testa que Source tem valor padrão
// Valida que NewCommand define Source como "procspy"
func TestCommand_DefaultSource(t *testing.T) {
	// Arrange & Act: Cria command
	cmd := NewCommand("user1", "test", "cmd", "output")

	// Assert: Valida source padrão
	if cmd.Source != "procspy" {
		t.Errorf("Source padrão = %s, esperado 'procspy'", cmd.Source)
	}
}

// TestCommand_ModifySource testa modificação do campo Source
// Valida que Source pode ser alterado após criação
func TestCommand_ModifySource(t *testing.T) {
	// Arrange: Cria command
	cmd := NewCommand("user1", "test", "cmd", "output")

	// Act: Modifica source
	cmd.Source = "Limit"

	// Assert: Valida modificação
	if cmd.Source != "Limit" {
		t.Errorf("Source após modificação = %s, esperado 'Limit'", cmd.Source)
	}

	// Valida que serialização reflete a mudança
	json := cmd.ToJson()
	if !strings.Contains(json, "Limit") {
		t.Error("JSON não contém Source modificado")
	}
}

// TestCommand_CreatedAtPrecision testa precisão do timestamp CreatedAt
// Valida que CreatedAt é definido com precisão adequada
func TestCommand_CreatedAtPrecision(t *testing.T) {
	// Arrange: Captura tempo antes
	before := time.Now()

	// Act: Cria command
	cmd := NewCommand("user1", "test", "cmd", "output")

	// Arrange: Captura tempo depois
	after := time.Now()

	// Assert: Valida que CreatedAt está entre before e after
	if cmd.CreatedAt.Before(before) {
		t.Error("CreatedAt está antes do tempo de criação")
	}

	if cmd.CreatedAt.After(after) {
		t.Error("CreatedAt está depois do tempo de criação")
	}

	// Valida que diferença é pequena (< 10ms)
	diff := cmd.CreatedAt.Sub(before)
	if diff > 10*time.Millisecond {
		t.Errorf("CreatedAt tem diferença muito grande: %v", diff)
	}
}

// TestCommand_MultipleInstances testa criação de múltiplos commands
// Valida que cada instância é independente
func TestCommand_MultipleInstances(t *testing.T) {
	// Arrange & Act: Cria múltiplos commands
	cmd1 := NewCommand("user1", "games", "cmd1", "out1")
	cmd2 := NewCommand("user2", "browsers", "cmd2", "out2")
	cmd3 := NewCommand("user3", "media", "cmd3", "out3")

	// Assert: Valida que são instâncias diferentes
	if cmd1 == cmd2 || cmd1 == cmd3 || cmd2 == cmd3 {
		t.Error("Commands deveriam ser instâncias diferentes")
	}

	// Valida que cada um tem seus próprios valores
	if cmd1.User == cmd2.User {
		t.Error("cmd1 e cmd2 não deveriam ter o mesmo User")
	}

	if cmd1.Name == cmd2.Name {
		t.Error("cmd1 e cmd2 não deveriam ter o mesmo Name")
	}

	// Valida que modificar um não afeta os outros
	cmd1.Source = "Modified"
	if cmd2.Source == "Modified" || cmd3.Source == "Modified" {
		t.Error("Modificar cmd1 afetou outros commands")
	}
}

// TestCommand_ToJson_ContainsAllFields testa que ToJson inclui todos os campos
// Valida que nenhum campo é omitido na serialização
func TestCommand_ToJson_ContainsAllFields(t *testing.T) {
	// Arrange: Cria command com todos os campos preenchidos
	cmd := &Command{
		User:        "user1",
		Name:        "games",
		CommandLine: "notify-send test",
		Return:      "Success",
		Source:      "Limit",
		CommandLog:  "Executed at 10:00",
		CreatedAt:   time.Now(),
	}

	// Act: Serializa
	json := cmd.ToJson()

	// Assert: Valida que todos os campos estão presentes
	requiredFields := []string{
		"user", "name", "command_line", "command_return",
		"source", "command_log", "created_at",
	}

	for _, field := range requiredFields {
		if !strings.Contains(json, field) {
			t.Errorf("JSON não contém campo obrigatório: %s", field)
		}
	}
}

// TestCommand_ToLog_IsCompact testa que ToLog é compacto
// Valida que não há indentação desnecessária
func TestCommand_ToLog_IsCompact(t *testing.T) {
	// Arrange: Cria command
	cmd := NewCommand("user1", "test", "echo test", "output")

	// Act: Serializa para log
	log := cmd.ToLog()

	// Assert: Valida formato compacto (sem newlines ou tabs extras)
	lines := strings.Split(log, "\n")
	if len(lines) > 1 {
		// JSON compacto pode ter uma linha
		t.Error("ToLog deveria ser uma única linha (compacto)")
	}

	// Valida que contém dados
	if len(log) == 0 {
		t.Error("ToLog retornou string vazia")
	}
}

// TestCommand_WithUnicodeCharacters testa Command com caracteres Unicode
// Valida que caracteres não-ASCII são tratados corretamente
func TestCommand_WithUnicodeCharacters(t *testing.T) {
	// Arrange: Cria command com Unicode
	cmd := NewCommand(
		"usuário",
		"jogos",
		"notify-send 'Tempo esgotado! 🎮'",
		"Sucesso ✓",
	)

	// Act: Serializa e desserializa
	json := cmd.ToJson()
	restored, err := CommandFromJson(json)

	// Assert: Valida que não houve erro
	if err != nil {
		t.Fatalf("Erro ao desserializar: %v", err)
	}

	// Valida que caracteres Unicode foram preservados
	if restored.User != cmd.User {
		t.Errorf("User com Unicode não preservado: %s != %s", restored.User, cmd.User)
	}

	if restored.Name != cmd.Name {
		t.Errorf("Name com Unicode não preservado: %s != %s", restored.Name, cmd.Name)
	}

	if !strings.Contains(restored.CommandLine, "🎮") {
		t.Error("Emoji não foi preservado no CommandLine")
	}

	if !strings.Contains(restored.Return, "✓") {
		t.Error("Símbolo Unicode não foi preservado no Return")
	}
}

// TestCommand_WithEmptyCommandLog testa Command com CommandLog vazio
// Valida que campo vazio é tratado corretamente
func TestCommand_WithEmptyCommandLog(t *testing.T) {
	// Arrange & Act: Cria command (CommandLog é vazio por padrão)
	cmd := NewCommand("user1", "test", "cmd", "output")

	// Assert: Valida que CommandLog está vazio
	if cmd.CommandLog != "" {
		t.Errorf("CommandLog deveria estar vazio, obteve: %s", cmd.CommandLog)
	}

	// Valida que serialização funciona
	json := cmd.ToJson()
	if json == "" {
		t.Error("ToJson retornou string vazia")
	}

	// Valida que campo está presente no JSON
	if !strings.Contains(json, "command_log") {
		t.Error("JSON não contém campo command_log")
	}
}

// TestCommand_ModifyCommandLog testa modificação do CommandLog
// Valida que CommandLog pode ser alterado após criação
func TestCommand_ModifyCommandLog(t *testing.T) {
	// Arrange: Cria command
	cmd := NewCommand("user1", "test", "cmd", "output")

	// Act: Modifica CommandLog
	cmd.CommandLog = "Log entry 1\nLog entry 2"

	// Assert: Valida modificação
	if !strings.Contains(cmd.CommandLog, "Log entry 1") {
		t.Error("CommandLog não foi modificado corretamente")
	}

	// Valida que serialização reflete a mudança
	json := cmd.ToJson()
	if !strings.Contains(json, "Log entry 1") {
		t.Error("JSON não contém CommandLog modificado")
	}
}

// TestCommandFromJson_WithMissingFields testa parsing com campos faltando
// Valida que campos opcionais podem estar ausentes
func TestCommandFromJson_WithMissingFields(t *testing.T) {
	// Arrange: JSON com apenas campos obrigatórios
	json := `{
		"user": "user1",
		"name": "test"
	}`

	// Act: Faz parsing
	result, err := CommandFromJson(json)

	// Assert: Valida que não houve erro
	if err != nil {
		t.Fatalf("Erro ao fazer parsing: %v", err)
	}

	// Valida que campos obrigatórios foram parseados
	if result.User != "user1" {
		t.Errorf("User = %s, esperado user1", result.User)
	}

	if result.Name != "test" {
		t.Errorf("Name = %s, esperado test", result.Name)
	}

	// Valida que campos opcionais estão vazios
	if result.CommandLine != "" {
		t.Error("CommandLine deveria estar vazio")
	}

	if result.Return != "" {
		t.Error("Return deveria estar vazio")
	}
}

// TestCommand_WithVeryLongCommandLine testa Command com linha de comando muito longa
// Valida que comandos longos são tratados corretamente
func TestCommand_WithVeryLongCommandLine(t *testing.T) {
	// Arrange: Cria comando muito longo
	longCmd := "bash -c '" + strings.Repeat("echo test; ", 1000) + "'"

	// Act: Cria command
	cmd := NewCommand("user1", "test", longCmd, "output")

	// Assert: Valida que comando foi armazenado
	if len(cmd.CommandLine) < 10000 {
		t.Error("CommandLine longo não foi armazenado corretamente")
	}

	// Valida que serialização funciona
	json := cmd.ToJson()
	if json == "" {
		t.Error("ToJson retornou string vazia com comando longo")
	}

	// Valida que pode ser desserializado
	restored, err := CommandFromJson(json)
	if err != nil {
		t.Fatalf("Erro ao desserializar comando longo: %v", err)
	}

	if restored.CommandLine != cmd.CommandLine {
		t.Error("CommandLine longo não foi preservado após round-trip")
	}
}

// TestCommand_CreatedAtBetweenCalls testa que CreatedAt é diferente entre chamadas
// Valida que cada command tem seu próprio timestamp
func TestCommand_CreatedAtBetweenCalls(t *testing.T) {
	// Arrange & Act: Cria dois commands com pequeno delay
	cmd1 := NewCommand("user1", "test", "cmd", "output")
	time.Sleep(1 * time.Millisecond)
	cmd2 := NewCommand("user1", "test", "cmd", "output")

	// Assert: Valida que timestamps são diferentes
	if !cmd2.CreatedAt.After(cmd1.CreatedAt) {
		t.Error("cmd2.CreatedAt deveria ser posterior a cmd1.CreatedAt")
	}

	// Valida que diferença é pequena mas mensurável
	diff := cmd2.CreatedAt.Sub(cmd1.CreatedAt)
	if diff < 1*time.Millisecond {
		t.Error("Diferença entre timestamps é muito pequena")
	}

	if diff > 100*time.Millisecond {
		t.Error("Diferença entre timestamps é muito grande")
	}
}

// TestCommand_SourceValues testa diferentes valores de Source
// Valida que Source pode ter diferentes valores
func TestCommand_SourceValues(t *testing.T) {
	sources := []string{"procspy", "Limit", "Warning", "Check", "Manual"}

	for _, source := range sources {
		// Arrange & Act: Cria command e modifica source
		cmd := NewCommand("user1", "test", "cmd", "output")
		cmd.Source = source

		// Assert: Valida que source foi definido
		if cmd.Source != source {
			t.Errorf("Source = %s, esperado %s", cmd.Source, source)
		}

		// Valida que serialização preserva o valor
		json := cmd.ToJson()
		if !strings.Contains(json, source) {
			t.Errorf("JSON não contém Source: %s", source)
		}
	}
}
