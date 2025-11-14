# Plano de Implementação

- [x] 1. Criar script make-changelog.sh com estrutura base e funções utilitárias
  - Criar arquivo `make-changelog.sh` na raiz do projeto
  - Adicionar shebang `#!/bin/bash` e `set -e`
  - Copiar definições de cores (RED, GREEN, YELLOW, BLUE, BOLD, NC) de build.sh
  - Copiar funções `log_msg()`, `format_duration()`, `start_timer()`, `end_timer()` de build.sh
  - Adicionar permissão de execução com `chmod +x make-changelog.sh`
  - _Requisitos: 3.3, 6.1, 6.2, 6.3, 6.4_

- [x] 2. Implementar funções de análise Git
  - [x] 2.1 Criar função `get_all_commits()`
    - Usar `git log --format="%H|%h|%ad|%s|%an" --date=format:"%Y-%m-%d %H:%M:%S"` (sem --reverse)
    - Retornar commits do mais recente ao mais antigo
    - Adicionar tratamento de erro se não for repositório Git
    - _Requisitos: 1.2, 1.3, 4.2, 4.3_

  - [x] 2.2 Criar função `get_commit_files()`
    - Usar `git show --name-status --format="" <hash>`
    - Retornar lista de arquivos modificados com status (A/M/D)
    - Tratar commits sem arquivos (merges vazios)
    - _Requisitos: 1.3, 2.1, 2.2_

  - [x] 2.3 Criar função `classify_change()`
    - Analisar nomes de arquivos para identificar padrões
    - Classificar como: Testes, Documentação, Configuração, Implementação, Outros
    - Usar grep para detectar extensões e caminhos (_test.go, .md, .sh, internal/, cmd/)
    - _Requisitos: 2.2, 2.3_

  - [x] 2.4 Criar função `generate_summary()`
    - Gerar resumo objetivo em PT-BR baseado no tipo de mudança
    - Contar arquivos por tipo (Go, testes, docs)
    - Focar em resultados alcançados ao invés de detalhes técnicos
    - Manter resumos concisos e didáticos
    - _Requisitos: 2.2, 2.5, 5.3, 5.4, 5.5_

- [x] 3. Implementar funções de geração de markdown
  - [x] 3.1 Criar função `generate_header()`
    - Gerar título "# CHANGELOG"
    - Adicionar descrição do documento
    - Incluir nota sobre geração automática
    - _Requisitos: 1.1, 3.1_

  - [x] 3.2 Criar função `generate_table()`
    - Gerar cabeçalho da tabela com colunas: Data/Hora, Commit, Mensagem, Resumo
    - Iterar sobre commits e adicionar linhas na tabela
    - Formatar commit hash como código inline com backticks
    - Escapar caracteres especiais (pipes) em mensagens
    - _Requisitos: 4.1, 4.2, 4.3, 4.4, 4.5_

  - [x] 3.3 Criar função `generate_mermaid()`
    - Gerar bloco de código mermaid com gitGraph
    - Adicionar cada commit com `commit id: "hash" tag: "data"`
    - Escapar aspas duplas em mensagens de commit
    - Garantir sintaxe Mermaid válida
    - _Requisitos: 4.6, 4.7, 4.8_

  - [x] 3.4 Criar função `generate_details()`
    - Gerar seção "## Detalhes dos Commits"
    - Para cada commit criar subsecção com hash curto e data/hora
    - Incluir: Mensagem, Autor, Tipo, Resumo, Arquivos modificados
    - Adicionar separador horizontal entre commits
    - _Requisitos: 5.1, 5.2, 5.3, 5.4, 5.5_

- [x] 4. Implementar fluxo principal do script
  - [x] 4.1 Adicionar validações iniciais
    - Verificar se está em repositório Git com `git rev-parse --git-dir`
    - Verificar se há commits com `git log -1`
    - Verificar permissão de escrita no diretório atual
    - Exibir mensagens de erro claras em PT-BR
    - _Requisitos: 1.1, 1.2_

  - [x] 4.2 Implementar coleta de dados
    - Iniciar timer "git_analysis"
    - Chamar `get_all_commits()` e armazenar resultado
    - Contar número de commits encontrados
    - Logar progresso com duração
    - _Requisitos: 1.2, 6.3, 6.5_

  - [x] 4.3 Implementar geração do documento
    - Iniciar timer "generation"
    - Redirecionar output das funções de geração para CHANGELOG.md (maiúsculas)
    - Chamar sequencialmente: generate_header, generate_mermaid, generate_table, generate_details
    - Logar sucesso com duração
    - _Requisitos: 1.1, 1.4, 3.1, 4.4, 4.8_

  - [x] 4.4 Adicionar finalização e logging
    - Calcular duração total com end_timer "total"
    - Exibir mensagem de sucesso com checkmark verde
    - Exibir duração total formatada
    - _Requisitos: 6.3, 6.5_

- [x] 5. Integrar make-changelog.sh com build.sh
  - [x] 5.1 Adicionar flag --no-changelog no build.sh
    - Adicionar variável `NO_CHANGELOG=false` no início
    - Adicionar case `--no-changelog` no parse de argumentos
    - Atualizar função `usage()` com documentação da nova flag
    - _Requisitos: 7.2, 7.3_

  - [x] 5.2 Adicionar chamada ao make-changelog.sh no build.sh
    - Adicionar execução antes das verificações de qualidade
    - Verificar se arquivo make-changelog.sh existe
    - Executar apenas se `NO_CHANGELOG=false`
    - Capturar falhas mas não abortar build (não crítico)
    - Logar início e resultado da geração
    - _Requisitos: 7.1, 7.4, 7.5_

