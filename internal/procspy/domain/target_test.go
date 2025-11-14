package domain

import (
	"strings"
	"testing"
	"time"
)

// TestNewTargetList testa a criação de uma nova lista de targets vazia
// Valida que a lista é inicializada corretamente com slice vazio
func TestNewTargetList(t *testing.T) {
	// Act: Cria nova lista de targets
	list := NewTargetList()

	// Assert: Valida que a lista foi criada corretamente
	if list == nil {
		t.Fatal("NewTargetList retornou nil")
	}

	if list.Targets == nil {
		t.Error("Targets slice não foi inicializado")
	}

	if len(list.Targets) != 0 {
		t.Errorf("Esperado lista vazia, obteve %d targets", len(list.Targets))
	}
}

// TestTargetListFromJson testa o parsing de JSON para TargetList
// Valida cenários: JSON válido, JSON inválido, JSON vazio
func TestTargetListFromJson(t *testing.T) {
	tests := []struct {
		name    string
		json    string
		wantErr bool
		wantLen int
	}{
		{
			name: "JSON válido com um target",
			json: `{
				"targets": [
					{
						"user": "test_user",
						"name": "games",
						"pattern": "steam|roblox",
						"kill": true
					}
				]
			}`,
			wantErr: false,
			wantLen: 1,
		},
		{
			name: "JSON válido com múltiplos targets",
			json: `{
				"targets": [
					{
						"user": "user1",
						"name": "games",
						"pattern": "steam",
						"kill": true
					},
					{
						"user": "user1",
						"name": "browsers",
						"pattern": "chrome|firefox",
						"kill": false
					}
				]
			}`,
			wantErr: false,
			wantLen: 2,
		},
		{
			name:    "JSON vazio",
			json:    `{"targets": []}`,
			wantErr: false,
			wantLen: 0,
		},
		{
			name:    "JSON inválido",
			json:    `{invalid json}`,
			wantErr: true,
			wantLen: 0,
		},
		{
			name:    "JSON com estrutura incorreta",
			json:    `{"wrong_field": []}`,
			wantErr: false,
			wantLen: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act: Faz parsing do JSON
			result, err := TargetListFromJson(tt.json)

			// Assert: Valida resultado
			if (err != nil) != tt.wantErr {
				t.Errorf("TargetListFromJson() erro = %v, esperava erro = %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if result == nil {
					t.Fatal("TargetListFromJson retornou nil sem erro")
				}

				if len(result.Targets) != tt.wantLen {
					t.Errorf("Esperado %d targets, obteve %d", tt.wantLen, len(result.Targets))
				}

				// Valida que weekdays foram configurados
				for _, target := range result.Targets {
					if target.Weekdays == nil || len(target.Weekdays) == 0 {
						t.Error("Weekdays não foram configurados para target")
					}
				}
			}
		})
	}
}

// TestTarget_Match testa o matching de processos com padrões regex
// Valida: match exato, match parcial, sem match, case insensitive
func TestTarget_Match(t *testing.T) {
	tests := []struct {
		name      string
		pattern   string
		value     string
		wantMatch bool
	}{
		{
			name:      "Match exato",
			pattern:   "steam",
			value:     "steam",
			wantMatch: true,
		},
		{
			name:      "Match parcial no início",
			pattern:   "steam",
			value:     "steam.exe",
			wantMatch: true,
		},
		{
			name:      "Match parcial no meio",
			pattern:   "steam",
			value:     "mysteamapp",
			wantMatch: true,
		},
		{
			name:      "Match com pipe (OR)",
			pattern:   "steam|roblox",
			value:     "roblox.exe",
			wantMatch: true,
		},
		{
			name:      "Sem match",
			pattern:   "steam",
			value:     "chrome.exe",
			wantMatch: false,
		},
		{
			name:      "Match case insensitive",
			pattern:   "(?i)steam",
			value:     "STEAM.EXE",
			wantMatch: true,
		},
		{
			name:      "Match com regex complexa",
			pattern:   "^(steam|roblox).*\\.exe$",
			value:     "steam.exe",
			wantMatch: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange: Cria target com pattern
			target := &Target{
				Pattern: tt.pattern,
			}

			// Act: Testa match
			result := target.Match(tt.value)

			// Assert: Valida resultado
			if result != tt.wantMatch {
				t.Errorf("Match() = %v, esperado %v para pattern '%s' e value '%s'",
					result, tt.wantMatch, tt.pattern, tt.value)
			}

			// Valida que regex foi compilada e armazenada
			if target.rgx == nil {
				t.Error("Regex não foi compilada e armazenada")
			}
		})
	}
}

