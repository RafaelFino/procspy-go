#!/bin/bash
# make-changelog.sh - Script de geração automática de CHANGELOG
# Analisa histórico Git e gera CHANGELOG.md estruturado

set -e

# Cores para output
RED='\033[0;31m'
GREEN='\033[0;32m'
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

# ============================================
# Funções de Análise Git
# ============================================

# Obtém todos os commits do mais recente ao mais antigo
# Retorna: hash|short_hash|datetime|message|author (um por linha)
get_all_commits() {
    # Verifica se está em repositório Git
    if ! git rev-parse --git-dir > /dev/null 2>&1; then
        log_msg "${RED}✗ Erro: Não é um repositório Git${NC}"
        exit 1
    fi
    
    # Verifica se há commits
    if ! git log -1 > /dev/null 2>&1; then
        log_msg "${YELLOW}⚠ Aviso: Repositório sem commits${NC}"
        exit 0
    fi
    
    # Lista commits do mais recente ao mais antigo com data/hora completa
    git log --format="%H|%h|%ad|%s|%an" --date=format:"%Y-%m-%d %H:%M:%S"
}

# Obtém arquivos modificados em um commit
# Parâmetro: hash do commit
# Retorna: lista de arquivos com status (A/M/D filename)
get_commit_files() {
    local hash=$1
    
    # Obtém lista de arquivos modificados
    local files=$(git show --name-status --format="" "$hash" 2>/dev/null)
    
    # Se não houver arquivos (merge vazio), retorna mensagem
    if [ -z "$files" ]; then
        echo "Commit vazio ou merge"
    else
        echo "$files"
    fi
}

# Classifica tipo de mudança baseado em arquivos modificados
# Parâmetro: lista de arquivos
# Retorna: tipo da mudança (Testes, Documentação, Configuração, Implementação, Outros)
classify_change() {
    local files=$1
    
    # Verifica padrões nos nomes de arquivos
    if echo "$files" | grep -q "_test\.go\|test\.sh"; then
        echo "Testes"
    elif echo "$files" | grep -q "\.md$\|docs/"; then
        echo "Documentação"
    elif echo "$files" | grep -q "\.sh$\|Makefile\|makefile\|\.yml$\|\.yaml$\|\.json$"; then
        echo "Configuração"
    elif echo "$files" | grep -q "\.go$\|internal/\|cmd/\|pkg/"; then
        echo "Implementação"
    else
        echo "Outros"
    fi
}

# Gera resumo objetivo baseado em arquivos modificados
# Parâmetros: lista de arquivos, tipo de mudança
# Retorna: resumo em PT-BR focado em resultados
generate_summary() {
    local files=$1
    local change_type=$2
    
    # Trata caso de commit vazio
    if [ "$files" = "Commit vazio ou merge" ]; then
        echo "Merge ou commit sem alterações de arquivos"
        return
    fi
    
    # Conta arquivos por tipo
    local go_files=$(echo "$files" | grep -c "\.go$" 2>/dev/null || echo "0")
    go_files=$(echo "$go_files" | tr -d '\n' | tr -d ' ')
    local test_files=$(echo "$files" | grep -c "_test\.go$" 2>/dev/null || echo "0")
    test_files=$(echo "$test_files" | tr -d '\n' | tr -d ' ')
    local doc_files=$(echo "$files" | grep -c "\.md$" 2>/dev/null || echo "0")
    doc_files=$(echo "$doc_files" | tr -d '\n' | tr -d ' ')
    local sh_files=$(echo "$files" | grep -c "\.sh$" 2>/dev/null || echo "0")
    sh_files=$(echo "$sh_files" | tr -d '\n' | tr -d ' ')
    local total_files=$(echo "$files" | wc -l)
    total_files=$(echo "$total_files" | tr -d '\n' | tr -d ' ')
    
    # Gera resumo baseado no tipo
    case "$change_type" in
        "Testes")
            if [ "$test_files" -gt 0 ]; then
                echo "Adicionados/atualizados testes ($test_files arquivo(s))"
            else
                echo "Atualização de scripts de teste ($sh_files arquivo(s))"
            fi
            ;;
        "Documentação")
            echo "Atualização de documentação ($doc_files arquivo(s))"
            ;;
        "Configuração")
            echo "Mudanças em configuração e scripts ($total_files arquivo(s))"
            ;;
        "Implementação")
            local impl_files=$((go_files - test_files))
            if [ $impl_files -gt 0 ]; then
                echo "Implementação de funcionalidades ($impl_files arquivo(s) Go)"
            else
                echo "Mudanças no código ($go_files arquivo(s))"
            fi
            ;;
        *)
            echo "Mudanças gerais no código ($total_files arquivo(s))"
            ;;
    esac
}

# ============================================
# Funções de Geração de Markdown
# ============================================

