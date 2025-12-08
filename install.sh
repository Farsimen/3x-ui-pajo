#!/bin/bash

red='\033[0;31m'
green='\033[0;32m'
blue='\033[0;34m'
yellow='\033[0;33m'
purple='\033[0;35m'
cyan='\033[0;36m'
plain='\033[0m'

cur_dir=$(pwd)

# Banner
print_banner() {
    clear
    echo -e "${purple}"
    cat << "EOF"
 ╔═══════════════════════════════════════════════════╗
 ║                                                   ║
 ║       3X-UI Pajo - Multi-Vendor Edition          ║
 ║     Role-Based Access Control (RBAC) System      ║
 ║                                                   ║
 ╚═══════════════════════════════════════════════════╝
EOF
    echo -e "${plain}"
}

# Check root
[[ $EUID -ne 0 ]] && echo -e "${red}Fatal error: ${plain} Please run this script with root privilege \n " && exit 1

# Check OS and set release variable
if [[ -f /etc/os-release ]]; then
    source /etc/os-release
    release=$ID
elif [[ -f /usr/lib/os-release ]]; then
    source /usr/lib/os-release
    release=$ID
else
    echo "Failed to check the system OS, please contact the author!" >&2
    exit 1
fi

arch() {
    case "$(uname -m)" in
        x86_64 | x64 | amd64) echo 'amd64' ;;
        i*86 | x86) echo '386' ;;
        armv8* | armv8 | arm64 | aarch64) echo 'arm64' ;;
        armv7* | armv7 | arm) echo 'armv7' ;;
        armv6* | armv6) echo 'armv6' ;;
        armv5* | armv5) echo 'armv5' ;;
        s390x) echo 's390x' ;;
        *) echo -e "${red}Unsupported CPU architecture! ${plain}" && exit 1 ;;
    esac
}

install_base() {
    echo -e "${cyan}📦 Installing base packages...${plain}"
    case "${release}" in
        ubuntu | debian | armbian)
            apt-get update && apt-get install -y -q wget curl tar tzdata git golang-go
        ;;
        fedora | amzn | virtuozzo | rhel | almalinux | rocky | ol)
            dnf -y update && dnf install -y -q wget curl tar tzdata git golang
        ;;
        centos)
            if [[ "${VERSION_ID}" =~ ^7 ]]; then
                yum -y update && yum install -y wget curl tar tzdata git golang
            else
                dnf -y update && dnf install -y -q wget curl tar tzdata git golang
            fi
        ;;
        arch | manjaro | parch)
            pacman -Syu && pacman -Syu --noconfirm wget curl tar tzdata git go
        ;;
        opensuse-tumbleweed | opensuse-leap)
            zypper refresh && zypper -q install -y wget curl tar timezone git go
        ;;
        alpine)
            apk update && apk add wget curl tar tzdata git go
        ;;
        *)
            apt-get update && apt-get install -y -q wget curl tar tzdata git golang-go
        ;;
    esac
    echo -e "${green}✓ Base packages installed${plain}"
}

get_server_ip() {
    local URL_lists=(
        "https://api4.ipify.org"
        "https://ipv4.icanhazip.com"
        "https://v4.api.ipinfo.io/ip"
    )
    local server_ip=""
    for ip_address in "${URL_lists[@]}"; do
        server_ip=$(curl -s --max-time 3 "${ip_address}" 2>/dev/null | tr -d '[:space:]')
        if [[ -n "${server_ip}" ]]; then
            echo "${server_ip}"
            return 0
        fi
    done
    echo "YOUR_SERVER_IP"
}

