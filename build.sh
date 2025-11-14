#!/bin/bash
# build.sh - Script de build do Procspy
# Compila todos os componentes para múltiplas arquiteturas e executa testes

set -e

# Cores para output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Configurações
BUILD_DIR="bin"
VERSION=$(git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_DATE=$(date -u +"%Y-%m-%d %H:%M:%S UTC")
LDFLAGS="-X 'main.buildDate=${BUILD_DATE}' -X 'main.version=${VERSION}'"

# Plataformas suportadas
PLATFORMS=(
    "linux/amd64"
    "linux/arm64"
    "linux/386"
    "windows/amd64"
    "windows/386"
    "darwin/amd64"
    "darwin/arm64"
)

echo -e "${GREEN}=== Procspy Build Script ===${NC}"
echo "Version: ${VERSION}"
echo "Build Date: ${BUILD_DATE}"
echo ""

# Função para exibir uso
usage() {
    echo "Uso: $0 [opções]"
    echo ""
    echo "Opções:"
    echo "  -h, --help          Exibe esta ajuda"
    echo "  -t, --test-only     Apenas executa testes"
    echo "  -b, --build-only    Apenas compila (pula testes)"
    echo "  -p, --platform OS/ARCH  Compila apenas para plataforma específica"
    echo "  -a, --all           Compila para todas as plataformas"
    echo "  -c, --clean         Limpa diretório de build antes"
    echo ""
    echo "Exemplos:"
    echo "  $0                  # Testa e compila para plataforma atual"
    echo "  $0 --all            # Testa e compila para todas as plataformas"
    echo "  $0 -p linux/amd64   # Compila apenas para Linux AMD64"
    echo "  $0 --test-only      # Apenas executa testes"
    exit 0
}



# Função para compilar um componente
build_component() {
    local component=$1
    local goos=$2
    local goarch=$3
    
    local output_name="${component}"
    if [ "$goos" = "windows" ]; then
        output_name="${component}.exe"
    fi
    
    local output_dir="${BUILD_DIR}/${goos}_${goarch}"
    local output_path="${output_dir}/${output_name}"
    
    mkdir -p "$output_dir"
    
    echo -n "  Building ${component} for ${goos}/${goarch}... "
    
    if GOOS=$goos GOARCH=$goarch go build \
        -ldflags "${LDFLAGS}" \
        -o "$output_path" \
        "cmd/${component}/main.go"; then
        echo -e "${GREEN}✓${NC}"
        return 0
    else
        echo -e "${RED}✗${NC}"
        return 1
    fi
}

# Função para compilar todos os componentes para uma plataforma
build_platform() {
    local platform=$1
    local goos=$(echo $platform | cut -d'/' -f1)
    local goarch=$(echo $platform | cut -d'/' -f2)
    
    echo -e "${YELLOW}=== Building for ${goos}/${goarch} ===${NC}"
    
    # Client
    build_component "client" "$goos" "$goarch" || return 1
    
    # Watcher
    build_component "watcher" "$goos" "$goarch" || return 1
    
    # Server (apenas Linux)
    if [ "$goos" = "linux" ]; then
        build_component "server" "$goos" "$goarch" || return 1
    fi
    
    echo ""
    return 0
}

# Função para verificar formatação do código
check_formatting() {
    echo -e "${YELLOW}Verificando formatação do código (go fmt)...${NC}"
    
    # Lista arquivos não formatados
    UNFORMATTED=$(gofmt -l . 2>&1 | grep -v "^vendor/" || true)
    
    if [ -n "$UNFORMATTED" ]; then
        echo -e "${RED}✗ Arquivos não formatados encontrados:${NC}"
        echo "$UNFORMATTED"
        echo ""
        echo -e "${YELLOW}Execute 'go fmt ./...' para corrigir${NC}"
        return 1
    else
        echo -e "${GREEN}✓ Todos os arquivos estão formatados corretamente${NC}"
        echo ""
        return 0
    fi
}

# Função para executar linter
run_linter() {
    echo -e "${YELLOW}Executando linter (go vet)...${NC}"
    
    if go vet ./...; then
        echo -e "${GREEN}✓ Nenhum problema encontrado pelo linter${NC}"
        echo ""
        return 0
    else
        echo -e "${RED}✗ Problemas encontrados pelo linter${NC}"
        echo ""
        return 1
    fi
}

# Função para executar verificações de qualidade
run_quality_checks() {
    echo -e "${YELLOW}=== Verificações de Qualidade ===${NC}"
    echo ""
    
    # Verifica formatação
    if ! check_formatting; then
        return 1
    fi
    
    # Executa linter
    if ! run_linter; then
        return 1
    fi
    
    echo -e "${GREEN}✓ Todas as verificações de qualidade passaram${NC}"
    echo ""
    return 0
}

# Função para limpar diretório de build
clean_build() {
    echo -e "${YELLOW}Limpando diretório de build...${NC}"
    rm -rf "$BUILD_DIR"
    rm -f coverage.out
    echo -e "${GREEN}✓ Diretório limpo${NC}"
    echo ""
}

# Parse de argumentos
TEST_ONLY=false
BUILD_ONLY=false
BUILD_ALL=false
CLEAN=false
SPECIFIC_PLATFORM=""

while [[ $# -gt 0 ]]; do
    case $1 in
        -h|--help)
            usage
            ;;
        -t|--test-only)
            TEST_ONLY=true
            shift
            ;;
        -b|--build-only)
            BUILD_ONLY=true
            shift
            ;;
        -a|--all)
            BUILD_ALL=true
            shift
            ;;
        -p|--platform)
            SPECIFIC_PLATFORM="$2"
            shift 2
            ;;
        -c|--clean)
            CLEAN=true
            shift
            ;;
        *)
            echo -e "${RED}Opção desconhecida: $1${NC}"
            usage
            ;;
    esac
