# Plano de Implementação - Testes Unitários

- [x] 1. Criar infraestrutura de testes
  - Criar script test.sh para execução automatizada de testes
  - Integrar test.sh com build.sh existente
  - Configurar flags --test-only e --build-only no build.sh
  - Adicionar relatório detalhado de coverage por pacote
  - Configurar meta de coverage de 99% com mínimos variáveis
  - _Requirements: 3.1, 3.2, 3.3, 3.4, 3.5, 4.1, 4.2, 4.3, 4.4, 6.1, 6.2, 6.5, 6.6_

- [x] 2. Implementar testes para pacote domain
  - [x] 2.1 Criar testes para domain/target.go
    - Testar NewTargetList, TargetListFromJson
    - Testar Target.Match com diferentes padrões regex
    - Testar Target.AddElapsed, SetElapsed, ResetElapsed
    - Testar Target.CheckLimit e CheckWarning
    - Testar serialização ToJson e ToLog
    - _Requirements: 1.1, 1.3, 5.1, 5.5, 7.4_

  - [x] 2.2 Criar testes para domain/match.go
    - Testar NewMatch com diferentes parâmetros
    - Testar MatchFromJson e MatchListFromJson
    - Testar serialização ToJson e ToLog
    - Validar casos de erro em parsing JSON
    - _Requirements: 1.1, 1.3, 5.1, 5.5, 7.4_

  - [x] 2.3 Criar testes para domain/command.go
    - Testar NewCommand com diferentes parâmetros
    - Testar CommandFromJson
    - Testar serialização ToJson e ToLog
    - Validar casos de erro em parsing JSON
    - _Requirements: 1.1, 1.3, 5.1, 5.5, 7.4_

- [x] 3. Implementar testes para pacote config
  - [x] 3.1 Criar testes para config/client.go
    - Testar NewConfig e SetDefaults
    - Testar ClientConfigFromJson com JSON válido e inválido
    - Testar ClientConfigFromFile com arquivo válido e inexistente
    - Testar ToJson
    - _Requirements: 1.1, 1.3, 5.1, 5.5_

  - [x] 3.2 Criar testes para config/server.go
    - Testar NewServer
    - Testar ServerConfigFromJson com JSON válido e inválido
    - Testar ServerConfigFromFile com arquivo válido e inexistente
    - Testar ToJson
    - _Requirements: 1.1, 1.3, 5.1, 5.5_

  - [x] 3.3 Criar testes para config/watcher.go
    - Testar parsing de configuração do watcher
    - Testar valores padrão
    - Testar serialização
    - _Requirements: 1.1, 1.3, 5.1, 5.5_

- [x] 4. Implementar testes para pacote storage
  - [x] 4.1 Criar testes para storage/dbconn.go
    - Testar NewDbConnection
    - Testar GetConn com banco em memória
    - Testar Close
    - Testar Exec com queries válidas e inválidas
    - _Requirements: 1.1, 1.3, 5.1, 5.5, 7.3_

  - [x] 4.2 Criar testes para storage/match.go
    - Testar Init (criação de tabela)
    - Testar InsertMatch com dados válidos
    - Testar GetMatches e GetMatchesInfo
    - Testar Close
    - Usar banco SQLite em memória (:memory:)
    - _Requirements: 1.1, 1.3, 5.1, 5.5, 7.3_

  - [x] 4.3 Criar testes para storage/command.go
    - Testar Init (criação de tabela)
    - Testar InsertCommand com dados válidos
    - Testar GetCommands
    - Testar Close
    - Usar banco SQLite em memória (:memory:)
    - _Requirements: 1.1, 1.3, 5.1, 5.5, 7.3_

- [x] 5. Implementar testes para pacote service
  - [x] 5.1 Criar testes para service/match.go
    - Testar NewMatch
    - Testar InsertMatch com validação de MATCH_MAX_ELAPSED
    - Testar GetMatches e GetMatchesInfo
    - Testar Close
    - Usar banco em memória para storage
    - _Requirements: 1.1, 1.3, 5.1, 5.5, 7.6_

  - [x] 5.2 Criar testes para service/command.go
    - Testar NewCommand
    - Testar InsertCommand
    - Testar GetCommands
    - Testar Close
    - Usar banco em memória para storage
    - _Requirements: 1.1, 1.3, 5.1, 5.5_

  - [x] 5.3 Criar testes para service/target.go
    - Testar busca de targets
    - Testar parsing de configurações
    - Testar cache de targets
    - _Requirements: 1.1, 1.3, 5.1, 5.5_

  - [x] 5.4 Criar testes para service/user.go
    - Testar validação de usuários
    - Testar busca de configurações de usuário
    - _Requirements: 1.1, 1.3, 5.1, 5.5_