// TestTarget_AddElapsed testa a adição de tempo decorrido
// Valida que elapsed é acumulado e remaining é calculado corretamente
func TestTarget_AddElapsed(t *testing.T) {
	// Arrange: Cria target com weekdays configurados
	target := &Target{
		Name:     "test",
		Weekdays: map[int]float64{int(time.Now().Weekday()): 1.0}, // 1 hora
		Elapsed:  0,
	}

	// Act: Adiciona tempo decorrido
	target.AddElapsed(1800.0) // 30 minutos

	// Assert: Valida elapsed
	if target.Elapsed != 1800.0 {
		t.Errorf("Elapsed = %.2f, esperado 1800.00", target.Elapsed)
	}

	// Valida remaining
	expectedRemaining := 3600.0 - 1800.0 // 1 hora - 30 minutos
	if target.Remaining != expectedRemaining {
		t.Errorf("Remaining = %.2f, esperado %.2f", target.Remaining, expectedRemaining)
	}

	// Act: Adiciona mais tempo
	target.AddElapsed(900.0) // mais 15 minutos

	// Assert: Valida acumulação
	if target.Elapsed != 2700.0 {
		t.Errorf("Elapsed após segunda adição = %.2f, esperado 2700.00", target.Elapsed)
	}
}

// TestTarget_SetElapsed testa a configuração direta de tempo decorrido
// Valida que elapsed é definido (não acumulado) e remaining é calculado
func TestTarget_SetElapsed(t *testing.T) {
	// Arrange: Cria target
	target := &Target{
		Name:     "test",
		Weekdays: map[int]float64{int(time.Now().Weekday()): 1.0},
		Elapsed:  1000.0, // Valor inicial
	}

	// Act: Define novo elapsed
	target.SetElapsed(2000.0)

	// Assert: Valida que foi substituído, não acumulado
	if target.Elapsed != 2000.0 {
		t.Errorf("Elapsed = %.2f, esperado 2000.00", target.Elapsed)
	}

	expectedRemaining := 3600.0 - 2000.0
	if target.Remaining != expectedRemaining {
		t.Errorf("Remaining = %.2f, esperado %.2f", target.Remaining, expectedRemaining)
	}
}

// TestTarget_ResetElapsed testa o reset de tempo decorrido
// Valida que elapsed volta a zero e remaining volta ao limite
func TestTarget_ResetElapsed(t *testing.T) {
	// Arrange: Cria target com tempo acumulado
	target := &Target{
		Name:     "test",
		Weekdays: map[int]float64{int(time.Now().Weekday()): 1.0},
		Elapsed:  2500.0,
	}

	// Act: Reseta elapsed
	target.ResetElapsed()

	// Assert: Valida reset
	if target.Elapsed != 0 {
		t.Errorf("Elapsed após reset = %.2f, esperado 0.00", target.Elapsed)
	}

	if target.Remaining != 3600.0 {
		t.Errorf("Remaining após reset = %.2f, esperado 3600.00", target.Remaining)
	}
}

