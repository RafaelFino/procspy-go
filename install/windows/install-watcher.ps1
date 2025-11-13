# install-watcher.ps1
# Script de instalação do Procspy Watcher para Windows usando NSSM
# Requer privilégios de Administrador
#
# Uso:
#   .\install-watcher.ps1 [-BinaryPath <path>] [-ConfigPath <path>] [-Help]

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
$ServiceName = "procspy-watcher"
$InstallPath = "C:\Program Files\Procspy"
$ExePath = "$InstallPath\procspy-watcher.exe"
$ConfigFile = "$InstallPath\watcher-config.json"
$LogPath = "$InstallPath\logs"

# Caminhos padrão
$DefaultBinaryPath = ".\bin\procspy-watcher.exe"
$DefaultConfigPath = ".\etc\watcher-config.json"

# Função para exibir ajuda
function Show-Help {
    Write-Host ""
    Write-Host "=== Instalação do Procspy Watcher ===" -ForegroundColor Green
    Write-Host ""
    Write-Host "Uso: .\install-watcher.ps1 [parâmetros]" -ForegroundColor Cyan
    Write-Host ""
    Write-Host "Parâmetros:" -ForegroundColor Yellow
    Write-Host "  -BinaryPath <path>    Caminho para o binário procspy-watcher.exe"
    Write-Host "                        Padrão: $DefaultBinaryPath"
    Write-Host ""
    Write-Host "  -ConfigPath <path>    Caminho para o arquivo watcher-config.json"
    Write-Host "                        Padrão: $DefaultConfigPath"
    Write-Host ""
    Write-Host "  -Help                 Exibe esta ajuda"
    Write-Host ""
    Write-Host "Exemplos:" -ForegroundColor Yellow
    Write-Host "  .\install-watcher.ps1"
    Write-Host "  .\install-watcher.ps1 -BinaryPath C:\Downloads\procspy-watcher.exe"
    Write-Host ""
    exit 0
}

if ($Help) {
    Show-Help
}

# Define caminhos padrão
if ([string]::IsNullOrEmpty($BinaryPath)) {
    $BinaryPath = $DefaultBinaryPath
}

if ([string]::IsNullOrEmpty($ConfigPath)) {
    $ConfigPath = $DefaultConfigPath
}

Write-Host ""
Write-Host "=== Instalação do Procspy Watcher ===" -ForegroundColor Green
Write-Host ""
Write-Host "Configuração:" -ForegroundColor Cyan
Write-Host "  Binário: $BinaryPath"
Write-Host "  Config:  $ConfigPath"
Write-Host ""

# Verifica NSSM
if (-not (Get-Command nssm -ErrorAction SilentlyContinue)) {
    Write-Host "Erro: NSSM não encontrado no PATH" -ForegroundColor Red
    Write-Host "Instale o NSSM primeiro: https://nssm.cc/download" -ForegroundColor Yellow
    exit 1
}

Write-Host "✓ NSSM encontrado" -ForegroundColor Green

# Verifica binário
if (-not (Test-Path $BinaryPath)) {
    Write-Host "Erro: Binário não encontrado em: $BinaryPath" -ForegroundColor Red
    Write-Host "Execute o build ou especifique o caminho correto" -ForegroundColor Yellow
    exit 1
}

Write-Host "✓ Binário encontrado: $BinaryPath" -ForegroundColor Green

# Verifica configuração
$ConfigExists = Test-Path $ConfigPath
if ($ConfigExists) {
    Write-Host "✓ Configuração encontrada: $ConfigPath" -ForegroundColor Green
} else {
    Write-Host "⚠ Configuração não encontrada, será criada uma padrão" -ForegroundColor Yellow
}

# Para serviço existente
if (Get-Service $ServiceName -ErrorAction SilentlyContinue) {
    Write-Host "Parando serviço existente..." -ForegroundColor Yellow
    nssm stop $ServiceName
    Start-Sleep -Seconds 2
}

# Cria diretórios
Write-Host "Criando diretórios..." -ForegroundColor Cyan
New-Item -ItemType Directory -Force -Path $InstallPath | Out-Null
New-Item -ItemType Directory -Force -Path $LogPath | Out-Null

# Copia binário
Write-Host "Copiando binário..." -ForegroundColor Cyan
Copy-Item $BinaryPath -Destination $ExePath -Force
Write-Host "✓ Binário instalado" -ForegroundColor Green

# Copia ou cria configuração
if ($ConfigExists) {
    Copy-Item $ConfigPath -Destination $ConfigFile -Force
    
    # Ajusta start_cmd para Windows
    $config = Get-Content $ConfigFile | ConvertFrom-Json
    $config.start_cmd = "nssm restart procspy-client"
    $config | ConvertTo-Json -Depth 10 | Out-File -FilePath $ConfigFile -Encoding UTF8
    
    Write-Host "✓ Configuração copiada e ajustada" -ForegroundColor Green
} else {
    $defaultConfig = @{
        log_path = "logs"
        interval = 10
        procspy_url = "http://localhost:8888/healthcheck"
        start_cmd = "nssm restart procspy-client"
    } | ConvertTo-Json -Depth 10

    $defaultConfig | Out-File -FilePath $ConfigFile -Encoding UTF8
    Write-Host "✓ Configuração padrão criada" -ForegroundColor Green
}

# Remove serviço existente
if (Get-Service $ServiceName -ErrorAction SilentlyContinue) {
    nssm remove $ServiceName confirm
    Start-Sleep -Seconds 2
}

# Instala serviço
Write-Host "Instalando serviço..." -ForegroundColor Cyan
nssm install $ServiceName $ExePath $ConfigFile

# Configura serviço
nssm set $ServiceName AppDirectory $InstallPath
nssm set $ServiceName DisplayName "Procspy Watcher"
nssm set $ServiceName Description "Procspy Watcher - Monitora e reinicia o Client automaticamente"
nssm set $ServiceName Start SERVICE_AUTO_START
nssm set $ServiceName AppRestartDelay 5000
nssm set $ServiceName AppStdout "$LogPath\nssm-watcher-stdout.log"
nssm set $ServiceName AppStderr "$LogPath\nssm-watcher-stderr.log"
nssm set $ServiceName AppRotateFiles 1
nssm set $ServiceName AppRotateBytes 1048576
nssm set $ServiceName DependOnService procspy-client

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
}

Write-Host ""
Write-Host "=== Instalação Completa ===" -ForegroundColor Green
Write-Host ""
Write-Host "Comandos Úteis:" -ForegroundColor Yellow
Write-Host "  nssm status $ServiceName"
Write-Host "  nssm restart $ServiceName"
Write-Host "  notepad `"$ConfigFile`""
Write-Host ""
Write-Host "✓ Procspy Watcher instalado com sucesso!" -ForegroundColor Green
