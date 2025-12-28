# 安装脚本使用指南

MCP Toolkit 提供了两个安装脚本，分别用于 Linux/macOS 和 Windows 平台。

## 📋 功能概览

两个脚本都支持以下功能：
- ✅ **安装**: 下载并安装最新版本或指定版本
- ✅ **更新**: 检查并更新到最新版本
- ✅ **卸载**: 完全卸载 MCP Toolkit
- ✅ **版本检查**: 显示已安装的版本
- ✅ **校验和验证**: 自动验证下载文件的完整性
- ✅ **智能更新**: 自动检测是否需要更新
- ✅ **备份恢复**: 安装新版本前自动备份旧版本

---

## 🐧 Linux/macOS 安装脚本

### 基本用法

#### 安装最新版本
```bash
curl -fsSL https://raw.githubusercontent.com/shibingli/mcp-toolkit/main/scripts/install.sh | bash
```

#### 下载脚本后使用
```bash
# 下载脚本
curl -fsSL https://raw.githubusercontent.com/shibingli/mcp-toolkit/main/scripts/install.sh -o install.sh
chmod +x install.sh

# 安装
./install.sh install

# 更新
./install.sh update

# 卸载
./install.sh uninstall

# 查看版本
./install.sh version

# 查看帮助
./install.sh help
```

### 高级用法

#### 安装指定版本
```bash
VERSION=v1.0.0 ./install.sh install
```

#### 安装到自定义目录
```bash
INSTALL_DIR=/usr/local/bin ./install.sh install
```

#### 使用自定义仓库
```bash
REPO=your-username/mcp-toolkit ./install.sh install
```

#### 跳过校验和验证
```bash
VERIFY_CHECKSUM=false ./install.sh install
```

#### 启用调试模式
```bash
DEBUG=true ./install.sh install
```

#### 组合使用
```bash
VERSION=v1.0.0 INSTALL_DIR=$HOME/bin DEBUG=true ./install.sh install
```

### 环境变量

| 变量 | 说明 | 默认值 |
|------|------|--------|
| `VERSION` | 要安装的版本 | 最新版本 |
| `INSTALL_DIR` | 安装目录 | `$HOME/.local/bin` |
| `REPO` | GitHub 仓库 | `shibingli/mcp-toolkit` |
| `VERIFY_CHECKSUM` | 是否验证校验和 | `true` |
| `DEBUG` | 启用调试输出 | `false` |

### 命令

| 命令 | 说明 |
|------|------|
| `install` | 安装（默认） |
| `uninstall` | 卸载 |
| `update` | 更新到最新版本 |
| `version` | 显示已安装版本 |
| `help` | 显示帮助信息 |

---

## 🪟 Windows 安装脚本

### 基本用法

#### 安装最新版本
```powershell
# 下载并运行
Invoke-WebRequest -Uri "https://raw.githubusercontent.com/shibingli/mcp-toolkit/main/scripts/install.ps1" -OutFile "install.ps1"
.\install.ps1
```

#### 一行命令安装
```powershell
iwr -useb https://raw.githubusercontent.com/shibingli/mcp-toolkit/main/scripts/install.ps1 | iex
```

#### 使用下载的脚本
```powershell
# 安装
.\install.ps1 install

# 更新
.\install.ps1 update

# 卸载
.\install.ps1 uninstall

# 查看版本
.\install.ps1 version

# 查看帮助
.\install.ps1 help
```

### 高级用法

#### 安装指定版本
```powershell
.\install.ps1 install -Version v1.0.0
```

#### 安装到自定义目录
```powershell
.\install.ps1 install -InstallDir "C:\Tools\mcp-toolkit"
```

#### 使用自定义仓库
```powershell
.\install.ps1 install -Repo "your-username/mcp-toolkit"
```

#### 跳过校验和验证
```powershell
.\install.ps1 install -SkipChecksum
```

#### 启用调试模式
```powershell
.\install.ps1 install -Debug
```

#### 组合使用
```powershell
.\install.ps1 install -Version v1.0.0 -InstallDir "C:\Tools" -Debug
```

### 参数

