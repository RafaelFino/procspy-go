# Requisitos - Implementação de Testes Unitários

## Introdução

Este documento especifica os requisitos para implementação de uma suíte completa de testes unitários para o projeto Procspy. O sistema de testes deve cobrir todos os componentes Go do projeto, garantir qualidade do código através de validação automatizada durante o processo de build, e fornecer documentação clara sobre o escopo e propósito de cada teste.

## Glossário

- **Sistema de Testes**: Conjunto de arquivos de teste unitário (`*_test.go`) e scripts de execução que validam o comportamento do código
- **Procspy**: Sistema distribuído de monitoramento e controle parental desenvolvido em Go
- **Client**: Componente que monitora processos localmente nos computadores das crianças
- **Server**: Componente que centraliza telemetria e fornece configurações
- **Watcher**: Componente que garante que o Client não seja desativado
- **Coverage**: Métrica que indica a porcentagem de código coberta por testes
- **test.sh**: Script bash responsável por executar todos os testes unitários
- **test.md**: Documento que descreve o escopo e propósito de cada arquivo de teste

## Requisitos

### Requisito 1: Cobertura de Testes para Arquivos Go

**User Story:** Como desenvolvedor, quero que todos os arquivos `.go` do projeto possuam arquivos de teste correspondentes, para garantir que todo o código seja validado.

#### Acceptance Criteria

1. WHEN o Sistema de Testes é executado, THE Sistema de Testes SHALL criar um arquivo `*_test.go` correspondente para cada arquivo `.go` que contenha funções ou métodos testáveis
2. WHEN um arquivo `.go` contém apenas definições de tipos sem lógica, THE Sistema de Testes SHALL documentar a ausência de testes no arquivo test.md com justificativa
3. THE Sistema de Testes SHALL incluir testes para todos os componentes principais: Client, Server, Watcher, handlers, services, storage, domain e config
4. THE Sistema de Testes SHALL garantir que cada arquivo de teste contenha pelo menos um caso de teste válido que execute sem erros

### Requisito 2: Documentação de Testes

**User Story:** Como desenvolvedor, quero ter documentação clara sobre o escopo de cada teste, para entender rapidamente o que está sendo testado e por quê.

#### Acceptance Criteria

1. THE Sistema de Testes SHALL criar um arquivo `test.md` na raiz do projeto que liste todos os arquivos de teste
2. WHEN o arquivo test.md é gerado, THE Sistema de Testes SHALL incluir para cada arquivo de teste: caminho do arquivo, escopo dos testes, e lista de funções/métodos testados
3. THE Sistema de Testes SHALL adicionar comentários em português em cada função de teste explicando o que está sendo testado e como
4. WHEN um teste valida múltiplos cenários, THE Sistema de Testes SHALL documentar cada cenário com comentários inline no código de teste

### Requisito 3: Script de Execução de Testes

**User Story:** Como desenvolvedor, quero executar todos os testes através de um script dedicado, para validar o código de forma consistente e automatizada.

#### Acceptance Criteria

1. THE Sistema de Testes SHALL criar um script `test.sh` na raiz do projeto que execute todos os testes unitários
2. WHEN o script test.sh é executado, THE Sistema de Testes SHALL executar `go test` em todos os pacotes do projeto
3. WHEN os testes são executados, THE Sistema de Testes SHALL gerar relatório de coverage em formato legível
4. IF algum teste falhar, THEN THE Sistema de Testes SHALL retornar código de saída não-zero e exibir mensagem de erro clara
5. THE Sistema de Testes SHALL exibir estatísticas de execução incluindo: número total de testes, testes passados, testes falhados, e porcentagem de coverage

### Requisito 4: Integração com Build

**User Story:** Como desenvolvedor, quero que os testes sejam executados automaticamente durante o build, para garantir que código com falhas não seja compilado.

#### Acceptance Criteria

