# 前置依赖 / Prerequisites

在开始构建和发布 MCP Toolkit 之前，请确保安装了以下工具。

Before building and releasing MCP Toolkit, make sure you have the following tools installed.

---

## 必需工具 / Required Tools

### 1. Go 语言环境 / Go Language Environment

**版本要求 / Version Required**: Go 1.25.5 或更高 / Go 1.25.5 or higher

#### 安装 / Installation

**Windows:**
```powershell
# 下载并安装 / Download and install
# 访问 https://go.dev/dl/
# 下载 go1.25.5.windows-amd64.msi 并安装

# 验证安装 / Verify installation
go version
```

**Linux:**
```bash
# 使用官方安装脚本 / Using official installation script
wget https://go.dev/dl/go1.25.5.linux-amd64.tar.gz
sudo rm -rf /usr/local/go
sudo tar -C /usr/local -xzf go1.25.5.linux-amd64.tar.gz

# 添加到 PATH / Add to PATH
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
source ~/.bashrc

# 验证安装 / Verify installation
go version
```

**macOS:**
```bash
# 使用 Homebrew / Using Homebrew
brew install go@1.25

# 或下载安装包 / Or download installer
# 访问 https://go.dev/dl/
# 下载 go1.25.5.darwin-amd64.pkg 或 go1.25.5.darwin-arm64.pkg

# 验证安装 / Verify installation
go version
```

### 2. Git

**版本要求 / Version Required**: Git 2.0 或更高 / Git 2.0 or higher

#### 安装 / Installation

**Windows:**
```powershell
# 下载并安装 Git for Windows
# https://git-scm.com/download/win

# 验证安装 / Verify installation
git --version
```

**Linux:**
```bash
# Ubuntu/Debian
sudo apt-get update
sudo apt-get install git

# CentOS/RHEL
sudo yum install git

# 验证安装 / Verify installation
git --version
```

**macOS:**
```bash
# 使用 Homebrew
brew install git

# 或使用 Xcode Command Line Tools
xcode-select --install

# 验证安装 / Verify installation
git --version
```

---

## 可选工具 / Optional Tools

### 3. Python (用于发布到 PyPI) / Python (for PyPI publishing)

**版本要求 / Version Required**: Python 3.8 或更高 / Python 3.8 or higher

#### 安装 / Installation

**Windows:**
```powershell
# 下载并安装 Python
# https://www.python.org/downloads/

# 验证安装 / Verify installation
python --version
```

**Linux:**
```bash
# Ubuntu/Debian
sudo apt-get update
sudo apt-get install python3 python3-pip

# CentOS/RHEL
sudo yum install python3 python3-pip

# 验证安装 / Verify installation
python3 --version
```

**macOS:**
```bash
# 使用 Homebrew
brew install python@3.12

# 验证安装 / Verify installation
python3 --version
```

#### 安装 Python 构建工具 / Install Python Build Tools

```bash
# 安装构建工具 / Install build tools
python -m pip install --upgrade pip
python -m pip install --upgrade build twine hatchling

# 验证安装 / Verify installation
python -m build --version
python -m twine --version
```

### 4. Make (可选，用于使用 Makefile) / Make (optional, for using Makefile)

**Windows:**
```powershell
# 使用 Chocolatey
choco install make

# 或使用 Scoop
scoop install make

# 或使用 Git Bash (已包含 make)
```

**Linux:**
```bash
# Ubuntu/Debian
sudo apt-get install build-essential

# CentOS/RHEL
sudo yum groupinstall "Development Tools"
```

**macOS:**
```bash
# 使用 Xcode Command Line Tools (已包含 make)
xcode-select --install

# 验证安装 / Verify installation
make --version
```

---

## 验证环境 / Verify Environment

运行以下命令验证所有工具是否正确安装：

Run the following commands to verify all tools are correctly installed:

```bash
# 检查 Go / Check Go
go version
# 期望输出 / Expected output: go version go1.25.5 ...

# 检查 Git / Check Git
git --version
# 期望输出 / Expected output: git version 2.x.x

# 检查 Python (可选) / Check Python (optional)
python --version
# 期望输出 / Expected output: Python 3.x.x

# 检查 Make (可选) / Check Make (optional)
make --version
# 期望输出 / Expected output: GNU Make x.x
```

---

## 项目依赖 / Project Dependencies

### Go 模块依赖 / Go Module Dependencies

项目的 Go 依赖会在首次构建时自动下载：

Project Go dependencies will be automatically downloaded on first build:

```bash
# 下载依赖 / Download dependencies
go mod download

# 验证依赖 / Verify dependencies
go mod verify
```

### Python 依赖 (仅用于发布) / Python Dependencies (for publishing only)

如果要发布到 PyPI，需要安装以下 Python 包：

If publishing to PyPI, install the following Python packages:

```bash
pip install build twine hatchling
```

---

## 环境变量配置 / Environment Variables Configuration

### Go 环境变量 / Go Environment Variables

```bash
# Linux/macOS
export GOPATH=$HOME/go
export PATH=$PATH:/usr/local/go/bin:$GOPATH/bin

# Windows (PowerShell)
$env:GOPATH = "$env:USERPROFILE\go"
$env:Path += ";C:\Go\bin;$env:GOPATH\bin"
```

### 构建标签 / Build Tags

对于 Linux 系统，需要设置额外的环境变量：

For Linux systems, additional environment variables are needed:

```bash
# Linux only
export GOEXPERIMENT="jsonv2,greenteagc"
```

---

## 常见问题 / Common Issues

### Go 版本过低 / Go Version Too Old

```bash
# 错误 / Error: go version go1.20.x is too old
# 解决方案 / Solution: 升级到 Go 1.25.5+
```

### Python build 模块未找到 / Python build module not found

```bash
# 错误 / Error: No module named build
# 解决方案 / Solution:
python -m pip install --upgrade build twine hatchling
```

### Make 命令未找到 / Make command not found

```bash
# 错误 / Error: make: command not found
# 解决方案 / Solution:
# Windows: 使用 Git Bash 或安装 make
# Linux: sudo apt-get install build-essential
# macOS: xcode-select --install
```

---

## 下一步 / Next Steps

环境配置完成后，可以继续：

After environment setup, you can proceed to:

1. **本地构建测试** / Local build test
   ```bash
   make build
   ```

2. **运行测试** / Run tests
   ```bash
   make test
   ```

3. **查看发布指南** / View release guide
   - [快速发布指南](QUICK_PUBLISH.md)
   - [详细发布指南](PUBLISH_GUIDE.md)

---

## 参考资源 / References

- [Go 官方网站](https://go.dev/)
- [Git 官方网站](https://git-scm.com/)
- [Python 官方网站](https://www.python.org/)
- [GNU Make 文档](https://www.gnu.org/software/make/)

