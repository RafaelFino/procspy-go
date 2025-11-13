# install-client.ps1
# Script de instalação do Procspy Client para Windows usando NSSM
# Requer privilégios de Administrador
#
# Uso:
#   .\install-client.ps1 [-BinaryPath <path>] [-ConfigPath <path>] [-Help]
#
# Parâmetros:
#   -BinaryPath    Caminho para o binário procspy-client.exe
#   -ConfigPath    Caminho para o arquivo config-client.json
#   -Help          Exibe esta ajuda

#Requires -RunAsAdministrator

[CmdletBinding()]
param(
    [Parameter(Mandatory=$false)]
    [string]$BinaryPath = "",
    
    [Parameter(Mandatory=$false)]
    [string]$ConfigPath = "",
    
    [Parameter(Mandatory=$false)]
    [switch]$Help
)

$ErrorActionPreference = "Stop"

# Configurações
$ServiceName = "procspy-client"
$InstallPath = "C:\Program Files\Procspy"
$ExePath = "$InstallPath\procspy-client.exe"
$ConfigFile = "$InstallPath\config-client.json"
$LogPath = "$InstallPath\logs"

# Caminhos padrão
$DefaultBinaryPath = ".\bin\procspy-client.exe"
$DefaultConfigPath = ".\etc\config-client.json"

# Função para exibir ajuda
function Show-Help {
    Write-Host ""
    Write-Host "=== Instalação do Procspy Client ===" -ForegroundColor Green
    Write-Host ""
    Write-Host "Uso: .\install-client.ps1 [parâmetros]" -ForegroundColor Cyan
    Write-Host ""
    Write-Host "Parâmetros:" -ForegroundColor Yellow
    Write-Host "  -BinaryPath <path>    Caminho para o binário procspy-client.exe"
    Write-Host "                        Padrão: $DefaultBinaryPath"
    Write-Host ""
    Write-Host "  -ConfigPath <path>    Caminho para o arquivo config-client.json"
    Write-Host "                        Padrão: $DefaultConfigPath"
    Write-Host ""
    Write-Host "  -Help                 Exibe esta ajuda"
    Write-Host ""
    Write-Host "Exemplos:" -ForegroundColor Yellow
    Write-Host "  # Instalação padrão"
    Write-Host "  .\install-client.ps1"
    Write-Host ""
    Write-Host "  # Instalação com binário customizado"
    Write-Host "  .\install-client.ps1 -BinaryPath C:\Downloads\procspy-client.exe"
    Write-Host ""
    Write-Host "  # Instalação com binário e config customizados"
    Write-Host "  .\install-client.ps1 -BinaryPath D:\builds\procspy-client.exe -ConfigPath D:\configs\config-client.json"
    Write-Host ""
    Write-Host "Notas:" -ForegroundColor Yellow
    Write-Host "  - Este script deve ser executado como Administrador"
    Write-Host "  - NSSM deve estar instalado e no PATH"
    Write-Host "  - O binário será copiado para: $InstallPath"
    Write-Host ""
    exit 0
}

# Exibe ajuda se solicitado
if ($Help) {
    Show-Help
}

# Define caminhos padrão se não foram especificados
if ([string]::IsNullOrEmpty($BinaryPath)) {
    $BinaryPath = $DefaultBinaryPath
}

if ([string]::IsNullOrEmpty($ConfigPath)) {
    $ConfigPath = $DefaultConfigPath
}

Write-Host ""
Write-Host "=== Instalação do Procspy Client ===" -ForegroundColor Green
Write-Host ""
Write-Host "Configuração:" -ForegroundColor Cyan
Write-Host "  Binário: $BinaryPath"
Write-Host "  Config:  $ConfigPath"
Write-Host ""

# Verifica se NSSM está disponível
if (-not (Get-Command nssm -ErrorAction SilentlyContinue)) {
    Write-Host "Erro: NSSM não encontrado no PATH" -ForegroundColor Red
    Write-Host ""
    Write-Host "Por favor, instale o NSSM primeiro:" -ForegroundColor Yellow
    Write-Host "  1. Baixe de: https://nssm.cc/download" -ForegroundColor Yellow
    Write-Host "  2. Extraia para C:\nssm\" -ForegroundColor Yellow
    Write-Host "  3. Adicione C:\nssm\win64 ao PATH" -ForegroundColor Yellow
    Write-Host ""
    Write-Host "Ou execute: .\get-nssm.sh (no Git Bash/WSL)" -ForegroundColor Yellow
    exit 1
}

Write-Host "✓ NSSM encontrado" -ForegroundColor Green

# Verifica se binário existe
if (-not (Test-Path $BinaryPath)) {
    Write-Host "Erro: Binário não encontrado em: $BinaryPath" -ForegroundColor Red
    Write-Host ""
    Write-Host "Soluções:" -ForegroundColor Yellow
    Write-Host "  1. Execute o build primeiro: .\build.sh" -ForegroundColor Yellow
    Write-Host "  2. Especifique o caminho correto: -BinaryPath <caminho>" -ForegroundColor Yellow
    Write-Host "  3. Use -Help para ver todas as opções" -ForegroundColor Yellow
    exit 1
}

Write-Host "✓ Binário encontrado: $BinaryPath" -ForegroundColor Green

