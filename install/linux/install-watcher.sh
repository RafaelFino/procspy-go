#!/bin/bash
# install-watcher.sh
# Script de instalação do Procspy Watcher para Linux usando systemd
# Requer privilégios de root
#
# Uso:
#   sudo ./install-watcher.sh [opções]
#
# Opções:
#   -b, --binary PATH       Caminho para o binário procspy-watcher
#   -c, --config PATH       Caminho para o arquivo watcher-config.json
#   -h, --help              Exibe esta ajuda

set -e

# Cores
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
CYAN='\033[0;36m'
NC='\033[0m'

# Configurações padrão
SERVICE_NAME="procspy-watcher"
INSTALL_DIR="/usr/local/bin"
CONFIG_DIR="/etc/procspy"
LOG_DIR="/var/log/procspy"
SERVICE_FILE="/etc/systemd/system/${SERVICE_NAME}.service"

# Caminhos padrão dos arquivos fonte
DEFAULT_BINARY_PATH="./bin/procspy-watcher"
DEFAULT_CONFIG_PATH="./etc/watcher-config.json"

# Variáveis para caminhos customizados
BINARY_PATH=""
CONFIG_PATH=""

# Função para exibir ajuda
show_help() {
    cat << EOF
${GREEN}=== Instalação do Procspy Watcher ===${NC}

Uso: sudo $0 [opções]

Opções:
  -b, --binary PATH       Caminho para o binário procspy-watcher
                          Padrão: $DEFAULT_BINARY_PATH
  
  -c, --config PATH       Caminho para o arquivo watcher-config.json
                          Padrão: $DEFAULT_CONFIG_PATH
  
  -h, --help              Exibe esta ajuda

Exemplos:
  # Instalação padrão
  sudo $0

  # Instalação com binário customizado
  sudo $0 --binary /tmp/procspy-watcher

  # Instalação com binário e config customizados
  sudo $0 -b /opt/builds/procspy-watcher -c /opt/configs/watcher-config.json

EOF
    exit 0
}

# Parse de argumentos
while [[ $# -gt 0 ]]; do
    case $1 in
        -b|--binary)
            BINARY_PATH="$2"
            shift 2
            ;;
        -c|--config)
            CONFIG_PATH="$2"
            shift 2
            ;;
        -h|--help)
            show_help
            ;;
        *)
            echo -e "${RED}Erro: Opção desconhecida: $1${NC}"
            echo "Use --help para ver as opções disponíveis"
            exit 1
            ;;
    esac
done

# Define caminhos padrão se não foram especificados
if [ -z "$BINARY_PATH" ]; then
    BINARY_PATH="$DEFAULT_BINARY_PATH"
fi

if [ -z "$CONFIG_PATH" ]; then
    CONFIG_PATH="$DEFAULT_CONFIG_PATH"
fi

echo -e "${GREEN}=== Instalação do Procspy Watcher ===${NC}"
echo ""
echo -e "${CYAN}Configuração:${NC}"
echo "  Binário: $BINARY_PATH"
echo "  Config:  $CONFIG_PATH"
echo ""

# Verifica se está rodando como root
if [ "$EUID" -ne 0 ]; then
    echo -e "${RED}Erro: Este script deve ser executado como root${NC}"
    echo "Execute: sudo $0"
    exit 1
fi

echo -e "${GREEN}✓ Privilégios de root verificados${NC}"

# Verifica se binário existe
if [ ! -f "$BINARY_PATH" ]; then
    echo -e "${RED}Erro: Binário não encontrado em: $BINARY_PATH${NC}"
    echo ""
    echo -e "${YELLOW}Soluções:${NC}"
    echo "  1. Execute o build primeiro: ./build.sh"
    echo "  2. Especifique o caminho correto: $0 --binary /caminho/para/procspy-watcher"
    exit 1
fi

echo -e "${GREEN}✓ Binário encontrado: $BINARY_PATH${NC}"

# Verifica se binário é executável
if [ ! -x "$BINARY_PATH" ]; then
    echo -e "${YELLOW}Aviso: Binário não é executável, ajustando permissões...${NC}"
    chmod +x "$BINARY_PATH"
fi

# Verifica se configuração existe
CONFIG_EXISTS=false
if [ -f "$CONFIG_PATH" ]; then
    CONFIG_EXISTS=true
    echo -e "${GREEN}✓ Arquivo de configuração encontrado: $CONFIG_PATH${NC}"
