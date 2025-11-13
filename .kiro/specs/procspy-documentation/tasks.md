# Implementation Plan

- [x] 1. Criar estrutura base do README.md
  - Criar arquivo README.md com estrutura de seções
  - Adicionar table of contents com links internos
  - Adicionar badges de tecnologias (Go, SQLite, Gin)
  - _Requirements: 1.1, 1.2, 1.3_

- [x] 2. Implementar seção de Visão Geral
  - [x] 2.1 Escrever descrição completa do propósito do sistema
    - Descrever o problema: controle parental de tempo de tela
    - Explicar a solução: monitoramento automático e aplicação de limites
    - Listar tipos de aplicações monitoradas (jogos, browsers, players)
    - Destacar benefícios para pais e famílias
    - _Requirements: 1.1, 1.4_

  - [x] 2.2 Documentar componentes principais e deployment
    - Descrever Client: instalado nos computadores das crianças
    - Descrever Server: hospedado em cloud para centralizar dados
    - Descrever Watcher: proteção contra fechamento do Client
    - Explicar a estratégia de deployment distribuído
    - _Requirements: 1.2_

  - [x] 2.3 Criar diagrama de arquitetura de alto nível
    - Implementar diagrama Mermaid mostrando Client, Server e Watcher
    - Mostrar múltiplos computadores de crianças conectados ao Server
    - Indicar comunicação HTTP/REST
    - Adicionar legenda explicativa
    - _Requirements: 1.3, 7.3_

  - [x] 2.4 Descrever capacidades principais do sistema
    - Monitoramento em tempo real de processos
    - Aplicação automática de limites de tempo
    - Sistema de avisos progressivos
    - Terminação automática de processos
    - Limites configuráveis por dia da semana
    - Proteção contra desativação
    - Telemetria centralizada
    - Geração de relatórios
    - _Requirements: 1.4, 1.5_

- [x] 3. Implementar seção de Arquitetura do Sistema
  - [x] 3.1 Criar diagrama de deployment
    - Implementar diagrama Mermaid mostrando distribuição física
    - Mostrar componentes em máquinas locais vs cloud
    - Indicar portas e protocolos de comunicação
    - _Requirements: 2.1, 7.1_

  - [x] 3.2 Criar diagrama de componentes internos
    - Implementar diagrama Mermaid com módulos internos
    - Mostrar cmd/, internal/, handlers, services, storage
    - Indicar dependências entre módulos
    - _Requirements: 2.1, 2.5_

  - [x] 3.3 Documentar comunicação entre componentes
    - Descrever protocolo HTTP/REST
    - Listar endpoints utilizados
    - Explicar fluxo de requisições
    - Documentar mecanismo de retry
    - _Requirements: 2.3, 2.4_

  - [x] 3.4 Criar diagrama de fluxo de dados
    - Implementar diagrama mostrando fluxo desde detecção até storage
    - Incluir buffers e mecanismos de retry
    - Mostrar persistência em SQLite
    - _Requirements: 2.4_

- [x] 4. Implementar seção de Componentes Detalhados
  - [x] 4.1 Documentar Client Component
    - Descrever função e responsabilidades
    - Listar tecnologias utilizadas (go-ps, gin-gonic)
    - Detalhar funcionalidades: scan, matching, limits, warnings, kill
    - Explicar buffer e retry mechanism
    - Documentar health check endpoint
    - Incluir exemplos de logs
    - _Requirements: 2.1, 8.5_

  - [x] 4.2 Documentar Server Component
    - Descrever função como centralizador de telemetria
    - Listar tecnologias utilizadas (gin-gonic, SQLite)
    - Detalhar API REST endpoints
    - Explicar armazenamento de dados
    - Documentar sistema de relatórios
    - Incluir exemplos de logs
    - _Requirements: 2.1, 8.5_

  - [x] 4.3 Documentar Watcher Component
    - Descrever função de proteção
    - Explicar mecanismo de detecção de falha
    - Documentar processo de restart do Client
    - Detalhar configuração de comandos de recuperação
    - Incluir exemplos de logs
    - _Requirements: 2.1, 5.5, 8.5_