// TestTarget_CheckLimit testa a verificação de limite excedido
// Valida cenários: abaixo do limite, no limite, acima do limite
func TestTarget_CheckLimit(t *testing.T) {
	tests := []struct {
		name        string
		elapsed     float64
		weekdayMult float64
		wantExceeded bool
	}{
		{
			name:        "Abaixo do limite",
			elapsed:     1800.0, // 30 minutos
			weekdayMult: 1.0,    // limite de 1 hora
			wantExceeded: false,
		},
		{
			name:        "Exatamente no limite",
			elapsed:     3600.0, // 1 hora
			weekdayMult: 1.0,    // limite de 1 hora
			wantExceeded: true,
		},
		{
			name:        "Acima do limite",
			elapsed:     4000.0, // mais de 1 hora
			weekdayMult: 1.0,    // limite de 1 hora
			wantExceeded: true,
		},
		{
			name:        "Limite zero (sem limite)",
			elapsed:     1000.0,
			weekdayMult: 0.0,
			wantExceeded: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange: Cria target
			target := &Target{
				Name:     "test",
				Weekdays: map[int]float64{int(time.Now().Weekday()): tt.weekdayMult},
				Elapsed:  tt.elapsed,
			}

			// Act: Verifica limite
			result := target.CheckLimit()

			// Assert: Valida resultado
			if result != tt.wantExceeded {
				t.Errorf("CheckLimit() = %v, esperado %v (elapsed=%.2f, limit=%.2f)",
					result, tt.wantExceeded, tt.elapsed, 3600.0*tt.weekdayMult)
			}
		})
	}
}

// TestTarget_CheckWarning testa a verificação de threshold de aviso
// Valida que aviso é disparado em 95% do limite
func TestTarget_CheckWarning(t *testing.T) {
	tests := []struct {
		name         string
		elapsed      float64
		weekdayMult  float64
		wantWarning  bool
	}{
		{
			name:        "Abaixo do threshold de aviso",
			elapsed:     3000.0, // 50 minutos (83% de 1 hora)
			weekdayMult: 1.0,
			wantWarning: false,
		},
		{
			name:        "No threshold de aviso (95%)",
			elapsed:     3421.0, // Ligeiramente acima de 95% de 1 hora
			weekdayMult: 1.0,
			wantWarning: true,
		},
		{
			name:        "Acima do threshold de aviso",
			elapsed:     3500.0, // 58.3 minutos (97% de 1 hora)
			weekdayMult: 1.0,
			wantWarning: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange: Cria target
			target := &Target{
				Name:     "test",
				Weekdays: map[int]float64{int(time.Now().Weekday()): tt.weekdayMult},
				Elapsed:  tt.elapsed,
			}

			// Act: Verifica warning
			result := target.CheckWarning()

			// Assert: Valida resultado
			if result != tt.wantWarning {
				t.Errorf("CheckWarning() = %v, esperado %v (elapsed=%.2f, warning_threshold=%.2f)",
					result, tt.wantWarning, tt.elapsed, 3600.0*tt.weekdayMult*DEFAULT_WARNING_ON)
			}
		})
	}
}

// TestTarget_ToJson testa a serialização de Target para JSON
// Valida que o JSON gerado é válido e contém os campos esperados
func TestTarget_ToJson(t *testing.T) {
	// Arrange: Cria target
	target := &Target{
		User:    "test_user",
		Name:    "games",
		Pattern: "steam|roblox",
		Kill:    true,
		Weekdays: map[int]float64{
			0: 2.0,
			1: 0.5,
		},
	}

	// Act: Serializa para JSON
	json := target.ToJson()

	// Assert: Valida que JSON não está vazio
	if json == "" {
		t.Error("ToJson retornou string vazia")
	}

	// Valida que contém campos esperados
	expectedFields := []string{"user", "name", "pattern", "kill", "weekdays"}
	for _, field := range expectedFields {
		if !strings.Contains(json, field) {
			t.Errorf("JSON não contém campo esperado: %s", field)
		}
	}
}

// TestTarget_ToLog testa a serialização compacta de Target
// Valida que o log é gerado sem indentação
func TestTarget_ToLog(t *testing.T) {
	// Arrange: Cria target
	target := &Target{
		User:    "test_user",
		Name:    "games",
		Pattern: "steam",
	}

	// Act: Serializa para log
	log := target.ToLog()

	// Assert: Valida que log não está vazio
	if log == "" {
		t.Error("ToLog retornou string vazia")
	}

	// Valida que não contém tabs (formato compacto)
	if strings.Contains(log, "\t") {
		t.Error("ToLog não deveria conter tabs (deve ser compacto)")
	}
}

