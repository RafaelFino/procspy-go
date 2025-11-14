package domain

import (
	"strings"
	"testing"
	"time"
)

// TestNewMatch testa a criação de um novo Match
// Valida que todos os campos são inicializados corretamente
func TestNewMatch(t *testing.T) {
	// Arrange: Define parâmetros
	user := "test_user"
	name := "games"
	pattern := "steam|roblox"
	match := "steam.exe"
	elapsed := 5.0

	// Act: Cria novo match
	result := NewMatch(user, name, pattern, match, elapsed)

	// Assert: Valida campos
	if result == nil {
		t.Fatal("NewMatch retornou nil")
	}

	if result.User != user {
		t.Errorf("User = %s, esperado %s", result.User, user)
	}

	if result.Name != name {
		t.Errorf("Name = %s, esperado %s", result.Name, name)
	}

	if result.Pattern != pattern {
		t.Errorf("Pattern = %s, esperado %s", result.Pattern, pattern)
	}

	if result.Match != match {
		t.Errorf("Match = %s, esperado %s", result.Match, match)
	}

	if result.Elapsed != elapsed {
		t.Errorf("Elapsed = %.2f, esperado %.2f", result.Elapsed, elapsed)
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

// TestMatch_ToJson testa a serialização de Match para JSON
// Valida que o JSON gerado é válido e formatado com indentação
func TestMatch_ToJson(t *testing.T) {
	// Arrange: Cria match
	match := NewMatch("user1", "games", "steam", "steam.exe", 10.5)

	// Act: Serializa para JSON
	json := match.ToJson()

	// Assert: Valida que JSON não está vazio
	if json == "" {
		t.Error("ToJson retornou string vazia")
	}

	// Valida que contém campos esperados
	expectedFields := []string{"user", "name", "pattern", "match", "elapsed"}
	for _, field := range expectedFields {
		if !strings.Contains(json, field) {
			t.Errorf("JSON não contém campo esperado: %s", field)
		}
	}

	// Valida que está indentado (contém espaços)
	if !strings.Contains(json, "  ") {
		t.Error("JSON não está indentado")
	}
}

// TestMatch_ToLog testa a serialização compacta de Match
// Valida que o log é gerado sem indentação
func TestMatch_ToLog(t *testing.T) {
	// Arrange: Cria match
	match := NewMatch("user1", "games", "steam", "steam.exe", 10.5)

	// Act: Serializa para log
	log := match.ToLog()

	// Assert: Valida que log não está vazio
	if log == "" {
		t.Error("ToLog retornou string vazia")
	}

	// Valida que contém campos esperados
	if !strings.Contains(log, "user1") || !strings.Contains(log, "games") {
		t.Error("ToLog não contém dados esperados")
	}

	// Valida que não contém indentação (formato compacto)
	if strings.Contains(log, "  ") {
		t.Error("ToLog não deveria conter indentação (deve ser compacto)")
	}
}

// TestMatchFromJson testa o parsing de JSON para Match
// Valida cenários: JSON válido, JSON inválido
func TestMatchFromJson(t *testing.T) {
	tests := []struct {
		name    string
		json    string
		wantErr bool
		wantUser string
		wantName string
	}{
		{
			name: "JSON válido completo",
			json: `{
				"user": "test_user",
				"name": "games",
				"pattern": "steam",
				"match": "steam.exe",
				"elapsed": 15.5
			}`,
			wantErr: false,
			wantUser: "test_user",
			wantName: "games",
		},
		{
			name: "JSON válido mínimo",
			json: `{
				"user": "user1",
				"name": "browsers"
			}`,
			wantErr: false,
			wantUser: "user1",
			wantName: "browsers",
		},
		{
			name:    "JSON inválido",
			json:    `{invalid json}`,
			wantErr: true,
		},
		{
			name:    "JSON vazio",
			json:    `{}`,
			wantErr: false,
			wantUser: "",
			wantName: "",
		},
		{
			name:    "String vazia",
			json:    ``,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act: Faz parsing do JSON
			result, err := MatchFromJson(tt.json)

			// Assert: Valida erro
			if (err != nil) != tt.wantErr {
				t.Errorf("MatchFromJson() erro = %v, esperava erro = %v", err, tt.wantErr)
				return
			}

			// Se não esperava erro, valida resultado
			if !tt.wantErr {
				if result == nil {
					t.Fatal("MatchFromJson retornou nil sem erro")
				}

				if result.User != tt.wantUser {
					t.Errorf("User = %s, esperado %s", result.User, tt.wantUser)
				}

				if result.Name != tt.wantName {
					t.Errorf("Name = %s, esperado %s", result.Name, tt.wantName)
				}
			}
		})
	}
}

// TestMatchListFromJson testa o parsing de JSON para MatchList
// Valida cenários: JSON válido com matches, JSON vazio, JSON inválido
func TestMatchListFromJson(t *testing.T) {
	tests := []struct {
		name    string
		json    string
		wantErr bool
		wantLen int
	}{
		{
			name: "JSON válido com matches",
			json: `{
				"matches": {
					"games": 3600.0,
					"browsers": 1800.0
				}
			}`,
			wantErr: false,
			wantLen: 2,
		},
		{
			name: "JSON vazio",
			json: `{
				"matches": {}
			}`,
			wantErr: false,
			wantLen: 0,
		},
		{
			name:    "JSON inválido",
			json:    `{invalid}`,
			wantErr: true,
			wantLen: 0,
		},
		{
			name: "JSON com um match",
			json: `{
				"matches": {
					"games": 7200.0
				}
			}`,
			wantErr: false,
			wantLen: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act: Faz parsing do JSON
			result, err := MatchListFromJson(tt.json)

			// Assert: Valida erro
			if (err != nil) != tt.wantErr {
				t.Errorf("MatchListFromJson() erro = %v, esperava erro = %v", err, tt.wantErr)
				return
			}

			// Se não esperava erro, valida resultado
			if !tt.wantErr {
				if result == nil {
					t.Fatal("MatchListFromJson retornou nil sem erro")
				}

				if result.Matches == nil {
					t.Fatal("Matches map não foi inicializado")
				}

				if len(result.Matches) != tt.wantLen {
					t.Errorf("Esperado %d matches, obteve %d", tt.wantLen, len(result.Matches))
				}
			}
		})
	}
}

