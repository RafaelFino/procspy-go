# Implementation Plan - Build Timing Metrics

- [x] 1. Implementar funções de medição de tempo
  - Criar função format_duration() que aceita nanosegundos e retorna string formatada
  - Implementar lógica de seleção de unidade (ns < 1ms, ms < 1s, s >= 1s)
  - Adicionar tratamento para valores com 2 casas decimais usando bc
  - Criar função start_timer() que armazena timestamp em variável global
  - Criar função end_timer() que calcula duração e retorna string formatada
  - Adicionar tratamento de erro para timer não iniciado
  - Testar funções com valores conhecidos (sleep)
  - _Requirements: 1.1, 1.2, 2.1, 2.2, 2.3, 2.4, 6.1, 6.2, 6.3_

- [x] 2. Integrar timing no build.sh
  - [x] 2.1 Adicionar timing à função auto_format()
    - Adicionar start_timer("format") no início
    - Adicionar end_timer("format") e exibir duração no sucesso
    - Adicionar end_timer("format") e exibir duração no erro
    - _Requirements: 1.1, 1.2, 1.3, 3.1, 5.1, 5.2, 5.3_
  
  - [x] 2.2 Adicionar timing à função run_quality_checks()
    - Adicionar start_timer("quality") no início
    - Adicionar end_timer("quality") e exibir duração no final
    - _Requirements: 1.1, 1.2, 1.3, 3.2, 5.1, 5.2, 5.3_
  
  - [x] 2.3 Adicionar timing à função build_component()
    - Adicionar start_timer com nome único por componente/plataforma
    - Adicionar end_timer e exibir duração nas mensagens de sucesso/erro
    - _Requirements: 1.1, 1.2, 1.3, 4.1, 5.1, 5.2, 5.3_
  
  - [x] 2.4 Adicionar timing à função build_platform()
    - Adicionar start_timer com nome único por plataforma
    - Adicionar end_timer e exibir duração total da plataforma
    - _Requirements: 1.1, 1.2, 1.3, 4.2, 4.3, 5.1, 5.2, 5.3_
  
  - [x] 2.5 Adicionar timing total do build
    - Adicionar start_timer("total") após parse de argumentos
    - Adicionar end_timer("total") antes do exit final
    - Exibir tempo total de execução no resumo final
    - _Requirements: 1.1, 1.2, 1.3, 3.5, 5.1, 5.2, 5.3_

- [x] 3. Integrar timing no test.sh
  - [x] 3.1 Adicionar funções de timing ao test.sh
    - Copiar format_duration() do build.sh
    - Copiar start_timer() do build.sh
    - Copiar end_timer() do build.sh
    - _Requirements: 1.4, 6.4_
  
  - [x] 3.2 Adicionar timing à função run_tests()
    - Adicionar start_timer("tests") no início
    - Adicionar end_timer("tests") e exibir duração no sucesso
    - Adicionar end_timer("tests") e exibir duração no erro
    - _Requirements: 1.4, 5.4_
  
  - [x] 3.3 Adicionar timing total do test.sh
    - Adicionar start_timer("total") no início do main
    - Adicionar end_timer("total") antes do exit
    - Exibir tempo total de execução
    - _Requirements: 1.4, 5.4_

- [x] 4. Organizar arquivos de coverage
  - [x] 4.1 Atualizar test.sh para usar diretório coverage
    - Adicionar variável COVERAGE_DIR="coverage"
    - Atualizar COVERAGE_FILE para "${COVERAGE_DIR}/coverage.out"
    - Atualizar COVERAGE_HTML para "${COVERAGE_DIR}/coverage.html"
    - Adicionar mkdir -p "$COVERAGE_DIR" no início
    - Atualizar comandos go test para usar novo caminho
    - Atualizar comandos go tool cover para usar novo caminho
    - _Requirements: 7.1, 7.2, 7.4, 7.5_
  
  - [x] 4.2 Atualizar build.sh para limpar diretório coverage
    - Atualizar função clean_build() para remover diretório coverage/
    - _Requirements: 7.6_
  
  - [x] 4.3 Mover arquivo COVERAGE_ANALYSIS.md
    - Mover COVERAGE_ANALYSIS.md da raiz para coverage/
    - Renomear para coverage_analysis.md (minúsculas)
    - _Requirements: 7.3_

