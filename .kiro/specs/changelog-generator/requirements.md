# Documento de Requisitos

## Introdução

Esta funcionalidade implementa um script bash (make-changelog.sh) que gera automaticamente um arquivo CHANGELOG.md analisando o histórico completo de commits do Git. O script usa comandos nativos do Git para percorrer todos os commits, analisa as diferenças de código para entender as mudanças reais, e gera um changelog estruturado e objetivo com diagrama Mermaid, tabela resumida e detalhamento dos commits. Os commits são exibidos do mais recente para o mais antigo. O script segue os mesmos padrões de logging dos scripts build.sh e test.sh existentes no projeto. Todos os arquivos markdown do projeto devem seguir a convenção de nomenclatura em maiúsculas (ex: CHANGELOG.md, README.md).

## Glossário

- **Script Gerador**: O script bash make-changelog.sh que analisa o histórico Git e gera o arquivo CHANGELOG.md
- **Repositório Git**: O repositório de controle de versão contendo o histórico de commits do projeto
- **Diff de Commit**: As diferenças de código introduzidas por um commit específico
- **Tabela Resumo**: Uma visão tabular no início do CHANGELOG.md contendo data, mensagem de commit e resumo das mudanças
- **Análise de Mudanças**: O processo de examinar diffs de código usando comandos git para entender quais resultados foram alcançados
- **Logging Padronizado**: Formato de log com timestamp e duração usado nos scripts build.sh e test.sh

## Requirements

### Requisito 1

**História de Usuário:** Como desenvolvedor, quero que o sistema gere automaticamente um arquivo CHANGELOG.md a partir do histórico Git, para que eu possa manter documentação do projeto sem esforço manual

#### Critérios de Aceitação

1. O Script Gerador DEVE criar ou atualizar um arquivo chamado "CHANGELOG.md" em maiúsculas
2. O Script Gerador DEVE percorrer o histórico de commits do Git e exibi-los do mais recente para o mais antigo
3. O Script Gerador DEVE analisar as diferenças de código de cada commit usando `git diff` e `git show --stat`
4. O Script Gerador DEVE gerar uma tabela resumo contendo colunas de data/hora, mensagem de commit e resumo das mudanças
5. O Script Gerador DEVE sobrescrever o arquivo CHANGELOG.md existente com o conteúdo completo atualizado

### Requisito 2

**História de Usuário:** Como desenvolvedor, quero que o changelog priorize a análise dos diffs de código ao invés das mensagens de commit, para que o changelog reflita os resultados reais alcançados

#### Critérios de Aceitação

1. QUANDO analisar um commit, O Script Gerador DEVE examinar as diferenças de código usando comandos git como fonte primária de verdade
2. O Script Gerador DEVE gerar resumos objetivos baseados nos arquivos modificados identificados por `git show --name-status`
3. O Script Gerador DEVE identificar o tipo de mudança analisando extensões e caminhos de arquivos (testes, docs, código, config)
4. O Script Gerador DEVE incluir a mensagem de commit na tabela resumo apenas para referência
5. O Script Gerador DEVE priorizar a análise do diff sobre o conteúdo da mensagem de commit ao determinar descrições

### Requisito 3

**História de Usuário:** Como desenvolvedor, quero que todos os arquivos markdown sigam a convenção de nomenclatura em maiúsculas, para que o projeto mantenha padrões consistentes de nomenclatura

#### Critérios de Aceitação

1. O Script Gerador DEVE nomear o arquivo de changelog como "CHANGELOG.md" em maiúsculas
2. O Script Gerador DEVE seguir a convenção do projeto onde todos os arquivos .md usam nomes em maiúsculas
3. O Script Gerador DEVE ser nomeado "make-changelog.sh" seguindo o padrão dos outros scripts do projeto

### Requisito 4

**História de Usuário:** Como desenvolvedor, quero que o changelog inclua uma tabela resumo estruturada e um gráfico visual de timeline, para que eu possa rapidamente visualizar todas as mudanças cronologicamente

#### Critérios de Aceitação

1. O Script Gerador DEVE criar uma tabela com colunas para Data/Hora, Mensagem de Commit e Resumo das Mudanças
2. O Script Gerador DEVE formatar data e hora no formato ISO completo (YYYY-MM-DD HH:MM:SS) usando `git log --date=format`
3. O Script Gerador DEVE ordenar as entradas da tabela cronologicamente do mais recente ao mais antigo
4. O Script Gerador DEVE incluir um diagrama Mermaid gitGraph representando uma timeline vertical dos commits
5. O Script Gerador DEVE exibir no diagrama Mermaid cada commit com sua data, ID curto e resumo
6. O Script Gerador DEVE posicionar o diagrama Mermaid como primeira seção após o cabeçalho do documento
7. O Script Gerador DEVE posicionar a tabela resumo logo após o diagrama Mermaid
8. O Script Gerador DEVE formatar a tabela usando sintaxe markdown

### Requisito 5

**História de Usuário:** Como desenvolvedor, quero descrições objetivas e didáticas das mudanças de cada commit, para que eu possa entender rapidamente o resultado alcançado sem excesso de detalhes

#### Critérios de Aceitação

1. QUANDO processar cada commit, O Script Gerador DEVE criar uma seção objetiva abaixo da tabela resumo
2. O Script Gerador DEVE agrupar mudanças por commit com cabeçalhos contendo hash curto e data/hora
3. O Script Gerador DEVE descrever os resultados das mudanças de forma didática e concisa
4. O Script Gerador DEVE ordenar as seções detalhadas do mais recente ao mais antigo
4. O Script Gerador DEVE focar no que foi alcançado (resultado) ao invés de listar todos os arquivos modificados
5. O Script Gerador DEVE manter descrições breves e diretas, evitando verbosidade excessiva

### Requisito 6

**História de Usuário:** Como desenvolvedor, quero que o script use o mesmo padrão de logging dos outros scripts do projeto, para manter consistência visual e facilitar debugging

#### Critérios de Aceitação

1. O Script Gerador DEVE usar as mesmas funções de logging (log_msg) dos scripts build.sh e test.sh
2. O Script Gerador DEVE incluir timestamps em todas as mensagens de log
3. O Script Gerador DEVE medir e exibir duração das operações principais usando start_timer e end_timer
4. O Script Gerador DEVE usar as mesmas cores e formatação (GREEN, RED, YELLOW, BLUE, BOLD, NC)
5. O Script Gerador DEVE exibir duração total da execução ao final

### Requisito 7

**História de Usuário:** Como desenvolvedor, quero que o changelog seja gerado automaticamente durante o build, para garantir que a documentação esteja sempre atualizada

#### Critérios de Aceitação

1. O script build.sh DEVE executar make-changelog.sh automaticamente antes de iniciar a compilação
2. O script build.sh DEVE aceitar um parâmetro --no-changelog para pular a geração do changelog
3. QUANDO o parâmetro --no-changelog for fornecido, O script build.sh DEVE pular a execução do make-changelog.sh
4. O script build.sh DEVE exibir mensagem de log indicando que a geração do changelog foi iniciada
5. SE o make-changelog.sh falhar, O script build.sh DEVE continuar com o build normalmente (não deve abortar)