# Gera cabeçalho do documento CHANGELOG
generate_header() {
    cat << 'EOF'
# CHANGELOG

Este documento contém o histórico completo de mudanças do projeto.
Gerado automaticamente a partir do histórico Git.

EOF
}

# Gera tabela resumo de mudanças
# Lê commits do stdin (formato: hash|short_hash|datetime|message|author)
generate_table() {
    echo "## Resumo de Mudanças"
    echo ""
    echo "| Data/Hora | Commit | Mensagem | Resumo das Mudanças |"
    echo "|-----------|--------|----------|---------------------|"
    
    # Para cada commit, adiciona linha na tabela
    while IFS='|' read -r hash short_hash datetime message author; do
        # Obtém arquivos e classifica mudança
        local files=$(get_commit_files "$hash")
        local change_type=$(classify_change "$files")
        local summary=$(generate_summary "$files" "$change_type")
        
        # Escapa pipes na mensagem para não quebrar tabela
        local safe_message=$(echo "$message" | sed 's/|/\\|/g')
        
        echo "| $datetime | \`$short_hash\` | $safe_message | $summary |"
    done
    
    echo ""
}

# Gera diagrama Mermaid gitGraph
# Lê commits do stdin (formato: hash|short_hash|datetime|message|author)
generate_mermaid() {
    echo "## Timeline de Commits"
    echo ""
    echo '```mermaid'
    echo "gitGraph"
    
    # Para cada commit, adiciona ao gitGraph
    while IFS='|' read -r hash short_hash datetime message author; do
        # Extrai apenas a data (YYYY-MM-DD) do datetime
        local date=$(echo "$datetime" | cut -d' ' -f1)
        # Escapa aspas duplas na mensagem
        local safe_message=$(echo "$message" | sed 's/"/\\"/g')
        echo "    commit id: \"$short_hash\" tag: \"$date\""
    done
    
    echo '```'
    echo ""
}

# Gera seções detalhadas para cada commit
# Lê commits do stdin (formato: hash|short_hash|datetime|message|author)
generate_details() {
    echo "## Detalhes dos Commits"
    echo ""
    
    while IFS='|' read -r hash short_hash datetime message author; do
        # Obtém informações do commit
        local files=$(get_commit_files "$hash")
        local change_type=$(classify_change "$files")
        local summary=$(generate_summary "$files" "$change_type")
        
        # Conta arquivos modificados
        if [ "$files" = "Commit vazio ou merge" ]; then
            local file_count=0
        else
            local file_count=$(echo "$files" | wc -l)
        fi
        
        # Gera seção do commit
        echo "### [\`$short_hash\`] - $datetime"
        echo ""
        echo "**Mensagem:** $message"
        echo ""
        echo "**Autor:** $author"
        echo ""
        echo "**Tipo:** $change_type"
        echo ""
        echo "**Resumo:** $summary"
        echo ""
        echo "**Arquivos modificados:** $file_count"
        echo ""
        echo "---"
        echo ""
    done
}

# ============================================
# Fluxo Principal
# ============================================

# Validações iniciais
validate_environment() {
    # Verifica se está em repositório Git
    if ! git rev-parse --git-dir > /dev/null 2>&1; then
        log_msg "${RED}✗ Erro: Não é um repositório Git${NC}"
        exit 1
    fi
    
    # Verifica se há commits
    if ! git log -1 > /dev/null 2>&1; then
        log_msg "${YELLOW}⚠ Aviso: Repositório sem commits${NC}"
        exit 0
    fi
    
    # Verifica permissão de escrita
    if [ ! -w . ]; then
        log_msg "${RED}✗ Erro: Sem permissão de escrita no diretório atual${NC}"
        exit 1
    fi
    
    # Avisa se vai sobrescrever arquivo existente
    if [ -f CHANGELOG.md ]; then
        log_msg "${YELLOW}Sobrescrevendo CHANGELOG.md existente${NC}"
    fi
}

# ============================================
# Execução Principal
# ============================================

log_msg "${GREEN}=== Gerador de CHANGELOG ===${NC}"
echo ""

# Inicia timer total
start_timer "total"

# Validações
validate_environment

# Coleta de dados
log_msg "${YELLOW}Analisando histórico Git...${NC}"
start_timer "git_analysis"

COMMITS=$(get_all_commits)
COMMIT_COUNT=$(echo "$COMMITS" | wc -l)

duration=$(end_timer "git_analysis")
log_msg "${GREEN}✓ Encontrados $COMMIT_COUNT commits${NC}" "$duration"
echo ""

# Geração do documento
log_msg "${YELLOW}Gerando CHANGELOG.md...${NC}"
start_timer "generation"

{
    generate_header
    echo "$COMMITS" | generate_mermaid
    echo "$COMMITS" | generate_table
    echo "$COMMITS" | generate_details
} > CHANGELOG.md

