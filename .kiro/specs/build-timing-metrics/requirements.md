# Requirements Document - Build Timing Metrics

## Introduction

Este documento define os requisitos para adicionar medição de tempo humanamente legível aos scripts de build e teste do Procspy. O objetivo é permitir que desenvolvedores avaliem facilmente quanto tempo cada etapa do processo leva (formatação, testes, builds para diferentes plataformas), exibindo os tempos em unidades apropriadas (ns, ms, s) de forma automática e legível.

## Glossary

- **Build Script**: O script build.sh responsável por compilar os binários do Procspy
- **Test Script**: O script test.sh responsável por executar os testes unitários
- **Timing Metric**: Medição de tempo de execução de uma operação específica
- **Human-Readable Format**: Formato de tempo que escolhe automaticamente a unidade apropriada (nanosegundos, milissegundos ou segundos)
- **Build Phase**: Uma etapa específica do processo de build (formatação, testes, compilação)
- **Platform Build**: Compilação para uma plataforma específica (OS/ARCH)

## Requirements

### Requirement 1: Medição de Tempo de Operações

**User Story:** Como desenvolvedor, quero ver quanto tempo cada operação do build leva, para que eu possa identificar gargalos no processo de compilação.

#### Acceptance Criteria

1. WHEN THE Build Script executa uma operação, THE Build Script SHALL registrar o tempo de início da operação
2. WHEN uma operação é concluída, THE Build Script SHALL calcular o tempo decorrido desde o início
3. WHEN THE Build Script exibe o resultado de uma operação, THE Build Script SHALL incluir o tempo de execução no formato humanamente legível
4. WHEN THE Test Script executa testes, THE Test Script SHALL medir e exibir o tempo total de execução dos testes

### Requirement 2: Formatação Automática de Unidades de Tempo

**User Story:** Como desenvolvedor, quero que os tempos sejam exibidos em unidades apropriadas automaticamente, para que eu não precise fazer conversões mentais.

#### Acceptance Criteria

1. WHEN o tempo de execução é menor que 1 milissegundo, THE Build Script SHALL exibir o tempo em nanosegundos com sufixo "ns"
2. WHEN o tempo de execução está entre 1 milissegundo e 1 segundo, THE Build Script SHALL exibir o tempo em milissegundos com sufixo "ms"
3. WHEN o tempo de execução é maior ou igual a 1 segundo, THE Build Script SHALL exibir o tempo em segundos com sufixo "s"
4. THE Build Script SHALL exibir valores numéricos com no máximo 2 casas decimais para legibilidade

### Requirement 3: Medição de Etapas do Build

**User Story:** Como desenvolvedor, quero ver o tempo de cada etapa do build separadamente, para que eu possa identificar qual etapa está mais lenta.

#### Acceptance Criteria

1. WHEN THE Build Script executa formatação automática, THE Build Script SHALL exibir o tempo gasto na formatação
2. WHEN THE Build Script executa verificações de qualidade, THE Build Script SHALL exibir o tempo gasto nas verificações
3. WHEN THE Build Script executa testes, THE Build Script SHALL exibir o tempo gasto nos testes
4. WHEN THE Build Script compila para uma plataforma, THE Build Script SHALL exibir o tempo gasto na compilação daquela plataforma
5. WHEN THE Build Script completa todas as operações, THE Build Script SHALL exibir o tempo total de execução

### Requirement 4: Medição de Builds por Plataforma

**User Story:** Como desenvolvedor, quero ver quanto tempo cada build de plataforma leva, para que eu possa avaliar o impacto de compilar para múltiplas plataformas.

#### Acceptance Criteria

1. WHEN THE Build Script compila um componente para uma plataforma, THE Build Script SHALL medir o tempo de compilação daquele componente
2. WHEN THE Build Script completa a compilação de todos os componentes de uma plataforma, THE Build Script SHALL exibir o tempo total para aquela plataforma
3. WHERE THE Build Script compila para múltiplas plataformas, THE Build Script SHALL exibir o tempo de cada plataforma individualmente
4. WHEN THE Build Script completa builds para todas as plataformas, THE Build Script SHALL exibir um resumo com tempos de todas as plataformas

