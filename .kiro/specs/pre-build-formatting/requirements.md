# Requirements Document

## Introduction

Este documento descreve os requisitos para implementar um processo automatizado de formatação de código Go que será executado antes de cada build do projeto procspy. O objetivo é garantir que todo o código esteja formatado de acordo com os padrões Go antes da compilação.

## Glossary

- **Build System**: O sistema responsável por compilar o código fonte do projeto procspy em binários executáveis
- **go fmt**: Ferramenta padrão do Go que formata código fonte de acordo com as convenções da linguagem
- **Pre-build Hook**: Um processo ou script que é executado automaticamente antes do processo de build iniciar
- **Build Script**: O script build.sh existente no projeto que compila os binários

## Requirements

### Requirement 1

**User Story:** Como desenvolvedor, eu quero que o código seja formatado automaticamente antes do build, para que eu não precise me preocupar em executar go fmt manualmente

#### Acceptance Criteria

1. WHEN THE Build System inicia o processo de build, THE Build System SHALL executar o comando "go fmt ./..." antes de compilar o código
2. THE Build System SHALL completar a formatação de todos os arquivos Go no projeto antes de prosseguir com a compilação
3. IF a formatação falhar, THEN THE Build System SHALL exibir uma mensagem de erro clara e interromper o processo de build
4. THE Build System SHALL registrar no log quando a formatação automática for executada com sucesso

### Requirement 2

**User Story:** Como desenvolvedor, eu quero que o processo de formatação seja integrado ao build.sh existente, para que eu não precise modificar meu fluxo de trabalho atual

#### Acceptance Criteria

1. THE Build System SHALL integrar o comando de formatação no script build.sh existente
2. THE Build System SHALL manter a compatibilidade com todos os comandos e opções existentes do build.sh
3. THE Build System SHALL executar a formatação antes de qualquer operação de compilação
4. THE Build System SHALL preservar o código de saída apropriado para indicar sucesso ou falha do processo completo

### Requirement 3

**User Story:** Como desenvolvedor, eu quero ter a opção de pular a formatação automática quando necessário, para que eu possa fazer builds rápidos durante debugging

#### Acceptance Criteria

1. WHERE uma flag de skip é fornecida, THE Build System SHALL permitir pular o processo de formatação automática
2. WHEN a flag "--no-fmt" é passada como argumento, THE Build System SHALL executar o build sem formatação prévia
3. THE Build System SHALL documentar a opção de skip na saída de ajuda do script

### Requirement 4

**User Story:** Como desenvolvedor, eu quero que todas as mensagens de log nos scripts de build e teste incluam timestamp, para que eu possa rastrear quando cada operação ocorreu

#### Acceptance Criteria

1. THE Build System SHALL incluir um timestamp no formato "[yyyy-mm-dd HH:MM:SS]" no início de cada mensagem de log
2. THE Build System SHALL exibir o timestamp em cor verde e em negrito
3. THE Build System SHALL aplicar o formato de timestamp em todas as mensagens do script build.sh
4. THE Build System SHALL aplicar o formato de timestamp em todas as mensagens do script test.sh
5. THE Build System SHALL manter a cor original da mensagem após o timestamp (verde para sucesso, vermelho para erro, amarelo para avisos)