else
    echo -e "${YELLOW}⚠ Arquivo de configuração não encontrado: $CONFIG_PATH${NC}"
    echo -e "${YELLOW}  Uma configuração padrão será criada${NC}"
fi

# Para serviço existente
if systemctl is-active --quiet $SERVICE_NAME; then
    echo -e "${YELLOW}Parando serviço existente...${NC}"
    systemctl stop $SERVICE_NAME
    sleep 1
fi

# Copia binário
echo -e "${CYAN}Copiando binário para $INSTALL_DIR...${NC}"
cp "$BINARY_PATH" "$INSTALL_DIR/procspy-watcher"
chmod +x "$INSTALL_DIR/procspy-watcher"
echo -e "${GREEN}✓ Binário instalado em $INSTALL_DIR/procspy-watcher${NC}"

# Cria diretórios
echo -e "${CYAN}Criando diretórios...${NC}"
mkdir -p "$CONFIG_DIR"
mkdir -p "$LOG_DIR"
chmod 755 "$LOG_DIR"
echo -e "${GREEN}✓ Diretórios criados${NC}"

# Copia ou cria configuração
if [ "$CONFIG_EXISTS" = true ]; then
    echo -e "${CYAN}Copiando configuração...${NC}"
    cp "$CONFIG_PATH" "$CONFIG_DIR/watcher-config.json"
    echo -e "${GREEN}✓ Configuração copiada para $CONFIG_DIR/watcher-config.json${NC}"
else
    echo -e "${CYAN}Criando configuração padrão...${NC}"
    cat > "$CONFIG_DIR/watcher-config.json" <<EOF
{
    "log_path": "$LOG_DIR",
    "interval": 10,
    "procspy_url": "http://localhost:8888/healthcheck",
    "start_cmd": "systemctl restart procspy-client"
}
EOF
    echo -e "${GREEN}✓ Configuração padrão criada em $CONFIG_DIR/watcher-config.json${NC}"
fi

# Cria arquivo de serviço systemd
echo -e "${CYAN}Criando arquivo de serviço systemd...${NC}"
cat > "$SERVICE_FILE" <<EOF
[Unit]
Description=Procspy Watcher - Monitora e reinicia o Client automaticamente
Documentation=https://github.com/seu-usuario/procspy
After=network.target procspy-client.service

[Service]
Type=simple
User=root
WorkingDirectory=$CONFIG_DIR
ExecStart=$INSTALL_DIR/procspy-watcher $CONFIG_DIR/watcher-config.json
Restart=always
RestartSec=5
StandardOutput=journal
StandardError=journal

# Limites de recursos
LimitNOFILE=65536

[Install]
WantedBy=multi-user.target
EOF

echo -e "${GREEN}✓ Arquivo de serviço criado em $SERVICE_FILE${NC}"

# Recarrega systemd
echo -e "${CYAN}Recarregando systemd...${NC}"
systemctl daemon-reload
echo -e "${GREEN}✓ Systemd recarregado${NC}"

# Habilita serviço
echo -e "${CYAN}Habilitando serviço para iniciar no boot...${NC}"
systemctl enable $SERVICE_NAME
echo -e "${GREEN}✓ Serviço habilitado${NC}"

# Inicia serviço
echo -e "${CYAN}Iniciando serviço...${NC}"
systemctl start $SERVICE_NAME
sleep 2

# Verifica status
if systemctl is-active --quiet $SERVICE_NAME; then
    echo -e "${GREEN}✓ Serviço iniciado com sucesso${NC}"
else
    echo -e "${YELLOW}⚠ Serviço instalado mas não está rodando${NC}"
    echo -e "${YELLOW}  Verifique os logs: journalctl -u $SERVICE_NAME -n 50${NC}"
fi

echo ""
echo -e "${GREEN}=== Instalação Completa ===${NC}"
echo ""
echo -e "${CYAN}Informações do Serviço:${NC}"
echo "  Nome:         $SERVICE_NAME"
echo "  Binário:      $INSTALL_DIR/procspy-watcher"
echo "  Configuração: $CONFIG_DIR/watcher-config.json"
echo "  Logs:         $LOG_DIR"
echo "  Service file: $SERVICE_FILE"
echo ""
echo -e "${CYAN}Comandos Úteis:${NC}"
echo "  systemctl status $SERVICE_NAME"
echo "  sudo systemctl restart $SERVICE_NAME"
echo "  journalctl -u $SERVICE_NAME -f"
echo "  sudo micro $CONFIG_DIR/watcher-config.json"
echo ""
echo -e "${GREEN}✓ Procspy Watcher instalado com sucesso!${NC}"
