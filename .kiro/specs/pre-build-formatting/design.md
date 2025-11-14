# Design Document - Pre-Build Formatting

## Overview

Esta feature adiciona duas melhorias aos scripts de build e teste do projeto procspy:

1. **Formatação automática**: Executa `go fmt ./...` automaticamente antes do build, garantindo que todo o código esteja formatado de acordo com os padrões Go antes da compilação
2. **Logging com timestamp**: Adiciona timestamps formatados em todas as mensagens de log dos scripts build.sh e test.sh para melhor rastreabilidade

A solução é integrada ao fluxo de build existente e oferece uma flag opcional para pular a formatação quando necessário.

## Architecture

### Current Build Flow

```
build.sh → Quality Checks (check_formatting + run_linter) → Tests → Build
```

### New Build Flow

```
build.sh → Auto Format → Quality Checks → Tests → Build
```

### Key Changes

1. **Nova função `auto_format()`**: Executa `go fmt ./...` e reporta resultados
2. **Nova flag `--no-fmt`**: Permite pular a formatação automática
3. **Integração no fluxo principal**: Chamada antes das verificações de qualidade
4. **Nova função `log_msg()`**: Adiciona timestamp formatado a todas as mensagens de log
5. **Atualização de todos os echo statements**: Substituir por chamadas a `log_msg()` em build.sh e test.sh

## Components and Interfaces

### 1. Log Message Function

**Função:** `log_msg()`

**Responsabilidades:**
- Adicionar timestamp formatado no início de cada mensagem
- Manter a cor original da mensagem
- Exibir timestamp em verde e negrito

**Assinatura:**
```bash
log_msg() {
    local message="$1"
    local timestamp=$(date '+%Y-%m-%d %H:%M:%S')
    echo -e "${GREEN}${BOLD}[${timestamp}]${NC} ${message}"
}
```

**Uso:**
```bash
# Antes
echo -e "${YELLOW}Executando testes...${NC}"

# Depois
log_msg "${YELLOW}Executando testes...${NC}"
```

**Notas de Implementação:**
- Adicionar variável `BOLD='\033[1m'` nas definições de cores
- A função preserva as cores passadas na mensagem
- O timestamp sempre aparece em verde e negrito, independente da cor da mensagem

### 2. Auto Format Function

**Função:** `auto_format()`

**Responsabilidades:**
- Executar o comando `go fmt ./...`
- Capturar e reportar erros
- Exibir feedback visual do processo
- Retornar código de saída apropriado

**Comportamento:**
```bash
auto_format() {
    log_msg "${YELLOW}Formatando código Go automaticamente...${NC}"
    
    if go fmt ./...; then
        log_msg "${GREEN}✓ Código formatado com sucesso${NC}"
        echo ""
        return 0
    else
        log_msg "${RED}✗ Erro ao formatar código${NC}"
        echo ""
        return 1
    fi
}
```

### 3. Flag --no-fmt

**Variável:** `NO_FMT` (boolean)

**Comportamento:**
- Default: `false` (formatação habilitada)
- Quando `true`: pula a etapa de formatação automática
- Parsing no bloco de argumentos existente

### 4. Integration Point

**Localização:** Após parse de argumentos, antes de `run_quality_checks()`

**Lógica:**
```bash
# Executa formatação automática se não for build-only e não tiver --no-fmt
if [ "$BUILD_ONLY" = false ] && [ "$NO_FMT" = false ]; then
    if ! auto_format; then
        echo -e "${RED}Build abortado devido a falha na formatação${NC}"
        exit 1
    fi
fi
```

## Data Models

Não há modelos de dados específicos para esta feature. A implementação trabalha apenas com:
- Variáveis de controle de fluxo (booleans)
- Códigos de saída de processos (integers)
- Strings para mensagens de output

## Error Handling

### Cenários de Erro

1. **go fmt falha**
   - Causa: Erro de sintaxe no código Go
   - Ação: Exibir mensagem de erro, abortar build, retornar exit code 1
   - Mensagem: "✗ Erro ao formatar código" + output do go fmt

2. **go fmt não encontrado**
   - Causa: Go não instalado ou não no PATH
   - Ação: Erro será capturado pelo if statement, abortar build
   - Mensagem: "✗ Erro ao formatar código" + mensagem do sistema

### Exit Codes

- `0`: Formatação bem-sucedida
- `1`: Erro na formatação (aborta build)

### Fallback Behavior

Não há fallback. Se a formatação falhar, o build é abortado para garantir qualidade do código. O usuário pode:
- Corrigir os erros de sintaxe
- Usar `--no-fmt` para pular a formatação temporariamente

## Testing Strategy

### Manual Testing

#### Testes de Formatação

1. **Teste de formatação bem-sucedida**
   - Executar `./build.sh` com código não formatado
   - Verificar que go fmt é executado
   - Verificar que build continua normalmente

2. **Teste de flag --no-fmt**
   - Executar `./build.sh --no-fmt`
   - Verificar que formatação é pulada
   - Verificar que build continua normalmente

3. **Teste de erro de formatação**
   - Criar arquivo Go com erro de sintaxe
   - Executar `./build.sh`
   - Verificar que build é abortado com mensagem clara

