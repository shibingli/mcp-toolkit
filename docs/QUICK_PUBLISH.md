# 快速发布指南 / Quick Publish Guide

这是一个快速参考指南，帮助你在 5 分钟内完成项目的发布配置。

This is a quick reference guide to help you complete the project release configuration in 5 minutes.

---

## 📋 发布前准备 / Pre-Release Preparation

### 0. 检查前置依赖 (必须) / Check Prerequisites (Required)

确保已安装以下工具：

Make sure the following tools are installed:

- ✅ Go 1.25.5+ (`go version`)
- ✅ Git 2.0+ (`git --version`)
- ✅ Python 3.8+ (可选，用于 PyPI 发布) (`python --version`)
- ✅ Make (可选，用于使用 Makefile) (`make --version`)

**如果缺少工具，请先查看 [前置依赖文档](PREREQUISITES.md) 进行安装。**

**If any tools are missing, see [Prerequisites](PREREQUISITES.md) for installation.**

#### 安装 Python 构建工具 (如果要发布到 PyPI)

Install Python build tools (if publishing to PyPI):

```bash
python -m pip install --upgrade pip
python -m pip install --upgrade build twine hatchling
```

### 1. 更新仓库信息 (必须) / Update Repository Info (Required)

将以下文件中的 `shibingli/mcp-toolkit` 替换为你的实际仓库地址：

Replace `shibingli/mcp-toolkit` with your actual repository in the following files:

```bash
# 使用 sed 批量替换 (Linux/macOS)
find . -type f \( -name "*.sh" -o -name "*.ps1" -o -name "*.py" -o -name "*.toml" -o -name "*.md" \) \
  -exec sed -i 's/shibingli\/mcp-toolkit/YOUR_USERNAME\/mcp-toolkit/g' {} +

# 或手动编辑以下文件:
# - scripts/install.sh (第 8 行)
# - scripts/install.ps1 (第 10 行)
# - python/mcp_toolkit_wrapper/installer.py (第 17 行)
# - pyproject.toml (URLs 部分)
# - 所有 docs/*.md 文件中的链接
```

### 2. 推送到 GitHub (必须) / Push to GitHub (Required)

```bash
# 创建 GitHub 仓库后
git remote add origin https://github.com/YOUR_USERNAME/mcp-toolkit.git
git add .
git commit -m "Initial commit with release configuration"
git push -u origin main
```

### 3. 配置 GitHub Actions 权限 (必须) / Configure GitHub Actions Permissions (Required)

1. 访问 `https://github.com/YOUR_USERNAME/mcp-toolkit/settings/actions`
2. 在 "Workflow permissions" 选择 **"Read and write permissions"**
3. 勾选 **"Allow GitHub Actions to create and approve pull requests"**

---

## 🚀 方式 1: 仅发布到 GitHub Releases (推荐开始)

这是最简单的方式，用户可以通过安装脚本下载。

This is the simplest way, users can download via installation scripts.

### 步骤 / Steps

```bash
# 1. 确保所有测试通过
make test

# 2. 创建版本标签
git tag -a v1.0.0 -m "Release version 1.0.0"

# 3. 推送标签 (这会自动触发 GitHub Actions)
git push origin v1.0.0

# 4. 等待 2-5 分钟，GitHub Actions 会自动:
#    - 构建所有平台的二进制文件
#    - 创建 GitHub Release
#    - 上传所有文件

# 5. 访问 https://github.com/YOUR_USERNAME/mcp-toolkit/releases
#    查看并完善 Release 说明
```

### 用户安装方式 / User Installation

```bash
# Linux/macOS
curl -fsSL https://raw.githubusercontent.com/YOUR_USERNAME/mcp-toolkit/main/scripts/install.sh | bash

# Windows
Invoke-WebRequest -Uri "https://raw.githubusercontent.com/YOUR_USERNAME/mcp-toolkit/main/scripts/install.ps1" -OutFile "install.ps1"
.\install.ps1
```

---

## 🎯 方式 2: 同时发布到 GitHub + PyPI (支持 uv)

如果希望用户可以通过 `uv tool install mcp-toolkit` 安装，需要额外配置 PyPI。

If you want users to install via `uv tool install mcp-toolkit`, additional PyPI configuration is needed.

### 额外步骤 / Additional Steps

#### A. 注册 PyPI 账号

1. 访问 https://pypi.org/account/register/
2. 创建账号并验证邮箱
3. 启用 2FA (强烈推荐)

#### B. 创建 API Token

1. 访问 https://pypi.org/manage/account/token/
2. 点击 "Add API token"
3. Token name: `mcp-toolkit-github-actions`
4. Scope: "Entire account" (或创建项目后选择特定项目)
5. 复制生成的 token (格式: `pypi-...`)

#### C. 配置 GitHub Secret

1. 访问 `https://github.com/YOUR_USERNAME/mcp-toolkit/settings/secrets/actions`
2. 点击 "New repository secret"
3. Name: `PYPI_API_TOKEN`
4. Value: 粘贴你的 PyPI token
5. 点击 "Add secret"

#### D. 发布

```bash
# 1. 创建并推送标签 (同方式 1)
git tag -a v1.0.0 -m "Release version 1.0.0"
git push origin v1.0.0

# 2. GitHub Actions 会自动:
#    - 构建并发布到 GitHub Releases
#    - 构建并发布 Python 包到 PyPI

# 3. 等待 5-10 分钟后，检查:
#    - GitHub: https://github.com/YOUR_USERNAME/mcp-toolkit/releases
#    - PyPI: https://pypi.org/project/mcp-toolkit/
```

