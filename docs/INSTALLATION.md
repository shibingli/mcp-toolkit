# MCP Toolkit 安装指南 / Installation Guide

本指南提供了 MCP Toolkit 的完整安装说明，包括多种安装方式、配置方法和故障排除。

This guide provides complete installation instructions for MCP Toolkit, including multiple installation methods,
configuration, and troubleshooting.

---

## 📋 目录 / Table of Contents

- [快速安装](#-快速安装--quick-installation)
- [详细安装方式](#-详细安装方式--installation-methods)
- [运行程序](#-运行程序--running-the-program)
- [配置 PATH](#-配置-path--configure-path)
- [验证安装](#-验证安装--verify-installation)
- [更新](#-更新--update)
- [卸载](#-卸载--uninstall)
- [安装脚本详细说明](#-安装脚本详细说明--installation-scripts)
- [故障排除](#-故障排除--troubleshooting)

---

## 🚀 快速安装 / Quick Installation

### 使用 uvx (最简单，无需安装)

```bash
# 直接运行，无需安装
uvx mcp-sandbox-toolkit --help
```

### 使用 uv (推荐)

```bash
# 安装 uv (如果还没有安装)
curl -LsSf https://astral.sh/uv/install.sh | sh

# 安装 MCP Toolkit
uv tool install mcp-sandbox-toolkit

# 运行
mcp-sandbox-toolkit --help
```

### 使用安装脚本

**Linux/macOS:**
```bash
curl -fsSL https://raw.githubusercontent.com/shibingli/mcp-toolkit/main/scripts/install.sh | bash
```

**Windows (PowerShell):**
```powershell
irm https://raw.githubusercontent.com/shibingli/mcp-toolkit/main/scripts/install.ps1 | iex
```

---

## 📦 详细安装方式 / Installation Methods

### 方式 1: 使用 uv (推荐)

`uv` 是一个快速的 Python 包管理器，提供最佳的安装体验。

#### 步骤 1: 安装 uv

```bash
# Linux/macOS
curl -LsSf https://astral.sh/uv/install.sh | sh

# Windows
powershell -c "irm https://astral.sh/uv/install.ps1 | iex"
```

#### 步骤 2: 安装 MCP Toolkit

```bash
# 安装
uv tool install mcp-sandbox-toolkit

# 查看已安装的工具
uv tool list

# 更新到最新版本
uv tool upgrade mcp-sandbox-toolkit

# 卸载
uv tool uninstall mcp-sandbox-toolkit
```

#### 优点

- ✅ 自动管理依赖和版本
- ✅ 隔离的环境，不污染系统
- ✅ 简单的更新和卸载
- ✅ 跨平台支持
- ✅ 速度快

---

### 方式 2: 使用 pip

```bash
# 安装
pip install mcp-sandbox-toolkit

# 更新
pip install --upgrade mcp-sandbox-toolkit

# 卸载
pip uninstall mcp-sandbox-toolkit
```

---

### 方式 3: 使用 pipx

`pipx` 是另一个流行的 Python 工具安装器，提供隔离环境。

```bash
# 安装 pipx (如果还没有安装)
python3 -m pip install --user pipx
python3 -m pipx ensurepath

# 安装 MCP Toolkit
pipx install mcp-sandbox-toolkit

# 更新
pipx upgrade mcp-sandbox-toolkit

# 卸载
pipx uninstall mcp-sandbox-toolkit
```

---

### 方式 4: 使用安装脚本

安装脚本提供了最简单的一键安装体验，支持自动下载、校验和配置。

#### Linux/macOS

```bash
# 方式 A: 直接运行
curl -fsSL https://raw.githubusercontent.com/shibingli/mcp-toolkit/main/scripts/install.sh | bash

# 方式 B: 下载后运行
curl -fsSL https://raw.githubusercontent.com/shibingli/mcp-toolkit/main/scripts/install.sh -o install.sh
chmod +x install.sh
./install.sh install
```

#### Windows (PowerShell)

```powershell
# 下载安装脚本
Invoke-WebRequest -Uri "https://raw.githubusercontent.com/shibingli/mcp-toolkit/main/scripts/install.ps1" -OutFile "install.ps1"

# 运行安装
.\install.ps1 install
```

**安装脚本支持的命令:**

| 命令          | 说明      |
|-------------|---------|
| `install`   | 安装（默认）  |
| `uninstall` | 卸载      |
| `update`    | 更新到最新版本 |
| `version`   | 显示已安装版本 |
| `help`      | 显示帮助信息  |

**高级用法:**

```bash
# Linux/macOS - 安装指定版本
VERSION=v1.0.0 ./install.sh install

# Linux/macOS - 安装到自定义目录
INSTALL_DIR=/usr/local/bin ./install.sh install

# Windows - 安装指定版本
.\install.ps1 install -Version v1.0.0

# Windows - 安装到自定义目录
.\install.ps1 install -InstallDir "C:\Tools\mcp-toolkit"
```

详细的安装脚本说明请参考 [安装脚本详细说明](#-安装脚本详细说明--installation-scripts) 部分。

---

### 方式 5: 手动下载

1. **访问 Releases 页面**

   https://github.com/shibingli/mcp-toolkit/releases

2. **下载对应平台的文件**

    - Windows (amd64): `mcp-toolkit-vX.X.X-windows-amd64.zip`
   - Windows (arm64): `mcp-toolkit-vX.X.X-windows-arm64.zip`
   - Linux (amd64): `mcp-toolkit-vX.X.X-linux-amd64.tar.gz`
   - Linux (arm64): `mcp-toolkit-vX.X.X-linux-arm64.tar.gz`
   - macOS (Intel): `mcp-toolkit-vX.X.X-darwin-amd64.tar.gz`
   - macOS (Apple Silicon): `mcp-toolkit-vX.X.X-darwin-arm64.tar.gz`

3. **解压文件**

   ```bash
   # Linux/macOS
   tar -xzf mcp-toolkit-*.tar.gz

   # Windows (PowerShell)
   Expand-Archive mcp-toolkit-*.zip
   ```

4. **移动到 PATH 目录**

   ```bash
   # Linux/macOS
   sudo mv mcp-toolkit-*/mcp-toolkit /usr/local/bin/

   # Windows: 手动移动到 C:\Program Files\mcp-toolkit\
   # 然后添加到 PATH 环境变量
   ```

---

### 方式 6: 从源码编译

```bash
# 克隆仓库
git clone https://github.com/shibingli/mcp-toolkit.git
cd mcp-toolkit

# 安装依赖
go mod download

# 编译
go build -tags="sonic" -o mcp-toolkit .

# 移动到 PATH 目录
sudo mv mcp-toolkit /usr/local/bin/
```

---

## 🏃 运行程序 / Running the Program

### 方式 1: 直接运行命令 (推荐)

安装后，可执行文件会被添加到以下位置：

- **Linux/macOS**: `~/.local/bin/mcp-sandbox-toolkit` 和 `~/.local/bin/mcp-toolkit`
- **Windows**: `%LOCALAPPDATA%\Programs\mcp-toolkit\mcp-sandbox-toolkit.exe` 和 `mcp-toolkit.exe`

如果 PATH 已正确配置，可以直接运行（两个命令都可以）：

```bash
# 使用完整名称
mcp-sandbox-toolkit --help
mcp-sandbox-toolkit --version

# 或使用短名称
mcp-toolkit --help
mcp-toolkit --version
```

### 方式 2: 使用 uvx (无需安装，推荐)

```bash
# 直接运行，无需安装
uvx mcp-sandbox-toolkit --help
uvx mcp-sandbox-toolkit --version

# 也可以使用短名称
uvx --from mcp-sandbox-toolkit mcp-toolkit --help
```

### 方式 3: 使用 uv tool run

```bash
# 使用 uv tool run
uv tool run mcp-sandbox-toolkit --help
```

---

## 🔧 配置 PATH / Configure PATH

如果安装后无法直接运行 `mcp-toolkit` 命令，需要将安装目录添加到 PATH。

### Linux/macOS

1. **确定你使用的 Shell**:
   ```bash
   echo $SHELL
   ```

2. **编辑配置文件**:

    - **Bash** (`~/.bashrc` 或 `~/.bash_profile`):
      ```bash
      echo 'export PATH="$HOME/.local/bin:$PATH"' >> ~/.bashrc
      source ~/.bashrc
      ```

    - **Zsh** (`~/.zshrc`):
      ```bash
      echo 'export PATH="$HOME/.local/bin:$PATH"' >> ~/.zshrc
      source ~/.zshrc
      ```

    - **Fish** (`~/.config/fish/config.fish`):
      ```bash
      echo 'set -gx PATH $HOME/.local/bin $PATH' >> ~/.config/fish/config.fish
      source ~/.config/fish/config.fish
      ```

3. **验证**:
   ```bash
   echo $PATH | grep ".local/bin"
   which mcp-toolkit
   ```

### Windows

#### 方式 1: PowerShell (临时)

```powershell
$env:Path += ";$env:LOCALAPPDATA\Programs\mcp-toolkit"
```

#### 方式 2: 永久添加到 PATH

1. 打开 **系统属性** → **高级** → **环境变量**
2. 在 **用户变量** 中找到 `Path`
3. 点击 **编辑** → **新建**
4. 添加: `%LOCALAPPDATA%\Programs\mcp-toolkit`
5. 点击 **确定** 保存
6. 重启终端

#### 方式 3: PowerShell 脚本 (永久)

```powershell
[Environment]::SetEnvironmentVariable(
    "Path",
    [Environment]::GetEnvironmentVariable("Path", "User") + ";$env:LOCALAPPDATA\Programs\mcp-toolkit",
    "User"
)
```

---

## ✅ 验证安装 / Verify Installation

```bash
# 检查版本
mcp-toolkit --version

# 查看帮助
mcp-toolkit --help

# 运行服务器
mcp-toolkit
```

---

## 🔄 更新 / Update

### 使用 uv

```bash
uv tool upgrade mcp-sandbox-toolkit
```

### 使用 pip

```bash
pip install --upgrade mcp-sandbox-toolkit
```

### 使用 pipx

```bash
pipx upgrade mcp-sandbox-toolkit
```

### 使用安装脚本

```bash
# Linux/macOS
./install.sh update

# Windows
.\install.ps1 update
```

---

## 🗑️ 卸载 / Uninstall

### 使用 uv

```bash
uv tool uninstall mcp-sandbox-toolkit
```

### 使用 pip

```bash
pip uninstall mcp-sandbox-toolkit
```

### 使用 pipx

```bash
pipx uninstall mcp-sandbox-toolkit
```

### 使用安装脚本

```bash
# Linux/macOS
./install.sh uninstall

# Windows
.\install.ps1 uninstall
```

### 手动卸载

```bash
# Linux/macOS
rm ~/.local/bin/mcp-sandbox-toolkit
rm ~/.local/bin/mcp-toolkit

# Windows
Remove-Item "$env:LOCALAPPDATA\Programs\mcp-toolkit" -Recurse
```

---

## 📜 安装脚本详细说明 / Installation Scripts

### 功能概览

安装脚本支持以下功能：

- ✅ **安装**: 下载并安装最新版本或指定版本
- ✅ **更新**: 检查并更新到最新版本
- ✅ **卸载**: 完全卸载 MCP Toolkit
- ✅ **版本检查**: 显示已安装的版本
- ✅ **校验和验证**: 自动验证下载文件的完整性
- ✅ **智能更新**: 自动检测是否需要更新
- ✅ **备份恢复**: 安装新版本前自动备份旧版本

### Linux/macOS 脚本

#### 环境变量

| 变量                | 说明        | 默认值                     |
|-------------------|-----------|-------------------------|
| `VERSION`         | 要安装的版本    | 最新版本                    |
| `INSTALL_DIR`     | 安装目录      | `$HOME/.local/bin`      |
| `REPO`            | GitHub 仓库 | `shibingli/mcp-toolkit` |
| `VERIFY_CHECKSUM` | 是否验证校验和   | `true`                  |
| `DEBUG`           | 启用调试输出    | `false`                 |

#### 使用示例

```bash
# 安装指定版本
VERSION=v1.0.0 ./install.sh install

# 安装到自定义目录
INSTALL_DIR=/usr/local/bin ./install.sh install

# 跳过校验和验证
VERIFY_CHECKSUM=false ./install.sh install

# 启用调试模式
DEBUG=true ./install.sh install

# 组合使用
VERSION=v1.0.0 INSTALL_DIR=$HOME/bin DEBUG=true ./install.sh install
```

### Windows 脚本

#### 参数

| 参数              | 说明        | 默认值                                   |
|-----------------|-----------|---------------------------------------|
| `-Version`      | 要安装的版本    | 最新版本                                  |
| `-InstallDir`   | 安装目录      | `%LOCALAPPDATA%\Programs\mcp-toolkit` |
| `-Repo`         | GitHub 仓库 | `shibingli/mcp-toolkit`               |
| `-SkipChecksum` | 跳过校验和验证   | `false`                               |
| `-Debug`        | 启用调试输出    | `false`                               |

#### 使用示例

```powershell
# 安装指定版本
.\install.ps1 install -Version v1.0.0

# 安装到自定义目录
.\install.ps1 install -InstallDir "C:\Tools\mcp-toolkit"

# 跳过校验和验证
.\install.ps1 install -SkipChecksum

# 启用调试模式
.\install.ps1 install -Debug

# 组合使用
.\install.ps1 install -Version v1.0.0 -InstallDir "C:\Tools" -Debug
```

### 智能更新检测

脚本会自动检测已安装的版本，如果已经是最新版本，则跳过安装：

```bash
# Linux/macOS
./install.sh install
# Output: Already up to date (version 1.0.0)

# Windows
.\install.ps1 install
# Output: Already up to date (version 1.0.0)
```

### 校验和验证

默认情况下，脚本会下载 `checksums.txt` 文件并验证下载的二进制文件：

```bash
[INFO] Downloading from https://github.com/.../mcp-toolkit-v1.0.0-linux-amd64.tar.gz...
[INFO] Verifying checksum...
[INFO] Checksum verified successfully
```

如果校验和不匹配，安装会失败：

```bash
[ERROR] Checksum verification failed!
Expected: abc123...
Actual: def456...
```

### 备份和恢复

安装新版本时，脚本会自动备份旧版本：

```bash
[INFO] Backing up old version...
[INFO] Installing new version...
[INFO] Installation completed successfully!
```

如果安装失败，可以手动恢复备份：

```bash
# Linux/macOS
cp ~/.local/bin/mcp-toolkit.backup ~/.local/bin/mcp-toolkit

# Windows
copy "%LOCALAPPDATA%\Programs\mcp-toolkit\mcp-toolkit.exe.backup" "%LOCALAPPDATA%\Programs\mcp-toolkit\mcp-toolkit.exe"
```

---

## 🐛 故障排除 / Troubleshooting

### 问题 1: 找不到 `mcp-toolkit` 命令

**症状**:

```bash
mcp-toolkit: command not found
```

**解决方案**:

1. **检查 PATH 是否包含安装目录**:
   ```bash
   # Linux/macOS
   echo $PATH | grep ".local/bin"

   # Windows
   echo $env:Path
   ```

2. **使用 uvx 代替**:
   ```bash
   uvx mcp-sandbox-toolkit --help
   ```

3. **使用完整路径运行**:
   ```bash
   # Linux/macOS
   ~/.local/bin/mcp-toolkit --help

   # Windows
   %LOCALAPPDATA%\Programs\mcp-toolkit\mcp-toolkit.exe --help
   ```

4. **重新加载 shell 配置**:
   ```bash
   # bash
   source ~/.bashrc

   # zsh
   source ~/.zshrc
   ```

5. **或者重启终端**

### 问题 2: 权限错误 (Linux/macOS)

**症状**:

```bash
[ERROR] Permission denied
```

**解决方案**:

```bash
# 方法 1: 给文件添加执行权限
chmod +x ~/.local/bin/mcp-toolkit

# 方法 2: 安装到用户目录（推荐）
INSTALL_DIR=$HOME/bin ./install.sh install

# 方法 3: 使用 sudo（不推荐）
sudo ./install.sh install
```

### 问题 3: Windows 安全警告

**症状**:
Windows Defender 或 SmartScreen 阻止运行

**解决方案**:

1. **方法 1: 解除锁定**
    - 右键点击 `mcp-toolkit.exe`
    - 选择 **属性** → **解除锁定**
    - 点击 **确定**

2. **方法 2: PowerShell 命令**:
   ```powershell
   Unblock-File "$env:LOCALAPPDATA\Programs\mcp-toolkit\mcp-toolkit.exe"
   ```

3. **方法 3: 添加到 Windows Defender 排除列表**

### 问题 4: 下载失败

**症状**:

```
[ERROR] Failed to download. Please check the URL and your internet connection.
```

**解决方案**:

1. **检查网络连接**
2. **检查 GitHub 是否可访问**
3. **尝试使用代理或 VPN**
4. **手动下载并安装**:
    - 访问 https://github.com/shibingli/mcp-toolkit/releases
    - 下载对应平台的文件
    - 手动解压并移动到 PATH 目录

### 问题 5: 校验和验证失败

**症状**:

```
[ERROR] Checksum verification failed!
```

**解决方案**:

1. **重新下载文件**
2. **检查网络是否稳定**
3. **如果确认文件正确，可以跳过验证**:
   ```bash
   # Linux/macOS
   VERIFY_CHECKSUM=false ./install.sh install

   # Windows
   .\install.ps1 install -SkipChecksum
   ```

### 问题 6: Python 版本不兼容

**症状**:

```
ERROR: Package 'mcp-sandbox-toolkit' requires a different Python: 3.7.0 not in '>=3.8'
```

**解决方案**:

1. **升级 Python 到 3.8 或更高版本**
2. **或使用安装脚本直接下载二进制文件**

### 问题 7: uvx 找不到可执行文件

**症状**:

```
An executable named `mcp-sandbox-toolkit` is not provided by package `mcp-sandbox-toolkit`.
```

**解决方案**:

这个问题已在最新版本中修复。如果仍然遇到，请：

1. **更新到最新版本**:
   ```bash
   uv tool upgrade mcp-sandbox-toolkit
   ```

2. **或使用短名称**:
   ```bash
   uvx --from mcp-sandbox-toolkit mcp-toolkit
   ```

---

## 📝 最佳实践 / Best Practices

### 1. 定期更新

建议定期检查更新：

```bash
# 使用 uv
uv tool upgrade mcp-sandbox-toolkit

# 使用安装脚本
./install.sh update  # Linux/macOS
.\install.ps1 update  # Windows
```

### 2. 使用版本锁定

在生产环境中，建议锁定特定版本：

```bash
# Linux/macOS
VERSION=v1.0.0 ./install.sh install

# Windows
.\install.ps1 install -Version v1.0.0
```

### 3. 验证安装

安装后验证：

```bash
mcp-toolkit --version
mcp-toolkit --help
```

### 4. 使用 uvx 进行测试

在安装前，可以使用 uvx 测试：

```bash
uvx mcp-sandbox-toolkit --help
```

---

## 🔗 相关链接 / Related Links

- **GitHub 仓库**: https://github.com/shibingli/mcp-toolkit
- **PyPI 包**: https://pypi.org/project/mcp-sandbox-toolkit/
- **GitHub Releases**: https://github.com/shibingli/mcp-toolkit/releases
- **问题反馈**: https://github.com/shibingli/mcp-toolkit/issues
- **版本管理说明**: [VERSION_MANAGEMENT.md](../VERSION_MANAGEMENT.md)

---

**最后更新**: 2025-12-28



