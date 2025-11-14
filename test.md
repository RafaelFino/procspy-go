# Documentação de Testes Unitários - Procspy

## Visão Geral

Este documento descreve o escopo e propósito de cada arquivo de teste no projeto Procspy, além de fornecer informações sobre compatibilidade cross-platform e requisitos para execução em diferentes sistemas operacionais.

### Estatísticas de Coverage

- **Coverage Total**: 64.9%
- **Pacotes Testados**: 8/8
- **Testes Totais**: 214

### Coverage por Pacote

| Pacote    | Coverage | Status | Meta    |
|-----------|----------|--------|---------|
| config    | 95.5%    | ✓      | 99%     |
| domain    | 85.7%    | ⚠      | 99%     |
| storage   | 74.8%    | ⚠      | 95%+    |
| watcher   | 73.8%    | ⚠      | 70%+    |
| service   | 68.9%    | ⚠      | 95%+    |
| client    | 54.0%    | ❌     | 70%+    |
| handlers  | 43.8%    | ❌     | 90%+    |
| server    | 35.6%    | ❌     | 90%+    |

**Legenda:**
- ✓ = Atingiu meta de coverage
- ⚠ = Coverage aceitável mas abaixo da meta
- ❌ = Coverage abaixo do mínimo aceitável

### Convenções

- Todos os testes seguem o padrão AAA (Arrange-Act-Assert)
- Testes de múltiplos cenários usam table-driven tests
- Comentários em português explicam o propósito de cada teste
- Testes são independentes e podem ser executados em qualquer ordem
- Uso de banco de dados SQLite em memória (`:memory:`) para testes de storage

## Compatibilidade Cross-Platform

### Sistemas Operacionais Suportados

Os testes foram projetados para serem executados em:
- **Linux** (testado e validado)
- **macOS** (compatível)
- **Windows** (compatível com Git Bash ou WSL)

### Requisitos por Sistema Operacional

#### Linux
- Go 1.16 ou superior
- Bash shell (padrão na maioria das distribuições)
- Utilitários padrão: `grep`, `awk`, `sed`, `wc`, `bc`
- SQLite3 (incluído no driver Go)

**Execução:**
```bash
./test.sh
```

#### macOS
- Go 1.16 ou superior
- Bash shell (padrão no macOS)
- Utilitários padrão: `grep`, `awk`, `sed`, `wc`, `bc`
- SQLite3 (incluído no driver Go)

**Execução:**
```bash
./test.sh
```

#### Windows

**Opção 1: Git Bash (Recomendado)**
- Go 1.16 ou superior
- Git for Windows (inclui Git Bash)
- SQLite3 (incluído no driver Go)

**Execução:**
```bash
# No Git Bash
./test.sh
```

**Opção 2: WSL (Windows Subsystem for Linux)**
- WSL 2 instalado
- Go 1.16 ou superior instalado no WSL
- Utilitários Linux padrão

**Execução:**
```bash
# No terminal WSL
./test.sh
```

**Opção 3: PowerShell (Alternativa)**
Se não for possível usar Git Bash ou WSL, execute diretamente com Go:
```powershell
# Executar todos os testes
go test -v -race -coverprofile=coverage.out ./...

# Ver coverage
go tool cover -func=coverage.out
```

### Características Cross-Platform dos Testes

#### ✅ Aspectos Compatíveis

1. **Código Go Puro**: Todos os testes usam apenas bibliotecas padrão do Go
2. **Banco de Dados em Memória**: Uso de SQLite `:memory:` evita dependências de filesystem
3. **Sem Comandos de Sistema**: Nenhum teste executa comandos específicos de OS
4. **Paths Relativos**: Todos os paths são relativos ou gerenciados pelo Go
5. **HTTP Testing**: Uso de `httptest` para simular servidores sem dependências de rede
6. **Sem Hardcoded Paths**: Nenhum path absoluto específico de OS nos testes

#### 🔍 Validações Realizadas