// TestTargetList_ToLog testa a serialização de TargetList
// Valida que a lista completa é serializada corretamente
func TestTargetList_ToLog(t *testing.T) {
	// Arrange: Cria lista com targets
	list := &TargetList{
		Targets: []*Target{
			{User: "user1", Name: "games", Pattern: "steam"},
			{User: "user1", Name: "browsers", Pattern: "chrome"},
		},
	}

	// Act: Serializa para log
	log := list.ToLog()

	// Assert: Valida que log não está vazio
	if log == "" {
		t.Error("ToLog retornou string vazia")
	}

	// Valida que contém ambos os targets
	if !strings.Contains(log, "games") || !strings.Contains(log, "browsers") {
		t.Error("ToLog não contém todos os targets")
	}
}

// TestTargetList_Hash testa a geração de hash da lista
// Valida que o hash é gerado e muda quando a lista muda
func TestTargetList_Hash(t *testing.T) {
	// Arrange: Cria duas listas diferentes
	list1 := &TargetList{
		Targets: []*Target{
			{
				User:     "user1",
				Name:     "games",
				Pattern:  "steam",
				Kill:     true,
				Weekdays: map[int]float64{0: 1.0},
			},
		},
	}

	list2 := &TargetList{
		Targets: []*Target{
			{
				User:     "user1",
				Name:     "browsers",
				Pattern:  "chrome",
				Kill:     false,
				Weekdays: map[int]float64{0: 1.0},
			},
		},
	}

	// Act: Gera hashes
	hash1 := list1.Hash()
	hash2 := list2.Hash()

	// Assert: Valida que hashes são diferentes
	if hash1 == hash2 {
		t.Error("Hashes de listas diferentes deveriam ser diferentes")
	}

	// Valida que hash não está vazio
	if hash1 == "" {
		t.Error("Hash não deveria estar vazio")
	}
}

// TestTarget_AddMatchInfo testa a adição de informações de match
// Valida que elapsed, first_match, last_match e ocurrences são atualizados
func TestTarget_AddMatchInfo(t *testing.T) {
	// Arrange: Cria target e match info
	target := &Target{
		Name:     "test",
		Weekdays: map[int]float64{int(time.Now().Weekday()): 1.0},
		Elapsed:  0,
	}

	matchInfo := &MatchInfo{
		Elapsed:    1500.0,
		FirstMatch: "2024-01-01 10:00:00",
		LastMatch:  "2024-01-01 10:25:00",
		Ocurrences: 5,
	}

	// Act: Adiciona match info
	target.AddMatchInfo(matchInfo)

	// Assert: Valida que informações foram adicionadas
	if target.Elapsed != 1500.0 {
		t.Errorf("Elapsed = %.2f, esperado 1500.00", target.Elapsed)
	}

	if target.FirstMatch != matchInfo.FirstMatch {
		t.Errorf("FirstMatch = %s, esperado %s", target.FirstMatch, matchInfo.FirstMatch)
	}

	if target.LastMatch != matchInfo.LastMatch {
		t.Errorf("LastMatch = %s, esperado %s", target.LastMatch, matchInfo.LastMatch)
	}

	if target.Ocurrences != matchInfo.Ocurrences {
		t.Errorf("Ocurrences = %d, esperado %d", target.Ocurrences, matchInfo.Ocurrences)
	}

	// Valida que remaining foi calculado
	expectedRemaining := 3600.0 - 1500.0
	if target.Remaining != expectedRemaining {
		t.Errorf("Remaining = %.2f, esperado %.2f", target.Remaining, expectedRemaining)
	}
}