duration=$(end_timer "generation")
log_msg "${GREEN}✓ CHANGELOG.md gerado com sucesso${NC}" "$duration"
echo ""

# Finalização
total_duration=$(end_timer "total")
log_msg "${GREEN}✓ Geração completa!${NC}" "$total_duration"

# ============================================
# Funções de Análise Git
# ============================================

# Obtém todos os commits do mais recente ao mais antigo
# Formato: hash|short_hash|datetime|message|author
get_all_commits() {
    # Verifica se está em repositório Git
    if ! git rev-parse --git-dir > /dev/null 2>&1; then
        log_msg "${RED}✗ Erro: Não é um repositório Git${NC}"
        exit 1
    fi
    
    # Verifica se há commits
    if ! git log -1 > /dev/null 2>&1; then
        log_msg "${YELLOW}⚠ Aviso: Repositório sem commits${NC}"
        exit 0
    fi
    
    # Lista commits do mais recente ao mais antigo (sem --reverse)
    git log --format="%H|%h|%ad|%s|%an" --date=format:"%Y-%m-%d %H:%M:%S"
}

# Obtém arquivos modificados em um commit
# Parâmetro: hash do commit
# Retorna: lista de arquivos com status (A/M/D)
get_commit_files() {
    local hash=$1
    
    # Obtém lista de arquivos modificados
    local files=$(git show --name-status --format="" "$hash" 2>/dev/null)
    
    # Trata commits sem arquivos (merges vazios)
    if [ -z "$files" ]; then
        echo "Commit vazio ou merge"
    else
        echo "$files"
    fi
}

# Classifica tipo de mudança baseado em arquivos modificados
# Parâmetro: lista de arquivos
# Retorna: tipo de mudança (Testes, Documentação, Configuração, Implementação, Outros)
classify_change() {
    local files=$1
    
    # Verifica padrões nos nomes de arquivos
    if echo "$files" | grep -q "_test\.go\|test\.sh"; then
        echo "Testes"
    elif echo "$files" | grep -q "\.md$\|docs/"; then
        echo "Documentação"
    elif echo "$files" | grep -q "\.sh$\|Makefile\|makefile\|\.yml$\|\.yaml$\|\.json$"; then
        echo "Configuração"
    elif echo "$files" | grep -q "\.go$\|internal/\|cmd/\|pkg/"; then
        echo "Implementação"
    else
        echo "Outros"
    fi
}

# Gera resumo objetivo baseado em arquivos modificados
# Parâmetros: lista de arquivos, tipo de mudança
# Retorna: resumo conciso em PT-BR
generate_summary() {
    local files=$1
    local change_type=$2
    
    # Trata caso de commit vazio
    if [ "$files" = "Commit vazio ou merge" ]; then
        echo "Merge ou commit sem alterações de arquivos"
        return
    fi
    
    # Conta arquivos por tipo (remove espaços em branco do resultado)
    local go_files=$(echo "$files" | grep -c "\.go$" 2>/dev/null | tr -d '[:space:]' || echo "0")
    local test_files=$(echo "$files" | grep -c "_test\.go$" 2>/dev/null | tr -d '[:space:]' || echo "0")
    local doc_files=$(echo "$files" | grep -c "\.md$" 2>/dev/null | tr -d '[:space:]' || echo "0")
    local sh_files=$(echo "$files" | grep -c "\.sh$" 2>/dev/null | tr -d '[:space:]' || echo "0")
    local total_files=$(echo "$files" | wc -l | tr -d '[:space:]')
    
    # Gera resumo baseado no tipo
    case "$change_type" in
        "Testes")
            if [ "$test_files" -gt 0 ]; then
                echo "Adicionados/atualizados testes ($test_files arquivo(s))"
            else
                echo "Atualização de testes ($total_files arquivo(s))"
            fi
            ;;
        "Documentação")
            echo "Atualização de documentação ($doc_files arquivo(s))"
            ;;
        "Configuração")
            if [ "$sh_files" -gt 0 ]; then
                echo "Atualização de scripts e configuração ($total_files arquivo(s))"
            else
                echo "Atualização de configuração ($total_files arquivo(s))"
            fi
            ;;
        "Implementação")
            echo "Implementação de funcionalidades ($go_files arquivo(s))"
            ;;
        *)
            echo "Mudanças gerais no código ($total_files arquivo(s))"
            ;;
    esac
}

# ============================================
# Funções de Geração de Markdown
# ============================================

# Gera cabeçalho do documento CHANGELOG
generate_header() {
    cat << 'EOF'
# CHANGELOG

Este documento contém o histórico completo de mudanças do projeto.
Gerado automaticamente a partir do histórico Git.

EOF
}