1. WHEN o script build.sh é executado sem flags especiais, THE Sistema de Testes SHALL executar test.sh antes da compilação
2. IF os testes falharem durante o build, THEN THE Sistema de Testes SHALL interromper o processo de build e exibir mensagem de erro
3. THE Sistema de Testes SHALL permitir que o build seja executado sem testes através da flag `--build-only`
4. THE Sistema de Testes SHALL permitir que apenas os testes sejam executados através da flag `--test-only`

### Requisito 5: Qualidade dos Testes

**User Story:** Como desenvolvedor, quero que os testes sejam bem estruturados e sigam boas práticas, para facilitar manutenção e compreensão.

#### Acceptance Criteria

1. THE Sistema de Testes SHALL utilizar a biblioteca padrão `testing` do Go para todos os testes
2. WHEN um teste valida comportamento de funções HTTP, THE Sistema de Testes SHALL utilizar `httptest` para criar servidores de teste
3. THE Sistema de Testes SHALL nomear funções de teste seguindo o padrão `Test<FunctionName>` ou `Test<FunctionName>_<Scenario>`
4. WHEN um teste requer setup ou teardown, THE Sistema de Testes SHALL implementar funções auxiliares claramente nomeadas
5. THE Sistema de Testes SHALL validar tanto casos de sucesso quanto casos de erro para cada função testada
6. THE Sistema de Testes SHALL evitar dependências externas nos testes, utilizando mocks quando necessário

### Requisito 6: Coverage Mínimo

**User Story:** Como desenvolvedor, quero garantir um nível elevado de cobertura de testes, para assegurar máxima qualidade do código.

#### Acceptance Criteria

1. THE Sistema de Testes SHALL calcular e exibir a porcentagem de coverage após execução dos testes
2. WHEN o coverage é calculado, THE Sistema de Testes SHALL gerar arquivo `coverage.out` na raiz do projeto
3. THE Sistema de Testes SHALL exibir warning se o coverage total for inferior a 99%
4. THE Sistema de Testes SHALL aceitar coverage mínimo de 70% apenas para componentes de alta complexidade onde 99% não seja viável
5. THE Sistema de Testes SHALL permitir visualização detalhada de coverage por arquivo através do comando `go tool cover`
6. THE Sistema de Testes SHALL exibir relatório detalhado de coverage por pacote identificando áreas abaixo da meta

### Requisito 7: Testes de Componentes Críticos

**User Story:** Como desenvolvedor, quero que componentes críticos do sistema tenham testes abrangentes, para garantir confiabilidade.

#### Acceptance Criteria

1. THE Sistema de Testes SHALL incluir testes para todas as funções públicas do pacote `client`
2. THE Sistema de Testes SHALL incluir testes para todos os handlers HTTP do pacote `server`
3. THE Sistema de Testes SHALL incluir testes para todas as operações de storage (database)
4. THE Sistema de Testes SHALL incluir testes para todas as funções de parsing e validação do pacote `domain`
5. THE Sistema de Testes SHALL incluir testes para o mecanismo de health check do `watcher`
6. WHEN testes envolvem operações de I/O, THE Sistema de Testes SHALL utilizar interfaces ou mocks para isolar a lógica de negócio

### Requisito 8: Execução Cross-Platform

**User Story:** Como desenvolvedor, quero que os testes sejam executáveis em diferentes sistemas operacionais, para garantir compatibilidade.

#### Acceptance Criteria

1. THE Sistema de Testes SHALL executar com sucesso em Linux, Windows e macOS
2. WHEN o script test.sh é executado em Windows, THE Sistema de Testes SHALL utilizar Git Bash ou WSL para compatibilidade
3. THE Sistema de Testes SHALL evitar uso de comandos específicos de sistema operacional nos testes
4. WHEN testes precisam validar comportamento específico de OS, THE Sistema de Testes SHALL utilizar build tags ou verificação de `runtime.GOOS`
