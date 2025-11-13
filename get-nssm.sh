#!/bin/bash
# get-nssm.sh - Script para baixar NSSM (Non-Sucking Service Manager)
# Baixa a versão mais recente do NSSM e extrai para /opt/nssm

set -e

# Cores para output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

# Configurações
NSSM_VERSION="2.24"
NSSM_URL="https://nssm.cc/release/nssm-${NSSM_VERSION}.zip"
DOWNLOAD_DIR="/tmp"
INSTALL_DIR="/opt"
NSSM_DIR="${INSTALL_DIR}/nssm"

echo -e "${GREEN}=== NSSM Download Script ===${NC}"
echo "Version: ${NSSM_VERSION}"
echo "Install Directory: ${NSSM_DIR}"
echo ""

# Verifica se está rodando no Windows (WSL ou Git Bash)
if [[ "$OSTYPE" == "msys" || "$OSTYPE" == "win32" || -n "$WINDIR" ]]; then
    echo -e "${YELLOW}Detectado ambiente Windows${NC}"
    INSTALL_DIR="C:/nssm"
    NSSM_DIR="${INSTALL_DIR}"
fi

# Verifica dependências
echo "Verificando dependências..."

if ! command -v curl &> /dev/null && ! command -v wget &> /dev/null; then
    echo -e "${RED}Erro: curl ou wget não encontrado${NC}"
    echo "Instale com: sudo apt install curl (ou wget)"
    exit 1
fi

if ! command -v unzip &> /dev/null; then
    echo -e "${RED}Erro: unzip não encontrado${NC}"
    echo "Instale com: sudo apt install unzip"
    exit 1
fi

echo -e "${GREEN}✓ Dependências OK${NC}"
echo ""

# Cria diretório de instalação
echo "Criando diretório de instalação..."
mkdir -p "$INSTALL_DIR"

# Baixa NSSM
echo -e "${YELLOW}Baixando NSSM ${NSSM_VERSION}...${NC}"
cd "$DOWNLOAD_DIR"

if command -v curl &> /dev/null; then
    curl -L -o "nssm-${NSSM_VERSION}.zip" "$NSSM_URL"
else
    wget -O "nssm-${NSSM_VERSION}.zip" "$NSSM_URL"
fi

echo -e "${GREEN}✓ Download completo${NC}"
echo ""

# Extrai arquivo
echo "Extraindo arquivo..."
unzip -q "nssm-${NSSM_VERSION}.zip"

# Move para diretório de instalação
echo "Instalando em ${NSSM_DIR}..."
rm -rf "$NSSM_DIR"
mv "nssm-${NSSM_VERSION}" "$NSSM_DIR"

# Limpa arquivos temporários
echo "Limpando arquivos temporários..."
rm "nssm-${NSSM_VERSION}.zip"

echo ""
echo -e "${GREEN}=== Instalação Completa ===${NC}"
echo ""
echo "NSSM instalado em: ${NSSM_DIR}"
echo ""
echo "Estrutura:"
echo "  ${NSSM_DIR}/"
echo "    ├── win32/nssm.exe    (32-bit)"
echo "    └── win64/nssm.exe    (64-bit)"
echo ""
echo "Para usar o NSSM:"
echo ""
echo "  Windows 64-bit:"
echo "    ${NSSM_DIR}/win64/nssm.exe install <service-name> <program>"
echo ""
echo "  Windows 32-bit:"
echo "    ${NSSM_DIR}/win32/nssm.exe install <service-name> <program>"
echo ""
echo "Ou adicione ao PATH:"
echo "  export PATH=\"\$PATH:${NSSM_DIR}/win64\""
echo ""
echo -e "${GREEN}✓ NSSM instalado com sucesso!${NC}"
