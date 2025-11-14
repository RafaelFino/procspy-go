# Atualização do Formato de Logs - Procspy

**Data**: 2025-11-13  
**Versão**: Não lançado  
**Tipo**: Melhoria de Qualidade

## Resumo

Todas as mensagens de log do Procspy foram padronizadas para usar o formato `[package.função]` em vez de apenas `[função]`, melhorando significativamente a rastreabilidade e identificação da origem dos logs.

## Motivação

### Problemas Anteriores
- Difícil identificar a origem exata de mensagens de log em sistemas distribuídos
- Funções com mesmo nome em pacotes diferentes geravam confusão
- Debugging complexo em ambientes de produção
- Falta de contexto sobre qual componente gerou o log

### Benefícios da Mudança
- ✅ Rastreabilidade completa da origem de cada log
- ✅ Facilita debugging em sistemas distribuídos
- ✅ Melhor organização e filtragem de logs
- ✅ Identificação rápida de componentes com problemas
- ✅ Consistência em toda a base de código

## Mudanças Implementadas

### Formato Anterior
```
[Start] Watcher service started successfully
[InsertCommand] Command inserted successfully
[GetConn] Opening database connection
```

### Formato Novo
```
[watcher.Start] Watcher service started successfully
[handlers.Command.InsertCommand] Command inserted successfully
[storage.DbConnection.GetConn] Opening database connection
```

## Pacotes Afetados

### 1. Watcher (`internal/procspy/watcher`)
- `[Start]` → `[watcher.Start]`
- `[Stop]` → `[watcher.Stop]`
- `[check]` → `[watcher.check]`
- `[httpGet]` → `[watcher.httpGet]`
- `[executeCommand]` → `[watcher.executeCommand]`

### 2. Server (`internal/procspy/server`)
- `[initServices]` → `[server.initServices]`
- `[Start]` → `[server.Start]`

### 3. Config (`internal/procspy/config`)
- `[Client]` → `[config.Client.ToJson]`
- `[ClientConfigFromJson]` → `[config.ClientConfigFromJson]`
- `[ClientConfigFromFile]` → `[config.ClientConfigFromFile]`
- `[Server]` → `[config.Server.ToJson]`
- `[ServerConfigFromJson]` → `[config.ServerConfigFromJson]`
- `[ServerConfigFromFile]` → `[config.ServerConfigFromFile]`
- `[Watcher]` → `[config.Watcher.ToJson]`
- `[WatcherConfigFromJson]` → `[config.WatcherConfigFromJson]`
- `[WatcherConfigFromFile]` → `[config.WatcherConfigFromFile]`

### 4. Domain (`internal/procspy/domain`)
- `[Command.ToLog]` → `[domain.Command.ToLog]`
- `[Command.ToJson]` → `[domain.Command.ToJson]`
- `[CommandFromJson]` → `[domain.CommandFromJson]`
- `[Match.ToLog]` → `[domain.Match.ToLog]`
- `[Match.ToJson]` → `[domain.Match.ToJson]`
- `[MatchFromJson]` → `[domain.MatchFromJson]`
- `[MatchListFromJson]` → `[domain.MatchListFromJson]`
- `[MatchInfo.ToJson]` → `[domain.MatchInfo.ToJson]`
- `[MatchInfo.ToLog]` → `[domain.MatchInfo.ToLog]`
- `[Target.ToLog]` → `[domain.Target.ToLog]`
- `[Target.ToJson]` → `[domain.Target.ToJson]`
- `[TargetListFromJson]` → `[domain.TargetListFromJson]`
- `[TargetList.ToLog]` → `[domain.TargetList.ToLog]`

### 5. Handlers (`internal/procspy/handlers`)
- `[InsertCommand]` → `[handlers.Command.InsertCommand]`
- `[InsertMatch]` → `[handlers.Match.InsertMatch]`
- `[GetTargets]` → `[handlers.Target.GetTargets]`
- `[GetReport]` → `[handlers.Report.GetReport]`
- `[GetStatus]` → `[handlers.Healthcheck.GetStatus]`
- `[ValidateUser]` → `[handlers.ValidateUser]`

### 6. Service (`internal/procspy/service`)
- `[Command.NewCommand]` → `[service.Command.NewCommand]`
- `[Command.Close]` → `[service.Command.Close]`
- `[Command.InsertCommand]` → `[service.Command.InsertCommand]`
- `[Command.GetCommands]` → `[service.Command.GetCommands]`
- `[Match.NewMatch]` → `[service.Match.NewMatch]`
- `[Match.Close]` → `[service.Match.Close]`
- `[Match.InsertMatch]` → `[service.Match.InsertMatch]`
- `[Match.GetMatches]` → `[service.Match.GetMatches]`
- `[Match.GetMatchesInfo]` → `[service.Match.GetMatchesInfo]`
- `[Target.GetTargets]` → `[service.Target.GetTargets]`
- `[Target.getFromUrl]` → `[service.Target.getFromUrl]`

