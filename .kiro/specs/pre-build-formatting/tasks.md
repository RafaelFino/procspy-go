# Implementation Plan

- [x] 1. Implementar função de logging com timestamp
  - Adicionar variável BOLD na seção de cores do build.sh
  - Criar função log_msg() que adiciona timestamp formatado em verde e negrito
  - Adicionar variável BOLD na seção de cores do test.sh
  - Criar função log_msg() no test.sh com mesma implementação
  - _Requirements: 4.1, 4.2, 4.5_

- [x] 2. Atualizar build.sh para usar log_msg
  - Substituir todos os echo -e por log_msg em todo o script build.sh
  - Preservar echo vazios (sem -e) para quebras de linha
  - Testar que todas as cores das mensagens são preservadas
  - _Requirements: 4.3, 4.5_

- [x] 3. Atualizar test.sh para usar log_msg
  - Substituir todos os echo -e por log_msg em todo o script test.sh
  - Preservar echo vazios (sem -e) para quebras de linha
  - Testar que todas as cores das mensagens são preservadas
  - _Requirements: 4.4, 4.5_

- [x] 4. Implementar função de formatação automática
  - Criar função auto_format() no build.sh que executa go fmt ./...
  - Implementar tratamento de erros com mensagens apropriadas usando log_msg
  - Adicionar retorno de código de saída apropriado
  - _Requirements: 1.1, 1.2, 1.3, 1.4_

- [x] 5. Adicionar flag --no-fmt ao build.sh
  - Adicionar variável NO_FMT=false na seção de parse de argumentos
  - Implementar case para --no-fmt no bloco while de parse
  - Atualizar função usage() para documentar a nova flag
  - _Requirements: 3.1, 3.2, 3.3_

- [x] 6. Integrar formatação automática no fluxo de build
  - Adicionar chamada a auto_format() antes de run_quality_checks()
  - Implementar lógica condicional: executar apenas se não for --build-only e não for --no-fmt
  - Adicionar tratamento de erro que aborta build se formatação falhar
  - _Requirements: 1.1, 2.1, 2.2, 2.3, 2.4_

- [x] 7. Validar implementação completa
  - Testar build.sh com código não formatado (deve formatar automaticamente)
  - Testar flag --no-fmt (deve pular formatação)
  - Testar todas as flags existentes para garantir compatibilidade
  - Verificar que timestamps aparecem em todas as mensagens de build.sh e test.sh
  - Verificar formato correto dos timestamps: [yyyy-mm-dd HH:MM:SS] em verde e negrito
  - Verificar que cores das mensagens são preservadas após timestamp
  - _Requirements: 1.1, 1.2, 1.3, 1.4, 2.1, 2.2, 2.3, 2.4, 3.1, 3.2, 3.3, 4.1, 4.2, 4.3, 4.4, 4.5_
