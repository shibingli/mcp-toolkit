#!/bin/bash
# 跨平台构建脚本 / Cross-platform build script
# 用于编译 Windows、Linux、macOS 的二进制文件

set -e

# 项目信息 / Project information
PROJECT_NAME="mcp-toolkit"
VERSION=${VERSION:-$(git describe --tags --always --dirty 2>/dev/null || echo "dev")}
BUILD_TIME=$(date -u '+%Y-%m-%d_%H:%M:%S')
GIT_COMMIT=$(git rev-parse --short HEAD 2>/dev/null || echo "unknown")

# 构建标签 / Build tags
BUILD_TAGS="sonic"

# 输出目录 / Output directory
OUTPUT_DIR="dist"
rm -rf ${OUTPUT_DIR}
mkdir -p ${OUTPUT_DIR}

# 构建信息 / Build information
LDFLAGS="-s -w -X main.Version=${VERSION} -X main.BuildTime=${BUILD_TIME} -X main.GitCommit=${GIT_COMMIT}"

echo "========================================="
echo "Building ${PROJECT_NAME} ${VERSION}"
echo "Build Time: ${BUILD_TIME}"
echo "Git Commit: ${GIT_COMMIT}"
echo "========================================="

# 构建函数 / Build function
build() {
    local GOOS=$1
    local GOARCH=$2
    local OUTPUT_NAME="${PROJECT_NAME}"
    
    if [ "$GOOS" = "windows" ]; then
        OUTPUT_NAME="${OUTPUT_NAME}.exe"
    fi
    
    local OUTPUT_PATH="${OUTPUT_DIR}/${PROJECT_NAME}-${GOOS}-${GOARCH}"
    mkdir -p ${OUTPUT_PATH}
    
    echo "Building for ${GOOS}/${GOARCH}..."
    
    # 设置环境变量 / Set environment variables
    export GOOS=${GOOS}
    export GOARCH=${GOARCH}
    export CGO_ENABLED=0
    
    # Linux 下需要额外的环境变量 / Additional environment variables for Linux
    if [ "$GOOS" = "linux" ]; then
        export GOEXPERIMENT="jsonv2,greenteagc"
    fi
    
    # 构建 / Build
    go build -tags="${BUILD_TAGS}" \
        -ldflags="${LDFLAGS}" \
        -o "${OUTPUT_PATH}/${OUTPUT_NAME}" \
        .
    
    # 复制 README 和 LICENSE / Copy README and LICENSE
    cp README.md ${OUTPUT_PATH}/
    cp LICENSE ${OUTPUT_PATH}/
    
    # 创建压缩包 / Create archive
    cd ${OUTPUT_DIR}
    if [ "$GOOS" = "windows" ]; then
        # Windows 使用 zip / Windows uses zip
        if command -v zip >/dev/null 2>&1; then
            zip -r "${PROJECT_NAME}-${VERSION}-${GOOS}-${GOARCH}.zip" "${PROJECT_NAME}-${GOOS}-${GOARCH}"
        else
            echo "Warning: zip command not found, skipping archive creation for ${GOOS}/${GOARCH}"
        fi
    else
        # Linux/macOS 使用 tar.gz / Linux/macOS uses tar.gz
        tar -czf "${PROJECT_NAME}-${VERSION}-${GOOS}-${GOARCH}.tar.gz" "${PROJECT_NAME}-${GOOS}-${GOARCH}"
    fi
    cd ..
    
    echo "✓ Built ${GOOS}/${GOARCH}"
}

# 构建所有平台 / Build all platforms
echo ""
echo "Building for all platforms..."
echo ""

# Windows
build windows amd64
build windows arm64

# Linux
build linux amd64
build linux arm64

# macOS
build darwin amd64
build darwin arm64

# 生成校验和文件 / Generate checksums
echo ""
echo "Generating checksums..."
cd ${OUTPUT_DIR}
sha256sum *.zip *.tar.gz 2>/dev/null > checksums.txt || true
cd ..

echo ""
echo "========================================="
echo "Build completed successfully!"
echo "Output directory: ${OUTPUT_DIR}"
echo "========================================="
echo ""
echo "Generated files:"
echo ""
echo "Windows packages (zip):"
ls -lh ${OUTPUT_DIR}/*windows*.zip 2>/dev/null || echo "  No Windows packages found"
echo ""
echo "Linux/macOS packages (tar.gz):"
ls -lh ${OUTPUT_DIR}/*{linux,darwin}*.tar.gz 2>/dev/null || echo "  No Linux/macOS packages found"
echo ""
echo "Checksums:"
ls -lh ${OUTPUT_DIR}/checksums.txt 2>/dev/null || echo "  No checksums file found"

