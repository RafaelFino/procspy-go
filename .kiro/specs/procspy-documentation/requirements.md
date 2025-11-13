# Requirements Document

## Introduction

Este documento define os requisitos para criar uma documentação técnica completa e detalhada do sistema Procspy. O Procspy é um sistema distribuído de monitoramento e controle parental desenvolvido em Go, projetado para monitorar o tempo de uso de aplicações específicas (jogos, browsers, players de vídeo) nos computadores das crianças e aplicar limites de tempo automaticamente. O sistema é composto por três componentes: Client (instalado nos computadores das crianças), Server (hospedado em cloud para centralizar telemetria e métricas), e Watcher (proteção contra fechamento do Client pelas crianças). A documentação deve incluir descrições arquiteturais com diagramas Mermaid, explicações detalhadas de componentes, fluxos de dados, guias de configuração e operação.

## Glossary

- **Procspy System**: O sistema completo de monitoramento de processos composto por três componentes principais (Client, Server, Watcher)
- **Client Component**: Componente cross-platform (Windows, Linux, macOS) que executa na máquina monitorada, responsável por escanear processos e enviar dados ao servidor, executando como serviço do sistema
- **Server Component**: Componente central que recebe dados dos clientes, armazena informações e fornece APIs REST
- **Watcher Component**: Componente de proteção cross-platform (Windows, Linux, macOS) instalado nos computadores das crianças que monitora a disponibilidade do Client Component e reinicia o serviço caso seja fechado, prevenindo que as crianças desabilitem o monitoramento, executando como serviço do sistema
- **NSSM**: Non-Sucking Service Manager, ferramenta para Windows que permite executar aplicações como serviços do Windows
- **systemctl**: Gerenciador de serviços do Linux (systemd) usado para controlar serviços
- **Installation Scripts**: Scripts automatizados localizados em ./install/ para instalação e configuração dos componentes como serviços em cada sistema operacional
- **Target**: Definição de um processo ou grupo de processos a serem monitorados, incluindo padrões de correspondência e limites de tempo
- **Match**: Registro de detecção de um processo alvo durante uma varredura
- **Command**: Registro de execução de comandos automáticos (avisos, limites, verificações)
- **Documentation System**: O conjunto completo de arquivos de documentação incluindo README, guias de arquitetura e diagramas
- **Mermaid Diagram**: Formato de diagramação baseado em texto usado para criar visualizações de arquitetura
- **REST API**: Interface de programação de aplicações baseada em HTTP usada para comunicação entre componentes
- **SQLite Database**: Sistema de banco de dados embutido usado para armazenar matches e comandos

## Requirements

### Requirement 1

**User Story:** Como pai/mãe ou administrador do sistema, eu quero uma documentação de visão geral do sistema, para que eu possa entender rapidamente o propósito, funcionamento e componentes do Procspy

#### Acceptance Criteria

1. THE Documentation System SHALL include a comprehensive overview section describing the parental control purpose
2. THE Documentation System SHALL describe all three main components (Client, Server, Watcher) and their deployment locations
3. THE Documentation System SHALL include a high-level architecture diagram using Mermaid format showing component distribution
4. THE Documentation System SHALL explain how the system monitors processes (games, browsers, video players) and enforces time limits
5. THE Documentation System SHALL describe the protection mechanism where Watcher prevents children from closing the Client

### Requirement 2

**User Story:** Como desenvolvedor, eu quero documentação detalhada da arquitetura do sistema, para que eu possa entender como os componentes interagem e se comunicam

#### Acceptance Criteria

1. THE Documentation System SHALL include a detailed architecture section with component descriptions
2. THE Documentation System SHALL include a sequence diagram showing the interaction flow between Client, Server, and Watcher
3. THE Documentation System SHALL document all REST API endpoints with their purposes
4. THE Documentation System SHALL describe the data flow from process detection to storage
5. THE Documentation System SHALL include a component diagram showing internal modules and their relationships

### Requirement 3

**User Story:** Como desenvolvedor, eu quero documentação dos modelos de dados, para que eu possa entender as estruturas de informação usadas no sistema

#### Acceptance Criteria

1. THE Documentation System SHALL document the Target data model with all fields and their purposes
2. THE Documentation System SHALL document the Match data model with all fields and their purposes
3. THE Documentation System SHALL document the Command data model with all fields and their purposes
4. THE Documentation System SHALL include a data model diagram showing relationships between entities
5. THE Documentation System SHALL describe the SQLite database schema

