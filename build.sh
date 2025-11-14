#!/bin/bash
# build.sh - Script de build do Procspy
# Compila todos os componentes para múltiplas arquiteturas e executa testes

set -e

# Cores para output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
BOLD='\033[1m'
NC='\033[0m' # No Color

# Configurações
BUILD_DIR="bin"
VERSION=$(git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_DATE=$(date -u +"%Y-%m-%d %H:%M:%S UTC")
LDFLAGS="-X 'main.buildDate=${BUILD_DATE}' -X 'main.version=${VERSION}'"

# Função para logging com timestamp e duração opcional
# Uso: log_msg "mensagem" [duração]
# Exemplo: log_msg "Teste completo" "1.5s"
log_msg() {
    local message="$1"
    local duration="$2"
    local timestamp=$(date '+%Y-%m-%d %H:%M:%S')
    
    if [ -n "$duration" ]; then
        # Com duração: [timestamp] duration (alinhado à direita, 12 chars) :: mensagem
        printf "${GREEN}${BOLD}[%s]${NC} ${BLUE}%12s${NC} :: %b\n" "$timestamp" "$duration" "$message"
    else
        # Sem duração: [timestamp] mensagem
        echo -e "${GREEN}${BOLD}[${timestamp}]${NC} ${message}"
    fi
}

# ============================================
# Funções de Medição de Tempo
# ============================================

# Formata duração em nanosegundos para formato legível
# Parâmetro: duração em nanosegundos
# Retorna: string formatada com unidade apropriada (ns, ms, s)
format_duration() {
    local ns=$1
    
    # Valida input
    if ! [[ "$ns" =~ ^[0-9]+$ ]]; then
        echo "N/A"
        return
    fi
    
    # Verifica se bc está disponível
    if ! command -v bc &> /dev/null; then
        # Fallback: usar apenas divisão inteira
        if [ $ns -lt 1000000 ]; then
            echo "${ns}ns"
        elif [ $ns -lt 1000000000 ]; then
            local ms=$((ns / 1000000))
            echo "${ms}ms"
        else
            local s=$((ns / 1000000000))
            echo "${s}s"
        fi
        return
    fi
    
    # Menos de 1ms: exibir em nanosegundos
    if [ $ns -lt 1000000 ]; then
        echo "${ns}ns"
    # Entre 1ms e 1s: exibir em milissegundos
    elif [ $ns -lt 1000000000 ]; then
        local ms=$(echo "scale=2; $ns / 1000000" | bc)
        echo "${ms}ms"
    # 1s ou mais: exibir em segundos
    else
        local s=$(echo "scale=2; $ns / 1000000000" | bc)
        echo "${s}s"
    fi
}

# Inicia medição de tempo para uma operação
# Parâmetro: nome da operação (usado como chave)
# Armazena tempo em nanosegundos em variável global
start_timer() {
    local timer_name=$1
    local var_name="TIMER_${timer_name}"
    eval "${var_name}=$(date +%s%N)"
}

# Finaliza medição e retorna duração formatada
# Parâmetro: nome da operação
# Retorna: string formatada (ex: "123ms", "1.5s")
end_timer() {
    local timer_name=$1
    local var_name="TIMER_${timer_name}"
    local start_time=$(eval echo \${${var_name}})
    
    # Verifica se timer foi iniciado
    if [ -z "$start_time" ]; then
        echo "N/A"
        return
    fi
    
    local end_time=$(date +%s%N)
    local duration=$((end_time - start_time))
    
    format_duration $duration
}

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

log_msg "${GREEN}=== Procspy Build Script ===${NC}"
log_msg "Version: ${VERSION}"
log_msg "Build Date: ${BUILD_DATE}"
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
    echo "  --no-fmt            Pula formatação automática do código"
    echo ""
    echo "Exemplos:"
    echo "  $0                  # Testa e compila para plataforma atual"
    echo "  $0 --all            # Testa e compila para todas as plataformas"
    echo "  $0 -p linux/amd64   # Compila apenas para Linux AMD64"
    echo "  $0 --test-only      # Apenas executa testes"
    echo "  $0 --no-fmt         # Compila sem formatar código"
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
    
    local timer_name="build_${component}_${goos}_${goarch}"
    start_timer "$timer_name"
    
    if GOOS=$goos GOARCH=$goarch go build \
        -ldflags "${LDFLAGS}" \
        -o "$output_path" \
        "cmd/${component}/main.go" 2>/dev/null; then
        local duration=$(end_timer "$timer_name")
        log_msg "  ${GREEN}✓${NC} Building ${component} for ${goos}/${goarch}" "$duration"
        return 0
    else
        local duration=$(end_timer "$timer_name")
        log_msg "  ${RED}✗${NC} Building ${component} for ${goos}/${goarch}" "$duration"
        return 1
    fi
}

