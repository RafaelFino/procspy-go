# Changelog

Todas as mudanças notáveis neste projeto serão documentadas neste arquivo.

O formato é baseado em [Keep a Changelog](https://keepachangelog.com/pt-BR/1.0.0/),
e este projeto adere ao [Semantic Versioning](https://semver.org/lang/pt-BR/).

## [Não Lançado]

### Alterado
- **Formato de Logs**: Todas as mensagens de log foram padronizadas para usar o formato `[package.função]` em vez de apenas `[função]`
  - Melhora a rastreabilidade e identificação da origem dos logs
  - Facilita debugging em sistemas distribuídos
  - Exemplos:
    - `[Start]` → `[watcher.Start]`
    - `[InsertCommand]` → `[handlers.Command.InsertCommand]`
    - `[GetConn]` → `[storage.DbConnection.GetConn]`
  - Pacotes afetados: watcher, server, config, domain, handlers, service, storage
  - Data: 2025-11-13

### Melhorado
- **Mensagens de Log**: Todas as mensagens de log foram revisadas para melhor clareza e contexto
  - Gramática corrigida
  - Contexto adicional incluído (URLs, valores, nomes de usuário, etc.)
  - Uso consistente de `%v` para formatação de erros
  - Mensagens mais descritivas e acionáveis
  - Data: 2025-11-13

## [Versões Anteriores]

Para histórico de versões anteriores, consulte os commits do Git.