- [x] 5. Implementar seção de Modelos de Dados
  - [x] 5.1 Documentar modelo Target
    - Listar e explicar todos os campos
    - Documentar pattern matching com regex
    - Explicar sistema de limites por dia da semana
    - Documentar comandos (check, warning, limit)
    - Incluir exemplos JSON
    - _Requirements: 3.1, 4.5_

  - [x] 5.2 Documentar modelo Match
    - Listar e explicar todos os campos
    - Explicar como matches são registrados
    - Documentar cálculo de tempo acumulado
    - Incluir exemplos JSON
    - _Requirements: 3.2_

  - [x] 5.3 Documentar modelo Command
    - Listar e explicar todos os campos
    - Documentar tipos de comandos (Check, Warning, Limit, Kill)
    - Explicar campo source
    - Incluir exemplos JSON
    - _Requirements: 3.3_

  - [x] 5.4 Criar diagrama de entidades
    - Implementar diagrama ER com Mermaid
    - Mostrar relacionamentos entre Target, Match e Command
    - Incluir cardinalidade
    - _Requirements: 3.4_

  - [x] 5.5 Documentar schema SQLite
    - Descrever tabelas do banco de dados
    - Documentar índices e constraints
    - Explicar estratégia de armazenamento
    - _Requirements: 3.5_

- [ ] 6. Implementar seção de API REST
  - [x] 6.1 Documentar endpoint GET /targets/:user
    - Descrever propósito e funcionamento
    - Documentar parâmetros
    - Incluir exemplo de request
    - Incluir exemplo de response
    - Documentar códigos de status HTTP
    - _Requirements: 2.3_

  - [x] 6.2 Documentar endpoint POST /match/:user
    - Descrever propósito e funcionamento
    - Documentar body da requisição
    - Incluir exemplo de request
    - Incluir exemplo de response
    - Documentar códigos de status HTTP
    - _Requirements: 2.3_

  - [x] 6.3 Documentar endpoint POST /command/:user
    - Descrever propósito e funcionamento
    - Documentar body da requisição
    - Incluir exemplo de request
    - Incluir exemplo de response
    - Documentar códigos de status HTTP
    - _Requirements: 2.3_

  - [x] 6.4 Documentar endpoint GET /report/:user
    - Descrever propósito e funcionamento
    - Documentar parâmetros opcionais
    - Incluir exemplo de response
    - Documentar códigos de status HTTP
    - _Requirements: 2.3_

  - [x] 6.5 Documentar endpoint GET /healthcheck
    - Descrever propósito
    - Incluir exemplo de response
    - Documentar uso pelo Watcher
    - _Requirements: 2.3_

- [x] 7. Implementar seção de Fluxos Operacionais
  - [x] 7.1 Criar flowchart do ciclo de monitoramento
    - Implementar diagrama Mermaid do loop principal
    - Mostrar scan de processos
    - Mostrar matching e cálculo de tempo
    - Mostrar envio de dados ao Server
    - _Requirements: 5.1_

  - [x] 7.2 Criar sequence diagram de aplicação de limites
    - Implementar diagrama Mermaid de sequência
    - Mostrar interação Client-Server-Database
    - Incluir verificação de limites
    - Mostrar terminação de processos
    - _Requirements: 5.2_

  - [x] 7.3 Criar flowchart do sistema de avisos
    - Implementar diagrama Mermaid de decisão
    - Mostrar threshold de warning
    - Mostrar execução de warning command
    - _Requirements: 5.3_

  - [x] 7.4 Criar sequence diagram de recuperação pelo Watcher
    - Implementar diagrama Mermaid de sequência
    - Mostrar verificação periódica do health check
    - Mostrar detecção de falha
    - Mostrar execução de comando de restart
    - _Requirements: 5.5_

  - [x] 7.5 Documentar mecanismo de buffer e retry
    - Explicar buffers de match e command
    - Documentar estratégia de retry
    - Explicar DLQ (Dead Letter Queue)
    - Incluir exemplos de logs
    - _Requirements: 5.4_

- [x] 8. Implementar seção de Configuração
  - [x] 8.1 Documentar configuração do Client
    - Listar todos os parâmetros do config-client.json
    - Explicar cada campo em detalhes
    - Incluir exemplo completo
    - Documentar valores recomendados
    - _Requirements: 4.1, 4.4_

  - [x] 8.2 Documentar configuração do Server
    - Listar todos os parâmetros do config-server.json
    - Explicar cada campo em detalhes
    - Documentar user_targets e URLs
    - Incluir exemplo completo
    - Documentar valores recomendados
    - _Requirements: 4.2, 4.4_

  - [x] 8.3 Documentar configuração do Watcher
    - Listar todos os parâmetros do watcher-config.json
    - Explicar cada campo em detalhes
    - Documentar start_cmd para diferentes sistemas operacionais
    - Incluir exemplos para Linux e Windows
    - _Requirements: 4.3, 4.4_

  - [x] 8.4 Documentar configuração de Targets
    - Explicar estrutura do arquivo user-targets.json
    - Documentar todos os campos disponíveis
    - Explicar pattern matching com regex
    - Documentar sistema de limites por weekday
    - Documentar comandos (check, warning, limit)
    - Incluir múltiplos exemplos práticos (games, browsers, video players)
    - _Requirements: 4.5, 8.4_