4. **Teste de compatibilidade**
   - Executar todas as flags existentes com formatação
   - Verificar que comportamento anterior é mantido
   - Testar: `--test-only`, `--build-only`, `--all`, `--clean`, `--platform`

5. **Teste de combinação de flags**
   - `./build.sh --no-fmt --build-only`
   - `./build.sh --no-fmt --all`
   - `./build.sh --clean --no-fmt`

#### Testes de Logging com Timestamp

1. **Teste de formato de timestamp em build.sh**
   - Executar `./build.sh`
   - Verificar que todas as mensagens têm formato `[yyyy-mm-dd HH:MM:SS] mensagem`
   - Verificar que timestamp está em verde e negrito
   - Verificar que cor da mensagem é preservada

2. **Teste de formato de timestamp em test.sh**
   - Executar `./test.sh`
   - Verificar que todas as mensagens têm formato `[yyyy-mm-dd HH:MM:SS] mensagem`
   - Verificar que timestamp está em verde e negrito
   - Verificar que cor da mensagem é preservada

3. **Teste de cores preservadas**
   - Verificar mensagens de sucesso (verde) mantêm cor após timestamp
   - Verificar mensagens de erro (vermelho) mantêm cor após timestamp
   - Verificar mensagens de aviso (amarelo) mantêm cor após timestamp
   - Verificar mensagens informativas (azul) mantêm cor após timestamp

4. **Teste de precisão de timestamp**
   - Executar script e verificar que timestamps são sequenciais
   - Verificar que formato de data/hora está correto

### Integration Testing

- Verificar que formatação ocorre antes de `check_formatting()`
- Verificar que `check_formatting()` passa após `auto_format()`
- Verificar que `--build-only` pula formatação (como esperado)

### Edge Cases

1. **Repositório sem arquivos Go**: go fmt deve completar sem erros
2. **Arquivos Go em vendor/**: go fmt ignora automaticamente
3. **Permissões de escrita**: go fmt precisa de permissão para modificar arquivos

## Implementation Notes

### Modificações no build.sh

1. **Adicionar variável `BOLD`** na seção de cores
2. **Adicionar função `log_msg()`** no início do arquivo, após definições de cores
3. **Substituir todos os `echo -e` por `log_msg`** em todo o script (exceto echo vazios)
4. **Adicionar função `auto_format()`** após `run_linter()`
5. **Adicionar variável `NO_FMT=false`** no início da seção de parse
6. **Adicionar case para --no-fmt** no bloco while de parse de argumentos
7. **Adicionar chamada a `auto_format()`** antes de `run_quality_checks()`
8. **Atualizar função `usage()`** para documentar a nova flag

### Modificações no test.sh

1. **Adicionar variável `BOLD`** na seção de cores
2. **Adicionar função `log_msg()`** no início do arquivo, após definições de cores
3. **Substituir todos os `echo -e` por `log_msg`** em todo o script (exceto echo vazios)

### Ordem de Execução

```
1. Parse de argumentos (incluindo --no-fmt)
2. Clean (se --clean)
3. Auto Format (se não --build-only e não --no-fmt) ← NOVO
4. Quality Checks (se não --build-only)
5. Tests (se não --build-only)
6. Build
```

### Backward Compatibility

- Todos os comandos existentes continuam funcionando
- Comportamento padrão adiciona formatação (melhoria, não breaking change)
- Flag --no-fmt permite restaurar comportamento anterior se necessário

## Design Decisions

### Por que go fmt em vez de gofmt?

`go fmt` é um wrapper para `gofmt` que opera em pacotes Go inteiros. É a ferramenta recomendada e mais conveniente para formatar projetos completos.

### Por que executar antes de check_formatting?

- `auto_format()` corrige problemas automaticamente
- `check_formatting()` apenas verifica e reporta
- Executar auto_format primeiro reduz falhas desnecessárias

### Por que abortar build em caso de erro?

Erros no `go fmt` geralmente indicam problemas de sintaxe graves que impediriam a compilação de qualquer forma. Abortar cedo fornece feedback mais claro.

### Por que --no-fmt em vez de --skip-fmt?

Consistência com outras flags do projeto que usam padrão `--no-*` (embora não haja exemplos no build.sh atual, é um padrão comum em ferramentas Unix).

### Por que não integrar ao makefile?

O makefile apenas chama build.sh. Toda a lógica deve estar no build.sh para manter a consistência e permitir uso direto do script.

### Por que usar função log_msg em vez de modificar echo diretamente?

- **Manutenibilidade**: Centraliza a lógica de timestamp em um único lugar
- **Consistência**: Garante formato uniforme em todas as mensagens
- **Flexibilidade**: Facilita mudanças futuras no formato de log
- **Legibilidade**: Código mais limpo e fácil de entender

### Por que timestamp em verde e negrito?

- **Verde**: Cor neutra que não interfere com as cores das mensagens (sucesso, erro, aviso)
- **Negrito**: Destaca o timestamp para fácil identificação visual
- **Consistência**: Mantém padrão visual uniforme em todo o output

### Por que aplicar em build.sh e test.sh?

Ambos os scripts são executados frequentemente durante o desenvolvimento e produzem output significativo. Ter timestamps consistentes em ambos facilita:
- Debugging de problemas
- Análise de performance
- Correlação de eventos entre build e testes
- Logs de CI/CD mais informativos