// TestMatchInfo_ToJson testa a serialização de MatchInfo para JSON
// Valida que o JSON é gerado corretamente com indentação
func TestMatchInfo_ToJson(t *testing.T) {
	// Arrange: Cria MatchInfo
	info := &MatchInfo{
		Elapsed:    3600.0,
		FirstMatch: "2024-01-01 10:00:00",
		LastMatch:  "2024-01-01 11:00:00",
		Ocurrences: 10,
	}

	// Act: Serializa para JSON
	json := info.ToJson()

	// Assert: Valida que JSON não está vazio
	if json == "" {
		t.Error("ToJson retornou string vazia")
	}

	// Valida que contém campos esperados
	expectedFields := []string{"elapsed", "first_match", "last_match", "ocurrences"}
	for _, field := range expectedFields {
		if !strings.Contains(json, field) {
			t.Errorf("JSON não contém campo esperado: %s", field)
		}
	}

	// Valida que está indentado
	if !strings.Contains(json, "  ") {
		t.Error("JSON não está indentado")
	}
}

// TestMatchInfo_ToLog testa a serialização compacta de MatchInfo
// Valida que o log é gerado sem indentação
func TestMatchInfo_ToLog(t *testing.T) {
	// Arrange: Cria MatchInfo
	info := &MatchInfo{
		Elapsed:    1800.0,
		FirstMatch: "2024-01-01 10:00:00",
		LastMatch:  "2024-01-01 10:30:00",
		Ocurrences: 5,
	}

	// Act: Serializa para log
	log := info.ToLog()

	// Assert: Valida que log não está vazio
	if log == "" {
		t.Error("ToLog retornou string vazia")
	}

	// Valida que não contém indentação (formato compacto)
	if strings.Contains(log, "  ") {
		t.Error("ToLog não deveria conter indentação (deve ser compacto)")
	}

	// Valida que contém valores esperados
	if !strings.Contains(log, "1800") {
		t.Error("ToLog não contém elapsed esperado")
	}
}

// TestMatch_WithZeroElapsed testa Match com elapsed zero
// Valida que valores zero são tratados corretamente
func TestMatch_WithZeroElapsed(t *testing.T) {
	// Arrange & Act: Cria match com elapsed zero
	match := NewMatch("user1", "test", "pattern", "match", 0.0)

	// Assert: Valida que elapsed é zero
	if match.Elapsed != 0.0 {
		t.Errorf("Elapsed = %.2f, esperado 0.00", match.Elapsed)
	}

	// Valida que serialização funciona com zero
	json := match.ToJson()
	if !strings.Contains(json, "elapsed") {
		t.Error("JSON não contém campo elapsed")
	}
}

