# Scripts de Instalação do Procspy

Este diretório contém scripts automatizados para instalação e desinstalação do Procspy em diferentes sistemas operacionais.

## Estrutura

```
install/
├── windows/          # Scripts PowerShell para Windows
│   ├── install-client.ps1
│   ├── install-watcher.ps1
│   ├── uninstall-client.ps1
│   └── uninstall-watcher.ps1
├── linux/            # Scripts Bash para Linux
│   ├── install-client.sh
│   ├── install-watcher.sh
│   ├── install-server.sh
│   ├── uninstall-client.sh
│   ├── uninstall-watcher.sh
│   └── uninstall-server.sh
└── README.md         # Este arquivo
```

## Instalação no Windows

### Pré-requisitos
- Windows 7 ou superior
- PowerShell 5.1 ou superior
- NSSM (Non-Sucking Service Manager)
- Privilégios de Administrador

### Instalar NSSM

```powershell
# Baixe de https://nssm.cc/download
# Ou use o script (no Git Bash/WSL):
./get-nssm.sh
```

### Instalar Client

```powershell
# Execute PowerShell como Administrador
cd procspy
.\install\windows\install-client.ps1
```

### Instalar Watcher

```powershell
.\install\windows\install-watcher.ps1
```

### Desinstalar

```powershell
.\install\windows\uninstall-client.ps1
.\install\windows\uninstall-watcher.ps1
```

## Instalação no Linux

### Pré-requisitos
- Linux com systemd
- Privilégios de root (sudo)
- Binários compilados em `./bin/`

### Instalar Client

```bash
chmod +x install/linux/install-client.sh
sudo ./install/linux/install-client.sh
```

### Instalar Watcher

```bash
chmod +x install/linux/install-watcher.sh
sudo ./install/linux/install-watcher.sh
```

### Instalar Server

```bash
chmod +x install/linux/install-server.sh
sudo ./install/linux/install-server.sh
```

### Desinstalar

```bash
sudo ./install/linux/uninstall-client.sh
sudo ./install/linux/uninstall-watcher.sh
sudo ./install/linux/uninstall-server.sh
```

## Configuração Pós-Instalação

### Windows

1. Edite as configurações:
   ```
   C:\Program Files\Procspy\config-client.json
   C:\Program Files\Procspy\watcher-config.json
   ```

2. Reinicie os serviços:
   ```powershell
   nssm restart procspy-client
   nssm restart procspy-watcher
   ```

### Linux

1. Edite as configurações:
   ```bash
   sudo micro /etc/procspy/config-client.json
   sudo micro /etc/procspy/watcher-config.json
   sudo micro /etc/procspy/config-server.json
   ```

2. Reinicie os serviços:
   ```bash
   sudo systemctl restart procspy-client
   sudo systemctl restart procspy-watcher
   sudo systemctl restart procspy-server
   ```

## Verificação

### Windows

```powershell
nssm status procspy-client
nssm status procspy-watcher
```

### Linux

```bash
systemctl status procspy-client
systemctl status procspy-watcher
systemctl status procspy-server
```

## Logs

### Windows
```
C:\Program Files\Procspy\logs\
```

### Linux
```bash
# Via journald
journalctl -u procspy-client -f

# Arquivos
/var/log/procspy/
```

## Troubleshooting

### Serviço não inicia

**Windows:**
```powershell
nssm status procspy-client
Get-Content "C:\Program Files\Procspy\logs\nssm-stderr.log"
```

**Linux:**
```bash
systemctl status procspy-client
journalctl -u procspy-client -n 50
```

### Permissões

Os serviços precisam rodar com privilégios elevados:
- Windows: SYSTEM ou Administrator
- Linux: root

### Firewall

**Linux Server:**
```bash
sudo ufw allow 8080/tcp
```

## Suporte

Para mais informações, consulte o README.md principal do projeto.