// TestTarget_setWeekdays testa a configuração automática de weekdays
// Valida que dias não configurados recebem valores padrão
func TestTarget_setWeekdays(t *testing.T) {
	// Arrange: Cria target sem weekdays
	target := &Target{
		Name: "test",
	}

	// Act: Configura weekdays
	target.setWeekdays()

	// Assert: Valida que todos os 7 dias foram configurados
	if len(target.Weekdays) != 7 {
		t.Errorf("Esperado 7 dias configurados, obteve %d", len(target.Weekdays))
	}

	// Valida valores padrão para dias de semana (1-5)
	for i := 1; i <= 5; i++ {
		if target.Weekdays[i] != DEFAULT_WEEKDAY_LIMIT {
			t.Errorf("Dia %d deveria ter limite %.2f, obteve %.2f",
				i, DEFAULT_WEEKDAY_LIMIT, target.Weekdays[i])
		}
	}

	// Valida valores padrão para fim de semana (0, 6)
	if target.Weekdays[0] != DEFAULT_WEEKEND_LIMIT {
		t.Errorf("Domingo deveria ter limite %.2f, obteve %.2f",
			DEFAULT_WEEKEND_LIMIT, target.Weekdays[0])
	}

	if target.Weekdays[6] != DEFAULT_WEEKEND_LIMIT {
		t.Errorf("Sábado deveria ter limite %.2f, obteve %.2f",
			DEFAULT_WEEKEND_LIMIT, target.Weekdays[6])
	}
}

// TestTarget_setWeekdays_PartialConfig testa weekdays com configuração parcial
// Valida que dias já configurados não são sobrescritos
func TestTarget_setWeekdays_PartialConfig(t *testing.T) {
	// Arrange: Cria target com alguns dias configurados
	target := &Target{
		Name: "test",
		Weekdays: map[int]float64{
			0: 3.0, // Domingo customizado
			1: 0.25, // Segunda customizada
		},
	}

	// Act: Configura weekdays
	target.setWeekdays()

	// Assert: Valida que valores customizados foram mantidos
	if target.Weekdays[0] != 3.0 {
		t.Errorf("Domingo customizado foi sobrescrito: %.2f", target.Weekdays[0])
	}

	if target.Weekdays[1] != 0.25 {
		t.Errorf("Segunda customizada foi sobrescrita: %.2f", target.Weekdays[1])
	}

	// Valida que outros dias receberam valores padrão
	if target.Weekdays[2] != DEFAULT_WEEKDAY_LIMIT {
		t.Errorf("Terça deveria ter valor padrão %.2f, obteve %.2f",
			DEFAULT_WEEKDAY_LIMIT, target.Weekdays[2])
	}
}

// TestTarget_CheckWarning_ZeroLimit testa CheckWarning com limite zero
// Valida que retorna false quando não há limite configurado
func TestTarget_CheckWarning_ZeroLimit(t *testing.T) {
	// Arrange: Cria target com limite zero
	target := &Target{
		Name:     "test",
		Weekdays: map[int]float64{int(time.Now().Weekday()): 0.0},
		Elapsed:  1000.0,
	}

	// Act: Verifica warning
	result := target.CheckWarning()

	// Assert: Valida que retorna false (sem limite, sem warning)
	if result != false {
		t.Error("CheckWarning deveria retornar false quando limite é zero")
	}
}

// TestTarget_Match_ReuseCompiledRegex testa que regex é reutilizada
// Valida que regex compilada é armazenada e reutilizada
func TestTarget_Match_ReuseCompiledRegex(t *testing.T) {
	// Arrange: Cria target
	target := &Target{
		Pattern: "steam",
	}

	// Act: Faz primeiro match (compila regex)
	target.Match("steam.exe")
	firstRegex := target.rgx

	// Act: Faz segundo match (reutiliza regex)
	target.Match("steam_api.dll")
	secondRegex := target.rgx

	// Assert: Valida que é a mesma instância de regex
	if firstRegex != secondRegex {
		t.Error("Regex deveria ser reutilizada, não recompilada")
	}
}

// TestTarget_getLimit_WithNegativeRemaining testa getLimit com remaining negativo
// Valida que remaining é resetado quando negativo
func TestTarget_getLimit_WithNegativeRemaining(t *testing.T) {
	// Arrange: Cria target com remaining negativo
	target := &Target{
		Name:      "test",
		Weekdays:  map[int]float64{int(time.Now().Weekday()): 1.0},
		Remaining: -100.0,
	}

	// Act: Chama getLimit
	limit := target.getLimit()

	// Assert: Valida que remaining foi resetado para o limite
	if target.Remaining != limit {
		t.Errorf("Remaining deveria ser resetado para %.2f, obteve %.2f", limit, target.Remaining)
	}

	if target.Remaining <= 0 {
		t.Error("Remaining não deveria ser negativo após getLimit")
	}
}

