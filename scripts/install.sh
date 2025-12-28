#!/bin/bash
# MCP Toolkit 安装脚本 / MCP Toolkit installation script
# 支持 Linux 和 macOS / Supports Linux and macOS

set -e

# 配置 / Configuration
REPO="${REPO:-shibingli/mcp-toolkit}"  # 可通过环境变量覆盖 / Can be overridden by environment variable
INSTALL_DIR="${INSTALL_DIR:-$HOME/.local/bin}"
BINARY_NAME="mcp-toolkit"
VERIFY_CHECKSUM="${VERIFY_CHECKSUM:-true}"  # 是否验证校验和 / Whether to verify checksum

# 颜色输出 / Color output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
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

debug() {
    if [ "${DEBUG:-false}" = "true" ]; then
        echo -e "${BLUE}[DEBUG]${NC} $1"
    fi
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
    debug "Fetching latest version from GitHub API..."
    local VERSION=$(curl -sSf "https://api.github.com/repos/${REPO}/releases/latest" 2>/dev/null | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')
    if [ -z "$VERSION" ]; then
        error "Failed to get latest version. Please check your internet connection or specify a version manually."
    fi
    debug "Latest version: $VERSION"
    echo "$VERSION"
}

# 获取已安装版本 / Get installed version
get_installed_version() {
    if [ -f "$INSTALL_DIR/${BINARY_NAME}" ]; then
        local VERSION=$("$INSTALL_DIR/${BINARY_NAME}" --version 2>/dev/null | head -n 1 | awk '{print $NF}' || echo "")
        echo "$VERSION"
    else
        echo ""
    fi
}

# 验证校验和 / Verify checksum
verify_checksum() {
    local FILE="$1"
    local EXPECTED_CHECKSUM="$2"

    if [ "$VERIFY_CHECKSUM" != "true" ]; then
        debug "Checksum verification disabled"
        return 0
    fi

    if [ -z "$EXPECTED_CHECKSUM" ]; then
        warn "No checksum provided, skipping verification"
        return 0
    fi

    info "Verifying checksum..."
    local ACTUAL_CHECKSUM=$(sha256sum "$FILE" | awk '{print $1}')

    if [ "$ACTUAL_CHECKSUM" = "$EXPECTED_CHECKSUM" ]; then
        info "Checksum verified successfully"
        return 0
    else
        error "Checksum verification failed!\nExpected: $EXPECTED_CHECKSUM\nActual: $ACTUAL_CHECKSUM"
    fi
}

# 下载并安装 / Download and install
install() {
    # 检查已安装版本 / Check installed version
    local INSTALLED_VERSION=$(get_installed_version)
    if [ -n "$INSTALLED_VERSION" ]; then
        info "Current installed version: $INSTALLED_VERSION"
    fi

    info "Detecting platform..."
    local PLATFORM=$(detect_platform)
    info "Platform: $PLATFORM"

    info "Getting latest version..."
    local VERSION=${VERSION:-$(get_latest_version)}
    info "Target version: $VERSION"

    # 检查是否需要更新 / Check if update is needed
    if [ -n "$INSTALLED_VERSION" ] && [ "$INSTALLED_VERSION" = "${VERSION#v}" ]; then
        info "Already up to date (version $INSTALLED_VERSION)"
        return 0
    fi

    local ARCHIVE_NAME="mcp-toolkit-${VERSION}-${PLATFORM}.tar.gz"
    local DOWNLOAD_URL="https://github.com/${REPO}/releases/download/${VERSION}/${ARCHIVE_NAME}"
    local CHECKSUMS_URL="https://github.com/${REPO}/releases/download/${VERSION}/checksums.txt"
    local TEMP_DIR=$(mktemp -d)

    trap "rm -rf ${TEMP_DIR}" EXIT

    info "Downloading from $DOWNLOAD_URL..."
    if ! curl -sSfL -o "${TEMP_DIR}/${ARCHIVE_NAME}" "$DOWNLOAD_URL"; then
        error "Failed to download. Please check the URL and your internet connection."
    fi

    # 下载并验证校验和 / Download and verify checksum
    if [ "$VERIFY_CHECKSUM" = "true" ]; then
        debug "Downloading checksums..."
        if curl -sSfL -o "${TEMP_DIR}/checksums.txt" "$CHECKSUMS_URL" 2>/dev/null; then
            local EXPECTED_CHECKSUM=$(grep "$ARCHIVE_NAME" "${TEMP_DIR}/checksums.txt" | awk '{print $1}')
            verify_checksum "${TEMP_DIR}/${ARCHIVE_NAME}" "$EXPECTED_CHECKSUM"
        else
            warn "Could not download checksums file, skipping verification"
        fi
    fi

    info "Extracting..."
    tar -xzf "${TEMP_DIR}/${ARCHIVE_NAME}" -C "${TEMP_DIR}"

    info "Installing to $INSTALL_DIR..."
    mkdir -p "$INSTALL_DIR"

    # 查找提取的目录 / Find extracted directory
    local EXTRACTED_DIR=$(find "${TEMP_DIR}" -type d -name "mcp-toolkit-*" -print -quit)
    if [ -z "$EXTRACTED_DIR" ]; then
        error "Failed to find extracted directory"
    fi

    debug "Extracted directory: $EXTRACTED_DIR"

    # 备份旧版本 / Backup old version
    if [ -f "$INSTALL_DIR/${BINARY_NAME}" ]; then
        info "Backing up old version..."
        cp "$INSTALL_DIR/${BINARY_NAME}" "$INSTALL_DIR/${BINARY_NAME}.backup"
    fi

    # 复制新版本 / Copy new version
    if [ ! -f "${EXTRACTED_DIR}/${BINARY_NAME}" ]; then
        error "Binary not found in extracted archive"
    fi

    cp "${EXTRACTED_DIR}/${BINARY_NAME}" "$INSTALL_DIR/"
    chmod +x "$INSTALL_DIR/${BINARY_NAME}"

    # 验证安装 / Verify installation
    if ! "$INSTALL_DIR/${BINARY_NAME}" --version >/dev/null 2>&1; then
        error "Installation verification failed. Binary may be corrupted."
    fi

    # 删除备份 / Remove backup
    if [ -f "$INSTALL_DIR/${BINARY_NAME}.backup" ]; then
        rm "$INSTALL_DIR/${BINARY_NAME}.backup"
    fi

    info "Installation completed successfully!"
    info "Binary installed to: $INSTALL_DIR/${BINARY_NAME}"

    # 显示版本信息 / Show version info
    local NEW_VERSION=$("$INSTALL_DIR/${BINARY_NAME}" --version 2>/dev/null | head -n 1 || echo "unknown")
    info "Installed version: $NEW_VERSION"

    # 检查是否在 PATH 中 / Check if in PATH
    if ! echo "$PATH" | grep -q "$INSTALL_DIR"; then
        warn "$INSTALL_DIR is not in your PATH"
        warn "Add the following line to your shell configuration file:"
        warn "  For bash: ~/.bashrc"
        warn "  For zsh: ~/.zshrc"
        echo ""
        echo "    export PATH=\"\$PATH:$INSTALL_DIR\""
        echo ""
    else
        info "Run '${BINARY_NAME} --help' to get started"
    fi
}

# 卸载 / Uninstall
uninstall() {
    if [ -f "$INSTALL_DIR/${BINARY_NAME}" ]; then
        local VERSION=$(get_installed_version)
        info "Uninstalling ${BINARY_NAME} ${VERSION}..."
        rm "$INSTALL_DIR/${BINARY_NAME}"
        info "Uninstalled successfully"
    else
        warn "Binary not found at $INSTALL_DIR/${BINARY_NAME}"
        warn "Nothing to uninstall"
    fi
}

# 更新 / Update
update() {
    local INSTALLED_VERSION=$(get_installed_version)
    if [ -z "$INSTALLED_VERSION" ]; then
        error "${BINARY_NAME} is not installed. Run 'install' first."
    fi

    info "Checking for updates..."
    local LATEST_VERSION=$(get_latest_version)

    if [ "$INSTALLED_VERSION" = "${LATEST_VERSION#v}" ]; then
        info "Already up to date (version $INSTALLED_VERSION)"
    else
        info "Update available: $INSTALLED_VERSION → ${LATEST_VERSION#v}"
        VERSION="$LATEST_VERSION" install
    fi
}

# 显示版本 / Show version
show_version() {
    local VERSION=$(get_installed_version)
    if [ -n "$VERSION" ]; then
        echo "${BINARY_NAME} version $VERSION"
        echo "Installed at: $INSTALL_DIR/${BINARY_NAME}"
    else
        echo "${BINARY_NAME} is not installed"
        exit 1
    fi
}

# 显示帮助 / Show help
show_help() {
    cat << EOF
MCP Toolkit Installation Script

Usage: $0 [COMMAND] [OPTIONS]

Commands:
    install     Install ${BINARY_NAME} (default)
    uninstall   Uninstall ${BINARY_NAME}
    update      Update ${BINARY_NAME} to the latest version
    version     Show installed version
    help        Show this help message

Environment Variables:
    VERSION         Specify version to install (e.g., v1.0.0)
    INSTALL_DIR     Installation directory (default: \$HOME/.local/bin)
    REPO            GitHub repository (default: shibingli/mcp-toolkit)
    VERIFY_CHECKSUM Verify checksum (default: true)
    DEBUG           Enable debug output (default: false)

Examples:
    # Install latest version
    $0 install

    # Install specific version
    VERSION=v1.0.0 $0 install

    # Install to custom directory
    INSTALL_DIR=/usr/local/bin $0 install

    # Update to latest version
    $0 update

    # Uninstall
    $0 uninstall

    # Show version
    $0 version

For more information, visit: https://github.com/${REPO}
EOF
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
        update)
            update
            ;;
        version)
            show_version
            ;;
        help|--help|-h)
            show_help
            ;;
        *)
            error "Unknown command: $1\nRun '$0 help' for usage information"
            ;;
    esac
}

main "$@"