# Gera tabela resumo de mudanças
# Lê commits do stdin no formato: hash|short_hash|datetime|message|author
generate_table() {
    echo "## Resumo de Mudanças"
    echo ""
    echo "| Data/Hora | Commit | Mensagem | Resumo das Mudanças |"
    echo "|-----------|--------|----------|---------------------|"
    
    # Para cada commit, adiciona linha na tabela
    while IFS='|' read -r hash short_hash datetime message author; do
        local files=$(get_commit_files "$hash")
        local change_type=$(classify_change "$files")
        local summary=$(generate_summary "$files" "$change_type")
        
        # Escapa pipes na mensagem para não quebrar tabela
        local safe_message=$(echo "$message" | sed 's/|/\\|/g')
        
        echo "| $datetime | \`$short_hash\` | $safe_message | $summary |"
    done
    
    echo ""
}

# Gera diagrama Mermaid gitGraph
# Lê commits do stdin no formato: hash|short_hash|datetime|message|author
generate_mermaid() {
    echo "## Timeline de Commits"
    echo ""
    echo '```mermaid'
    echo "gitGraph"
    
    # Para cada commit, adiciona ao gitGraph
    while IFS='|' read -r hash short_hash datetime message author; do
        # Extrai apenas a data (YYYY-MM-DD) do datetime
        local date=$(echo "$datetime" | cut -d' ' -f1)
        
        # Escapa aspas duplas na mensagem
        local safe_message=$(echo "$message" | sed 's/"/\\"/g')
        
        echo "    commit id: \"$short_hash\" tag: \"$date\""
    done
    
    echo '```'
    echo ""
}

# Gera seções detalhadas para cada commit
# Lê commits do stdin no formato: hash|short_hash|datetime|message|author
generate_details() {
    echo "## Detalhes dos Commits"
    echo ""
    
    while IFS='|' read -r hash short_hash datetime message author; do
        local files=$(get_commit_files "$hash")
        local change_type=$(classify_change "$files")
        local summary=$(generate_summary "$files" "$change_type")
        local file_count=$(echo "$files" | wc -l)
        
        echo "### [\`$short_hash\`] - $datetime"
        echo ""
        echo "**Mensagem:** $message"
        echo ""
        echo "**Autor:** $author"
        echo ""
        echo "**Tipo:** $change_type"
        echo ""
        echo "**Resumo:** $summary"
        echo ""
        echo "**Arquivos modificados:** $file_count"
        echo ""
        echo "---"
        echo ""
    done
}

# ============================================
# Fluxo Principal
# ============================================

log_msg "${GREEN}=== Gerador de CHANGELOG ===${NC}"
echo ""

# Validações iniciais
start_timer "validation"

# Verifica se está em repositório Git
if ! git rev-parse --git-dir > /dev/null 2>&1; then
    log_msg "${RED}✗ Erro: Não é um repositório Git${NC}"
    log_msg "${YELLOW}Execute este script dentro de um repositório Git${NC}"
    exit 1
fi

# Verifica se há commits
if ! git log -1 > /dev/null 2>&1; then
    log_msg "${YELLOW}⚠ Aviso: Repositório sem commits${NC}"
    log_msg "${YELLOW}Nada a fazer${NC}"
    exit 0
fi

# Verifica permissão de escrita
if [ ! -w . ]; then
    log_msg "${RED}✗ Erro: Sem permissão de escrita no diretório atual${NC}"
    exit 1
fi

duration=$(end_timer "validation")
log_msg "${GREEN}✓ Validações concluídas${NC}" "$duration"
echo ""

# Coleta de dados do Git
log_msg "${YELLOW}Analisando histórico Git...${NC}"
start_timer "git_analysis"

commits=$(get_all_commits)
commit_count=$(echo "$commits" | wc -l)

duration=$(end_timer "git_analysis")
log_msg "${GREEN}✓ Encontrados $commit_count commits${NC}" "$duration"
echo ""

# Geração do documento CHANGELOG.md
log_msg "${YELLOW}Gerando CHANGELOG.md...${NC}"
start_timer "generation"

# Avisa se vai sobrescrever arquivo existente
if [ -f CHANGELOG.md ]; then
    log_msg "${YELLOW}Sobrescrevendo CHANGELOG.md existente${NC}"
fi

# Gera documento na ordem: header -> mermaid -> table -> details
{
    generate_header
    echo "$commits" | generate_mermaid
    echo "$commits" | generate_table
    echo "$commits" | generate_details
} > CHANGELOG.md

duration=$(end_timer "generation")
log_msg "${GREEN}✓ CHANGELOG.md gerado com sucesso${NC}" "$duration"
echo ""

# Finalização
total_duration=$(end_timer "total")
log_msg "${GREEN}✓ Geração completa!${NC}" "$total_duration"
log_msg "Arquivo gerado: ${BOLD}CHANGELOG.md${NC}"