- [x] 6. Implementar testes para pacote handlers
  - [x] 6.1 Criar testes para handlers/healthcheck.go
    - Testar NewHealthcheck
    - Testar GetStatus com httptest.ResponseRecorder
    - Validar response status 200 e JSON correto
    - _Requirements: 1.1, 1.3, 5.2, 5.5, 7.2, 7.5_

  - [x] 6.2 Criar testes para handlers/match.go
    - Testar NewMatch
    - Testar InsertMatch com request válido
    - Testar InsertMatch com request inválido
    - Testar validação de usuário
    - Usar httptest para simular requests
    - _Requirements: 1.1, 1.3, 5.2, 5.5, 7.2_

  - [x] 6.3 Criar testes para handlers/command.go
    - Testar NewCommand
    - Testar InsertCommand com request válido
    - Testar InsertCommand com request inválido
    - Testar validação de usuário
    - Usar httptest para simular requests
    - _Requirements: 1.1, 1.3, 5.2, 5.5, 7.2_

  - [x] 6.4 Criar testes para handlers/target.go
    - Testar NewTarget
    - Testar GetTargets com usuário válido
    - Testar GetTargets com usuário inválido
    - Testar integração com service layer
    - Usar httptest para simular requests
    - _Requirements: 1.1, 1.3, 5.2, 5.5, 7.2_

  - [x] 6.5 Criar testes para handlers/report.go
    - Testar NewReport
    - Testar GetReport com diferentes filtros
    - Testar FormatInterval
    - Usar httptest para simular requests
    - _Requirements: 1.1, 1.3, 5.2, 5.5, 7.2_

  - [x] 6.6 Criar testes para handlers/util.go
    - Testar ValidateUser com usuário válido e inválido
    - Testar funções auxiliares de handlers
    - _Requirements: 1.1, 1.3, 5.5_

- [x] 7. Implementar testes para pacote client
  - [x] 7.1 Criar testes para client/client.go
    - Testar NewSpy
    - Testar updateTargets com httptest server
    - Testar postMatch e postCommand
    - Testar consumeBuffers
    - Testar kill com processos mock
    - Testar roundFloat e executeCommand
    - _Requirements: 1.1, 1.3, 5.1, 5.5, 7.1_

- [x] 8. Implementar testes para pacote server
  - [x] 8.1 Criar testes para server/server.go
    - Testar NewServer
    - Testar initServices
    - Validar criação de handlers e services
    - _Requirements: 1.1, 1.3, 5.1, 5.5, 7.2_

- [x] 9. Implementar testes para pacote watcher
  - [x] 9.1 Criar testes para watcher/watcher.go
    - Testar NewWatcher
    - Testar check com httptest server (client up)
    - Testar check com httptest server (client down)
    - Testar executeCommand
    - _Requirements: 1.1, 1.3, 5.1, 5.5, 7.5_

- [x] 10. Criar documentação de testes
  - Criar arquivo test.md na raiz do projeto
  - Documentar escopo de cada arquivo de teste
  - Listar funções testadas por arquivo
  - Incluir estatísticas de coverage
  - Adicionar convenções e boas práticas
  - _Requirements: 2.1, 2.2, 2.3, 2.4_

- [x] 11. Analisar coverage atual e identificar gaps
  - Executar test.sh e gerar relatório de coverage detalhado
  - Gerar relatório HTML de coverage por arquivo
  - Identificar todos os arquivos com coverage abaixo de 99%
  - Listar linhas específicas não cobertas em cada arquivo
  - Priorizar arquivos por impacto (domain, config, storage primeiro)
  - _Requirements: 6.1, 6.2, 6.5, 6.6_