- [x] 6. Testar script manualmente
  - [x] 6.1 Teste básico de execução
    - Executar `./make-changelog.sh` diretamente
    - Verificar se CHANGELOG.md foi criado
    - Verificar estrutura do documento (cabeçalho, tabela, mermaid, detalhes)
    - Verificar ordenação cronológica dos commits
    - _Requisitos: 1.1, 1.2, 4.3_

  - [x] 6.2 Teste de integração com build.sh
    - Executar `./build.sh` sem flags
    - Verificar se changelog foi gerado automaticamente
    - Verificar logs de execução
    - Executar `./build.sh --no-changelog`
    - Verificar se geração foi pulada
    - _Requisitos: 7.1, 7.2, 7.3_

  - [x] 6.3 Validar renderização do CHANGELOG.md
    - Abrir CHANGELOG.md no GitHub ou visualizador markdown
    - Verificar renderização correta da tabela
    - Verificar renderização do diagrama Mermaid gitGraph
    - Verificar formatação das seções detalhadas
    - _Requisitos: 4.5, 4.6, 4.7_

  - [x] 6.4 Validar conteúdo e qualidade
    - Verificar se todos os commits estão presentes
    - Verificar se datas estão no formato correto (YYYY-MM-DD)
    - Verificar se resumos são objetivos e didáticos
    - Verificar se classificações fazem sentido
    - Ajustar lógica de classificação ou resumos se necessário
    - _Requisitos: 2.3, 2.5, 4.2, 5.3, 5.5_

  - [x] 6.5 Testar tratamento de erros
    - Testar execução fora de repositório Git
    - Testar em diretório sem permissão de escrita
    - Testar em repositório sem commits
    - Verificar mensagens de erro claras em PT-BR
    - _Requisitos: 1.1, 1.2_

- [x] 7. Atualizar script para nova ordenação e formato de data/hora
  - [x] 7.1 Atualizar função `get_all_commits()`
    - Remover flag `--reverse` do comando git log
    - Alterar formato de data de `--date=short` para `--date=format:"%Y-%m-%d %H:%M:%S"`
    - Garantir que commits sejam retornados do mais recente ao mais antigo
    - _Requisitos: 1.2, 4.2, 4.3_

  - [x] 7.2 Atualizar função `generate_table()`
    - Alterar cabeçalho da coluna de "Data" para "Data/Hora"
    - Ajustar largura da coluna para acomodar formato completo
    - Garantir que datetime seja exibido no formato YYYY-MM-DD HH:MM:SS
    - _Requisitos: 4.1, 4.2_

  - [x] 7.3 Atualizar função `generate_mermaid()`
    - Extrair apenas a data (YYYY-MM-DD) do campo datetime para o tag
    - Manter formato curto no diagrama para melhor visualização
    - _Requisitos: 4.4, 4.5_

  - [x] 7.4 Atualizar função `generate_details()`
    - Alterar cabeçalho de seção para incluir data/hora completa
    - Formato: `### [\`hash\`] - YYYY-MM-DD HH:MM:SS`
    - _Requisitos: 5.1, 5.2_

  - [x] 7.5 Atualizar ordem de geração no fluxo principal
    - Alterar ordem de chamada das funções de geração
    - Nova ordem: generate_header, generate_mermaid, generate_table, generate_details
    - Garantir que Mermaid apareça primeiro, seguido da tabela
    - _Requisitos: 4.6, 4.7, 4.8_

- [x] 8. Garantir nomenclatura em maiúsculas para arquivos markdown
  - [x] 8.1 Verificar geração de CHANGELOG.md
    - Confirmar que o arquivo gerado é CHANGELOG.md (maiúsculas)
    - Verificar todas as referências ao arquivo no código
    - _Requisitos: 1.1, 3.1, 3.2_

  - [x] 8.2 Criar arquivo README.md na raiz (se não existir)
    - Criar README.md em maiúsculas na raiz do projeto
    - Adicionar seção de Changelog apontando para CHANGELOG.md
    - Incluir link direto: `[CHANGELOG.md](CHANGELOG.md)`
    - _Requisitos: 3.1, 3.2_

  - [x] 8.3 Atualizar documentação
    - Verificar que todos os arquivos .md no projeto seguem convenção de maiúsculas
    - Atualizar referências em comentários e logs se necessário
    - _Requisitos: 3.1, 3.2, 3.3_

- [x] 9. Testar implementação completa com novas mudanças
  - [x] 9.1 Teste de ordenação
    - Executar make-changelog.sh e verificar ordem dos commits
    - Confirmar que commits mais recentes aparecem primeiro
    - Verificar ordem em: tabela, mermaid e detalhes
    - _Requisitos: 1.2, 4.3, 5.4_

  - [x] 9.2 Teste de formato de data/hora
    - Verificar formato YYYY-MM-DD HH:MM:SS na tabela
    - Verificar formato YYYY-MM-DD HH:MM:SS nos detalhes
    - Verificar formato YYYY-MM-DD no diagrama Mermaid
    - _Requisitos: 4.1, 4.2, 4.5_

  - [x] 9.3 Teste de estrutura do documento
    - Verificar ordem: Cabeçalho → Mermaid → Tabela → Detalhes
    - Confirmar que Mermaid é a primeira seção após o cabeçalho
    - Verificar renderização completa no visualizador markdown
    - _Requisitos: 4.6, 4.7, 4.8_

  - [x] 9.4 Teste de nomenclatura
    - Confirmar que CHANGELOG.md está em maiúsculas
    - Verificar que README.md existe e aponta para CHANGELOG.md
    - Validar links e referências
    - _Requisitos: 3.1, 3.2_

