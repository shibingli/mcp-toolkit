# 🚀 快速开始 - 发布配置 / Quick Start - Release Setup

欢迎！这个文档将帮助你在 **10 分钟内** 完成 MCP Toolkit 的发布配置。

Welcome! This document will help you complete the MCP Toolkit release setup in **10 minutes**.

---

## ✅ 第一步：检查环境 / Step 1: Check Environment

### 必需工具 / Required Tools

运行以下命令检查：

```bash
# 检查 Go (必需)
go version
# 期望: go version go1.25.5 或更高

# 检查 Git (必需)
git --version
# 期望: git version 2.x.x

# 检查 Python (可选，用于 PyPI 发布)
python --version
# 期望: Python 3.8 或更高
```

### 如果缺少工具 / If Tools Are Missing

查看详细安装指南：[docs/PREREQUISITES.md](docs/PREREQUISITES.md)

**快速安装 Python 构建工具：**

```bash
python -m pip install --upgrade pip
python -m pip install --upgrade build twine hatchling
```

---

## 📝 第二步：更新仓库信息 / Step 2: Update Repository Info

### 方式 A: 批量替换（推荐）/ Method A: Batch Replace (Recommended)

**Linux/macOS:**
```bash
# 将 your-username 替换为你的 GitHub 用户名
find . -type f \( -name "*.sh" -o -name "*.ps1" -o -name "*.py" -o -name "*.toml" -o -name "*.md" \) \
  -exec sed -i 's/your-username/YOUR_GITHUB_USERNAME/g' {} +
```

**Windows (PowerShell):**
```powershell
# 将 your-username 替换为你的 GitHub 用户名
Get-ChildItem -Recurse -Include *.sh,*.ps1,*.py,*.toml,*.md | 
  ForEach-Object {
    (Get-Content $_.FullName) -replace 'your-username', 'YOUR_GITHUB_USERNAME' | 
    Set-Content $_.FullName
  }
```

### 方式 B: 手动编辑 / Method B: Manual Edit

编辑以下关键文件：

1. `scripts/install.sh` (第 8 行)
2. `scripts/install.ps1` (第 10 行)
3. `python/mcp_toolkit_wrapper/installer.py` (第 17 行)
4. `pyproject.toml` (URLs 部分)

---

## 🔧 第三步：推送到 GitHub / Step 3: Push to GitHub

```bash
# 1. 创建 GitHub 仓库
# 访问 https://github.com/new
# 创建名为 mcp-toolkit 的仓库

# 2. 关联远程仓库
git remote add origin https://github.com/YOUR_USERNAME/mcp-toolkit.git

# 3. 推送代码
git add .
git commit -m "Add release configuration"
git push -u origin main
```

---

## ⚙️ 第四步：配置 GitHub Actions / Step 4: Configure GitHub Actions

1. 访问仓库设置页面：
   ```
   https://github.com/YOUR_USERNAME/mcp-toolkit/settings/actions
   ```

2. 在 "Workflow permissions" 部分：
   - ✅ 选择 **"Read and write permissions"**
   - ✅ 勾选 **"Allow GitHub Actions to create and approve pull requests"**

3. 点击 "Save" 保存

---

## 🎯 第五步：测试构建 / Step 5: Test Build

### 本地测试 / Local Test

```bash
# 运行测试
go test -v ./...

# 本地构建
go build -tags="sonic" -o mcp-toolkit .

# 或使用 Makefile
make test
make build
```

### 测试 Python 包构建（可选）/ Test Python Package Build (Optional)

```bash
# 构建 Python 包
python -m build

# 检查包
python -m twine check dist/*

# 应该看到: PASSED
```

---

## 🚀 第六步：创建第一个发布 / Step 6: Create First Release

### 方式 1: 仅发布到 GitHub Releases（推荐开始）

