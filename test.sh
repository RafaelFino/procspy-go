#!/bin/bash
# test.sh - Script de execução de testes unitários do Procspy
# Executa todos os testes e gera relatórios de coverage

set -e

# Cores para output
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
BOLD='\033[1m'
NC='\033[0m' # No Color

# Função para logging com timestamp
log_msg() {
    local message="$1"
    local timestamp=$(date '+%Y-%m-%d %H:%M:%S')
    echo -e "${GREEN}${BOLD}[${timestamp}]${NC} ${message}"
}

# Configurações
COVERAGE_FILE="coverage.out"
COVERAGE_HTML="coverage.html"
MIN_COVERAGE=70

# Função para exibir cabeçalho
print_header() {
    log_msg "${BLUE}=== Procspy - Testes Unitários ===${NC}"
    echo ""
}

# Função para executar testes
run_tests() {
    log_msg "${YELLOW}Executando testes...${NC}"
    echo ""
    
    # Verifica se há arquivos de teste
    if ! find . -name "*_test.go" -type f | grep -q .; then
        log_msg "${YELLOW}Aviso: Nenhum arquivo de teste encontrado${NC}"
        return 0
    fi
    
    # Executa testes com coverage e race detector
    if go test -v -race -coverprofile="$COVERAGE_FILE" ./...; then
        echo ""
        log_msg "${GREEN}✓ Todos os testes passaram${NC}"
        return 0
    else
        echo ""
        log_msg "${RED}✗ Alguns testes falharam${NC}"
        return 1
    fi
}

# Função para gerar relatório de coverage
generate_coverage() {
    if [ ! -f "$COVERAGE_FILE" ]; then
        log_msg "${YELLOW}Arquivo de coverage não encontrado${NC}"
        return 0
    fi
    
    echo ""
    log_msg "${YELLOW}=== Relatório de Coverage ===${NC}"
    echo ""
    
    # Calcula coverage total
    COVERAGE=$(go tool cover -func="$COVERAGE_FILE" | grep total | awk '{print $3}')
    COVERAGE_NUM=$(echo "$COVERAGE" | sed 's/%//')
    
    log_msg "Coverage Total: ${GREEN}${COVERAGE}${NC}"
    echo ""
    
    # Verifica se atingiu o mínimo
    if (( $(echo "$COVERAGE_NUM < $MIN_COVERAGE" | bc -l) )); then
        log_msg "${YELLOW}⚠ Warning: Coverage abaixo do mínimo recomendado (${MIN_COVERAGE}%)${NC}"
    else
        log_msg "${GREEN}✓ Coverage acima do mínimo recomendado (${MIN_COVERAGE}%)${NC}"
    fi
    
    echo ""
}

# Função para exibir estatísticas
display_stats() {
    log_msg "${YELLOW}=== Estatísticas de Execução ===${NC}"
    echo ""
    
    # Conta testes
    TOTAL_TESTS=$(go test -list . ./... 2>/dev/null | grep -c "^Test" || echo "0")
    
    log_msg "Total de testes: ${TOTAL_TESTS}"
    
    if [ -f "$COVERAGE_FILE" ]; then
        # Conta pacotes testados
        PACKAGES=$(go list ./... | wc -l)
        log_msg "Pacotes testados: ${PACKAGES}"
    fi
    
    echo ""
}

# Função para gerar HTML de coverage (opcional)
generate_html() {
    if [ -f "$COVERAGE_FILE" ]; then
        log_msg "${BLUE}Para visualizar coverage detalhado, execute:${NC}"
        echo "  go tool cover -html=$COVERAGE_FILE"
        echo ""
    fi
}

# Main
main() {
    print_header
    
    # Executa testes
    if ! run_tests; then
        exit 1
    fi
    
    # Gera relatório de coverage
    generate_coverage
    
    # Exibe estatísticas
    display_stats
    
    # Informação sobre HTML
    generate_html
    
    log_msg "${GREEN}✓ Execução de testes concluída com sucesso!${NC}"
    exit 0
}

# Executa main
main