// TestMatch_WithLargeElapsed testa Match com valores grandes de elapsed
// Valida que valores grandes são tratados corretamente
func TestMatch_WithLargeElapsed(t *testing.T) {
	// Arrange & Act: Cria match com elapsed grande
	largeElapsed := 86400.0 // 24 horas em segundos
	match := NewMatch("user1", "test", "pattern", "match", largeElapsed)

	// Assert: Valida que elapsed foi armazenado corretamente
	if match.Elapsed != largeElapsed {
		t.Errorf("Elapsed = %.2f, esperado %.2f", match.Elapsed, largeElapsed)
	}

	// Valida que serialização funciona
	json := match.ToJson()
	if !strings.Contains(json, "86400") {
		t.Error("JSON não contém valor correto de elapsed")
	}
}

// TestMatchInfo_WithEmptyStrings testa MatchInfo com strings vazias
// Valida que campos vazios são tratados corretamente
func TestMatchInfo_WithEmptyStrings(t *testing.T) {
	// Arrange: Cria MatchInfo com strings vazias
	info := &MatchInfo{
		Elapsed:    100.0,
		FirstMatch: "",
		LastMatch:  "",
		Ocurrences: 0,
	}

	// Act: Serializa
	json := info.ToJson()
	log := info.ToLog()

	// Assert: Valida que serialização funciona
	if json == "" {
		t.Error("ToJson retornou string vazia")
	}

	if log == "" {
		t.Error("ToLog retornou string vazia")
	}

	// Valida que contém elapsed
	if !strings.Contains(json, "100") {
		t.Error("JSON não contém elapsed")
	}
}

// TestMatch_Serialization_RoundTrip testa serialização e desserialização
// Valida que dados são preservados após round-trip JSON
func TestMatch_Serialization_RoundTrip(t *testing.T) {
	// Arrange: Cria match original
	original := NewMatch("user1", "games", "steam|roblox", "steam.exe", 25.5)

	// Act: Serializa e desserializa
	json := original.ToJson()
	restored, err := MatchFromJson(json)

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

	if restored.Pattern != original.Pattern {
		t.Errorf("Pattern não preservado: %s != %s", restored.Pattern, original.Pattern)
	}

	if restored.Match != original.Match {
		t.Errorf("Match não preservado: %s != %s", restored.Match, original.Match)
	}

	if restored.Elapsed != original.Elapsed {
		t.Errorf("Elapsed não preservado: %.2f != %.2f", restored.Elapsed, original.Elapsed)
	}
}

// TestMatchList_WithMultipleMatches testa MatchList com múltiplos matches
// Valida que todos os matches são armazenados corretamente
func TestMatchList_WithMultipleMatches(t *testing.T) {
	// Arrange: Cria JSON com múltiplos matches
	json := `{
		"matches": {
			"games": 3600.0,
			"browsers": 1800.0,
			"media": 900.0
		}
	}`

	// Act: Faz parsing
	result, err := MatchListFromJson(json)

	// Assert: Valida resultado
	if err != nil {
		t.Fatalf("Erro ao fazer parsing: %v", err)
	}

	if len(result.Matches) != 3 {
		t.Errorf("Esperado 3 matches, obteve %d", len(result.Matches))
	}

	// Valida valores específicos
	if result.Matches["games"] != 3600.0 {
		t.Errorf("games = %.2f, esperado 3600.00", result.Matches["games"])
	}

	if result.Matches["browsers"] != 1800.0 {
		t.Errorf("browsers = %.2f, esperado 1800.00", result.Matches["browsers"])
	}

	if result.Matches["media"] != 900.0 {
		t.Errorf("media = %.2f, esperado 900.00", result.Matches["media"])
	}
}

// TestMatch_WithNegativeElapsed testa Match com elapsed negativo
// Valida que valores negativos são aceitos (pode representar ajustes)
func TestMatch_WithNegativeElapsed(t *testing.T) {
	// Arrange & Act: Cria match com elapsed negativo
	match := NewMatch("user1", "test", "pattern", "match", -10.0)

	// Assert: Valida que elapsed negativo foi armazenado
	if match.Elapsed != -10.0 {
		t.Errorf("Elapsed = %.2f, esperado -10.00", match.Elapsed)
	}

	// Valida que serialização funciona
	json := match.ToJson()
	if !strings.Contains(json, "-10") {
		t.Error("JSON não contém elapsed negativo")
	}
}