### Requirement 5: Integração com Output Existente

**User Story:** Como desenvolvedor, quero que as métricas de tempo sejam integradas naturalmente ao output existente, para que eu não perca informações importantes.

#### Acceptance Criteria

1. THE Build Script SHALL preservar todas as mensagens de log existentes
2. THE Build Script SHALL adicionar informações de tempo sem quebrar o formato atual das mensagens
3. THE Build Script SHALL manter as cores e formatação existentes do output
4. THE Test Script SHALL adicionar métricas de tempo sem alterar o formato do relatório de coverage

### Requirement 6: Função Reutilizável de Formatação de Tempo

**User Story:** Como desenvolvedor, quero uma função reutilizável para formatar tempos, para que o código seja consistente e fácil de manter.

#### Acceptance Criteria

1. THE Build Script SHALL implementar uma função format_duration que aceita tempo em nanosegundos
2. THE format_duration function SHALL retornar uma string formatada com valor e unidade apropriada
3. THE format_duration function SHALL ser utilizada em todas as medições de tempo do Build Script
4. THE Test Script SHALL implementar a mesma função format_duration para consistência

### Requirement 7: Organização de Arquivos de Coverage

**User Story:** Como desenvolvedor, quero que todos os arquivos de coverage sejam organizados em uma pasta dedicada, para que a raiz do projeto permaneça limpa e organizada.

#### Acceptance Criteria

1. THE Test Script SHALL criar um diretório ./coverage na raiz do projeto se não existir
2. THE Test Script SHALL armazenar todos os arquivos de coverage (*.out, *.html) no diretório ./coverage
3. THE Test Script SHALL mover o arquivo COVERAGE_ANALYSIS.md para o diretório ./coverage se existir na raiz
4. THE Build Script SHALL referenciar arquivos de coverage no diretório ./coverage quando necessário
5. THE Test Script SHALL atualizar todas as referências a arquivos de coverage para usar o caminho ./coverage
6. THE Build Script SHALL limpar arquivos de coverage do diretório ./coverage quando usar a flag --clean

### Requirement 8: Padronização de Nomes de Arquivos Markdown

**User Story:** Como desenvolvedor, quero que todos os arquivos markdown na raiz sigam uma convenção de nomenclatura consistente, para que o projeto tenha uma estrutura organizada e previsível.

#### Acceptance Criteria

1. THE Project SHALL renomear todos os arquivos .md na raiz do projeto para usar letras minúsculas
2. WHEN um arquivo .md tem nome em maiúsculas, THE Project SHALL renomear para minúsculas mantendo underscores
3. THE Project SHALL atualizar todas as referências internas aos arquivos renomeados
4. THE README.md SHALL ser renomeado para readme.md

### Requirement 9: Documentação de Referência no README

**User Story:** Como desenvolvedor, quero que o README aponte para documentação importante do projeto, para que eu possa encontrar facilmente informações sobre testes e compatibilidade.

#### Acceptance Criteria

1. THE readme.md SHALL incluir uma seção de documentação com links para test.md
2. THE readme.md SHALL incluir link para cross_platform_testing.md na seção de documentação
3. THE readme.md SHALL incluir breve descrição do conteúdo de cada documento referenciado
4. THE readme.md SHALL organizar links de documentação em uma seção claramente identificada

### Requirement 10: Qualidade de Mensagens de Log

**User Story:** Como desenvolvedor, quero que todas as mensagens de log sejam gramaticalmente corretas e didáticas, para que eu possa entender facilmente os eventos do sistema.

#### Acceptance Criteria

1. THE Project SHALL revisar todas as mensagens de log em arquivos .go para correção gramatical
2. THE Project SHALL garantir que mensagens de log sejam didáticas e expliquem claramente o evento
3. THE Project SHALL manter a formatação e estética já adotada nas mensagens de log existentes
4. THE Project SHALL usar linguagem consistente em todas as mensagens de log
5. THE Project SHALL incluir contexto suficiente em mensagens de erro para facilitar debugging
