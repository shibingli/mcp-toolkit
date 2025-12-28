#!/bin/bash
# MCP Toolkit 安装脚本 / MCP Toolkit installation script
# 支持 Linux 和 macOS / Supports Linux and macOS

set -e

# 配置 / Configuration
REPO="shibingli/mcp-toolkit"  # 请替换为你的 GitHub 仓库 / Replace with your GitHub repository
INSTALL_DIR="${INSTALL_DIR:-$HOME/.local/bin}"
BINARY_NAME="mcp-toolkit"

# 颜色输出 / Color output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

error() {
    echo -e "${RED}[ERROR]${NC} $1"
    exit 1
}

# 检测操作系统和架构 / Detect OS and architecture
detect_platform() {
    local OS=$(uname -s | tr '[:upper:]' '[:lower:]')
    local ARCH=$(uname -m)
    
    case "$OS" in
        linux*)
            OS="linux"
            ;;
        darwin*)
            OS="darwin"
            ;;
        *)
            error "Unsupported operating system: $OS"
            ;;
    esac
    
    case "$ARCH" in
        x86_64|amd64)
            ARCH="amd64"
            ;;
        aarch64|arm64)
            ARCH="arm64"
            ;;
        *)
            error "Unsupported architecture: $ARCH"
            ;;
    esac
    
    echo "${OS}-${ARCH}"
}

# 获取最新版本 / Get latest version
get_latest_version() {
    local VERSION=$(curl -s "https://api.github.com/repos/${REPO}/releases/latest" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')
    if [ -z "$VERSION" ]; then
        error "Failed to get latest version"
    fi
    echo "$VERSION"
}

# 下载并安装 / Download and install
install() {
    info "Detecting platform..."
    local PLATFORM=$(detect_platform)
    info "Platform: $PLATFORM"
    
    info "Getting latest version..."
    local VERSION=${VERSION:-$(get_latest_version)}
    info "Version: $VERSION"
    
    local DOWNLOAD_URL="https://github.com/${REPO}/releases/download/${VERSION}/mcp-toolkit-${VERSION}-${PLATFORM}.tar.gz"
    local TEMP_DIR=$(mktemp -d)
    
    info "Downloading from $DOWNLOAD_URL..."
    if ! curl -L -o "${TEMP_DIR}/mcp-toolkit.tar.gz" "$DOWNLOAD_URL"; then
        error "Failed to download"
    fi
    
    info "Extracting..."
    tar -xzf "${TEMP_DIR}/mcp-toolkit.tar.gz" -C "${TEMP_DIR}"
    
    info "Installing to $INSTALL_DIR..."
    mkdir -p "$INSTALL_DIR"
    
    local EXTRACTED_DIR="${TEMP_DIR}/mcp-toolkit-${PLATFORM#*-}"
    if [ ! -d "$EXTRACTED_DIR" ]; then
        EXTRACTED_DIR=$(find "${TEMP_DIR}" -type d -name "mcp-toolkit-*" | head -n 1)
    fi
    
    cp "${EXTRACTED_DIR}/${BINARY_NAME}" "$INSTALL_DIR/"
    chmod +x "$INSTALL_DIR/${BINARY_NAME}"
    
    rm -rf "${TEMP_DIR}"
    
    info "Installation completed!"
    info "Binary installed to: $INSTALL_DIR/${BINARY_NAME}"
    
    # 检查是否在 PATH 中 / Check if in PATH
    if ! echo "$PATH" | grep -q "$INSTALL_DIR"; then
        warn "$INSTALL_DIR is not in your PATH"
        warn "Add the following line to your ~/.bashrc or ~/.zshrc:"
        echo ""
        echo "    export PATH=\"\$PATH:$INSTALL_DIR\""
        echo ""
    fi
    
    info "Run '${BINARY_NAME} --help' to get started"
}

# 卸载 / Uninstall
uninstall() {
    if [ -f "$INSTALL_DIR/${BINARY_NAME}" ]; then
        rm "$INSTALL_DIR/${BINARY_NAME}"
        info "Uninstalled successfully"
    else
        warn "Binary not found at $INSTALL_DIR/${BINARY_NAME}"
    fi
}

# 主函数 / Main function
main() {
    case "${1:-install}" in
        install)
            install
            ;;
        uninstall)
            uninstall
            ;;
        *)
            echo "Usage: $0 {install|uninstall}"
            exit 1
            ;;
    esac
}

main "$@"