// TestTarget_getLimit_MissingWeekday testa getLimit quando dia não está configurado
// Valida que usa valor padrão de dia de semana
func TestTarget_getLimit_MissingWeekday(t *testing.T) {
	// Arrange: Cria target sem o dia atual configurado
	target := &Target{
		Name:     "test",
		Weekdays: map[int]float64{}, // Vazio
	}

	// Act: Chama getLimit
	limit := target.getLimit()

	// Assert: Valida que usou valor padrão
	expectedLimit := DEFAULT_BASE_LIMIT * DEFAULT_WEEKDAY_LIMIT
	if limit != expectedLimit {
		t.Errorf("Limit = %.2f, esperado %.2f (valor padrão)", limit, expectedLimit)
	}
}

// TestTargetList_Hash_EmptyList testa Hash com lista vazia
// Valida que hash de lista vazia é string vazia
func TestTargetList_Hash_EmptyList(t *testing.T) {
	// Arrange: Cria lista vazia
	list := NewTargetList()

	// Act: Gera hash
	hash := list.Hash()

	// Assert: Valida que hash é string vazia
	if hash != "" {
		t.Errorf("Hash de lista vazia deveria ser '', obteve '%s'", hash)
	}
}

// TestTargetList_Hash_SameContent testa que mesmo conteúdo gera mesmo hash
// Valida consistência do hash
func TestTargetList_Hash_SameContent(t *testing.T) {
	// Arrange: Cria duas listas com mesmo conteúdo
	list1 := &TargetList{
		Targets: []*Target{
			{
				User:     "user1",
				Name:     "games",
				Pattern:  "steam",
				Kill:     true,
				Weekdays: map[int]float64{0: 1.0},
			},
		},
	}

	list2 := &TargetList{
		Targets: []*Target{
			{
				User:     "user1",
				Name:     "games",
				Pattern:  "steam",
				Kill:     true,
				Weekdays: map[int]float64{0: 1.0},
			},
		},
	}

	// Act: Gera hashes
	hash1 := list1.Hash()
	hash2 := list2.Hash()

	// Assert: Valida que hashes são iguais
	if hash1 != hash2 {
		t.Error("Hashes de listas com mesmo conteúdo deveriam ser iguais")
	}
}

// TestTarget_AddElapsed_MultipleAdditions testa múltiplas adições de elapsed
// Valida acumulação correta ao longo de várias chamadas
func TestTarget_AddElapsed_MultipleAdditions(t *testing.T) {
	// Arrange: Cria target
	target := &Target{
		Name:     "test",
		Weekdays: map[int]float64{int(time.Now().Weekday()): 2.0}, // 2 horas
		Elapsed:  0,
	}

	// Act: Adiciona elapsed várias vezes
	target.AddElapsed(600.0)  // 10 min
	target.AddElapsed(900.0)  // 15 min
	target.AddElapsed(1200.0) // 20 min

	// Assert: Valida acumulação total
	expectedElapsed := 2700.0 // 45 min total
	if target.Elapsed != expectedElapsed {
		t.Errorf("Elapsed = %.2f, esperado %.2f", target.Elapsed, expectedElapsed)
	}

	// Valida remaining
	expectedRemaining := 7200.0 - 2700.0
	if target.Remaining != expectedRemaining {
		t.Errorf("Remaining = %.2f, esperado %.2f", target.Remaining, expectedRemaining)
	}
}

// TestTarget_SetElapsed_OverwritesPrevious testa que SetElapsed sobrescreve
// Valida que valor anterior é substituído, não acumulado
func TestTarget_SetElapsed_OverwritesPrevious(t *testing.T) {
	// Arrange: Cria target com elapsed inicial
	target := &Target{
		Name:     "test",
		Weekdays: map[int]float64{int(time.Now().Weekday()): 1.0},
		Elapsed:  5000.0,
	}

	// Act: Define novo elapsed (menor que o anterior)
	target.SetElapsed(500.0)

	// Assert: Valida que foi substituído
	if target.Elapsed != 500.0 {
		t.Errorf("Elapsed = %.2f, esperado 500.00 (deveria sobrescrever)", target.Elapsed)
	}
}

