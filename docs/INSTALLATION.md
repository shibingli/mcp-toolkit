# 安装指南 / Installation Guide

MCP Toolkit 提供多种安装方式，选择最适合你的方式。

MCP Toolkit provides multiple installation methods. Choose the one that suits you best.

---

## 快速安装 / Quick Installation

### 使用 uv (推荐) / Using uv (Recommended)

```bash
# 安装 uv (如果还没有安装) / Install uv (if not already installed)
curl -LsSf https://astral.sh/uv/install.sh | sh

# 使用 uv 安装 MCP Toolkit / Install MCP Toolkit using uv
uv tool install mcp-toolkit

# 或者使用 pipx / Or use pipx
pipx install mcp-toolkit
```

### 使用安装脚本 / Using Installation Script

#### Linux/macOS

```bash
curl -fsSL https://raw.githubusercontent.com/shibingli/mcp-toolkit/main/scripts/install.sh | bash
```

#### Windows (PowerShell)

```powershell
Invoke-WebRequest -Uri "https://raw.githubusercontent.com/shibingli/mcp-toolkit/main/scripts/install.ps1" -OutFile "install.ps1"
.\install.ps1
```

---

## 详细安装方式 / Detailed Installation Methods

### 方式 1: 使用 uv 工具 (推荐) / Method 1: Using uv Tool (Recommended)

`uv` 是一个快速的 Python 包管理器，可以轻松安装和管理工具。

`uv` is a fast Python package manager that makes it easy to install and manage tools.

#### 步骤 1: 安装 uv / Step 1: Install uv

```bash
# Linux/macOS
curl -LsSf https://astral.sh/uv/install.sh | sh

# Windows
powershell -c "irm https://astral.sh/uv/install.ps1 | iex"
```

#### 步骤 2: 安装 MCP Toolkit / Step 2: Install MCP Toolkit

```bash
# 使用 uv tool install
uv tool install mcp-toolkit

# 查看已安装的工具
uv tool list

# 更新到最新版本
uv tool upgrade mcp-toolkit

# 卸载
uv tool uninstall mcp-toolkit
```

#### 优点 / Advantages

- ✅ 自动管理依赖和版本 / Automatic dependency and version management
- ✅ 隔离的环境，不污染系统 / Isolated environment, no system pollution
- ✅ 简单的更新和卸载 / Simple update and uninstall
- ✅ 跨平台支持 / Cross-platform support

---

### 方式 2: 使用 pipx / Method 2: Using pipx

`pipx` 是另一个流行的 Python 工具安装器。

`pipx` is another popular Python tool installer.

```bash
# 安装 pipx (如果还没有安装) / Install pipx (if not already installed)
python3 -m pip install --user pipx
python3 -m pipx ensurepath

# 安装 MCP Toolkit
pipx install mcp-toolkit

# 更新
pipx upgrade mcp-toolkit

# 卸载
pipx uninstall mcp-toolkit
```

---

### 方式 3: 使用安装脚本 / Method 3: Using Installation Script

#### Linux/macOS

```bash
# 方式 A: 直接运行 / Method A: Direct run
curl -fsSL https://raw.githubusercontent.com/shibingli/mcp-toolkit/main/scripts/install.sh | bash

# 方式 B: 下载后运行 / Method B: Download and run
wget https://raw.githubusercontent.com/shibingli/mcp-toolkit/main/scripts/install.sh
chmod +x install.sh
./install.sh

# 指定版本 / Specify version
VERSION=v1.0.0 ./install.sh

# 指定安装目录 / Specify installation directory
INSTALL_DIR=/usr/local/bin ./install.sh

# 卸载 / Uninstall
./install.sh uninstall
```

#### Windows (PowerShell)