- ✅ Nenhum uso de `exec.Command` com comandos específicos de OS
- ✅ Nenhum uso de paths absolutos Unix (`/tmp/`, `/usr/`, etc.)
- ✅ Nenhum uso de paths absolutos Windows (`C:\`, etc.)
- ✅ Nenhum uso de `runtime.GOOS` para lógica condicional
- ✅ Nenhuma dependência de variáveis de ambiente específicas de OS
- ✅ Nenhum uso de separadores de path hardcoded

### Script test.sh

O script `test.sh` usa utilitários Unix padrão:
- `grep`: Busca de padrões
- `awk`: Processamento de texto
- `sed`: Substituição de strings
- `wc`: Contagem de linhas
- `bc`: Cálculos matemáticos

**Nota para Windows**: Estes utilitários estão disponíveis no Git Bash e WSL.

### Executando Testes Específicos

#### Testar um pacote específico
```bash
# Linux/macOS/Git Bash
go test -v ./internal/procspy/domain

# PowerShell
go test -v ./internal/procspy/domain
```

#### Testar uma função específica
```bash
# Linux/macOS/Git Bash
go test -v -run TestTarget_Match ./internal/procspy/domain

# PowerShell
go test -v -run TestTarget_Match ./internal/procspy/domain
```

#### Gerar coverage HTML
```bash
# Linux/macOS/Git Bash/PowerShell
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

## Cobertura de Testes por Pacote

### internal/procspy/client (Coverage: 54.0%)

#### client_test.go
**Escopo**: Testes do componente Client que monitora processos localmente

**Funções Testadas**:
- `TestNewSpy`: Valida criação de instância Spy com configuração
- `TestSpy_IsEnabled`: Valida verificação de habilitação do spy
- `TestRoundFloat`: Valida arredondamento de valores float
  - Cenário: 2 casas decimais
  - Cenário: 1 casa decimal
  - Cenário: 0 casas decimais
  - Cenário: Valor exato
- Testes adicionais de funções auxiliares e formatação

**Limitações**:
- Função `Run()` não testada (loop infinito bloqueante)
- Funções de scan de processos não testadas (dependem de OS)

### internal/procspy/config (Coverage: 95.5%)

#### client_test.go
**Escopo**: Testes de parsing e validação de configuração do client

**Funções Testadas**:
- `TestNewConfig`: Valida criação de configuração vazia
- `TestClient_SetDefaults`: Valida aplicação de valores padrão
  - Cenário: Todos os campos vazios
  - Cenário: Interval abaixo do mínimo
  - Cenário: Valores válidos não são alterados
- `TestClient_ToJson`: Valida serialização para JSON
- `TestClientConfigFromJson`: Valida parsing de JSON
  - Cenário: JSON válido completo
  - Cenário: JSON mínimo com defaults
  - Cenário: JSON inválido
  - Cenário: JSON com interval abaixo do mínimo
- `TestClientConfigFromFile`: Valida leitura de arquivo de configuração
- `TestClientConfigFromFile_FileNotFound`: Valida tratamento de arquivo inexistente
- `TestClientConfigFromFile_InvalidJSON`: Valida tratamento de JSON inválido
- `TestClient_Serialization_RoundTrip`: Valida serialização e deserialização
- `TestClient_DebugOmitEmpty`: Valida omissão de campo debug quando false

#### server_test.go
**Escopo**: Testes de parsing e validação de configuração do server

**Funções Testadas**:
- `TestNewServer`: Valida criação de configuração de servidor
- `TestServer_ToJson`: Valida serialização para JSON
- `TestServerConfigFromJson`: Valida parsing de JSON
  - Cenário: JSON válido
  - Cenário: JSON inválido
- `TestServerConfigFromFile`: Valida leitura de arquivo de configuração
- `TestServerConfigFromFile_NotFound`: Valida tratamento de arquivo inexistente

#### watcher_test.go
**Escopo**: Testes de parsing e validação de configuração do watcher

**Funções Testadas**:
- `TestNewWatcher`: Valida criação de configuração do watcher
- `TestWatcher_SetDefaults`: Valida aplicação de valores padrão
  - Cenário: Todos os campos vazios
  - Cenário: Interval abaixo do mínimo
  - Cenário: Valores válidos não são alterados
- `TestWatcher_ToJson`: Valida serialização para JSON
- `TestWatcherConfigFromJson`: Valida parsing de JSON
  - Cenário: JSON válido
  - Cenário: JSON inválido
- `TestWatcherConfigFromFile`: Valida leitura de arquivo de configuração
- `TestWatcherConfigFromFile_NotFound`: Valida tratamento de arquivo inexistente

### internal/procspy/domain (Coverage: 85.7%)

#### command_test.go
**Escopo**: Testes do modelo Command e operações de serialização

**Funções Testadas**:
- `TestNewCommand`: Valida criação de comando
- `TestCommand_ToJson`: Valida serialização para JSON
- `TestCommand_ToLog`: Valida formatação para log
- `TestCommandFromJson`: Valida parsing de JSON
  - Cenário: JSON válido completo
  - Cenário: JSON válido mínimo
  - Cenário: JSON inválido
  - Cenário: JSON vazio
  - Cenário: String vazia
  - Cenário: JSON com campos extras
- `TestCommand_WithEmptyStrings`: Valida comando com strings vazias
- `TestCommand_WithLongStrings`: Valida comando com strings longas
- `TestCommand_Serialization_RoundTrip`: Valida serialização e deserialização
- `TestCommand_WithSpecialCharacters`: Valida comando com caracteres especiais
- `TestCommand_DefaultSource`: Valida source padrão
- `TestCommand_ModifySource`: Valida modificação de source
- `TestCommand_CreatedAtPrecision`: Valida precisão de timestamp
- `TestCommand_MultipleInstances`: Valida múltiplas instâncias independentes

#### match_test.go
**Escopo**: Testes do modelo Match e operações de serialização

**Funções Testadas**:
- `TestNewMatch`: Valida criação de match
- `TestMatch_ToJson`: Valida serialização para JSON
- `TestMatch_ToLog`: Valida formatação para log
- `TestMatchFromJson`: Valida parsing de JSON
  - Cenário: JSON válido completo
  - Cenário: JSON válido mínimo
  - Cenário: JSON inválido
  - Cenário: JSON vazio
  - Cenário: String vazia
- `TestMatchListFromJson`: Valida parsing de lista de matches
  - Cenário: JSON válido com matches
  - Cenário: JSON vazio
  - Cenário: JSON inválido
  - Cenário: JSON com um match
- `TestMatchInfo_ToJson`: Valida serialização de MatchInfo
- `TestMatchInfo_ToLog`: Valida formatação de MatchInfo para log
- `TestMatch_WithZeroElapsed`: Valida match com elapsed zero
- `TestMatch_WithLargeElapsed`: Valida match com elapsed grande
- `TestMatchInfo_WithEmptyStrings`: Valida MatchInfo com strings vazias
- `TestMatch_Serialization_RoundTrip`: Valida serialização e deserialização
- `TestMatchList_WithMultipleMatches`: Valida lista com múltiplos matches

#### target_test.go
**Escopo**: Testes do modelo Target e operações relacionadas

**Funções Testadas**:
- `TestNewTargetList`: Valida criação de lista vazia de targets
- `TestTargetListFromJson`: Valida parsing de JSON para TargetList
  - Cenário: JSON válido com um target
  - Cenário: JSON válido com múltiplos targets
  - Cenário: JSON vazio
  - Cenário: JSON inválido
  - Cenário: JSON com estrutura incorreta
- `TestTarget_Match`: Valida matching de processos com regex
  - Cenário: Match exato
  - Cenário: Match parcial no início
  - Cenário: Match parcial no meio
  - Cenário: Match com pipe (OR)
  - Cenário: Sem match
  - Cenário: Match case insensitive
  - Cenário: Match com regex complexa
- `TestTarget_AddElapsed`: Valida acumulação de tempo
- `TestTarget_SetElapsed`: Valida definição de tempo
- `TestTarget_ResetElapsed`: Valida reset de tempo
- `TestTarget_CheckLimit`: Valida verificação de limite
  - Cenário: Abaixo do limite
  - Cenário: Exatamente no limite
  - Cenário: Acima do limite
  - Cenário: Limite zero (sem limite)
- `TestTarget_CheckWarning`: Valida verificação de aviso
  - Cenário: Abaixo do threshold de aviso
  - Cenário: No threshold de aviso (95%)
  - Cenário: Acima do threshold de aviso
- `TestTarget_ToJson`: Valida serialização para JSON
- `TestTarget_ToLog`: Valida formatação para log
- `TestTargetList_ToLog`: Valida formatação de lista para log
- `TestTargetList_Hash`: Valida geração de hash da lista
- `TestTarget_AddMatchInfo`: Valida adição de informação de match
- `TestTarget_setWeekdays`: Valida configuração de dias da semana
- `TestTarget_setWeekdays_PartialConfig`: Valida configuração parcial de dias

### internal/procspy/handlers (Coverage: 43.8%)

#### command_test.go
**Escopo**: Testes do handler de comandos HTTP

**Funções Testadas**:
- `TestNewCommand`: Valida criação do handler de comando

#### healthcheck_test.go
**Escopo**: Testes do handler de health check

**Funções Testadas**:
- `TestNewHealthcheck`: Valida criação do handler de healthcheck
- `TestHealthcheck_GetStatus`: Valida endpoint de status
  - Valida response status 200
  - Valida JSON de resposta

#### match_test.go
**Escopo**: Testes do handler de matches HTTP

**Funções Testadas**:
- `TestNewMatch`: Valida criação do handler de match

#### report_test.go
**Escopo**: Testes do handler de relatórios

**Funções Testadas**:
- `TestNewReport`: Valida criação do handler de report
- `TestFormatInterval`: Valida formatação de intervalos de tempo
  - Cenário: 1 segundo
  - Cenário: 60 segundos
  - Cenário: 3600 segundos

#### target_test.go
**Escopo**: Testes do handler de targets HTTP

**Funções Testadas**:
- `TestNewTarget`: Valida criação do handler de target

#### util_test.go
**Escopo**: Testes de funções utilitárias dos handlers

**Funções Testadas**:
- `TestValidateUser_ValidUser`: Valida usuário existente
- `TestValidateUser_InvalidUser`: Valida usuário inexistente

### internal/procspy/server (Coverage: 35.6%)

#### server_test.go
**Escopo**: Testes de inicialização e configuração do servidor

**Funções Testadas**:
- `TestNewServer`: Valida criação de servidor
  - Valida inicialização de services
  - Valida criação de handlers

### internal/procspy/service (Coverage: 68.9%)

#### command_test.go
**Escopo**: Testes do service layer de comandos

**Funções Testadas**:
- `TestNewCommand`: Valida criação do service de comando
- `TestCommand_InsertCommand`: Valida inserção de comando
- `TestCommand_GetCommands`: Valida busca de comandos por usuário

#### match_test.go
**Escopo**: Testes do service layer de matches

**Funções Testadas**:
- `TestNewMatch`: Valida criação do service de match
- `TestMatch_InsertMatch`: Valida inserção de match
- `TestMatch_InsertMatch_MaxElapsed`: Valida validação de MATCH_MAX_ELAPSED
- `TestMatch_GetMatches`: Valida busca de matches

#### target_test.go
**Escopo**: Testes do service layer de targets

**Funções Testadas**:
- `TestNewTarget`: Valida criação do service de target
- `TestTarget_GetTargets_NoUser`: Valida busca sem usuário

#### user_test.go
**Escopo**: Testes do service layer de usuários

**Funções Testadas**:
- `TestNewUsers`: Valida criação do service de users
- `TestUsers_GetUsers`: Valida busca de usuários
- `TestUsers_Exists`: Valida verificação de existência de usuário

### internal/procspy/storage (Coverage: 74.8%)

#### command_test.go
**Escopo**: Testes de persistência de comandos

**Funções Testadas**:
- `TestNewCommand`: Valida criação do storage de comando
- `TestCommand_InsertCommand`: Valida inserção no banco
- `TestCommand_GetCommands`: Valida busca de comandos
- `TestCommand_Close`: Valida fechamento de conexão

**Nota**: Usa banco SQLite em memória (`:memory:`)

#### dbconn_test.go
**Escopo**: Testes de conexão com banco de dados

**Funções Testadas**:
- `TestNewDbConnection`: Valida criação de conexão
- `TestDbConnection_GetConn`: Valida obtenção de conexão
- `TestDbConnection_Close`: Valida fechamento de conexão
- `TestDbConnection_Exec`: Valida execução de queries
- `TestDbConnection_makeDBPath`: Valida criação de path do banco

**Nota**: Usa banco SQLite em memória (`:memory:`)

#### match_test.go
**Escopo**: Testes de persistência de matches

**Funções Testadas**:
- `TestNewMatch`: Valida criação do storage de match
- `TestMatch_InsertMatch`: Valida inserção no banco
- `TestMatch_GetMatches`: Valida busca de matches agregados
- `TestMatch_GetMatchesInfo`: Valida busca de informações detalhadas
- `TestMatch_Close`: Valida fechamento de conexão

**Nota**: Usa banco SQLite em memória (`:memory:`)

### internal/procspy/watcher (Coverage: 73.8%)

#### watcher_test.go
**Escopo**: Testes do componente Watcher que monitora o Client

**Funções Testadas**:
- `TestNewWatcher`: Valida criação de instância Watcher
- `TestWatcher_check`: Valida verificação de health check
  - Cenário: Procspy up (status 200)
  - Cenário: Procspy down sem comando de start
  - Cenário: Procspy down com comando de start
  - Cenário: Erro de conexão
- `TestWatcher_Stop`: Valida parada do watcher
- `TestExecuteCommand`: Valida execução de comandos
  - Cenário: Comando inválido

**Limitações**:
- Função `Run()` não testada (loop infinito bloqueante)

## Troubleshooting

### Problema: "bc: command not found" no Linux/macOS

**Solução**:
```bash
# Ubuntu/Debian
sudo apt-get install bc

# macOS
brew install bc

# Ou execute testes diretamente com Go
go test -v ./...
```

### Problema: Script não executa no Windows

**Solução**: Use Git Bash ou WSL, ou execute diretamente:
```powershell
go test -v -race -coverprofile=coverage.out ./...
```

### Problema: "permission denied" ao executar test.sh

**Solução**:
```bash
chmod +x test.sh
./test.sh
```

### Problema: Testes falham com "database is locked"

**Causa**: Testes usam banco em memória, não deve ocorrer

**Solução**: Verifique se não há processos do Go travados:
```bash
# Linux/macOS
pkill -9 go

# Windows PowerShell
taskkill /F /IM go.exe
```

## Integração com CI/CD

### GitHub Actions (Exemplo)

```yaml
name: Tests
on: [push, pull_request]
jobs:
  test:
    strategy:
      matrix:
        os: [ubuntu-latest, macos-latest, windows-latest]
        go: ['1.16', '1.17', '1.18']
    runs-on: ${{ matrix.os }}
    steps:
      - uses: actions/checkout@v2
      - uses: actions/setup-go@v2
        with:
          go-version: ${{ matrix.go }}
      - name: Run tests
        run: go test -v -race -coverprofile=coverage.out ./...
      - name: Upload coverage
        uses: codecov/codecov-action@v2
        with:
          files: ./coverage.out
```

## Estratégias Utilizadas para Atingir Coverage

### Pacote config (95.5%)
- ✅ Testes completos de parsing JSON
- ✅ Testes de validação de defaults
- ✅ Testes de serialização round-trip
- ✅ Testes de edge cases (valores negativos, strings vazias, etc.)
- ⚠ Limitação: Linhas de tratamento de erro em `json.MarshalIndent` não são testáveis com structs normais (requerem tipos customizados que falham no marshalling)

### Pacote domain (85.7%)
- ✅ Testes completos de modelos de dados
- ✅ Testes de serialização JSON
- ✅ Testes de regex matching
- ✅ Testes de validação de limites
- ⚠ Algumas linhas de log não cobertas (não críticas)

### Pacote storage (74.8%)
- ✅ Uso de banco SQLite em memória para testes
- ✅ Testes de CRUD completo
- ✅ Testes de queries complexas
- ⚠ Algumas linhas de tratamento de erro de conexão não cobertas

### Pacote watcher (73.8%)
- ✅ Testes de health check com httptest
- ✅ Testes de execução de comandos
- ✅ Testes de cenários de erro
- ⚠ Função `Run()` não testada (loop infinito bloqueante)

### Pacote service (68.9%)
- ✅ Testes de lógica de negócio
- ✅ Testes de validação de dados
- ⚠ Algumas integrações com storage não totalmente cobertas

### Pacote client (54.0%)
- ✅ Testes de funções auxiliares
- ✅ Testes de arredondamento e formatação
- ❌ Função `Run()` não testada (loop infinito com scan de processos)
- ❌ Funções de scan de processos não testadas (dependem de OS)

### Pacote handlers (43.8%)
- ✅ Testes de criação de handlers
- ✅ Testes básicos de endpoints
- ❌ Muitos cenários de erro HTTP não cobertos
- ❌ Validações de request body não totalmente cobertas

### Pacote server (35.6%)
- ✅ Testes de inicialização de servidor
- ✅ Testes de criação de services e handlers
- ❌ Função `Start()` não testada (servidor HTTP bloqueante com signal handling)

## Limitações Conhecidas de Coverage

### Funções Não Testáveis em Testes Unitários

1. **Loops Infinitos Bloqueantes**:
   - `client.Run()`: Loop infinito que escaneia processos
   - `watcher.Run()`: Loop infinito que monitora health check
   - `server.Start()`: Servidor HTTP bloqueante com signal handling

2. **Interações com Sistema Operacional**:
   - Scan de processos em `client.go`
   - Execução de comandos do sistema
   - Tratamento de sinais (SIGINT, SIGTERM)

3. **Erros de Marshalling JSON**:
   - Linhas de tratamento de erro em `ToJson()` que só falham com tipos customizados

### Recomendações para Testes de Integração

Para cobrir as funções não testáveis em testes unitários, recomenda-se:

1. **Testes de Integração E2E**:
   - Iniciar servidor real em porta aleatória
   - Executar client por tempo limitado
   - Validar comunicação entre componentes

2. **Testes Manuais**:
   - Validar comportamento em diferentes sistemas operacionais
   - Testar shutdown graceful com sinais reais
   - Validar scan de processos reais

## Próximos Passos

Para melhorar a cobertura de testes:

1. **Handlers** (43.8% → 90%): 
   - Adicionar testes HTTP completos com httptest
   - Testar todos os cenários de erro (400, 401, 500)
   - Testar validação de request body
   - Testar diferentes content-types

2. **Client** (54.0% → 70%):
   - Adicionar testes de funções auxiliares restantes
   - Documentar limitações de teste de `Run()`
   - Considerar refatoração para tornar mais testável

3. **Service** (68.9% → 95%):
   - Adicionar testes de cenários de erro
   - Melhorar cobertura de integrações com storage

4. **Server** (35.6% → 90%):
   - Documentar limitações de teste de `Start()`
   - Considerar refatoração para separar lógica de roteamento

## Referências

- [Go Testing Package](https://pkg.go.dev/testing)
- [Go Coverage Tool](https://go.dev/blog/cover)
- [Table Driven Tests](https://github.com/golang/go/wiki/TableDrivenTests)
- [httptest Package](https://pkg.go.dev/net/http/httptest)
