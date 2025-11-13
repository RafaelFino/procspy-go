# uninstall-client.ps1
# Script de desinstalação do Procspy Client para Windows

#Requires -RunAsAdministrator

$ErrorActionPreference = "Stop"

$ServiceName = "procspy-client"
$InstallPath = "C:\Program Files\Procspy"

Write-Host "=== Desinstalação do Procspy Client ===" -ForegroundColor Yellow
Write-Host ""

# Verifica se serviço existe
if (-not (Get-Service $ServiceName -ErrorAction SilentlyContinue)) {
    Write-Host "Serviço $ServiceName não encontrado" -ForegroundColor Yellow
} else {
    # Para serviço
    Write-Host "Parando serviço..." -ForegroundColor Cyan
    nssm stop $ServiceName
    Start-Sleep -Seconds 2
    
    # Remove serviço
    Write-Host "Removendo serviço..." -ForegroundColor Cyan
    nssm remove $ServiceName confirm
    Write-Host "✓ Serviço removido" -ForegroundColor Green
}

# Remove arquivos (opcional)
$removeFiles = Read-Host "Deseja remover os arquivos de $InstallPath? (S/N)"
if ($removeFiles -eq "S" -or $removeFiles -eq "s") {
    if (Test-Path $InstallPath) {
        Write-Host "Removendo arquivos..." -ForegroundColor Cyan
        Remove-Item -Path "$InstallPath\procspy-client.exe" -Force -ErrorAction SilentlyContinue
        Remove-Item -Path "$InstallPath\config-client.json" -Force -ErrorAction SilentlyContinue
        Write-Host "✓ Arquivos removidos" -ForegroundColor Green
    }
}

Write-Host ""
Write-Host "✓ Procspy Client desinstalado" -ForegroundColor Green
