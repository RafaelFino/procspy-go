#!/bin/bash
# uninstall-watcher.sh - Desinstalação do Procspy Watcher

set -e
SERVICE_NAME="procspy-watcher"

if [ "$EUID" -ne 0 ]; then
    echo "Erro: Execute como root"
    exit 1
fi

echo "=== Desinstalando Procspy Watcher ==="

systemctl stop $SERVICE_NAME 2>/dev/null || true
systemctl disable $SERVICE_NAME 2>/dev/null || true
rm -f /etc/systemd/system/${SERVICE_NAME}.service
systemctl daemon-reload
rm -f /usr/local/bin/procspy-watcher

echo "✓ Procspy Watcher desinstalado"