```powershell
# 下载安装脚本 / Download installation script
Invoke-WebRequest -Uri "https://raw.githubusercontent.com/shibingli/mcp-toolkit/main/scripts/install.ps1" -OutFile "install.ps1"

# 运行安装 / Run installation
.\install.ps1

# 指定版本 / Specify version
.\install.ps1 -Version "v1.0.0"

# 指定安装目录 / Specify installation directory
.\install.ps1 -InstallDir "C:\Program Files\mcp-toolkit"

# 卸载 / Uninstall
.\install.ps1 uninstall
```

---

### 方式 4: 手动下载 / Method 4: Manual Download

1. **访问 Releases 页面 / Visit Releases page**
   
   https://github.com/shibingli/mcp-toolkit/releases

2. **下载对应平台的文件 / Download file for your platform**
   
   - Windows (amd64): `mcp-toolkit-vX.X.X-windows-amd64.zip`
   - Windows (arm64): `mcp-toolkit-vX.X.X-windows-arm64.zip`
   - Linux (amd64): `mcp-toolkit-vX.X.X-linux-amd64.tar.gz`
   - Linux (arm64): `mcp-toolkit-vX.X.X-linux-arm64.tar.gz`
   - macOS (Intel): `mcp-toolkit-vX.X.X-darwin-amd64.tar.gz`
   - macOS (Apple Silicon): `mcp-toolkit-vX.X.X-darwin-arm64.tar.gz`

3. **解压文件 / Extract file**
   
   ```bash
   # Linux/macOS
   tar -xzf mcp-toolkit-*.tar.gz
   
   # Windows (PowerShell)
   Expand-Archive mcp-toolkit-*.zip
   ```

4. **移动到 PATH 目录 / Move to PATH directory**
   
   ```bash
   # Linux/macOS
   sudo mv mcp-toolkit-*/mcp-toolkit /usr/local/bin/
   
   # Windows: 手动移动到 C:\Program Files\mcp-toolkit\
   # 然后添加到 PATH 环境变量
   ```

---

### 方式 5: 从源码编译 / Method 5: Build from Source

```bash
# 克隆仓库 / Clone repository
git clone https://github.com/shibingli/mcp-toolkit.git
cd mcp-toolkit

# 安装依赖 / Install dependencies
go mod download

# 编译 / Build
go build -tags="sonic" -o mcp-toolkit .

# 移动到 PATH 目录 / Move to PATH directory
sudo mv mcp-toolkit /usr/local/bin/
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

## 配置 PATH / Configure PATH

### Linux/macOS

添加到 `~/.bashrc` 或 `~/.zshrc`:

```bash
export PATH="$PATH:$HOME/.local/bin"
```

然后重新加载配置:

```bash
source ~/.bashrc  # 或 source ~/.zshrc
```

### Windows

1. 打开"系统属性" → "高级" → "环境变量"
2. 在"用户变量"中找到"Path"
3. 添加安装目录，例如: `C:\Users\YourName\AppData\Local\Programs\mcp-toolkit`

或使用 PowerShell:

```powershell
$env:Path += ";$env:LOCALAPPDATA\Programs\mcp-toolkit"
```

---

## 更新 / Update

### 使用 uv/pipx

```bash
# uv
uv tool upgrade mcp-toolkit

# pipx
pipx upgrade mcp-toolkit
```

### 使用安装脚本

重新运行安装脚本即可:

```bash
./install.sh
```

---

## 卸载 / Uninstall

### 使用 uv/pipx

```bash
# uv
uv tool uninstall mcp-toolkit

# pipx
pipx uninstall mcp-toolkit
```

### 使用安装脚本

```bash
./install.sh uninstall
```

### 手动卸载

```bash
# Linux/macOS
rm ~/.local/bin/mcp-toolkit

# Windows
Remove-Item "$env:LOCALAPPDATA\Programs\mcp-toolkit\mcp-toolkit.exe"
```

---

## 故障排查 / Troubleshooting

参见 [RELEASE.md](RELEASE.md#故障排查--troubleshooting)

