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

# Configurações
COVERAGE_DIR="coverage"
COVERAGE_FILE="${COVERAGE_DIR}/coverage.out"
COVERAGE_HTML="${COVERAGE_DIR}/coverage.html"
MIN_COVERAGE=70

# Cria diretório de coverage
mkdir -p "$COVERAGE_DIR"

# Função para exibir cabeçalho
print_header() {
    log_msg "${BLUE}=== Procspy - Testes Unitários ===${NC}"
    echo ""
}

# Função para executar testes
run_tests() {
    start_timer "tests"
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
        local duration=$(end_timer "tests")
        log_msg "${GREEN}✓ Todos os testes passaram${NC}" "$duration"
        return 0
    else
        echo ""
        local duration=$(end_timer "tests")
        log_msg "${RED}✗ Alguns testes falharam${NC}" "$duration"
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
    start_timer "total"
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
    
    local total_duration=$(end_timer "total")
    log_msg "${GREEN}✓ Execução de testes concluída com sucesso!${NC}" "$total_duration"
    exit 0
}

# Executa main
main
