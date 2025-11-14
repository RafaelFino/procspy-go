# Análise de Coverage - Procspy

**Data da Análise**: 2025-11-13  
**Coverage Total Atual**: 35.9%  
**Meta de Coverage**: 99% (mínimo 70% para componentes complexos)

## Resumo por Pacote

| Pacote | Coverage Atual | Meta | Gap | Prioridade | Status |
|--------|----------------|------|-----|------------|--------|
| **config** | 95.5% | 99% | 3.5% | Alta | ✅ Próximo da meta |
| **domain** | 84.9% | 99% | 14.1% | Alta | ⚠️ Precisa melhorar |
| **storage** | 69.5% | 95% | 25.5% | Alta | ⚠️ Precisa melhorar |
| **service** | 47.3% | 95% | 47.7% | Alta | ❌ Crítico |
| **server** | 35.6% | 90% | 54.4% | Média | ❌ Crítico |
| **handlers** | 10.2% | 90% | 79.8% | Alta | ❌ Crítico |
| **watcher** | 4.8% | 70% | 65.2% | Média | ❌ Crítico |
| **client** | 2.2% | 70% | 67.8% | Alta | ❌ Crítico |
| **cmd/*** | 0.0% | N/A | N/A | Baixa | ⏸️ Não testável (main) |

## Análise Detalhada por Pacote

### 1. Config (95.5% → Meta: 99%)

**Status**: ✅ Muito próximo da meta

**Funções com coverage < 100%**:
- `Client.ToJson()`: 75.0%
- `Server.ToJson()`: 75.0%
- `Watcher.ToJson()`: 75.0%

**Ações Necessárias**:
- Adicionar testes para casos de erro em serialização JSON
- Testar comportamento com campos nil/vazios
- Testar valores extremos

**Estimativa**: 2-3 testes adicionais por arquivo

---

### 2. Domain (84.9% → Meta: 99%)

**Status**: ⚠️ Precisa de melhorias

**Funções com coverage < 100%**:
- `Command.ToLog()`: 60.0%
- `Command.ToJson()`: 60.0%
- `Match.ToLog()`: 60.0%
- `Match.ToJson()`: 60.0%
- `MatchInfo.ToJson()`: 60.0%
- `MatchInfo.ToLog()`: 60.0%
- `Target.ToLog()`: 60.0%
- `Target.ToJson()`: 75.0%
- `TargetList.ToLog()`: 60.0%
- `Target.CheckWarning()`: 75.0%

**Ações Necessárias**:
- Adicionar testes para casos de erro em serialização
- Testar todos os branches de CheckWarning
- Testar valores extremos e edge cases
- Adicionar testes para campos opcionais

**Estimativa**: 10-15 testes adicionais

---

### 3. Storage (69.5% → Meta: 95%)

**Status**: ⚠️ Precisa de melhorias significativas

**Funções com coverage < 100%**:
- `Command.NewCommand()`: 66.7%
- `Command.Init()`: 62.5%
- `Command.Close()`: 50.0%
- `Command.InsertCommand()`: 80.0%
- `Command.GetCommands()`: 73.3%
- `Match.NewMatch()`: 66.7%
- `Match.Init()`: 62.5%
- `Match.Close()`: 50.0%
- `DbConnection.GetConn()`: 77.8%
- `DbConnection.Close()`: 80.0%
- `DbConnection.Exec()`: 57.9%

**Ações Necessárias**:
- Testar cenários de erro de conexão
- Testar comportamento com banco fechado
- Testar queries inválidas
- Testar transações e rollback
- Adicionar testes para Init com erros
- Testar Close múltiplas vezes
- Testar operações após Close

**Estimativa**: 20-25 testes adicionais

---

### 4. Service (47.3% → Meta: 95%)

**Status**: ❌ Crítico - Precisa de trabalho extensivo

**Funções com coverage 0%**:
- `Command.Close()`: 0.0%
- `Match.Close()`: 0.0%
- `Match.GetMatchesInfo()`: 0.0%
- `Target.getFromUrl()`: 0.0%

**Funções com coverage < 100%**:
- `Command.NewCommand()`: 71.4%
- `Match.NewMatch()`: 71.4%
- `Match.GetMatches()`: 75.0%
- `Target.GetTargets()`: 26.3%

**Ações Necessárias**:
- Adicionar testes para Close() em ambos os services
- Adicionar testes para GetMatchesInfo()
- Adicionar testes para getFromUrl() com mock HTTP
- Melhorar cobertura de GetTargets com diferentes cenários
- Testar cenários de erro de storage
- Testar validação de MATCH_MAX_ELAPSED
- Testar cache de targets

**Estimativa**: 25-30 testes adicionais

---

### 5. Handlers (10.2% → Meta: 90%)

**Status**: ❌ Crítico - Quase sem cobertura

**Funções com coverage 0%**:
- `Command.InsertCommand()`: 0.0%
- `Match.InsertMatch()`: 0.0%
- `Report.GetReport()`: 0.0%
- `Target.GetTargets()`: 0.0%

**Ações Necessárias**:
- Adicionar testes completos para InsertCommand:
  - Request válido
  - Request inválido (JSON malformado)
  - Body vazio
  - Usuário inválido
  - Erro de service layer
  - Diferentes métodos HTTP
  
- Adicionar testes completos para InsertMatch:
  - Request válido
  - Request inválido
  - Validação de campos
  - Erro de service layer
  
- Adicionar testes completos para GetReport:
  - Diferentes filtros (user, interval, target)
  - Combinações de filtros
  - Filtros inválidos
  - Erro de service layer
  
- Adicionar testes completos para GetTargets:
  - Usuário válido
  - Usuário inválido
  - Erro de service layer
  - Response vazio

**Estimativa**: 40-50 testes adicionais

---

### 6. Server (35.6% → Meta: 90%)

**Status**: ❌ Crítico

**Funções com coverage 0%**:
- `Server.Start()`: 0.0%

**Ações Necessárias**:
- Adicionar testes para Start():
  - Inicialização bem-sucedida
  - Erro ao iniciar servidor HTTP
  - Erro ao criar services
  - Validação de configuração
  - Teste de rotas registradas

**Estimativa**: 10-15 testes adicionais

---

### 7. Client (2.2% → Meta: 70%)

**Status**: ❌ Crítico - Componente complexo

**Funções com coverage 0%**:
- `startHttpServer()`: 0.0%
- `stopHttpServer()`: 0.0%
- `httpGet()`: 0.0%
- `httpPost()`: 0.0%
- `updateTargets()`: 0.0%
- `postMatch()`: 0.0%
- `postCommand()`: 0.0%
- `consumeBuffers()`: 0.0%
- `run()`: 0.0%
- `kill()`: 0.0%
- `Start()`: 0.0%
- `Stop()`: 0.0%
- `executeCommand()`: 0.0%

**Ações Necessárias**:
- Adicionar testes com httptest para funções HTTP
- Mockar interações com sistema operacional
- Testar ciclo de vida (Start/Stop)
- Testar consumeBuffers com diferentes cenários
- Testar kill com processos mockados
- Documentar limitações de coverage devido a interações com OS

**Nota**: Este componente tem alta complexidade devido a interações com processos do sistema operacional. Meta de 70% é aceitável.

**Estimativa**: 30-40 testes adicionais

---

### 8. Watcher (4.8% → Meta: 70%)

**Status**: ❌ Crítico - Componente complexo

**Funções não cobertas**: Maioria das funções principais

**Ações Necessárias**:
- Adicionar testes para check() com httptest
- Testar cenários de client up/down
- Testar executeCommand
- Testar ciclo de vida do watcher
- Documentar limitações de coverage devido a interações com OS

**Nota**: Este componente tem alta complexidade devido a interações com processos do sistema operacional. Meta de 70% é aceitável.

**Estimativa**: 20-25 testes adicionais

---

## Plano de Ação Priorizado

### Fase 1: Melhorias Rápidas (Coverage fácil de aumentar)
1. **Config** (95.5% → 99%): ~2-3 testes
2. **Domain** (84.9% → 99%): ~10-15 testes

**Impacto**: +8% no coverage total  
**Esforço**: Baixo  
**Tempo estimado**: 2-3 horas

### Fase 2: Componentes Críticos de Negócio
3. **Storage** (69.5% → 95%): ~20-25 testes
4. **Service** (47.3% → 95%): ~25-30 testes

**Impacto**: +25% no coverage total  
**Esforço**: Médio  
**Tempo estimado**: 6-8 horas

### Fase 3: Camada de Apresentação
5. **Handlers** (10.2% → 90%): ~40-50 testes
6. **Server** (35.6% → 90%): ~10-15 testes

**Impacto**: +35% no coverage total  
**Esforço**: Alto  
**Tempo estimado**: 8-10 horas

### Fase 4: Componentes Complexos
7. **Client** (2.2% → 70%): ~30-40 testes
8. **Watcher** (4.8% → 70%): ~20-25 testes

**Impacto**: +15% no coverage total  
**Esforço**: Alto (devido a complexidade)  
**Tempo estimado**: 10-12 horas

---

## Estimativa Total

- **Testes adicionais necessários**: ~160-200 testes
- **Tempo total estimado**: 26-33 horas
- **Coverage final esperado**: 85-95% (total)

## Próximos Passos

1. ✅ Executar análise de coverage (CONCLUÍDO)
2. ⏭️ Começar Fase 1: Melhorar config e domain
3. ⏭️ Continuar com Fase 2: Storage e Service
4. ⏭️ Prosseguir com Fase 3: Handlers e Server
5. ⏭️ Finalizar com Fase 4: Client e Watcher
6. ⏭️ Validar coverage final e documentar

## Notas Importantes

- Os pacotes `cmd/*` não precisam de testes (são apenas entry points)
- Client e Watcher têm meta reduzida (70%) devido à complexidade de interação com OS
- Priorizar testes de casos de erro e edge cases
- Usar mocks para isolar dependências externas
- Documentar limitações conhecidas de coverage