// TestMatch_CreatedAtIsRecent testa que CreatedAt é definido corretamente
// Valida que timestamp está próximo do momento de criação
func TestMatch_CreatedAtIsRecent(t *testing.T) {
	// Arrange: Captura tempo antes
	before := time.Now()

	// Act: Cria match
	match := NewMatch("user1", "test", "pattern", "match", 5.0)

	// Arrange: Captura tempo depois
	after := time.Now()

	// Assert: Valida que CreatedAt está no intervalo
	if match.CreatedAt.Before(before) || match.CreatedAt.After(after) {
		t.Error("CreatedAt não está no intervalo esperado")
	}
}

// TestMatchFromJson_WithCreatedAt testa parsing de JSON com created_at
// Valida que timestamp é preservado durante desserialização
func TestMatchFromJson_WithCreatedAt(t *testing.T) {
	// Arrange: JSON com created_at
	json := `{
		"user": "user1",
		"name": "games",
		"created_at": "2024-01-01T10:00:00Z"
	}`

	// Act: Faz parsing
	result, err := MatchFromJson(json)

	// Assert: Valida resultado
	if err != nil {
		t.Fatalf("Erro ao fazer parsing: %v", err)
	}

	// Valida que CreatedAt foi parseado
	if result.CreatedAt.IsZero() {
		t.Error("CreatedAt não foi parseado do JSON")
	}
}

// TestMatchFromJson_WithAllFields testa parsing com todos os campos
// Valida que todos os campos opcionais são parseados corretamente
func TestMatchFromJson_WithAllFields(t *testing.T) {
	// Arrange: JSON completo
	json := `{
		"user": "user1",
		"name": "games",
		"pattern": "steam",
		"match": "steam.exe",
		"elapsed": 123.45,
		"first_match": "2024-01-01 10:00:00",
		"last_match": "2024-01-01 10:02:03",
		"ocurrences": 5
	}`

	// Act: Faz parsing
	result, err := MatchFromJson(json)

	// Assert: Valida todos os campos
	if err != nil {
		t.Fatalf("Erro ao fazer parsing: %v", err)
	}

	if result.User != "user1" {
		t.Errorf("User = %s, esperado user1", result.User)
	}

	if result.Name != "games" {
		t.Errorf("Name = %s, esperado games", result.Name)
	}

	if result.Pattern != "steam" {
		t.Errorf("Pattern = %s, esperado steam", result.Pattern)
	}

	if result.Match != "steam.exe" {
		t.Errorf("Match = %s, esperado steam.exe", result.Match)
	}

	if result.Elapsed != 123.45 {
		t.Errorf("Elapsed = %.2f, esperado 123.45", result.Elapsed)
	}

	if result.FirstMatch != "2024-01-01 10:00:00" {
		t.Errorf("FirstMatch = %s, esperado 2024-01-01 10:00:00", result.FirstMatch)
	}

	if result.LastMatch != "2024-01-01 10:02:03" {
		t.Errorf("LastMatch = %s, esperado 2024-01-01 10:02:03", result.LastMatch)
	}

	if result.Ocurrences != 5 {
		t.Errorf("Ocurrences = %d, esperado 5", result.Ocurrences)
	}
}

// TestMatchInfo_WithZeroOccurrences testa MatchInfo com zero ocorrências
// Valida que zero é um valor válido
func TestMatchInfo_WithZeroOccurrences(t *testing.T) {
	// Arrange: Cria MatchInfo com zero ocorrências
	info := &MatchInfo{
		Elapsed:    0.0,
		FirstMatch: "",
		LastMatch:  "",
		Ocurrences: 0,
	}

	// Act: Serializa
	json := info.ToJson()
	log := info.ToLog()

	// Assert: Valida que serialização funciona
	if json == "" {
		t.Error("ToJson retornou string vazia")
	}

	if log == "" {
		t.Error("ToLog retornou string vazia")
	}

	// Valida que contém ocurrences com valor zero
	if !strings.Contains(json, "ocurrences") {
		t.Error("JSON não contém campo ocurrences")
	}
}

// TestMatchInfo_WithHighOccurrences testa MatchInfo com muitas ocorrências
// Valida que valores grandes são tratados corretamente
func TestMatchInfo_WithHighOccurrences(t *testing.T) {
	// Arrange: Cria MatchInfo com muitas ocorrências
	info := &MatchInfo{
		Elapsed:    86400.0, // 24 horas
		FirstMatch: "2024-01-01 00:00:00",
		LastMatch:  "2024-01-02 00:00:00",
		Ocurrences: 999999,
	}

	// Act: Serializa
	json := info.ToJson()

	// Assert: Valida que valores grandes são serializados
	if !strings.Contains(json, "999999") {
		t.Error("JSON não contém ocurrences alto")
	}

	if !strings.Contains(json, "86400") {
		t.Error("JSON não contém elapsed alto")
	}
}

