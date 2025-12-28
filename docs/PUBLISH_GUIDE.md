# 发布到公开仓库指南 / Publishing to Public Repository Guide

本指南详细说明如何将 MCP Toolkit 发布到 GitHub 和 PyPI，以便用户可以通过 `uv` 或其他包管理器下载和安装。

This guide explains how to publish MCP Toolkit to GitHub and PyPI so users can download and install it via `uv` or other package managers.

---

## 目录 / Table of Contents

1. [发布到 GitHub Releases](#1-发布到-github-releases)
2. [发布到 PyPI](#2-发布到-pypi-可选)
3. [配置自动化发布](#3-配置自动化发布)
4. [用户安装方式](#4-用户安装方式)

---

## 1. 发布到 GitHub Releases

### 步骤 1: 准备仓库 / Step 1: Prepare Repository

```bash
# 1. 创建 GitHub 仓库 / Create GitHub repository
# 在 GitHub 上创建新仓库，例如: shibingli/mcp-toolkit

# 2. 克隆或关联仓库 / Clone or link repository
git remote add origin https://github.com/shibingli/mcp-toolkit.git

# 3. 推送代码 / Push code
git add .
git commit -m "Initial commit"
git push -u origin main
```

### 步骤 2: 更新配置文件 / Step 2: Update Configuration Files

需要更新以下文件中的仓库信息：

Update repository information in the following files:

1. **scripts/install.sh** (第 8 行)
   ```bash
   REPO="shibingli/mcp-toolkit"  # 替换为你的仓库
   ```

2. **scripts/install.ps1** (第 10 行)
   ```powershell
   $Repo = "shibingli/mcp-toolkit"  # 替换为你的仓库
   ```

3. **python/mcp_toolkit_wrapper/installer.py** (第 17 行)
   ```python
   REPO = "shibingli/mcp-toolkit"  # 替换为你的仓库
   ```

4. **pyproject.toml** (URLs 部分)
   ```toml
   [project.urls]
   Homepage = "https://github.com/shibingli/mcp-toolkit"
   Repository = "https://github.com/shibingli/mcp-toolkit"
   Issues = "https://github.com/shibingli/mcp-toolkit/issues"
   ```

5. **docs/RELEASE.md** 和其他文档中的链接

### 步骤 3: 创建第一个发布 / Step 3: Create First Release

```bash
# 1. 确保所有测试通过 / Ensure all tests pass
make test

# 2. 构建所有平台 / Build all platforms
make build-all

# 3. 创建版本标签 / Create version tag
git tag -a v1.0.0 -m "Release version 1.0.0"

# 4. 推送标签 / Push tag
git push origin v1.0.0
```

### 步骤 4: GitHub Actions 自动发布 / Step 4: GitHub Actions Auto-Release

推送标签后，GitHub Actions 会自动：

After pushing the tag, GitHub Actions will automatically:

1. 构建所有平台的二进制文件 / Build binaries for all platforms
2. 生成 SHA256 校验和 / Generate SHA256 checksums
3. 创建 GitHub Release / Create GitHub Release
4. 上传所有构建产物 / Upload all build artifacts

查看进度：`https://github.com/shibingli/mcp-toolkit/actions`

### 步骤 5: 完善 Release 页面 / Step 5: Enhance Release Page

1. 访问 `https://github.com/shibingli/mcp-toolkit/releases`
2. 编辑自动创建的 Release
3. 添加详细的发布说明：
   - 主要变更 / Major changes
   - 新功能 / New features
   - Bug 修复 / Bug fixes
   - 升级说明 / Upgrade notes

---

## 2. 发布到 PyPI (可选)

如果希望用户可以通过 `uv tool install mcp-toolkit` 安装，需要发布到 PyPI。

If you want users to install via `uv tool install mcp-toolkit`, publish to PyPI.

#### 前置要求 / Prerequisites

确保已安装 Python 构建工具：

Make sure Python build tools are installed:

```bash
python -m pip install --upgrade pip
python -m pip install --upgrade build twine hatchling
```

详见 [前置依赖文档](PREREQUISITES.md)。

See [Prerequisites](PREREQUISITES.md) for details.

### 步骤 1: 注册 PyPI 账号 / Step 1: Register PyPI Account

1. 访问 https://pypi.org/account/register/
2. 创建账号并验证邮箱 / Create account and verify email
3. 启用 2FA (推荐) / Enable 2FA (recommended)

### 步骤 2: 创建 API Token / Step 2: Create API Token

1. 访问 https://pypi.org/manage/account/token/
2. 创建新的 API token
3. 保存 token (只显示一次) / Save token (shown only once)

### 步骤 3: 配置 GitHub Secrets / Step 3: Configure GitHub Secrets

1. 访问 `https://github.com/shibingli/mcp-toolkit/settings/secrets/actions`
2. 添加新的 secret:
   - Name: `PYPI_API_TOKEN`
   - Value: 你的 PyPI API token

### 步骤 4: 添加 PyPI 发布工作流 / Step 4: Add PyPI Release Workflow

创建 `.github/workflows/publish-pypi.yml`:

```yaml
name: Publish to PyPI

on:
  release:
    types: [published]

jobs:
  publish:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      
      - name: Set up Python
        uses: actions/setup-python@v5
        with:
          python-version: '3.11'
      
      - name: Install dependencies
        run: |
          python -m pip install --upgrade pip
          pip install build twine
      
      - name: Build package
        run: python -m build
      
      - name: Publish to PyPI
        env:
          TWINE_USERNAME: __token__
          TWINE_PASSWORD: ${{ secrets.PYPI_API_TOKEN }}
        run: python -m twine upload dist/*
```

### 步骤 5: 测试发布 / Step 5: Test Publishing

先发布到 TestPyPI 测试：

Test by publishing to TestPyPI first:

```bash
# 1. 注册 TestPyPI 账号 / Register TestPyPI account
# https://test.pypi.org/account/register/

# 2. 构建包 / Build package
python -m build

# 3. 上传到 TestPyPI / Upload to TestPyPI
python -m twine check dist/*
python -m twine upload --repository testpypi dist/*

# 4. 测试安装 / Test installation
uv tool install --index-url https://test.pypi.org/simple/ mcp-toolkit
```

### 步骤 6: 正式发布 / Step 6: Official Release

```bash
# 创建 Release 时会自动触发 PyPI 发布
git tag -a v1.0.0 -m "Release v1.0.0"
git push origin v1.0.0
```

---

## 3. 配置自动化发布

### GitHub Actions 工作流说明 / GitHub Actions Workflow Explanation

项目已包含两个工作流：

The project includes two workflows:

1. **`.github/workflows/release.yml`** - 构建和发布到 GitHub Releases
2. **`.github/workflows/publish-pypi.yml`** (可选) - 发布到 PyPI

### 触发条件 / Trigger Conditions

- **自动触发**: 推送以 `v` 开头的标签 (例如 `v1.0.0`)
- **手动触发**: 在 GitHub Actions 页面手动运行

### 权限配置 / Permission Configuration

确保 GitHub Actions 有足够的权限：

Ensure GitHub Actions has sufficient permissions:

1. 访问 `Settings` → `Actions` → `General`
2. 在 "Workflow permissions" 中选择 "Read and write permissions"
3. 勾选 "Allow GitHub Actions to create and approve pull requests"

---

## 4. 用户安装方式

发布后，用户可以通过以下方式安装：

After publishing, users can install via:

### 方式 1: 使用 uv (推荐)

```bash
# 安装 uv
curl -LsSf https://astral.sh/uv/install.sh | sh

# 安装 MCP Toolkit
uv tool install mcp-toolkit

# 运行
mcp-toolkit --help
```

### 方式 2: 使用安装脚本

```bash
# Linux/macOS
curl -fsSL https://raw.githubusercontent.com/shibingli/mcp-toolkit/main/scripts/install.sh | bash

# Windows (PowerShell)
Invoke-WebRequest -Uri "https://raw.githubusercontent.com/shibingli/mcp-toolkit/main/scripts/install.ps1" -OutFile "install.ps1"
.\install.ps1
```

### 方式 3: 手动下载

1. 访问 https://github.com/shibingli/mcp-toolkit/releases
2. 下载对应平台的压缩包
3. 解压并添加到 PATH

---

## 5. 发布后验证

### 检查清单 / Checklist

- [ ] GitHub Release 页面正常显示 / GitHub Release page displays correctly
- [ ] 所有平台的二进制文件可下载 / All platform binaries downloadable
- [ ] 校验和文件正确 / Checksum file correct
- [ ] 安装脚本可用 / Installation scripts work
- [ ] PyPI 页面正常 (如果发布) / PyPI page normal (if published)
- [ ] `uv tool install` 可用 / `uv tool install` works
- [ ] 文档链接正确 / Documentation links correct

### 测试安装 / Test Installation

```bash
# 测试 GitHub 安装脚本
curl -fsSL https://raw.githubusercontent.com/shibingli/mcp-toolkit/main/scripts/install.sh | bash
mcp-toolkit --help

# 测试 uv 安装 (如果发布到 PyPI)
uv tool install mcp-toolkit
mcp-toolkit --help
```

---

## 6. 维护和更新

### 发布新版本 / Release New Version

```bash
# 1. 更新代码和文档
# 2. 更新 CHANGELOG.md
# 3. 更新版本号 (pyproject.toml, installer.py)
# 4. 提交更改
git add .
git commit -m "Prepare for v1.1.0 release"
git push

# 5. 创建新标签
git tag -a v1.1.0 -m "Release v1.1.0"
git push origin v1.1.0

# 6. 等待 GitHub Actions 完成
# 7. 验证发布
```

### 修复发布问题 / Fix Release Issues

如果发现问题需要重新发布：

If issues are found and re-release is needed:

```bash
# 1. 删除远程标签
git push --delete origin v1.0.0

# 2. 删除本地标签
git tag -d v1.0.0

# 3. 删除 GitHub Release (在网页上操作)

# 4. 修复问题后重新发布
git tag -a v1.0.0 -m "Release v1.0.0 (fixed)"
git push origin v1.0.0
```

---

## 7. 常见问题 / FAQ

### Q: GitHub Actions 构建失败怎么办？

**A:** 检查以下几点：
1. Go 版本是否正确 (1.25.5+)
2. 依赖是否都能下载
3. 构建脚本是否有执行权限
4. 查看 Actions 日志获取详细错误信息

### Q: 如何更新已发布的 Release？

**A:** 
1. 在 GitHub Release 页面点击 "Edit release"
2. 更新描述或上传新文件
3. 保存更改

### Q: 用户报告下载失败怎么办?

**A:**
1. 检查文件是否正确上传
2. 检查文件大小和校验和
3. 确认 GitHub Release 是公开的
4. 检查网络连接问题

---

## 相关资源 / Related Resources

- [GitHub Actions 文档](https://docs.github.com/en/actions)
- [PyPI 发布指南](https://packaging.python.org/tutorials/packaging-projects/)
- [语义化版本规范](https://semver.org/)
- [Keep a Changelog](https://keepachangelog.com/)

