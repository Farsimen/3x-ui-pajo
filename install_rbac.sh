#!/bin/bash

# 3X-UI RBAC Installation Script
# Multi-Vendor Support with Role-Based Access Control

set -e

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

echo -e "${GREEN}================================${NC}"
echo -e "${GREEN}3X-UI RBAC Installer${NC}"
echo -e "${GREEN}Multi-Vendor Panel with RBAC${NC}"
echo -e "${GREEN}================================${NC}"
echo ""

# Check if running as root
if [[ $EUID -ne 0 ]]; then
   echo -e "${RED}Error: This script must be run as root${NC}" 
   exit 1
fi

# Check OS
if [[ ! -f /etc/os-release ]]; then
    echo -e "${RED}Error: Cannot detect OS${NC}"
    exit 1
fi

source /etc/os-release
if [[ "$ID" != "ubuntu" && "$ID" != "debian" ]]; then
    echo -e "${YELLOW}Warning: This script is tested on Ubuntu/Debian only${NC}"
fi

echo -e "${GREEN}[1/6] Installing dependencies...${NC}"
apt-get update -qq
apt-get install -y wget curl git sqlite3 gcc build-essential > /dev/null 2>&1

echo -e "${GREEN}[2/6] Installing Go...${NC}"
if ! command -v go &> /dev/null; then
    cd /tmp
    wget -q https://go.dev/dl/go1.23.4.linux-amd64.tar.gz
    rm -rf /usr/local/go
    tar -C /usr/local -xzf go1.23.4.linux-amd64.tar.gz
    echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
    export PATH=$PATH:/usr/local/go/bin
    rm go1.23.4.linux-amd64.tar.gz
fi

echo -e "${GREEN}[3/6] Cloning 3X-UI RBAC...${NC}"
cd /tmp
rm -rf 3x-ui-rbac
git clone -q -b rbac-implementation https://github.com/Farsimen/3x-ui-pajo.git 3x-ui-rbac
cd 3x-ui-rbac

echo -e "${GREEN}[4/6] Building application (this may take a few minutes)...${NC}"
CGO_ENABLED=1 go build -v -o x-ui main.go > /dev/null 2>&1

if [[ ! -f x-ui ]]; then
    echo -e "${RED}Error: Build failed${NC}"
    exit 1
fi

echo -e "${GREEN}[5/6] Installing...${NC}"

# Stop existing service if any
systemctl stop x-ui 2>/dev/null || true

# Create installation directory
mkdir -p /usr/local/x-ui

# Copy files
cp x-ui /usr/local/x-ui/
chmod +x /usr/local/x-ui/x-ui

# Copy xray binary if exists
if [[ -d bin ]]; then
    cp -r bin /usr/local/x-ui/
fi

# Setup systemd service
cat > /etc/systemd/system/x-ui.service << 'EOF'
[Unit]
Description=3X-UI RBAC - Multi-Vendor Xray Panel
After=network.target

[Service]
Type=simple
User=root
WorkingDirectory=/usr/local/x-ui
ExecStart=/usr/local/x-ui/x-ui
Restart=on-failure
RestartSec=5s
LimitNOFILE=1048576

[Install]
WantedBy=multi-user.target
EOF

systemctl daemon-reload
systemctl enable x-ui

echo -e "${GREEN}[6/6] Starting service...${NC}"
systemctl start x-ui
sleep 3

if systemctl is-active --quiet x-ui; then
    echo ""
    echo -e "${GREEN}================================${NC}"
    echo -e "${GREEN}Installation Successful!${NC}"
    echo -e "${GREEN}================================${NC}"
    echo ""
    echo -e "Default Login:"
    echo -e "  URL: ${YELLOW}http://YOUR_IP:2087${NC}"
    echo -e "  Username: ${YELLOW}admin${NC}"
    echo -e "  Password: ${YELLOW}admin${NC}"
    echo ""
    echo -e "${YELLOW}IMPORTANT: Change default password immediately!${NC}"
    echo ""
    echo -e "RBAC API Endpoints:"
    echo -e "  POST /vendor/create  - Create vendor"
    echo -e "  GET  /vendor/list    - List vendors"
    echo -e "  DEL  /vendor/delete/:id - Delete vendor"
    echo -e "  POST /vendor/grant   - Grant inbound access"
    echo -e "  POST /vendor/revoke  - Revoke inbound access"
    echo ""
    echo -e "Logs: ${YELLOW}journalctl -u x-ui -f${NC}"
    echo -e "Status: ${YELLOW}systemctl status x-ui${NC}"
    echo ""
else
    echo -e "${RED}Error: Service failed to start${NC}"
    echo -e "Check logs: journalctl -u x-ui -n 50"
    exit 1
fi