- [x] 12. Melhorar coverage do pacote domain para 99%
  - [x] 12.1 Analisar gaps de coverage em domain/target.go
    - Identificar branches não testados
    - Identificar edge cases não cobertos
    - _Requirements: 1.1, 5.5, 6.4, 7.4_
  
  - [x] 12.2 Adicionar testes para atingir 99% em domain/target_test.go
    - Adicionar testes para todos os branches não cobertos
    - Adicionar testes para edge cases (nil, vazio, valores extremos)
    - Adicionar testes para casos de erro não cobertos
    - _Requirements: 1.1, 5.5, 6.4, 7.4_
  
  - [x] 12.3 Analisar gaps de coverage em domain/match.go
    - Identificar branches não testados
    - Identificar edge cases não cobertos
    - _Requirements: 1.1, 5.5, 6.4, 7.4_
  
  - [x] 12.4 Adicionar testes para atingir 99% em domain/match_test.go
    - Adicionar testes para todos os branches não cobertos
    - Adicionar testes para edge cases
    - Adicionar testes para casos de erro não cobertos
    - _Requirements: 1.1, 5.5, 6.4, 7.4_
  
  - [x] 12.5 Analisar gaps de coverage em domain/command.go
    - Identificar branches não testados
    - Identificar edge cases não cobertos
    - _Requirements: 1.1, 5.5, 6.4, 7.4_
  
  - [x] 12.6 Adicionar testes para atingir 99% em domain/command_test.go
    - Adicionar testes para todos os branches não cobertos
    - Adicionar testes para edge cases
    - Adicionar testes para casos de erro não cobertos
    - _Requirements: 1.1, 5.5, 6.4, 7.4_

- [x] 13. Melhorar coverage do pacote config para 99%
  - [x] 13.1 Analisar e melhorar coverage de config/client_test.go
    - Identificar gaps de coverage
    - Adicionar testes para branches não cobertos
    - Testar todos os cenários de erro possíveis
    - _Requirements: 1.1, 5.5, 6.4_
  
  - [x] 13.2 Analisar e melhorar coverage de config/server_test.go
    - Identificar gaps de coverage
    - Adicionar testes para branches não cobertos
    - Testar todos os cenários de erro possíveis
    - _Requirements: 1.1, 5.5, 6.4_
  
  - [x] 13.3 Analisar e melhorar coverage de config/watcher_test.go
    - Identificar gaps de coverage
    - Adicionar testes para branches não cobertos
    - Testar todos os cenários de erro possíveis
    - _Requirements: 1.1, 5.5, 6.4_

- [x] 14. Melhorar coverage do pacote storage para 95%+
  - [x] 14.1 Analisar e melhorar coverage de storage/dbconn_test.go
    - Identificar gaps de coverage
    - Adicionar testes para cenários de erro de conexão
    - Testar comportamento com queries inválidas
    - Testar comportamento com conexão fechada
    - _Requirements: 1.1, 5.5, 6.4, 7.3_
  
  - [x] 14.2 Analisar e melhorar coverage de storage/match_test.go
    - Identificar gaps de coverage
    - Adicionar testes para cenários de erro
    - Testar operações com dados inválidos
    - Testar queries complexas
    - _Requirements: 1.1, 5.5, 6.4, 7.3_
  
  - [x] 14.3 Analisar e melhorar coverage de storage/command_test.go
    - Identificar gaps de coverage
    - Adicionar testes para cenários de erro
    - Testar operações com dados inválidos
    - _Requirements: 1.1, 5.5, 6.4, 7.3_

- [x] 15. Melhorar coverage do pacote service para 95%+
  - [x] 15.1 Analisar e melhorar coverage de service/match_test.go
    - Identificar gaps de coverage
    - Adicionar testes para validação de MATCH_MAX_ELAPSED
    - Testar cenários de erro de storage
    - _Requirements: 1.1, 5.5, 6.4, 7.6_
  
  - [x] 15.2 Analisar e melhorar coverage de service/command_test.go
    - Identificar gaps de coverage
    - Adicionar testes para cenários de erro
    - _Requirements: 1.1, 5.5, 6.4_
  
  - [x] 15.3 Analisar e melhorar coverage de service/target_test.go
    - Identificar gaps de coverage
    - Adicionar testes para cache de targets
    - Testar cenários de erro de parsing
    - _Requirements: 1.1, 5.5, 6.4_
  
  - [x] 15.4 Analisar e melhorar coverage de service/user_test.go
    - Identificar gaps de coverage
    - Adicionar testes para validação de usuários
    - Testar cenários de erro
    - _Requirements: 1.1, 5.5, 6.4_