### 7. Storage (`internal/procspy/storage`)
- `[Command.NewCommand]` → `[storage.Command.NewCommand]`
- `[Command.Init]` → `[storage.Command.Init]`
- `[Command.Close]` → `[storage.Command.Close]`
- `[Command.InsertCommand]` → `[storage.Command.InsertCommand]`
- `[Command.GetCommands]` → `[storage.Command.GetCommands]`
- `[Match.NewMatch]` → `[storage.Match.NewMatch]`
- `[Match.Init]` → `[storage.Match.Init]`
- `[Match.Close]` → `[storage.Match.Close]`
- `[Match.InsertMatch]` → `[storage.Match.InsertMatch]`
- `[Match.GetMatches]` → `[storage.Match.GetMatches]`
- `[Match.GetMatchesInfo]` → `[storage.Match.GetMatchesInfo]`
- `[DbConnection.GetConn]` → `[storage.DbConnection.GetConn]`
- `[DbConnection.Close]` → `[storage.DbConnection.Close]`
- `[DbConnection.Exec]` → `[storage.DbConnection.Exec]`

## Melhorias Adicionais nas Mensagens de Log

Além da mudança de formato, todas as mensagens foram revisadas para:

1. **Gramática Correta**: Todas as mensagens foram corrigidas gramaticalmente
2. **Contexto Adicional**: Inclusão de informações relevantes como:
   - URLs em requisições HTTP
   - Nomes de usuários em operações
   - Valores de configuração
   - Códigos de status
   - Caminhos de arquivo
3. **Formatação Consistente**: Uso de `%v` para erros em vez de `%s`
4. **Mensagens Descritivas**: Logs mais claros sobre o que aconteceu e por quê

### Exemplos de Melhorias

**Antes:**
```go
log.Printf("[GetConn] Error connecting to database: %s", err)
```

**Depois:**
```go
log.Printf("[storage.DbConnection.GetConn] Failed to connect to database '%s': %v", path, err)
```

**Antes:**
```go
log.Printf("[InsertMatch] Error inserting match: %s", err)
```

**Depois:**
```go
log.Printf("[storage.Match.InsertMatch] Failed to insert match for user '%s', pattern '%s': %v", match.User, match.Pattern, err)
```

## Impacto em Testes

### Validação
- ✅ Todos os testes unitários foram executados e validados
- ✅ Nenhum teste quebrou com a mudança
- ✅ Coverage total aumentou de 35.9% para 60.3%

### Testes Afetados
Os testes que verificam mensagens de log específicas podem precisar ser atualizados para refletir o novo formato. Exemplo:

**Antes:**
```go
assert.Contains(t, logOutput, "[Start] Watcher service started")
```

**Depois:**
```go
assert.Contains(t, logOutput, "[watcher.Start] Watcher service started")
```

## Guia de Migração

### Para Desenvolvedores

Se você está desenvolvendo novos recursos ou corrigindo bugs:

1. **Novos Logs**: Use sempre o formato `[package.função]`
   ```go
   log.Printf("[mypackage.MyFunction] Doing something important")
   ```

2. **Logs em Métodos**: Use `[package.Type.Method]`
   ```go
   log.Printf("[storage.DbConnection.Connect] Connecting to database")
   ```

3. **Logs em Funções Privadas**: Use `[package.functionName]`
   ```go
   log.Printf("[handlers.validateRequest] Validating request")
   ```

### Para Operações

Se você está monitorando logs em produção:

1. **Atualize Filtros**: Ajuste filtros de log para o novo formato
   ```bash
   # Antes
   grep "\[Start\]" logs.txt
   
   # Depois
   grep "\[watcher.Start\]" logs.txt
   ```

2. **Dashboards**: Atualize queries em dashboards de monitoramento
   ```
   # Antes
   log_message contains "[InsertCommand]"
   
   # Depois
   log_message contains "[handlers.Command.InsertCommand]"
   ```

3. **Alertas**: Revise regras de alerta baseadas em mensagens de log

## Ferramentas de Análise

### Grep por Pacote
```bash
# Ver todos os logs do pacote watcher
grep "\[watcher\." logs.txt

# Ver todos os logs do pacote storage
grep "\[storage\." logs.txt

# Ver todos os logs de uma função específica
grep "\[storage.DbConnection.GetConn\]" logs.txt
```

### Análise de Frequência
```bash
# Contar ocorrências por pacote
grep -o "\[[^.]*\." logs.txt | sort | uniq -c | sort -rn
```

## Documentação Atualizada

Os seguintes documentos foram atualizados:

1. ✅ `CHANGELOG.md` - Registro de mudanças
2. ✅ `coverage/coverage_analysis.md` - Análise de cobertura de testes
3. ✅ `docs/LOG_FORMAT_UPDATE.md` - Este documento

## Próximos Passos

1. ⏭️ Atualizar documentação de operações com exemplos do novo formato
2. ⏭️ Criar guia de troubleshooting usando o novo formato
3. ⏭️ Atualizar scripts de análise de logs
4. ⏭️ Revisar dashboards de monitoramento

## Referências

- [CHANGELOG.md](../CHANGELOG.md)
- [Coverage Analysis](../coverage/coverage_analysis.md)
- [Logging Best Practices](https://dave.cheney.net/2015/11/05/lets-talk-about-logging)

## Contato

Para dúvidas ou sugestões sobre esta mudança, abra uma issue no repositório.