- [x] 9. Implementar seção de Suporte Cross-Platform
  - [x] 9.1 Documentar plataformas suportadas
    - Listar Windows (7, 8, 10, 11, Server)
    - Listar Linux (Ubuntu, Debian, CentOS, RHEL)
    - Listar macOS (10.12+)
    - _Requirements: 8.1_

  - [x] 9.2 Documentar process killing cross-platform
    - Explicar como Go lida com processos em diferentes OS
    - Documentar uso de os.FindProcess() e Process.Kill()
    - Explicar implementação específica para Windows (TerminateProcess)
    - Explicar implementação específica para Linux/macOS (SIGKILL)
    - _Requirements: 8.2_

  - [x] 9.3 Documentar serviços do sistema
    - Documentar NSSM para Windows
    - Documentar systemd/systemctl para Linux
    - Documentar launchd para macOS (opcional)
    - Explicar configuração de restart automático
    - _Requirements: 8.3, 8.4_

  - [x] 9.4 Documentar estrutura de scripts de instalação
    - Documentar diretório ./install/
    - Listar scripts para Windows (PowerShell)
    - Listar scripts para Linux (Bash)
    - Listar arquivos .service para systemd
    - _Requirements: 8.5, 8.6_

- [x] 10. Implementar seção de Instalação e Deployment
  - [x] 10.1 Documentar processo de build
    - Incluir comandos para build de todos os componentes
    - Documentar dependências necessárias (Go 1.24+)
    - Incluir instruções para cross-compilation
    - Documentar build para Windows, Linux e macOS
    - _Requirements: 9.1_

  - [x] 10.2 Documentar instalação no Windows
    - Documentar download e instalação do NSSM
    - Criar guia passo-a-passo usando install-client.ps1
    - Criar guia passo-a-passo usando install-watcher.ps1
    - Incluir exemplo completo do script PowerShell
    - Documentar verificação via services.msc
    - Documentar comandos nssm (status, start, stop, restart)
    - _Requirements: 8.4, 9.2, 9.3_

  - [x] 10.3 Documentar instalação no Linux (Client/Watcher)
    - Criar guia passo-a-passo usando install-client.sh
    - Criar guia passo-a-passo usando install-watcher.sh
    - Incluir exemplo completo do script Bash
    - Incluir exemplo completo do arquivo .service
    - Documentar localização de binários (/usr/local/bin/)
    - Documentar localização de configs (/etc/procspy/)
    - Documentar verificação via systemctl status
    - Documentar comandos systemctl (start, stop, restart, enable)
    - _Requirements: 8.5, 9.2, 9.3_

  - [x] 10.4 Documentar deployment do Server em cloud (Linux)
    - Criar guia passo-a-passo usando install-server.sh
    - Incluir exemplo completo do script Bash
    - Incluir exemplo completo do arquivo .service
    - Documentar configuração de firewall (ufw)
    - Documentar configuração de proxy reverso (nginx/apache)
    - Documentar configuração de HTTPS
    - Documentar verificação via systemctl status
    - _Requirements: 9.4, 9.6_

  - [x] 10.5 Documentar scripts de desinstalação
    - Documentar uninstall-client.ps1 para Windows
    - Documentar uninstall-watcher.ps1 para Windows
    - Documentar uninstall-client.sh para Linux
    - Documentar uninstall-watcher.sh para Linux
    - Documentar uninstall-server.sh para Linux
    - _Requirements: 8.7, 9.3_