| 参数 | 说明 | 默认值 |
|------|------|--------|
| `-Version` | 要安装的版本 | 最新版本 |
| `-InstallDir` | 安装目录 | `%LOCALAPPDATA%\Programs\mcp-toolkit` |
| `-Repo` | GitHub 仓库 | `shibingli/mcp-toolkit` |
| `-SkipChecksum` | 跳过校验和验证 | `false` |
| `-Debug` | 启用调试输出 | `false` |

### 命令

| 命令 | 说明 |
|------|------|
| `install` | 安装（默认） |
| `uninstall` | 卸载 |
| `update` | 更新到最新版本 |
| `version` | 显示已安装版本 |
| `help` | 显示帮助信息 |

---

## 🔧 功能详解

### 1. 智能更新检测

脚本会自动检测已安装的版本，如果已经是最新版本，则跳过安装：

```bash
# Linux/macOS
./install.sh install
# Output: Already up to date (version 1.0.0)

# Windows
.\install.ps1 install
# Output: Already up to date (version 1.0.0)
```

### 2. 校验和验证

默认情况下，脚本会下载 `checksums.txt` 文件并验证下载的二进制文件：

```bash
# 验证过程
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

### 3. 备份和恢复

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

### 4. PATH 管理

#### Linux/macOS
如果安装目录不在 PATH 中，脚本会提示添加：
```bash
[WARN] /home/user/.local/bin is not in your PATH
[WARN] Add the following line to your shell configuration file:
[WARN]   For bash: ~/.bashrc
[WARN]   For zsh: ~/.zshrc

    export PATH="$PATH:/home/user/.local/bin"
```

#### Windows
脚本会询问是否自动添加到 PATH：
```powershell
[WARN] C:\Users\user\AppData\Local\Programs\mcp-toolkit is not in your PATH

Do you want to add it to your PATH? (Y/N)
```

---

## 🐛 故障排查

### 问题 1: 下载失败

**错误信息**:
```
[ERROR] Failed to download. Please check the URL and your internet connection.
```

**解决方法**:
1. 检查网络连接
2. 检查 GitHub 是否可访问
3. 尝试使用代理或 VPN
4. 手动下载并安装

### 问题 2: 校验和验证失败

**错误信息**:
```
[ERROR] Checksum verification failed!
```

**解决方法**:
1. 重新下载文件
2. 检查网络是否稳定
3. 如果确认文件正确，可以跳过验证：
   ```bash
   # Linux/macOS
   VERIFY_CHECKSUM=false ./install.sh install
   
   # Windows
   .\install.ps1 install -SkipChecksum
   ```

### 问题 3: 权限错误

**Linux/macOS 错误**:
```
[ERROR] Permission denied
```

**解决方法**:
```bash
# 方法 1: 使用 sudo（不推荐）
sudo ./install.sh install

# 方法 2: 安装到用户目录（推荐）
INSTALL_DIR=$HOME/bin ./install.sh install
```

**Windows 错误**:
```
Access denied
```

**解决方法**:
以管理员身份运行 PowerShell

### 问题 4: 找不到命令

**错误**:
```bash
mcp-toolkit: command not found
```

**解决方法**:
1. 检查 PATH 设置
2. 重新加载 shell 配置：
   ```bash
   # bash
   source ~/.bashrc
   
   # zsh
   source ~/.zshrc
   ```
3. 或者重启终端

---

## 📝 最佳实践

### 1. 定期更新

建议定期检查更新：
```bash
# Linux/macOS
./install.sh update

# Windows
.\install.ps1 update
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

### 4. 备份配置

在更新前备份配置文件（如果有）：
```bash
# 备份配置
cp ~/.config/mcp-toolkit/config.yaml ~/.config/mcp-toolkit/config.yaml.backup

# 更新
./install.sh update

# 如果需要，恢复配置
cp ~/.config/mcp-toolkit/config.yaml.backup ~/.config/mcp-toolkit/config.yaml
```

---

## 🔗 相关链接

- [快速开始指南](../GET_STARTED.md)
- [完整安装指南](INSTALLATION.md)
- [GitHub Releases](https://github.com/shibingli/mcp-toolkit/releases)

