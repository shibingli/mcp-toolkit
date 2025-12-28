# MCP Toolkit 安装和使用指南

## 📦 安装方式

### 方式 1: 使用 uv (推荐)

```bash
# 安装 uv (如果还没有安装)
curl -LsSf https://astral.sh/uv/install.sh | sh

# 安装 MCP Toolkit
uv tool install mcp-sandbox-toolkit
```

### 方式 2: 使用 pip

```bash
pip install mcp-sandbox-toolkit
```

### 方式 3: 使用 pipx

```bash
# 安装 pipx (如果还没有安装)
python -m pip install --user pipx
python -m pipx ensurepath

# 安装 MCP Toolkit
pipx install mcp-sandbox-toolkit
```

### 方式 4: 使用安装脚本

**Linux/macOS:**

```bash
curl -fsSL https://raw.githubusercontent.com/shibingli/mcp-toolkit/main/scripts/install.sh | bash
```

**Windows (PowerShell):**

```powershell
irm https://raw.githubusercontent.com/shibingli/mcp-toolkit/main/scripts/install.ps1 | iex
```

---

## 🚀 运行程序

### 方式 1: 直接运行命令 (推荐)

安装后，可执行文件会被添加到以下位置：

- **Linux/macOS**: `~/.local/bin/mcp-toolkit`
- **Windows**: `%LOCALAPPDATA%\Programs\mcp-toolkit\mcp-toolkit.exe`

如果 PATH 已正确配置，可以直接运行：

```bash
mcp-toolkit --help
mcp-toolkit --version
```

### 方式 2: 使用 uvx (无需安装)

```bash
# 直接运行，无需安装
uvx mcp-sandbox-toolkit --help
uvx mcp-sandbox-toolkit --version
```

### 方式 3: 使用 uv tool run

```bash
# 使用 uv tool run
uv tool run mcp-sandbox-toolkit --help
```

---

## 🔧 配置 PATH 环境变量

如果安装后无法直接运行 `mcp-toolkit` 命令，需要将安装目录添加到 PATH：

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

## ✅ 验证安装

```bash
# 检查版本
mcp-toolkit --version

# 查看帮助
mcp-toolkit --help

# 运行服务器
mcp-toolkit
```

---

## 🔄 更新

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

---

## 🗑️ 卸载

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

---

## 🐛 故障排除

### 问题 1: 找不到 `mcp-toolkit` 命令

**解决方案**:

1. 检查 PATH 是否包含安装目录
2. 使用 `uvx mcp-sandbox-toolkit` 代替
3. 使用完整路径运行

### 问题 2: 权限错误 (Linux/macOS)

**解决方案**:

```bash
chmod +x ~/.local/bin/mcp-toolkit
```

### 问题 3: Windows 安全警告

**解决方案**:

1. 右键点击 `mcp-toolkit.exe`
2. 选择 **属性** → **解除锁定**
3. 或在 PowerShell 中运行:
   ```powershell
   Unblock-File "$env:LOCALAPPDATA\Programs\mcp-toolkit\mcp-toolkit.exe"
   ```

---

## 📚 更多信息

- **GitHub**: https://github.com/shibingli/mcp-toolkit
- **PyPI**: https://pypi.org/project/mcp-sandbox-toolkit/
- **文档**: https://github.com/shibingli/mcp-toolkit/blob/main/README.md