### Requirement 4

**User Story:** Como administrador de sistema, eu quero documentação de configuração, para que eu possa instalar e configurar corretamente todos os componentes

#### Acceptance Criteria

1. THE Documentation System SHALL document all configuration parameters for the Client Component
2. THE Documentation System SHALL document all configuration parameters for the Server Component
3. THE Documentation System SHALL document all configuration parameters for the Watcher Component
4. THE Documentation System SHALL include examples of configuration files with explanations
5. THE Documentation System SHALL document the Target configuration format with all available options

### Requirement 5

**User Story:** Como desenvolvedor ou administrador, eu quero documentação dos fluxos operacionais, para que eu possa entender como o sistema funciona em diferentes cenários

#### Acceptance Criteria

1. THE Documentation System SHALL include a flowchart showing the Client monitoring cycle with process scanning
2. THE Documentation System SHALL include a flowchart showing the limit enforcement process including automatic process termination
3. THE Documentation System SHALL include a flowchart showing the warning system operation
4. THE Documentation System SHALL document the buffer and retry mechanism for failed API calls
5. THE Documentation System SHALL include a sequence diagram showing how Watcher detects Client failure and restarts the service

### Requirement 6

**User Story:** Como desenvolvedor, eu quero documentação da estrutura de diretórios e organização do código, para que eu possa navegar facilmente no projeto

#### Acceptance Criteria

1. THE Documentation System SHALL include a complete directory structure overview
2. THE Documentation System SHALL describe the purpose of each main directory
3. THE Documentation System SHALL document the organization of internal packages
4. THE Documentation System SHALL explain the separation between cmd, internal, and configuration directories
5. THE Documentation System SHALL document the location and purpose of data and log directories
6. THE Documentation System SHALL document the ./install/ directory containing installation scripts for each operating system

### Requirement 7

**User Story:** Como desenvolvedor, eu quero diagramas visuais da arquitetura, para que eu possa visualizar rapidamente a estrutura do sistema

#### Acceptance Criteria

1. THE Documentation System SHALL include a deployment diagram showing how components are distributed
2. THE Documentation System SHALL include a network diagram showing communication paths
3. THE Documentation System SHALL use Mermaid syntax for all diagrams to enable version control
4. THE Documentation System SHALL include legends or annotations explaining diagram elements
5. THE Documentation System SHALL ensure all diagrams are properly formatted and renderable

### Requirement 8

**User Story:** Como desenvolvedor ou administrador, eu quero documentação sobre suporte cross-platform e scripts de instalação, para que eu possa entender como o sistema funciona em diferentes sistemas operacionais

#### Acceptance Criteria

1. THE Documentation System SHALL document cross-platform support for Windows, Linux and macOS
2. THE Documentation System SHALL document how process killing works across different operating systems
3. THE Documentation System SHALL document service installation using NSSM on Windows
4. THE Documentation System SHALL document service installation using systemctl on Linux
5. THE Documentation System SHALL document the installation scripts in ./install/ directory
6. THE Documentation System SHALL include examples of installation scripts for Windows (install-windows.ps1 or install-windows.bat)
7. THE Documentation System SHALL include examples of installation scripts for Linux (install-linux.sh)
8. THE Documentation System SHALL document systemd service files for Linux
9. THE Documentation System SHALL document NSSM commands for Windows service management

### Requirement 9

**User Story:** Como pai/mãe ou administrador, eu quero documentação de uso e operação, para que eu possa instalar, configurar e monitorar o sistema nos computadores das crianças

#### Acceptance Criteria

1. THE Documentation System SHALL document how to build the system using the provided scripts
2. THE Documentation System SHALL document the cross-platform installation process on children's computers (Client + Watcher) for Windows, Linux and macOS
3. THE Documentation System SHALL document step-by-step usage of installation scripts from ./install/ directory
4. THE Documentation System SHALL document the Server deployment process in cloud environment as a systemctl service
5. THE Documentation System SHALL document how to configure targets (games, browsers, video players) with time limits
6. THE Documentation System SHALL document how to monitor system health, view logs and access telemetry reports
7. THE Documentation System SHALL include troubleshooting guide for common scenarios across different operating systems
