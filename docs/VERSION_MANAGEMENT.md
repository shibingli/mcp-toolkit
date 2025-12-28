# 版本管理说明 / Version Management Guide

## 📋 版本管理策略

本项目采用 **Git Tag 驱动的自动版本管理** 策略，版本号由 GitHub Release 标签自动控制。

This project uses **Git Tag-driven automatic version management**, where version numbers are automatically controlled by
GitHub Release tags.

---

## 🔄 版本号流程

### 1. 源代码中的版本号

**pyproject.toml**:

```toml
[project]
name = "mcp-sandbox-toolkit"
version = "0.0.0"  # Placeholder, will be updated by GitHub Actions
```

**python/mcp_toolkit_wrapper/__init__.py**:

```python
import importlib.metadata

try:
    # 从已安装的包元数据获取版本
    __version__ = importlib.metadata.version("mcp-sandbox-toolkit")
except importlib.metadata.PackageNotFoundError:
    # 开发环境的后备版本
    __version__ = "0.0.0-dev"
```

### 2. GitHub Actions 自动更新

当推送 Git 标签时（如 `v1.2.0`），GitHub Actions 会：

1. **提取版本号**:
   ```bash
   VERSION=${GITHUB_REF#refs/tags/v}  # 移除 'v' 前缀
   ```

2. **更新 pyproject.toml**:
   ```bash
   sed -i "s/^version = .*/version = \"$VERSION\"/" pyproject.toml
   ```

3. **构建 Python 包**:
   ```bash
   python -m build
   ```

4. **发布到 PyPI**:
    - 包的版本号为标签版本（如 `1.2.0`）
    - 用户安装后，`importlib.metadata.version()` 会返回正确的版本号

---

## 🚀 发布新版本

### 步骤 1: 确保代码已提交

```bash
git add -A
git commit -m "Your commit message"
git push origin main
```

### 步骤 2: 创建并推送标签

```bash
# 创建标签（使用语义化版本号）
git tag -a v1.2.0 -m "Release v1.2.0"

# 推送标签到远程仓库
git push origin v1.2.0
```

### 步骤 3: GitHub Actions 自动执行

推送标签后，GitHub Actions 会自动：

- ✅ 运行所有测试
- ✅ 构建所有平台的二进制文件
- ✅ 更新 `pyproject.toml` 中的版本号
- ✅ 构建 Python 包
- ✅ 发布到 PyPI
- ✅ 创建 GitHub Release

### 步骤 4: 验证发布

```bash
# 检查 GitHub Release
# https://github.com/shibingli/mcp-toolkit/releases

# 检查 PyPI
# https://pypi.org/project/mcp-sandbox-toolkit/

# 测试安装
pip install mcp-sandbox-toolkit==1.2.0
python -c "import mcp_toolkit_wrapper; print(mcp_toolkit_wrapper.__version__)"
```

---

## 📝 版本号规范

遵循 [语义化版本 2.0.0](https://semver.org/lang/zh-CN/) 规范：

```
MAJOR.MINOR.PATCH

例如: 1.2.3
```

- **MAJOR**: 不兼容的 API 修改
- **MINOR**: 向下兼容的功能性新增
- **PATCH**: 向下兼容的问题修正

### 示例

- `v1.0.0` - 首个稳定版本
- `v1.1.0` - 添加新功能
- `v1.1.1` - 修复 bug
- `v2.0.0` - 重大更新，可能不兼容

---

## 🔍 版本号查询

### 查看当前版本

**已安装的包**:

```python
import mcp_toolkit_wrapper
print(mcp_toolkit_wrapper.__version__)
```

**命令行**:

```bash
python -c "import mcp_toolkit_wrapper; print(mcp_toolkit_wrapper.__version__)"
```

**pip**:

```bash
pip show mcp-sandbox-toolkit
```

### 查看可用版本

```bash
pip index versions mcp-sandbox-toolkit
```

---

## 🛠️ 开发环境

在开发环境中（未安装包），版本号会显示为 `0.0.0-dev`：

```python
>>> import mcp_toolkit_wrapper
>>> mcp_toolkit_wrapper.__version__
'0.0.0-dev'
```

这是正常的，因为包还没有被安装，无法从元数据中获取版本号。

---

## ⚠️ 注意事项

### 1. 不要手动修改版本号

❌ **错误做法**:

```toml
# pyproject.toml
version = "1.2.0"  # 手动修改
```

✅ **正确做法**:

```bash
# 通过 Git 标签控制版本
git tag -a v1.2.0 -m "Release v1.2.0"
git push origin v1.2.0
```

### 2. 标签格式

- ✅ 正确: `v1.2.0`, `v2.0.0-beta.1`
- ❌ 错误: `1.2.0`, `version-1.2.0`, `release-1.2.0`

### 3. 删除错误的标签

如果推送了错误的标签：

```bash
# 删除本地标签
git tag -d v1.2.0

# 删除远程标签
git push origin :refs/tags/v1.2.0
```

### 4. 版本号一致性

确保以下版本号一致：

- Git 标签版本
- PyPI 包版本
- GitHub Release 版本
- 二进制文件版本

这些都由 GitHub Actions 自动保证一致性。

---

## 📊 版本历史

查看所有版本：

```bash
# 查看所有标签
git tag -l

# 查看标签详情
git show v1.2.0

# 查看版本历史
git log --oneline --decorate --tags
```

---

## 🔗 相关链接

- **GitHub Releases**: https://github.com/shibingli/mcp-toolkit/releases
- **PyPI**: https://pypi.org/project/mcp-sandbox-toolkit/
- **语义化版本**: https://semver.org/lang/zh-CN/
- **GitHub Actions**: https://github.com/shibingli/mcp-toolkit/actions

---

**最后更新**: 2025-12-28