- [x] 5. Padronizar nomes de arquivos markdown
  - Renomear CROSS_PLATFORM_TESTING.md para cross_platform_testing.md
  - Renomear PLATFORM_COMPATIBILITY_SUMMARY.md para platform_compatibility_summary.md
  - Renomear README.md para readme.md
  - Atualizar referências nos scripts se houver
  - _Requirements: 8.1, 8.2, 8.3, 8.4_

- [x] 6. Adicionar seção de documentação no readme.md
  - Criar seção "Documentação" no readme.md
  - Adicionar subseção "Testes" com link para test.md
  - Adicionar subseção "Compatibilidade" com links para cross_platform_testing.md e platform_compatibility_summary.md
  - Adicionar subseção "Coverage" com link para coverage/coverage_analysis.md
  - Incluir descrição breve de cada documento
  - _Requirements: 9.1, 9.2, 9.3, 9.4_

- [x] 7. Revisar mensagens de log em arquivos Go
  - [x] 7.1 Revisar logs do pacote client
    - Revisar internal/procspy/client/client.go
    - Corrigir gramática e melhorar clareza das mensagens
    - Adicionar contexto relevante (URLs, valores, erros)
    - Manter formatação existente ([função] mensagem)
    - _Requirements: 10.1, 10.2, 10.3, 10.4, 10.5_
  
  - [x] 7.2 Revisar logs do pacote watcher
    - Revisar internal/procspy/watcher/watcher.go
    - Corrigir gramática e melhorar clareza das mensagens
    - Adicionar contexto relevante
    - Manter formatação existente
    - _Requirements: 10.1, 10.2, 10.3, 10.4, 10.5_
  
  - [x] 7.3 Revisar logs do pacote server
    - Revisar internal/procspy/server/server.go
    - Corrigir gramática e melhorar clareza das mensagens
    - Adicionar contexto relevante
    - Manter formatação existente
    - _Requirements: 10.1, 10.2, 10.3, 10.4, 10.5_
  
  - [x] 7.4 Revisar logs do pacote config
    - Revisar internal/procspy/config/*.go
    - Corrigir gramática e melhorar clareza das mensagens
    - Adicionar contexto relevante (nomes de arquivo, valores)
    - Manter formatação existente
    - _Requirements: 10.1, 10.2, 10.3, 10.4, 10.5_
  
  - [x] 7.5 Revisar logs do pacote domain
    - Revisar internal/procspy/domain/*.go
    - Corrigir gramática e melhorar clareza das mensagens
    - Adicionar contexto relevante
    - Manter formatação existente
    - _Requirements: 10.1, 10.2, 10.3, 10.4, 10.5_
  
  - [x] 7.6 Revisar logs do pacote handlers
    - Revisar internal/procspy/handlers/*.go
    - Corrigir gramática e melhorar clareza das mensagens
    - Adicionar contexto relevante (status codes, endpoints)
    - Manter formatação existente
    - _Requirements: 10.1, 10.2, 10.3, 10.4, 10.5_
  
  - [x] 7.7 Revisar logs do pacote service
    - Revisar internal/procspy/service/*.go
    - Corrigir gramática e melhorar clareza das mensagens
    - Adicionar contexto relevante
    - Manter formatação existente
    - _Requirements: 10.1, 10.2, 10.3, 10.4, 10.5_
  
  - [x] 7.8 Revisar logs do pacote storage
    - Revisar internal/procspy/storage/*.go
    - Corrigir gramática e melhorar clareza das mensagens
    - Adicionar contexto relevante (queries, erros de DB)
    - Manter formatação existente
    - _Requirements: 10.1, 10.2, 10.3, 10.4, 10.5_

- [x] 8. Validar implementação completa
  - Executar build.sh --test-only e verificar tempos exibidos
  - Executar build.sh --all e verificar tempos por plataforma
  - Verificar que arquivos de coverage estão em ./coverage/
  - Verificar que todos os arquivos .md estão em minúsculas
  - Verificar seção de documentação no readme.md
  - Executar testes e verificar que mensagens de log estão corretas
  - Validar que cores e formatação são preservadas
  - _Requirements: 1.1, 1.2, 1.3, 1.4, 2.1, 2.2, 2.3, 2.4, 3.1, 3.2, 3.3, 3.4, 3.5, 4.1, 4.2, 4.3, 4.4, 5.1, 5.2, 5.3, 5.4, 7.1, 7.2, 7.3, 7.4, 7.5, 7.6, 8.1, 8.2, 8.3, 8.4, 9.1, 9.2, 9.3, 9.4, 10.1, 10.2, 10.3, 10.4, 10.5_