# Função para compilar todos os componentes para uma plataforma
build_platform() {
    local platform=$1
    local goos=$(echo $platform | cut -d'/' -f1)
    local goarch=$(echo $platform | cut -d'/' -f2)
    
    local timer_name="platform_${goos}_${goarch}"
    start_timer "$timer_name"
    
    log_msg "${YELLOW}=== Building for ${goos}/${goarch} ===${NC}"
    
    # Client
    build_component "client" "$goos" "$goarch" || return 1
    
    # Watcher
    build_component "watcher" "$goos" "$goarch" || return 1
    
    # Server (apenas Linux)
    if [ "$goos" = "linux" ]; then
        build_component "server" "$goos" "$goarch" || return 1
    fi
    
    local duration=$(end_timer "$timer_name")
    log_msg "${GREEN}✓ Platform ${goos}/${goarch} build complete${NC}" "$duration"
    echo ""
    return 0
}

# Função para verificar formatação do código
check_formatting() {
    log_msg "${YELLOW}Verificando formatação do código (go fmt)...${NC}"
    
    # Lista arquivos não formatados
    UNFORMATTED=$(gofmt -l . 2>&1 | grep -v "^vendor/" || true)
    
    if [ -n "$UNFORMATTED" ]; then
        log_msg "${RED}✗ Arquivos não formatados encontrados:${NC}"
        echo "$UNFORMATTED"
        echo ""
        log_msg "${YELLOW}Execute 'go fmt ./...' para corrigir${NC}"
        return 1
    else
        log_msg "${GREEN}✓ Todos os arquivos estão formatados corretamente${NC}"
        echo ""
        return 0
    fi
}

# Função para executar linter
run_linter() {
    log_msg "${YELLOW}Executando linter (go vet)...${NC}"
    
    if go vet ./...; then
        log_msg "${GREEN}✓ Nenhum problema encontrado pelo linter${NC}"
        echo ""
        return 0
    else
        log_msg "${RED}✗ Problemas encontrados pelo linter${NC}"
        echo ""
        return 1
    fi
}

# Função para formatação automática
auto_format() {
    start_timer "format"
    log_msg "${YELLOW}Formatando código Go automaticamente...${NC}"
    
    if go fmt ./...; then
        local duration=$(end_timer "format")
        log_msg "${GREEN}✓ Código formatado com sucesso${NC}" "$duration"
        echo ""
        return 0
    else
        local duration=$(end_timer "format")
        log_msg "${RED}✗ Erro ao formatar código${NC}" "$duration"
        echo ""
        return 1
    fi
}

# Função para executar verificações de qualidade
run_quality_checks() {
    start_timer "quality"
    log_msg "${YELLOW}=== Verificações de Qualidade ===${NC}"
    echo ""
    
    # Verifica formatação
    if ! check_formatting; then
        return 1
    fi
    
    # Executa linter
    if ! run_linter; then
        return 1
    fi
    
    local duration=$(end_timer "quality")
    log_msg "${GREEN}✓ Todas as verificações de qualidade passaram${NC}" "$duration"
    echo ""
    return 0
}

# Função para limpar diretório de build
clean_build() {
    log_msg "${YELLOW}Limpando diretório de build...${NC}"
    rm -rf "$BUILD_DIR"
    rm -rf coverage/
    log_msg "${GREEN}✓ Diretório limpo${NC}"
    echo ""
}

# Parse de argumentos
TEST_ONLY=false
BUILD_ONLY=false
BUILD_ALL=false
CLEAN=false
NO_FMT=false
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
        --no-fmt)
            NO_FMT=true
            shift
            ;;
        *)
            log_msg "${RED}Opção desconhecida: $1${NC}"
            usage
            ;;
    esac
done

# Inicia timer total
start_timer "total"

# Limpa se solicitado
if [ "$CLEAN" = true ]; then
    clean_build
fi

# Executa formatação automática se não for build-only e não tiver --no-fmt
if [ "$BUILD_ONLY" = false ] && [ "$NO_FMT" = false ]; then
    if ! auto_format; then
        log_msg "${RED}Build abortado devido a falha na formatação${NC}"
        exit 1
    fi
fi

# Executa verificações de qualidade se não for build-only
if [ "$BUILD_ONLY" = false ]; then
    if ! run_quality_checks; then
        log_msg "${RED}Build abortado devido a falhas nas verificações de qualidade${NC}"
        exit 1
    fi
fi

# Executa testes se não for build-only
if [ "$BUILD_ONLY" = false ]; then
    if ! ./test.sh; then
        log_msg "${RED}Build abortado devido a falhas nos testes${NC}"
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
    log_msg "${YELLOW}=== Building for all platforms ===${NC}"
    echo ""
    
    for platform in "${PLATFORMS[@]}"; do
        build_platform "$platform" || {
            log_msg "${RED}Build failed for $platform${NC}"
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
log_msg "${GREEN}=== Build Completo ===${NC}"
log_msg "Binários gerados em: ${BUILD_DIR}/"
echo ""
log_msg "Estrutura:"
tree -L 2 "$BUILD_DIR" 2>/dev/null || find "$BUILD_DIR" -type f

echo ""
TOTAL_DURATION=$(end_timer "total")
log_msg "${GREEN}✓ Build concluído com sucesso!${NC}" "$TOTAL_DURATION"