- [x] 11. Implementar seção de Operação e Monitoramento
  - [x] 11.1 Documentar gerenciamento de serviços no Windows
    - Documentar comandos nssm status
    - Documentar comandos nssm start/stop/restart
    - Documentar comandos nssm remove
    - Incluir exemplos práticos
    - _Requirements: 8.4, 9.2_

  - [x] 11.2 Documentar gerenciamento de serviços no Linux
    - Documentar comandos systemctl status
    - Documentar comandos systemctl start/stop/restart
    - Documentar comandos systemctl enable/disable
    - Incluir exemplos práticos
    - _Requirements: 8.5, 9.2_

  - [x] 11.3 Documentar inicialização manual (para testes)
    - Incluir comandos para iniciar cada componente manualmente
    - Documentar opções de linha de comando
    - Incluir exemplos de output inicial
    - _Requirements: 9.2_

  - [x] 11.4 Documentar monitoramento de logs
    - Explicar estrutura dos logs
    - Documentar localização dos arquivos de log (Windows vs Linux)
    - Incluir exemplos de mensagens importantes
    - Explicar rotação de logs
    - _Requirements: 9.5, 9.8_

  - [x] 11.5 Documentar verificação de saúde
    - Explicar como verificar se componentes estão rodando
    - Documentar uso dos health check endpoints
    - Incluir comandos curl de exemplo
    - Documentar verificação específica por OS
    - _Requirements: 9.3_

  - [x] 11.6 Documentar acesso a relatórios
    - Explicar como acessar relatórios de uso
    - Incluir exemplos de requisições
    - Mostrar interpretação dos dados
    - _Requirements: 9.6_

  - [x] 11.7 Criar guia de troubleshooting cross-platform
    - Documentar problema: Client não inicia (Windows/Linux)
    - Documentar problema: Watcher não detecta Client
    - Documentar problema: Server não recebe dados
    - Documentar problema: Processos não são terminados
    - Documentar problema: Problemas de conectividade
    - Documentar problema: Serviço não inicia automaticamente
    - Documentar problema: Permissões insuficientes
    - Incluir soluções específicas por OS
    - Incluir comandos de diagnóstico
    - _Requirements: 9.7, 9.9_

- [x] 12. Implementar seção de Estrutura do Projeto
  - [x] 12.1 Criar árvore de diretórios
    - Documentar estrutura completa do projeto incluindo ./install/
    - Explicar propósito de cada diretório
    - Documentar estrutura completa do projeto
    - Explicar propósito de cada diretório
    - Documentar convenções de organização
    - _Requirements: 6.1, 6.2_

  - [x] 12.2 Documentar organização de pacotes internos
    - Explicar separação cmd/ vs internal/
    - Documentar estrutura de domain/
    - Documentar estrutura de handlers/
    - Documentar estrutura de services/
    - Documentar estrutura de storage/
    - _Requirements: 6.3, 6.4_

  - [x] 12.3 Documentar diretórios de runtime
    - Explicar diretório logs/
    - Explicar diretório data/
    - Documentar permissões necessárias
    - _Requirements: 6.5_

- [x] 13. Adicionar seções complementares
  - [x] 13.1 Criar seção de Tecnologias Utilizadas
    - Listar Go e versão
    - Listar bibliotecas principais (gin-gonic, go-ps, SQLite)
    - Adicionar badges no topo do README
    - _Requirements: 1.1_

  - [x] 13.2 Criar seção de Características e Capacidades
    - Listar todas as capacidades do sistema
    - Destacar diferenciais (proteção Watcher, limites por dia, cross-platform)
    - Incluir casos de uso práticos
    - _Requirements: 1.4, 8.1_

  - [x] 13.3 Criar seção de Segurança
    - Documentar considerações de segurança
    - Explicar proteção contra desativação
    - Documentar permissões necessárias
    - Documentar execução como serviço privilegiado
    - _Requirements: 9.7_

  - [x] 13.4 Adicionar seção de FAQ
    - Responder perguntas comuns
    - Incluir dicas de configuração
    - Documentar limitações conhecidas
    - Incluir perguntas sobre cross-platform
    - _Requirements: 9.7_

  - [x] 13.5 Adicionar seção de Contribuição e Licença
    - Documentar como contribuir (se aplicável)
    - Incluir informações de licença
    - Adicionar informações de contato/suporte
    - _Requirements: 1.1_

- [x] 14. Revisar e finalizar documentação
  - [x] 14.1 Revisar todos os diagramas Mermaid
    - Verificar sintaxe correta
    - Testar renderização no GitHub
    - Ajustar estilos e cores
    - _Requirements: 7.3, 7.4, 7.5_

  - [x] 14.2 Revisar formatação Markdown
    - Verificar hierarquia de headers
    - Verificar links internos
    - Verificar code blocks e syntax highlighting
    - Verificar listas e tabelas
    - _Requirements: 7.2_

  - [x] 14.3 Revisar conteúdo textual
    - Verificar clareza e completude
    - Corrigir erros de português
    - Garantir tom acessível
    - Verificar consistência de termos
    - Verificar informações cross-platform
    - _Requirements: 1.1, 1.2, 1.3, 1.4, 8.1_

  - [x] 14.4 Adicionar table of contents
    - Criar índice completo no início
    - Adicionar links para todas as seções
    - Testar navegação
    - _Requirements: 7.2_

  - [x] 14.5 Fazer backup do README original
    - Renomear README.md atual para README.old.md
    - Substituir com novo README.md
    - _Requirements: 1.1_