interactive_config() {
    print_banner
    echo -e "${cyan}╔═══════════════════════════════════════════════════╗${plain}"
    echo -e "${cyan}║           Interactive Setup Wizard                ║${plain}"
    echo -e "${cyan}╚═══════════════════════════════════════════════════╝${plain}"
    echo ""
    
    # Panel Port
    echo -e "${yellow}┌─ Panel Port${plain}"
    while true; do
        read -rp "$(echo -e "${yellow}└─> Enter port [1-65535] (Default: 2053): ${plain}")" config_port
        config_port=${config_port:-2053}
        if [[ "$config_port" =~ ^[0-9]+$ ]] && [ "$config_port" -ge 1 ] && [ "$config_port" -le 65535 ]; then
            echo -e "${green}    ✓ Port set to: ${config_port}${plain}"
            break
        else
            echo -e "${red}    ✗ Invalid port. Please enter 1-65535${plain}"
        fi
    done
    echo ""
    
    # Admin Username
    echo -e "${yellow}┌─ Admin Username${plain}"
    while true; do
        read -rp "$(echo -e "${yellow}└─> Enter username (min 3 chars): ${plain}")" config_username
        if [[ -z "$config_username" ]]; then
            echo -e "${red}    ✗ Username cannot be empty${plain}"
            continue
        fi
        if [[ ${#config_username} -ge 3 ]]; then
            echo -e "${green}    ✓ Username: ${config_username}${plain}"
            break
        else
            echo -e "${red}    ✗ Username too short (min 3 characters)${plain}"
        fi
    done
    echo ""
    
    # Admin Password
    echo -e "${yellow}┌─ Admin Password${plain}"
    while true; do
        read -rsp "$(echo -e "${yellow}└─> Enter password (min 4 chars): ${plain}")" config_password
        echo ""
        if [[ -z "$config_password" ]]; then
            echo -e "${red}    ✗ Password cannot be empty${plain}"
            continue
        fi
        if [[ ${#config_password} -ge 4 ]]; then
            echo -e "${green}    ✓ Password set (hidden)${plain}"
            break
        else
            echo -e "${red}    ✗ Password too short (min 4 characters)${plain}"
        fi
    done
    echo ""
    
    # Installation Path
    echo -e "${yellow}┌─ Installation Path${plain}"
    read -rp "$(echo -e "${yellow}└─> Enter path (Default: /usr/local/x-ui): ${plain}")" install_path
    install_path=${install_path:-/usr/local/x-ui}
    echo -e "${green}    ✓ Path: ${install_path}${plain}"
    echo ""
    
    # Confirmation
    echo -e "${cyan}╔═══════════════════════════════════════════════════╗${plain}"
    echo -e "${cyan}║           Configuration Summary                   ║${plain}"
    echo -e "${cyan}╚═══════════════════════════════════════════════════╝${plain}"
    echo -e "${blue}  Port:${plain}      ${green}${config_port}${plain}"
    echo -e "${blue}  Username:${plain}  ${green}${config_username}${plain}"
    echo -e "${blue}  Password:${plain}  ${green}$(echo ${config_password} | sed 's/./*/g')${plain}"
    echo -e "${blue}  Path:${plain}      ${green}${install_path}${plain}"
    echo ""
    
    read -rp "$(echo -e "${yellow}Continue with installation? [Y/n]: ${plain}")" confirm
    confirm=${confirm:-Y}
    if [[ "$confirm" != "y" && "$confirm" != "Y" ]]; then
        echo -e "${red}✗ Installation cancelled${plain}"
        exit 0
    fi
    echo ""
}

install_x-ui_pajo() {
    echo -e "${cyan}╔═══════════════════════════════════════════════════╗${plain}"
    echo -e "${cyan}║         Installing 3X-UI Pajo RBAC Edition        ║${plain}"
    echo -e "${cyan}╚═══════════════════════════════════════════════════╝${plain}"
    echo ""
    
    # Stop existing service
    if systemctl is-active --quiet x-ui 2>/dev/null; then
        echo -e "${yellow}⚠  Stopping existing x-ui service...${plain}"
        systemctl stop x-ui
    fi
    
    # Create installation directory
    echo -e "${cyan}📁 Creating installation directory...${plain}"
    mkdir -p "${install_path}"
    cd "${install_path}" || exit 1
    
    # Backup existing installation
    if [[ -f "x-ui" ]]; then
        echo -e "${yellow}📦 Backing up existing installation...${plain}"
        cp x-ui x-ui.backup.$(date +%Y%m%d_%H%M%S) 2>/dev/null || true
    fi
    
    # Clone repository
    echo -e "${cyan}⬇️  Downloading from GitHub...${plain}"
    if [[ -d ".git" ]]; then
        rm -rf .git
    fi
    
    rm -rf temp_clone 2>/dev/null || true
    git clone --depth 1 -b rbac-implementation https://github.com/Farsimen/3x-ui-pajo.git temp_clone
    if [[ $? -ne 0 ]]; then
        echo -e "${red}✗ Failed to clone repository${plain}"
        echo -e "${yellow}  Please check your internet connection and try again${plain}"
        exit 1
    fi
    
    # Move files
    echo -e "${cyan}📋 Extracting files...${plain}"
    cp -r temp_clone/* . 2>/dev/null || true
    cp -r temp_clone/.github . 2>/dev/null || true
    rm -rf temp_clone
    
    # Build
    echo -e "${cyan}🔨 Building application (this may take a few minutes)...${plain}"
    export GO111MODULE=on
    go mod download
    go build -ldflags="-s -w" -o x-ui main.go
    if [[ $? -ne 0 ]]; then
        echo -e "${red}✗ Build failed${plain}"
        echo -e "${yellow}  Please check Go installation: go version${plain}"
        exit 1
    fi
    
    chmod +x x-ui
    chmod +x x-ui.sh 2>/dev/null || true
    
    # Set permissions for bin directory
    if [[ -d "bin" ]]; then
        chmod +x bin/* 2>/dev/null || true
    fi
    
    echo -e "${green}✓ Build completed successfully${plain}"
    
    # Configure panel
    echo -e "${cyan}⚙️  Configuring panel settings...${plain}"
    ./x-ui setting -username "${config_username}" -password "${config_password}" -port "${config_port}" 2>/dev/null || true
    
    # Install CLI tool
    echo -e "${cyan}🔧 Installing CLI tool...${plain}"
    cp -f x-ui.sh /usr/bin/x-ui 2>/dev/null || cat > /usr/bin/x-ui << 'EOFCLI'
#!/bin/bash
/usr/local/x-ui/x-ui "$@"
EOFCLI
    chmod +x /usr/bin/x-ui
    
    # Create systemd service
    echo -e "${cyan}🔧 Creating systemd service...${plain}"
    cat > /etc/systemd/system/x-ui.service <<EOF
[Unit]
Description=3X-UI Pajo - Multi-Vendor Xray Panel with RBAC
Documentation=https://github.com/Farsimen/3x-ui-pajo
After=network.target nss-lookup.target

[Service]
Type=simple
User=root
WorkingDirectory=${install_path}
ExecStart=${install_path}/x-ui
Restart=on-failure
RestartSec=5s
LimitNOFILE=1048576

[Install]
WantedBy=multi-user.target
EOF
    
    # Reload and start service
    echo -e "${cyan}🚀 Starting service...${plain}"
    systemctl daemon-reload
    systemctl enable x-ui >/dev/null 2>&1
    systemctl restart x-ui
    
    # Wait for service to start
    sleep 3
    
    if systemctl is-active --quiet x-ui; then
        echo -e "${green}✓ Service started successfully!${plain}"
    else
        echo -e "${yellow}⚠  Service may not have started. Check: journalctl -u x-ui -n 50${plain}"
    fi
    
    echo ""
}

show_access_info() {
    local server_ip=$(get_server_ip)
    
    clear
    print_banner
    
    echo -e "${green}"
    cat << "EOF"
 ╔═══════════════════════════════════════════════════╗
 ║                                                   ║
 ║       ✨ Installation Completed Successfully! ✨   ║
 ║                                                   ║
 ╚═══════════════════════════════════════════════════╝
EOF
    echo -e "${plain}"
    echo ""
    
    echo -e "${cyan}╔═══════════════════════════════════════════════════╗${plain}"
    echo -e "${cyan}║              🌐 Access Information                ║${plain}"
    echo -e "${cyan}╚═══════════════════════════════════════════════════╝${plain}"
    echo ""
    echo -e "  ${yellow}📡 Panel URL (HTTP):${plain}"
    echo -e "     ${green}http://${server_ip}:${config_port}${plain}"
    echo ""
    echo -e "  ${yellow}🔒 Panel URL (HTTPS):${plain}"
    echo -e "     ${green}https://${server_ip}:${config_port}${plain}"
    echo ""
    echo -e "  ${yellow}👤 Admin Username:${plain}"
    echo -e "     ${green}${config_username}${plain}"
    echo ""
    echo -e "  ${yellow}🔑 Admin Password:${plain}"
    echo -e "     ${green}${config_password}${plain}"
    echo ""
    echo -e "  ${yellow}📂 Installation Path:${plain}"
    echo -e "     ${green}${install_path}${plain}"
    echo ""
    
    echo -e "${cyan}╔═══════════════════════════════════════════════════╗${plain}"
    echo -e "${cyan}║              ⭐ RBAC Features Enabled             ║${plain}"
    echo -e "${cyan}╚═══════════════════════════════════════════════════╝${plain}"
    echo ""
    echo -e "  ${green}✓${plain} Role-Based Access Control"
    echo -e "  ${green}✓${plain} Multi-Vendor Support (up to 50+ vendors)"
    echo -e "  ${green}✓${plain} Client Ownership Tracking"
    echo -e "  ${green}✓${plain} Isolated Vendor Dashboard"
    echo -e "  ${green}✓${plain} Admin Vendor Management"
    echo ""
    
    echo -e "${cyan}╔═══════════════════════════════════════════════════╗${plain}"
    echo -e "${cyan}║              📚 Useful Commands                   ║${plain}"
    echo -e "${cyan}╚═══════════════════════════════════════════════════╝${plain}"
    echo ""
    echo -e "  ${blue}x-ui${plain}                 - Open control menu"
    echo -e "  ${blue}x-ui start${plain}           - Start panel service"
    echo -e "  ${blue}x-ui stop${plain}            - Stop panel service"
    echo -e "  ${blue}x-ui restart${plain}         - Restart panel service"
    echo -e "  ${blue}x-ui status${plain}          - Check service status"
    echo -e "  ${blue}x-ui log${plain}             - View panel logs"
    echo -e "  ${blue}systemctl status x-ui${plain} - Detailed service status"
    echo ""
    
    echo -e "${cyan}╔═══════════════════════════════════════════════════╗${plain}"
    echo -e "${cyan}║              ⚠️  Important Security Notes          ║${plain}"
    echo -e "${cyan}╚═══════════════════════════════════════════════════╝${plain}"
    echo ""
    echo -e "  ${yellow}1.${plain} Change your password after first login"
    echo -e "  ${yellow}2.${plain} Configure firewall: ${green}ufw allow ${config_port}/tcp${plain}"
    echo -e "  ${yellow}3.${plain} Enable HTTPS with SSL certificate in panel settings"
    echo -e "  ${yellow}4.${plain} Keep your system and panel updated regularly"
    echo ""
    
    echo -e "${cyan}╔═══════════════════════════════════════════════════╗${plain}"
    echo -e "${cyan}║              📖 Documentation & Support           ║${plain}"
    echo -e "${cyan}╚═══════════════════════════════════════════════════╝${plain}"
    echo ""
    echo -e "  ${purple}🌐 GitHub:${plain}"
    echo -e "     ${blue}https://github.com/Farsimen/3x-ui-pajo${plain}"
    echo ""
    echo -e "  ${purple}📚 Documentation:${plain}"
    echo -e "     ${blue}https://github.com/Farsimen/3x-ui-pajo/tree/rbac-implementation${plain}"
    echo ""
    echo -e "  ${purple}📋 Guide:${plain}"
    echo -e "     ${blue}- IMPLEMENTATION_STEPS.md (Persian/English)${plain}"
    echo -e "     ${blue}- VENDOR_PERMISSIONS.md${plain}"
    echo ""
    
    echo -e "${green}═════════════════════════════════════════════════════${plain}"
    echo -e "         ${purple}Made with ❤️  by Farsimen for the community${plain}"
    echo -e "${green}═════════════════════════════════════════════════════${plain}"
    echo ""
    
    echo -e "${yellow}💡 Tip: Save this information in a secure place!${plain}"
    echo ""
}

# Main execution
main() {
    print_banner
    echo -e "${cyan}Checking system requirements...${plain}"
    echo -e "${green}✓ OS: ${release}${plain}"
    echo -e "${green}✓ Architecture: $(arch)${plain}"
    echo ""
    sleep 1
    
    install_base
    echo ""
    interactive_config
    install_x-ui_pajo
    show_access_info
}

# Run
main