```bash
# 1. 创建版本标签
git tag -a v1.0.0 -m "Release version 1.0.0"

# 2. 推送标签（这会触发自动构建和发布）
git push origin v1.0.0

# 3. 等待 2-5 分钟，访问 Releases 页面
# https://github.com/YOUR_USERNAME/mcp-toolkit/releases
```

### 方式 2: 同时发布到 GitHub + PyPI

**额外步骤：配置 PyPI**

1. 注册 PyPI 账号：https://pypi.org/account/register/
2. 创建 API Token：https://pypi.org/manage/account/token/
3. 添加 GitHub Secret：
   - 访问：`https://github.com/YOUR_USERNAME/mcp-toolkit/settings/secrets/actions`
   - 点击 "New repository secret"
   - Name: `PYPI_API_TOKEN`
   - Value: 粘贴你的 PyPI token
   - 点击 "Add secret"

**然后创建发布：**

```bash
git tag -a v1.0.0 -m "Release version 1.0.0"
git push origin v1.0.0
```

---

## 🎉 完成！/ Done!

### 验证发布 / Verify Release

1. **GitHub Release**
   - 访问：`https://github.com/YOUR_USERNAME/mcp-toolkit/releases`
   - 应该看到 v1.0.0 版本和所有平台的二进制文件

2. **PyPI**（如果发布）
   - 访问：`https://pypi.org/project/mcp-toolkit/`
   - 应该看到 1.0.0 版本

### 用户安装方式 / User Installation

**方式 1: 使用 uv（如果发布到 PyPI）**
```bash
uv tool install mcp-toolkit
```

**方式 2: 使用安装脚本**
```bash
# Linux/macOS
curl -fsSL https://raw.githubusercontent.com/YOUR_USERNAME/mcp-toolkit/main/scripts/install.sh | bash

# Windows
Invoke-WebRequest -Uri "https://raw.githubusercontent.com/YOUR_USERNAME/mcp-toolkit/main/scripts/install.ps1" -OutFile "install.ps1"
.\install.ps1
```

**方式 3: 手动下载**
- 访问 Releases 页面下载对应平台的二进制文件

---

## 📚 下一步 / Next Steps

- 📖 阅读 [发布检查清单](docs/RELEASE_CHECKLIST.md)
- 📖 查看 [详细发布指南](docs/PUBLISH_GUIDE.md)
- 📖 了解 [用户安装方式](docs/INSTALLATION.md)

---

## ❓ 遇到问题？/ Having Issues?

### 常见问题 / Common Issues

1. **GitHub Actions 失败**
   - 检查 Actions 权限设置
   - 查看 Actions 日志：`https://github.com/YOUR_USERNAME/mcp-toolkit/actions`

2. **Python 构建失败**
   ```bash
   # 重新安装构建工具
   python -m pip install --upgrade build twine hatchling
   ```

3. **找不到 Go 命令**
   - 确保 Go 1.25.5+ 已安装
   - 查看 [前置依赖文档](docs/PREREQUISITES.md)

### 获取帮助 / Get Help

- 📖 查看 [前置依赖文档](docs/PREREQUISITES.md)
- 📖 查看 [故障排查](docs/RELEASE.md#故障排查--troubleshooting)
- 🐛 提交 Issue：`https://github.com/YOUR_USERNAME/mcp-toolkit/issues`

---

## 📊 发布流程总结 / Release Process Summary

```
1. 检查环境 ✅
   ↓
2. 更新仓库信息 ✅
   ↓
3. 推送到 GitHub ✅
   ↓
4. 配置 Actions 权限 ✅
   ↓
5. 测试构建 ✅
   ↓
6. 创建标签并推送 ✅
   ↓
7. 自动构建和发布 🚀
   ↓
8. 用户可以安装 🎉
```

---

**🎊 恭喜！你已经完成了发布配置！**

**🎊 Congratulations! You've completed the release setup!**

现在只需创建标签并推送，GitHub Actions 会自动完成剩余工作！

Now just create a tag and push, GitHub Actions will automatically do the rest!

