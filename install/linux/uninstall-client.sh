#!/bin/bash
# uninstall-client.sh - Desinstalação do Procspy Client

set -e
SERVICE_NAME="procspy-client"

if [ "$EUID" -ne 0 ]; then
    echo "Erro: Execute como root"
    exit 1
fi

echo "=== Desinstalando Procspy Client ==="

systemctl stop $SERVICE_NAME 2>/dev/null || true
systemctl disable $SERVICE_NAME 2>/dev/null || true
rm -f /etc/systemd/system/${SERVICE_NAME}.service
systemctl daemon-reload
rm -f /usr/local/bin/procspy-client

echo "✓ Procspy Client desinstalado"
echo "Configurações mantidas em /etc/procspy/"
echo "Logs mantidos em /var/log/procspy/"
