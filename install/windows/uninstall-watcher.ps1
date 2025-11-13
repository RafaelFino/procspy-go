# uninstall-watcher.ps1
# Script de desinstalação do Procspy Watcher para Windows

#Requires -RunAsAdministrator

$ErrorActionPreference = "Stop"

$ServiceName = "procspy-watcher"
$InstallPath = "C:\Program Files\Procspy"

Write-Host "=== Desinstalação do Procspy Watcher ===" -ForegroundColor Yellow
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
        Remove-Item -Path "$InstallPath\procspy-watcher.exe" -Force -ErrorAction SilentlyContinue
        Remove-Item -Path "$InstallPath\watcher-config.json" -Force -ErrorAction SilentlyContinue
        Write-Host "✓ Arquivos removidos" -ForegroundColor Green
    }
}

Write-Host ""
Write-Host "✓ Procspy Watcher desinstalado" -ForegroundColor Green
