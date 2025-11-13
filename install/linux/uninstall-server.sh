#!/bin/bash
# uninstall-server.sh - Desinstalação do Procspy Server

set -e
SERVICE_NAME="procspy-server"

if [ "$EUID" -ne 0 ]; then
    echo "Erro: Execute como root"
    exit 1
fi

echo "=== Desinstalando Procspy Server ==="

systemctl stop $SERVICE_NAME 2>/dev/null || true
systemctl disable $SERVICE_NAME 2>/dev/null || true
rm -f /etc/systemd/system/${SERVICE_NAME}.service
systemctl daemon-reload
rm -f /usr/local/bin/procspy-server

echo "✓ Procspy Server desinstalado"
echo "Banco de dados mantido em /var/lib/procspy/"