### 用户安装方式 / User Installation

```bash
# 方式 A: 使用 uv (推荐)
uv tool install mcp-toolkit

# 方式 B: 使用 pipx
pipx install mcp-toolkit

# 方式 C: 使用安装脚本 (同方式 1)
curl -fsSL https://raw.githubusercontent.com/YOUR_USERNAME/mcp-toolkit/main/scripts/install.sh | bash
```

---

## 🧪 测试发布 (可选)

在正式发布前，可以先测试：

Before official release, you can test:

### 测试 GitHub Release

```bash
# 使用 dev 标签测试
git tag -a v0.0.1-dev -m "Test release"
git push origin v0.0.1-dev

# 检查 Actions: https://github.com/YOUR_USERNAME/mcp-toolkit/actions
# 检查 Release: https://github.com/YOUR_USERNAME/mcp-toolkit/releases

# 测试完成后删除
git push --delete origin v0.0.1-dev
git tag -d v0.0.1-dev
```

### 测试 PyPI (使用 TestPyPI)

```bash
# 1. 注册 TestPyPI: https://test.pypi.org/account/register/
# 2. 创建 API token
# 3. 添加 GitHub Secret: TEST_PYPI_API_TOKEN

# 4. 手动触发工作流
# 访问: https://github.com/YOUR_USERNAME/mcp-toolkit/actions/workflows/publish-pypi.yml
# 点击 "Run workflow"
# 勾选 "Publish to TestPyPI"

# 5. 测试安装
uv tool install --index-url https://test.pypi.org/simple/ mcp-toolkit
```

---

## 📝 发布检查清单 / Release Checklist

发布前快速检查：

Quick check before release:

- [ ] 所有测试通过 (`make test`)
- [ ] 代码已格式化 (`make fmt`)
- [ ] 更新了 CHANGELOG.md
- [ ] 更新了版本号 (pyproject.toml, installer.py)
- [ ] 更新了仓库地址 (所有文件)
- [ ] GitHub Actions 权限已配置
- [ ] (可选) PyPI API Token 已配置

---

## 🔄 更新发布 / Update Release

发布新版本：

Release new version:

```bash
# 1. 更新代码
git add .
git commit -m "Add new features"
git push

# 2. 更新 CHANGELOG.md
# 3. 更新版本号

# 4. 创建新标签
git tag -a v1.1.0 -m "Release v1.1.0"
git push origin v1.1.0

# 5. 等待自动发布完成
```

---

## ❌ 回滚发布 / Rollback Release

如果发现问题需要回滚：

If issues are found and rollback is needed:

```bash
# 1. 删除远程标签
git push --delete origin v1.0.0

# 2. 删除本地标签
git tag -d v1.0.0

# 3. 在 GitHub 上删除 Release
# 访问: https://github.com/YOUR_USERNAME/mcp-toolkit/releases
# 点击对应 Release 的 "Delete"

# 4. (如果发布到 PyPI) 联系 PyPI 支持删除版本
# 或发布修复版本 v1.0.1
```

---

## 📊 监控发布 / Monitor Release

### GitHub Release

- **Releases**: https://github.com/YOUR_USERNAME/mcp-toolkit/releases
- **Actions**: https://github.com/YOUR_USERNAME/mcp-toolkit/actions
- **Insights**: https://github.com/YOUR_USERNAME/mcp-toolkit/pulse

### PyPI (如果发布)

- **项目页面**: https://pypi.org/project/mcp-toolkit/
- **下载统计**: https://pypistats.org/packages/mcp-toolkit
- **安全扫描**: https://pypi.org/project/mcp-toolkit/#security

---

## 🆘 常见问题 / Common Issues

### GitHub Actions 失败

```bash
# 查看日志
https://github.com/YOUR_USERNAME/mcp-toolkit/actions

# 常见原因:
# 1. 权限不足 -> 检查 Actions 权限设置
# 2. 构建失败 -> 检查 Go 版本和依赖
# 3. 标签格式错误 -> 确保以 'v' 开头 (v1.0.0)
```

### PyPI 发布失败

```bash
# 常见原因:
# 1. Token 无效 -> 重新创建并更新 Secret
# 2. 版本已存在 -> PyPI 不允许覆盖，需要新版本号
# 3. 包名冲突 -> 更改项目名称
```

### 用户安装失败

```bash
# 检查:
# 1. Release 是否公开
# 2. 文件是否正确上传
# 3. 校验和是否匹配
# 4. 安装脚本中的仓库地址是否正确
```

---

## 📚 相关文档 / Related Documentation

- [详细发布指南](PUBLISH_GUIDE.md)
- [发布检查清单](RELEASE_CHECKLIST.md)
- [安装指南](INSTALLATION.md)
- [发布文档](RELEASE.md)

---

## ✅ 完成！/ Done!

现在你的项目已经配置好发布流程，只需：

Now your project is configured for release, just:

1. 创建标签: `git tag -a v1.0.0 -m "Release v1.0.0"`
2. 推送标签: `git push origin v1.0.0`
3. 等待自动发布完成！

用户可以通过以下方式安装：

Users can install via:

```bash
# 使用 uv (如果发布到 PyPI)
uv tool install mcp-toolkit

# 使用安装脚本
curl -fsSL https://raw.githubusercontent.com/YOUR_USERNAME/mcp-toolkit/main/scripts/install.sh | bash
```

🎉 祝发布顺利！/ Happy releasing!