# Verifica se configuração existe
$ConfigExists = Test-Path $ConfigPath
if ($ConfigExists) {
    Write-Host "✓ Arquivo de configuração encontrado: $ConfigPath" -ForegroundColor Green
} else {
    Write-Host "⚠ Arquivo de configuração não encontrado: $ConfigPath" -ForegroundColor Yellow
    Write-Host "  Uma configuração padrão será criada" -ForegroundColor Yellow
}

# Para serviço existente
if (Get-Service $ServiceName -ErrorAction SilentlyContinue) {
    Write-Host "Parando serviço existente..." -ForegroundColor Yellow
    nssm stop $ServiceName
    Start-Sleep -Seconds 2
}

# Cria diretório de instalação
Write-Host "Criando diretório de instalação..." -ForegroundColor Cyan
New-Item -ItemType Directory -Force -Path $InstallPath | Out-Null
New-Item -ItemType Directory -Force -Path $LogPath | Out-Null

# Copia binário
Write-Host "Copiando binário..." -ForegroundColor Cyan
Copy-Item $BinaryPath -Destination $ExePath -Force
Write-Host "✓ Binário instalado em $ExePath" -ForegroundColor Green

# Copia ou cria configuração
if ($ConfigExists) {
    Write-Host "Copiando configuração..." -ForegroundColor Cyan
    Copy-Item $ConfigPath -Destination $ConfigFile -Force
    Write-Host "✓ Configuração copiada para $ConfigFile" -ForegroundColor Green
} else {
    Write-Host "Criando configuração padrão..." -ForegroundColor Cyan
    $defaultConfig = @{
        user = $env:USERNAME
        log_path = "logs"
        debug = $false
        interval = 5
        server_url = "https://seu-servidor.com/procspy"
        api_host = "localhost"
        api_port = 8888
    } | ConvertTo-Json -Depth 10

    $defaultConfig | Out-File -FilePath $ConfigFile -Encoding UTF8
    Write-Host "✓ Configuração padrão criada em $ConfigFile" -ForegroundColor Green
    Write-Host "  IMPORTANTE: Edite a configuração antes de iniciar" -ForegroundColor Yellow
    Write-Host "  Execute: notepad `"$ConfigFile`"" -ForegroundColor Yellow
}

# Remove serviço existente se houver
if (Get-Service $ServiceName -ErrorAction SilentlyContinue) {
    Write-Host "Removendo serviço existente..." -ForegroundColor Yellow
    nssm remove $ServiceName confirm
    Start-Sleep -Seconds 2
}

# Instala serviço com NSSM
Write-Host "Instalando serviço..." -ForegroundColor Cyan
nssm install $ServiceName $ExePath $ConfigFile

# Configura serviço
Write-Host "Configurando serviço..." -ForegroundColor Cyan
nssm set $ServiceName AppDirectory $InstallPath
nssm set $ServiceName DisplayName "Procspy Client"
nssm set $ServiceName Description "Procspy - Monitoramento de processos para controle parental"
nssm set $ServiceName Start SERVICE_AUTO_START
nssm set $ServiceName AppRestartDelay 5000
nssm set $ServiceName AppStdout "$LogPath\nssm-stdout.log"
nssm set $ServiceName AppStderr "$LogPath\nssm-stderr.log"
nssm set $ServiceName AppRotateFiles 1
nssm set $ServiceName AppRotateBytes 1048576

Write-Host "✓ Serviço configurado" -ForegroundColor Green

# Inicia serviço
Write-Host "Iniciando serviço..." -ForegroundColor Cyan
nssm start $ServiceName
Start-Sleep -Seconds 3

# Verifica status
$status = nssm status $ServiceName
if ($status -eq "SERVICE_RUNNING") {
    Write-Host "✓ Serviço iniciado com sucesso" -ForegroundColor Green
} else {
    Write-Host "⚠ Serviço instalado mas não está rodando" -ForegroundColor Yellow
    Write-Host "  Status: $status" -ForegroundColor Yellow
    Write-Host "  Verifique os logs em: $LogPath" -ForegroundColor Yellow
}

Write-Host ""
Write-Host "=== Instalação Completa ===" -ForegroundColor Green
Write-Host ""
Write-Host "Informações do Serviço:" -ForegroundColor Cyan
Write-Host "  Nome:         $ServiceName"
Write-Host "  Instalado em: $InstallPath"
Write-Host "  Configuração: $ConfigFile"
Write-Host "  Logs:         $LogPath"
Write-Host ""
Write-Host "Comandos Úteis:" -ForegroundColor Yellow
Write-Host "  nssm status $ServiceName      # Verificar status"
Write-Host "  nssm start $ServiceName       # Iniciar serviço"
Write-Host "  nssm stop $ServiceName        # Parar serviço"
Write-Host "  nssm restart $ServiceName     # Reiniciar serviço"
Write-Host "  nssm edit $ServiceName        # Editar configuração"
Write-Host ""
Write-Host "  notepad `"$ConfigFile`"       # Editar config"
Write-Host "  nssm restart $ServiceName     # Após editar"
Write-Host ""
Write-Host "✓ Procspy Client instalado com sucesso!" -ForegroundColor Green