// TestTarget_AddMatchInfo_AccumulatesElapsed testa que AddMatchInfo acumula
// Valida que elapsed é acumulado, não substituído
func TestTarget_AddMatchInfo_AccumulatesElapsed(t *testing.T) {
	// Arrange: Cria target com elapsed inicial
	target := &Target{
		Name:     "test",
		Weekdays: map[int]float64{int(time.Now().Weekday()): 1.0},
		Elapsed:  1000.0,
	}

	matchInfo := &MatchInfo{
		Elapsed:    500.0,
		FirstMatch: "2024-01-01 10:00:00",
		LastMatch:  "2024-01-01 10:08:20",
		Ocurrences: 3,
	}

	// Act: Adiciona match info
	target.AddMatchInfo(matchInfo)

	// Assert: Valida que elapsed foi acumulado
	if target.Elapsed != 1500.0 {
		t.Errorf("Elapsed = %.2f, esperado 1500.00 (acumulado)", target.Elapsed)
	}
}

// TestTargetListFromJson_WithWeekdaysConfig testa parsing com weekdays customizados
// Valida que weekdays do JSON são preservados e complementados
func TestTargetListFromJson_WithWeekdaysConfig(t *testing.T) {
	// Arrange: JSON com weekdays parcialmente configurados
	json := `{
		"targets": [
			{
				"user": "user1",
				"name": "games",
				"pattern": "steam",
				"weekdays": {
					"0": 3.0,
					"6": 2.5
				}
			}
		]
	}`

	// Act: Faz parsing
	result, err := TargetListFromJson(json)

	// Assert: Valida resultado
	if err != nil {
		t.Fatalf("Erro ao fazer parsing: %v", err)
	}

	if len(result.Targets) != 1 {
		t.Fatalf("Esperado 1 target, obteve %d", len(result.Targets))
	}

	target := result.Targets[0]

	// Valida que weekdays customizados foram preservados
	if target.Weekdays[0] != 3.0 {
		t.Errorf("Domingo = %.2f, esperado 3.0", target.Weekdays[0])
	}

	if target.Weekdays[6] != 2.5 {
		t.Errorf("Sábado = %.2f, esperado 2.5", target.Weekdays[6])
	}

	// Valida que outros dias foram preenchidos com padrão
	if len(target.Weekdays) != 7 {
		t.Errorf("Esperado 7 dias configurados, obteve %d", len(target.Weekdays))
	}
}

// TestTarget_CheckLimit_EdgeCases testa CheckLimit com casos extremos
// Valida comportamento com valores extremos
func TestTarget_CheckLimit_EdgeCases(t *testing.T) {
	tests := []struct {
		name         string
		elapsed      float64
		weekdayMult  float64
		wantExceeded bool
	}{
		{
			name:         "Elapsed muito pequeno",
			elapsed:      0.001,
			weekdayMult:  1.0,
			wantExceeded: false,
		},
		{
			name:         "Elapsed muito grande",
			elapsed:      999999.0,
			weekdayMult:  1.0,
			wantExceeded: true,
		},
		{
			name:         "Limite muito pequeno",
			elapsed:      1.0,
			weekdayMult:  0.0001,
			wantExceeded: true,
		},
		{
			name:         "Limite muito grande",
			elapsed:      1000.0,
			weekdayMult:  100.0,
			wantExceeded: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			target := &Target{
				Name:     "test",
				Weekdays: map[int]float64{int(time.Now().Weekday()): tt.weekdayMult},
				Elapsed:  tt.elapsed,
			}

			result := target.CheckLimit()

			if result != tt.wantExceeded {
				t.Errorf("CheckLimit() = %v, esperado %v", result, tt.wantExceeded)
			}
		})
	}
}