- [x] 16. Melhorar coverage do pacote handlers para 90%+
  - [x] 16.1 Analisar e melhorar coverage de handlers/healthcheck_test.go
    - Identificar gaps de coverage
    - Adicionar testes para diferentes métodos HTTP
    - Testar comportamento com headers variados
    - _Requirements: 1.1, 5.2, 5.5, 7.2, 7.5_
  
  - [x] 16.2 Analisar e melhorar coverage de handlers/match_test.go
    - Identificar gaps de coverage
    - Adicionar testes para todos os cenários de erro
    - Testar validação de body vazio
    - Testar diferentes content-types
    - _Requirements: 1.1, 5.2, 5.5, 7.2_
  
  - [x] 16.3 Analisar e melhorar coverage de handlers/command_test.go
    - Identificar gaps de coverage
    - Adicionar testes para todos os cenários de erro
    - Testar validação de parâmetros
    - _Requirements: 1.1, 5.2, 5.5, 7.2_
  
  - [x] 16.4 Analisar e melhorar coverage de handlers/target_test.go
    - Identificar gaps de coverage
    - Adicionar testes para cenários de erro de service
    - Testar query parameters variados
    - _Requirements: 1.1, 5.2, 5.5, 7.2_
  
  - [x] 16.5 Analisar e melhorar coverage de handlers/report_test.go
    - Identificar gaps de coverage
    - Adicionar testes para todos os filtros possíveis
    - Testar FormatInterval com diferentes inputs
    - _Requirements: 1.1, 5.2, 5.5, 7.2_
  
  - [x] 16.6 Analisar e melhorar coverage de handlers/util_test.go
    - Identificar gaps de coverage
    - Adicionar testes para todos os cenários de ValidateUser
    - _Requirements: 1.1, 5.5_

- [x] 17. Melhorar coverage do pacote server para 90%+
  - [x] 17.1 Analisar e melhorar coverage de server/server_test.go
    - Identificar gaps de coverage
    - Adicionar testes para inicialização completa
    - Testar cenários de erro na criação de services
    - _Requirements: 1.1, 5.5, 7.2_

- [x] 18. Melhorar coverage dos pacotes client e watcher (meta: 70%+)
  - [x] 18.1 Analisar e melhorar coverage de client/client_test.go
    - Identificar gaps de coverage viáveis
    - Adicionar testes para cenários não cobertos
    - Documentar limitações de coverage devido a interações com OS
    - _Requirements: 1.1, 5.5, 7.1_
  
  - [x] 18.2 Analisar e melhorar coverage de watcher/watcher_test.go
    - Identificar gaps de coverage viáveis
    - Adicionar testes para cenários não cobertos
    - Documentar limitações de coverage devido a interações com OS
    - _Requirements: 1.1, 5.5, 7.5_

- [x] 19. Validar cobertura final e qualidade
  - Executar test.sh e validar que todos os testes passam
  - Gerar relatório de coverage final
  - Verificar se coverage alvo de 99% foi atingido (ou mínimos aceitáveis)
  - Gerar relatório detalhado por pacote
  - Validar que build.sh integra corretamente com testes
  - Testar flags --test-only e --build-only
  - _Requirements: 3.5, 4.1, 4.2, 4.3, 4.4, 6.1, 6.2, 6.3, 6.4, 6.5, 6.6_

- [x] 20. Atualizar documentação de testes
  - Atualizar test.md com estatísticas finais de coverage
  - Documentar estratégias usadas para atingir 99%
  - Listar limitações conhecidas de coverage
  - Adicionar guia de manutenção de coverage alto
  - _Requirements: 2.1, 2.2, 2.3, 2.4_

- [x] 21. Validar compatibilidade cross-platform
  - Executar testes em Linux
  - Validar que testes não usam comandos específicos de OS
  - Documentar requisitos para execução em Windows/macOS
  - _Requirements: 8.1, 8.2, 8.3, 8.4_