// TestMatchList_WithZeroMatches testa MatchList vazia
// Valida que lista vazia é tratada corretamente
func TestMatchList_WithZeroMatches(t *testing.T) {
	// Arrange: Cria JSON com matches vazio
	json := `{"matches": {}}`

	// Act: Faz parsing
	result, err := MatchListFromJson(json)

	// Assert: Valida resultado
	if err != nil {
		t.Fatalf("Erro ao fazer parsing: %v", err)
	}

	if result.Matches == nil {
		t.Fatal("Matches map não foi inicializado")
	}

	if len(result.Matches) != 0 {
		t.Errorf("Esperado 0 matches, obteve %d", len(result.Matches))
	}
}

// TestMatchList_WithNullMatches testa MatchList com null
// Valida tratamento de campo null
func TestMatchList_WithNullMatches(t *testing.T) {
	// Arrange: Cria JSON com matches null
	json := `{"matches": null}`

	// Act: Faz parsing
	result, err := MatchListFromJson(json)

	// Assert: Valida que não houve erro
	if err != nil {
		t.Fatalf("Erro ao fazer parsing: %v", err)
	}

	// Valida que Matches é nil (comportamento esperado do JSON unmarshaling)
	if result.Matches != nil {
		t.Error("Matches deveria ser nil quando JSON contém null")
	}
}

// TestMatch_ToJson_ContainsAllFields testa que ToJson inclui todos os campos
// Valida que nenhum campo é omitido na serialização
func TestMatch_ToJson_ContainsAllFields(t *testing.T) {
	// Arrange: Cria match com todos os campos preenchidos
	match := &Match{
		User:       "user1",
		Name:       "games",
		Pattern:    "steam",
		Match:      "steam.exe",
		Elapsed:    100.5,
		CreatedAt:  time.Now(),
		FirstMatch: "2024-01-01 10:00:00",
		LastMatch:  "2024-01-01 10:01:40",
		Ocurrences: 10,
	}

	// Act: Serializa
	json := match.ToJson()

	// Assert: Valida que todos os campos estão presentes
	requiredFields := []string{
		"user", "name", "pattern", "match", "elapsed",
		"created_at", "first_match", "last_match", "ocurrences",
	}

	for _, field := range requiredFields {
		if !strings.Contains(json, field) {
			t.Errorf("JSON não contém campo obrigatório: %s", field)
		}
	}
}

// TestMatchInfo_ToLog_IsCompact testa que ToLog é compacto
// Valida que não há indentação desnecessária
func TestMatchInfo_ToLog_IsCompact(t *testing.T) {
	// Arrange: Cria MatchInfo
	info := &MatchInfo{
		Elapsed:    100.0,
		FirstMatch: "2024-01-01 10:00:00",
		LastMatch:  "2024-01-01 10:01:40",
		Ocurrences: 5,
	}

	// Act: Serializa para log
	log := info.ToLog()

	// Assert: Valida formato compacto
	if strings.Contains(log, "\n") {
		t.Error("ToLog não deveria conter newlines (deve ser compacto)")
	}

	if strings.Contains(log, "\t") {
		t.Error("ToLog não deveria conter tabs (deve ser compacto)")
	}

	// Valida que contém dados
	if len(log) == 0 {
		t.Error("ToLog retornou string vazia")
	}
}

// TestMatch_MultipleInstances testa criação de múltiplos matches
// Valida que cada instância é independente
func TestMatch_MultipleInstances(t *testing.T) {
	// Arrange & Act: Cria múltiplos matches
	match1 := NewMatch("user1", "games", "steam", "steam.exe", 10.0)
	match2 := NewMatch("user2", "browsers", "chrome", "chrome.exe", 20.0)
	match3 := NewMatch("user3", "media", "vlc", "vlc.exe", 30.0)

	// Assert: Valida que são instâncias diferentes
	if match1 == match2 || match1 == match3 || match2 == match3 {
		t.Error("Matches deveriam ser instâncias diferentes")
	}

	// Valida que cada um tem seus próprios valores
	if match1.User == match2.User {
		t.Error("match1 e match2 não deveriam ter o mesmo User")
	}

	if match1.Elapsed == match2.Elapsed {
		t.Error("match1 e match2 não deveriam ter o mesmo Elapsed")
	}

	// Valida que modificar um não afeta os outros
	match1.Elapsed = 999.0
	if match2.Elapsed == 999.0 || match3.Elapsed == 999.0 {
		t.Error("Modificar match1 afetou outros matches")
	}
}