done

# Limpa se solicitado
if [ "$CLEAN" = true ]; then
    clean_build
fi

# Executa verificações de qualidade se não for build-only
if [ "$BUILD_ONLY" = false ]; then
    if ! run_quality_checks; then
        echo -e "${RED}Build abortado devido a falhas nas verificações de qualidade${NC}"
        exit 1
    fi
fi

# Executa testes se não for build-only
if [ "$BUILD_ONLY" = false ]; then
    if ! ./test.sh; then
        echo -e "${RED}Build abortado devido a falhas nos testes${NC}"
        exit 1
    fi
    echo ""
fi

# Se for test-only, para aqui
if [ "$TEST_ONLY" = true ]; then
    exit 0
fi

# Cria diretório de build
mkdir -p "$BUILD_DIR"

# Determina o que compilar
if [ -n "$SPECIFIC_PLATFORM" ]; then
    # Plataforma específica
    build_platform "$SPECIFIC_PLATFORM"
elif [ "$BUILD_ALL" = true ]; then
    # Todas as plataformas
    echo -e "${YELLOW}=== Building for all platforms ===${NC}"
    echo ""
    
    for platform in "${PLATFORMS[@]}"; do
        build_platform "$platform" || {
            echo -e "${RED}Build failed for $platform${NC}"
            exit 1
        }
    done
else
    # Apenas plataforma atual
    CURRENT_OS=$(go env GOOS)
    CURRENT_ARCH=$(go env GOARCH)
    build_platform "${CURRENT_OS}/${CURRENT_ARCH}"
fi

# Resumo
echo -e "${GREEN}=== Build Completo ===${NC}"
echo "Binários gerados em: ${BUILD_DIR}/"
echo ""
echo "Estrutura:"
tree -L 2 "$BUILD_DIR" 2>/dev/null || find "$BUILD_DIR" -type f

echo ""
echo -e "${GREEN}✓ Build concluído com sucesso!${NC}"
