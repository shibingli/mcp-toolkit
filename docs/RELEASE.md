# 发布指南 / Release Guide

本文档介绍如何构建、打包和发布 MCP Toolkit 到 GitHub Releases，以及如何通过 `uv` 工具安装。

This document describes how to build, package, and release MCP Toolkit to GitHub Releases, and how to install it using the `uv` tool.

---

## 目录 / Table of Contents

- [本地构建](#本地构建--local-build)
- [GitHub Actions 自动发布](#github-actions-自动发布--github-actions-auto-release)
- [安装方式](#安装方式--installation-methods)
- [使用 uv 工具](#使用-uv-工具--using-uv-tool)

---

## 本地构建 / Local Build

### Linux/macOS

```bash
# 赋予执行权限 / Grant execute permission
chmod +x scripts/build.sh

# 构建所有平台 / Build all platforms
./scripts/build.sh

# 指定版本号构建 / Build with specific version
VERSION=v1.0.0 ./scripts/build.sh
```

### Windows

```powershell
# 使用 PowerShell 或 CMD
.\scripts\build.bat

# 指定版本号构建 / Build with specific version
set VERSION=v1.0.0
.\scripts\build.bat
```

### 构建输出 / Build Output

构建完成后，会在 `dist/` 目录下生成以下文件：

After building, the following files will be generated in the `dist/` directory:

```
dist/
├── mcp-toolkit-v1.0.0-windows-amd64.zip
├── mcp-toolkit-v1.0.0-windows-arm64.zip
├── mcp-toolkit-v1.0.0-linux-amd64.tar.gz
├── mcp-toolkit-v1.0.0-linux-arm64.tar.gz
├── mcp-toolkit-v1.0.0-darwin-amd64.tar.gz
├── mcp-toolkit-v1.0.0-darwin-arm64.tar.gz
└── checksums.txt
```

---

## GitHub Actions 自动发布 / GitHub Actions Auto Release

### 配置步骤 / Configuration Steps

1. **推送代码到 GitHub**

```bash
git add .
git commit -m "Add release workflow"
git push origin main
```

2. **创建并推送标签 / Create and push tag**

```bash
# 创建标签 / Create tag
git tag -a v1.0.0 -m "Release version 1.0.0"

# 推送标签 / Push tag
git push origin v1.0.0
```

3. **自动构建和发布 / Auto build and release**

推送标签后，GitHub Actions 会自动：
- 构建所有平台的二进制文件
- 生成 SHA256 校验和
- 创建 GitHub Release
- 上传所有构建产物

After pushing the tag, GitHub Actions will automatically:
- Build binaries for all platforms
- Generate SHA256 checksums
- Create GitHub Release
- Upload all build artifacts

### 手动触发 / Manual Trigger

也可以在 GitHub Actions 页面手动触发工作流：

You can also manually trigger the workflow on the GitHub Actions page:

1. 进入仓库的 Actions 页面 / Go to the repository's Actions page
2. 选择 "Release" 工作流 / Select the "Release" workflow
3. 点击 "Run workflow" / Click "Run workflow"

---

## 安装方式 / Installation Methods

### 方式 1: 使用安装脚本 / Method 1: Using Installation Script

#### Linux/macOS

```bash
# 下载并运行安装脚本 / Download and run installation script
curl -fsSL https://raw.githubusercontent.com/shibingli/mcp-toolkit/main/scripts/install.sh | bash

# 或者手动下载后安装 / Or download manually and install
wget https://raw.githubusercontent.com/shibingli/mcp-toolkit/main/scripts/install.sh
chmod +x install.sh
./install.sh

# 指定版本安装 / Install specific version
VERSION=v1.0.0 ./install.sh

# 指定安装目录 / Specify installation directory
INSTALL_DIR=/usr/local/bin ./install.sh

# 卸载 / Uninstall
./install.sh uninstall
```

#### Windows

```powershell
# 下载并运行安装脚本 / Download and run installation script
Invoke-WebRequest -Uri "https://raw.githubusercontent.com/shibingli/mcp-toolkit/main/scripts/install.ps1" -OutFile "install.ps1"
.\install.ps1

# 指定版本安装 / Install specific version
.\install.ps1 -Version "v1.0.0"

# 指定安装目录 / Specify installation directory
.\install.ps1 -InstallDir "C:\Program Files\mcp-toolkit"

# 卸载 / Uninstall
.\install.ps1 uninstall
```

### 方式 2: 手动下载 / Method 2: Manual Download

1. 访问 [Releases 页面](https://github.com/shibingli/mcp-toolkit/releases)
2. 下载对应平台的压缩包
3. 解压到目标目录
4. 将二进制文件路径添加到 PATH 环境变量

---

## 使用 uv 工具 / Using uv Tool

`uv` 是一个现代化的 Python 包管理器，也支持安装其他工具。

`uv` is a modern Python package manager that also supports installing other tools.

### 前提条件 / Prerequisites

首先安装 `uv`：

First, install `uv`:

```bash
# Linux/macOS
curl -LsSf https://astral.sh/uv/install.sh | sh

# Windows
powershell -c "irm https://astral.sh/uv/install.ps1 | iex"
```

### 通过 uv 安装 MCP Toolkit / Install MCP Toolkit via uv

目前 `uv` 主要用于 Python 工具，对于 Go 二进制文件，建议使用以下方式：

Currently, `uv` is mainly used for Python tools. For Go binaries, the following methods are recommended:

#### 方式 A: 使用 uv tool install (如果项目发布到 PyPI)

如果将 MCP Toolkit 包装为 Python 包并发布到 PyPI：

If MCP Toolkit is wrapped as a Python package and published to PyPI:

```bash
uv tool install mcp-toolkit
```

#### 方式 B: 使用 pipx (推荐用于 Python 包装的工具)

```bash
pipx install mcp-toolkit
```

#### 方式 C: 直接使用安装脚本 (推荐)

```bash
# 这是最直接的方式 / This is the most direct way
curl -fsSL https://raw.githubusercontent.com/shibingli/mcp-toolkit/main/scripts/install.sh | bash
```

---

## 验证安装 / Verify Installation

```bash
# 检查版本 / Check version
mcp-toolkit --help

# 运行测试 / Run test
mcp-toolkit -transport stdio
```

---

## 更新 / Update

```bash
# 重新运行安装脚本即可更新到最新版本
# Re-run the installation script to update to the latest version
./install.sh
```

---

## 故障排查 / Troubleshooting

### 问题 1: 找不到命令 / Command not found

**解决方案 / Solution:**

确保安装目录在 PATH 中：

Make sure the installation directory is in PATH:

```bash
# Linux/macOS
export PATH="$PATH:$HOME/.local/bin"

# Windows (PowerShell)
$env:Path += ";$env:LOCALAPPDATA\Programs\mcp-toolkit"
```

### 问题 2: 权限被拒绝 / Permission denied

**解决方案 / Solution:**

```bash
# Linux/macOS
chmod +x ~/.local/bin/mcp-toolkit
```

### 问题 3: 下载失败 / Download failed

**解决方案 / Solution:**

检查网络连接，或手动从 GitHub Releases 页面下载。

Check network connection, or manually download from GitHub Releases page.

---

## 相关链接 / Related Links

- [GitHub Repository](https://github.com/shibingli/mcp-toolkit)
- [GitHub Releases](https://github.com/shibingli/mcp-toolkit/releases)
- [Documentation](https://github.com/shibingli/mcp-toolkit/blob/main/README.md)

