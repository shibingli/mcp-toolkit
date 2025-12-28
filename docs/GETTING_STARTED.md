# MCP Toolkit 快速开始指南 / Getting Started Guide

本指南将帮助你快速开始使用 MCP Toolkit，包括安装、配置、开发和发布。

This guide will help you quickly get started with MCP Toolkit, including installation, configuration, development, and
publishing.

---

## 📋 目录 / Table of Contents

- [用户快速开始](#-用户快速开始--quick-start-for-users)
- [开发者快速开始](#-开发者快速开始--quick-start-for-developers)
- [前置依赖](#-前置依赖--prerequisites)
- [从源码构建](#-从源码构建--build-from-source)
- [发布项目](#-发布项目--publish-project)
- [常见问题](#-常见问题--faq)

---

## 🚀 用户快速开始 / Quick Start for Users

### 最简单的方式：使用 uvx (无需安装)

```bash
# 直接运行，无需安装
uvx mcp-sandbox-toolkit --help

# 启动服务器
uvx mcp-sandbox-toolkit
```

### 方式 2: 使用 uv 安装

```bash
# 安装 uv (如果还没有安装)
curl -LsSf https://astral.sh/uv/install.sh | sh

# 安装 MCP Toolkit
uv tool install mcp-sandbox-toolkit

# 运行
mcp-sandbox-toolkit --help
```

### 方式 3: 使用安装脚本

**Linux/macOS:**

```bash
curl -fsSL https://raw.githubusercontent.com/shibingli/mcp-toolkit/main/scripts/install.sh | bash
```

**Windows (PowerShell):**

```powershell
irm https://raw.githubusercontent.com/shibingli/mcp-toolkit/main/scripts/install.ps1 | iex
```

### 运行服务器

```bash
# Stdio 模式 (默认)
mcp-sandbox-toolkit

# HTTP 模式
mcp-sandbox-toolkit -transport http -http-port 8080

# SSE 模式
mcp-sandbox-toolkit -transport sse -sse-port 8081

# 指定沙箱目录
mcp-sandbox-toolkit -sandbox /path/to/sandbox
```

**详细安装说明请参考**: [INSTALLATION.md](INSTALLATION.md)

---

## 👨‍💻 开发者快速开始 / Quick Start for Developers

### 1. 克隆项目

```bash
git clone https://github.com/shibingli/mcp-toolkit.git
cd mcp-toolkit
```

### 2. 安装依赖

```bash
go mod download
```

### 3. 编译

```bash
# 使用 Sonic (推荐，最高性能)
go build -tags="sonic" -o mcp-toolkit .

# 或使用 Makefile
make build
```

### 4. 运行

```bash
# Windows
.\mcp-toolkit.exe

# Linux/macOS
./mcp-toolkit
```

### 5. 运行测试

```bash
# 运行所有测试
go test -v ./...

# 查看测试覆盖率
go test -v -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# 或使用 Makefile
make test
make test-cover
```

---

## 📦 前置依赖 / Prerequisites

### 必需工具 / Required Tools

| 工具         | 版本要求    | 用途           | 检查命令               |
|------------|---------|--------------|--------------------|
| **Go**     | 1.25.5+ | 编译和构建        | `go version`       |
| **Git**    | 2.0+    | 版本控制         | `git --version`    |
| **Python** | 3.8+    | PyPI 发布 (可选) | `python --version` |
| **Make**   | 任意版本    | 构建工具 (可选)    | `make --version`   |

### 安装 Go

**Windows:**

```powershell
# 访问 https://go.dev/dl/
# 下载 go1.25.5.windows-amd64.msi 并安装

# 验证安装
go version
```

**Linux:**

```bash
wget https://go.dev/dl/go1.25.5.linux-amd64.tar.gz
sudo rm -rf /usr/local/go
sudo tar -C /usr/local -xzf go1.25.5.linux-amd64.tar.gz

# 添加到 PATH
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
source ~/.bashrc

# 验证安装
go version
```

**macOS:**

```bash
# 使用 Homebrew
brew install go@1.25

# 验证安装
go version
```

### 安装 Python 构建工具 (可选，用于 PyPI 发布)

```bash
python -m pip install --upgrade pip
python -m pip install --upgrade build twine hatchling
```

---

## 🔨 从源码构建 / Build from Source

### 基本构建

```bash
# 克隆项目
git clone https://github.com/shibingli/mcp-toolkit.git
cd mcp-toolkit

# 安装依赖
go mod download

# 编译 (使用 Sonic，推荐)
go build -tags="sonic" -o mcp-toolkit .

# 运行
./mcp-toolkit --help
```

### 使用 Makefile

```bash
# 查看所有可用命令
make help

# 构建当前平台
make build

# 构建所有平台
make build-all

# 运行测试
make test

# 查看测试覆盖率
make test-cover

# 格式化代码
make fmt

# 清理构建产物
make clean
```

### 构建所有平台

```bash
# 使用 Makefile
make build-all

# 或使用脚本
./scripts/build.sh  # Linux/macOS
.\scripts\build.bat  # Windows
```

构建产物将保存在 `dist/` 目录：

- `mcp-toolkit-vX.X.X-windows-amd64.zip`
- `mcp-toolkit-vX.X.X-linux-amd64.tar.gz`
- `mcp-toolkit-vX.X.X-darwin-arm64.tar.gz`
- 等等...

---

## 📤 发布项目 / Publish Project

### 快速发布 (5 分钟)

#### 步骤 1: 配置 GitHub Secrets

在 GitHub 仓库设置中添加以下 Secrets：

1. **Settings** → **Secrets and variables** → **Actions** → **New repository secret**

2. 添加 `PYPI_API_TOKEN` (可选，用于发布到 PyPI):
    - 访问 https://pypi.org/manage/account/token/
    - 创建新 Token
    - 复制 Token 并添加到 GitHub Secrets

#### 步骤 2: 更新仓库信息

将以下文件中的 `shibingli/mcp-toolkit` 替换为你的仓库地址：

```bash
# 使用 sed 批量替换 (Linux/macOS)
find . -type f \( -name "*.sh" -o -name "*.ps1" -o -name "*.py" -o -name "*.toml" -o -name "*.md" \) \
  -exec sed -i 's/shibingli\/mcp-toolkit/YOUR_USERNAME\/mcp-toolkit/g' {} +

# Windows (PowerShell)
Get-ChildItem -Recurse -Include *.sh,*.ps1,*.py,*.toml,*.md |
  ForEach-Object {
    (Get-Content $_.FullName) -replace 'shibingli/mcp-toolkit', 'YOUR_USERNAME/mcp-toolkit' |
    Set-Content $_.FullName
  }
```

或手动编辑以下文件：

- `scripts/install.sh` (第 8 行)
- `scripts/install.ps1` (第 11 行)
- `python/mcp_toolkit_wrapper/installer.py` (第 18 行)
- `pyproject.toml` (URLs 部分)

#### 步骤 3: 创建并推送标签

```bash
# 1. 提交所有更改
git add -A
git commit -m "Prepare for release v1.0.0"
git push origin main

# 2. 创建标签
git tag -a v1.0.0 -m "Release v1.0.0"

# 3. 推送标签 (这将触发自动发布)
git push origin v1.0.0
```

#### 步骤 4: 等待自动发布

GitHub Actions 将自动：

- ✅ 运行所有测试
- ✅ 构建所有平台的二进制文件
- ✅ 生成校验和文件
- ✅ 创建 GitHub Release
- ✅ 上传所有构建产物
- ✅ 构建并发布 Python 包到 PyPI (如果配置了 Token)

查看发布进度：

```
https://github.com/YOUR_USERNAME/mcp-toolkit/actions
```

#### 步骤 5: 验证发布

```bash
# 检查 GitHub Release
# https://github.com/YOUR_USERNAME/mcp-toolkit/releases

# 检查 PyPI (如果发布了)
# https://pypi.org/project/mcp-sandbox-toolkit/

# 测试安装
uv tool install mcp-sandbox-toolkit
mcp-sandbox-toolkit --help
```

### 发布检查清单

发布前确认：

- [ ] ✅ 所有测试通过 (`make test`)
- [ ] ✅ 代码已格式化 (`make fmt`)
- [ ] ✅ 更新了 CHANGELOG.md
- [ ] ✅ 更新了版本号相关文档
- [ ] ✅ 仓库信息已更新
- [ ] ✅ GitHub Secrets 已配置
- [ ] ✅ 本地构建成功 (`make build-all`)

### 版本管理

版本号由 Git Tag 自动控制：

```bash
# 创建新版本
git tag -a v1.2.0 -m "Release v1.2.0"
git push origin v1.2.0

# 查看所有版本
git tag -l

# 删除错误的标签
git tag -d v1.2.0
git push origin :refs/tags/v1.2.0
```

**详细版本管理说明**: [VERSION_MANAGEMENT.md](../VERSION_MANAGEMENT.md)

---

## 🎯 常见任务 / Common Tasks

### 本地开发

```bash
# 运行开发模式
make dev

# 或直接运行
go run -tags="sonic" . -sandbox ./sandbox
```

### 运行测试

```bash
# 运行所有测试
make test

# 运行特定包的测试
go test -v ./pkg/server/...

# 运行测试并查看覆盖率
make test-cover
```

### 格式化代码

```bash
# 格式化所有代码
make fmt

# 或使用 gofmt
gofmt -w .
```

### 清理构建产物

```bash
# 清理所有构建产物
make clean

# 手动清理
rm -rf dist/ build/ *.exe mcp-toolkit
```

### 更新依赖

```bash
# 下载依赖
go mod download

# 整理依赖
go mod tidy

# 更新所有依赖
go get -u ./...
go mod tidy
```

---

## ❓ 常见问题 / FAQ

### Q1: 如何选择 JSON 库？

**A**: 默认使用 Sonic (最高性能)。编译时使用 `-tags="sonic"` 标签。

```bash
# Sonic (推荐)
go build -tags="sonic" -o mcp-toolkit .

# 标准库
go build -o mcp-toolkit .
```

### Q2: 如何配置沙箱目录？

**A**: 使用 `-sandbox` 参数：

```bash
mcp-toolkit -sandbox /path/to/sandbox
```

### Q3: 如何更改监听端口？

**A**: 使用对应的端口参数：

```bash
# HTTP
mcp-toolkit -transport http -http-port 8080

# SSE
mcp-toolkit -transport sse -sse-port 8081
```

### Q4: 如何启用调试日志？

**A**: 设置环境变量：

```bash
# Linux/macOS
export LOG_LEVEL=debug
./mcp-toolkit

# Windows
set LOG_LEVEL=debug
mcp-toolkit.exe
```

### Q5: 构建失败怎么办？

**A**: 检查以下几点：

1. Go 版本是否 >= 1.25.5
2. 依赖是否已下载 (`go mod download`)
3. 是否使用了正确的编译标签
4. 查看详细错误信息

### Q6: 如何发布新版本？

**A**: 只需创建并推送 Git 标签：

```bash
git tag -a v1.2.0 -m "Release v1.2.0"
git push origin v1.2.0
```

GitHub Actions 会自动完成所有发布流程。

### Q7: 如何跳过 PyPI 发布？

**A**: 不配置 `PYPI_API_TOKEN` Secret，或者在 `.github/workflows/release.yml` 中注释掉 PyPI 发布步骤。

### Q8: 安装后找不到命令？

**A**: 需要配置 PATH 环境变量。详见 [INSTALLATION.md](INSTALLATION.md#-配置-path--configure-path)

---

## 📚 相关文档 / Related Documentation

### 用户文档

- **[INSTALLATION.md](INSTALLATION.md)** - 详细安装指南
- **[README.md](../README.md)** - 项目概览
- **[EXAMPLE.md](EXAMPLE.md)** - 使用示例

### 开发文档

- **[TRANSPORT.md](TRANSPORT.md)** - 传输协议说明
- **[COMMAND_EXECUTION.md](COMMAND_EXECUTION.md)** - 命令执行
- **[RECOVERY.md](RECOVERY.md)** - Panic 恢复机制

### 发布文档

- **[VERSION_MANAGEMENT.md](../VERSION_MANAGEMENT.md)** - 版本管理
- **[RELEASE.md](RELEASE.md)** - 发布说明
- **[RELEASE_CHECKLIST.md](RELEASE_CHECKLIST.md)** - 发布检查清单

### GitHub 相关

- **[GITHUB_ACTIONS.md](GITHUB_ACTIONS.md)** - GitHub Actions 配置

---

## 🆘 需要帮助？/ Need Help?

### 获取帮助

- 📖 查看文档: https://github.com/shibingli/mcp-toolkit/tree/main/docs
- 🐛 报告问题: https://github.com/shibingli/mcp-toolkit/issues
- 💬 讨论: https://github.com/shibingli/mcp-toolkit/discussions

### 贡献

欢迎贡献代码、文档或报告问题！

---

**最后更新**: 2025-12-28

