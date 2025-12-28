# GitHub Actions 工作流说明

本项目包含三个主要的 GitHub Actions 工作流，用于自动化构建、测试和发布流程。

## 📋 工作流概览

### 1. CI (持续集成)
**文件**: `.github/workflows/ci.yml`

**触发条件**:
- Push 到 `main` 或 `develop` 分支
- Pull Request 到 `main` 或 `develop` 分支
- 手动触发

**功能**:
- ✅ 在多个操作系统上运行测试 (Ubuntu, Windows, macOS)
- ✅ 代码质量检查 (golangci-lint)
- ✅ 测试覆盖率报告
- ✅ 构建验证
- ✅ Python 包验证

**使用场景**: 每次代码提交时自动运行，确保代码质量

---

### 2. Release (发布)
**文件**: `.github/workflows/release.yml`

**触发条件**:
- Push 标签 (格式: `v*`, 例如 `v1.0.0`)
- 手动触发（可指定版本号）

**功能**:
- ✅ 构建 6 个平台的二进制文件
  - Windows (amd64, arm64) - zip 格式
  - Linux (amd64, arm64) - tar.gz 格式
  - macOS (amd64, arm64) - tar.gz 格式
- ✅ 生成 SHA256 校验和
- ✅ 运行测试并生成覆盖率报告
- ✅ 创建 GitHub Release
- ✅ 上传所有构建产物
- ✅ 自动生成发布说明

**使用方法**:

#### 方式 1: 使用 Git 标签（推荐）
```bash
# 1. 更新 CHANGELOG.md

# 2. 提交更改
git add .
git commit -m "Prepare for release v1.0.0"
git push

# 3. 创建并推送标签
git tag -a v1.0.0 -m "Release version 1.0.0"
git push origin v1.0.0

# 4. 等待 GitHub Actions 完成（约 2-5 分钟）
# 5. 访问 https://github.com/YOUR_USERNAME/mcp-toolkit/releases
```

#### 方式 2: 手动触发
1. 访问 GitHub Actions 页面
2. 选择 "Release" 工作流
3. 点击 "Run workflow"
4. 输入版本号（例如: `v1.0.0`）
5. 点击 "Run workflow"

**输出**:
- GitHub Release 页面包含所有二进制文件
- 自动生成的发布说明
- 安装说明

---

### 3. Publish to PyPI (发布到 PyPI)
**文件**: `.github/workflows/publish-pypi.yml`

**触发条件**:
- GitHub Release 发布时自动触发
- 手动触发（可选择 TestPyPI 或 PyPI）

**功能**:
- ✅ 构建 Python 包 (wheel 和 source distribution)
- ✅ 验证包的完整性
- ✅ 发布到 PyPI 或 TestPyPI
- ✅ 自动更新版本号

**前置要求**:

1. **注册 PyPI 账号**
   - 访问 https://pypi.org/account/register/

2. **创建 API Token**
   - 访问 https://pypi.org/manage/account/token/
   - 创建一个新的 API token
   - 复制 token（只显示一次）

3. **添加 GitHub Secret**
   - 访问 `https://github.com/YOUR_USERNAME/mcp-toolkit/settings/secrets/actions`
   - 点击 "New repository secret"
   - Name: `PYPI_API_TOKEN`
   - Value: 粘贴你的 PyPI token
   - 点击 "Add secret"

**使用方法**:

#### 自动发布（推荐）
当你创建 GitHub Release 时，会自动触发 PyPI 发布：
```bash
git tag -a v1.0.0 -m "Release version 1.0.0"
git push origin v1.0.0
# Release 工作流会创建 GitHub Release
# PyPI 工作流会自动发布到 PyPI
```

#### 手动发布
1. 访问 GitHub Actions 页面
2. 选择 "Publish to PyPI" 工作流
3. 点击 "Run workflow"
4. 选择选项:
   - **Publish to TestPyPI**: 勾选以发布到测试环境
   - **Version**: 输入版本号（留空使用 git tag）
5. 点击 "Run workflow"

---

## 🔧 配置说明

### 必需的 GitHub Secrets

| Secret 名称 | 用途 | 必需 |
|------------|------|------|
| `PYPI_API_TOKEN` | PyPI 发布 | 可选* |

*如果不发布到 PyPI，可以不配置

### GitHub Actions 权限设置

1. 访问 `https://github.com/shibingli/mcp-toolkit/settings/actions`
2. 在 "Workflow permissions" 部分
3. 选择 "Read and write permissions"
4. 勾选 "Allow GitHub Actions to create and approve pull requests"
5. 点击 "Save"

---

## 📊 工作流状态徽章

在 README.md 中添加状态徽章：

```markdown
![CI](https://github.com/YOUR_USERNAME/mcp-toolkit/workflows/CI/badge.svg)
![Release](https://github.com/YOUR_USERNAME/mcp-toolkit/workflows/Release/badge.svg)
```

---

## 🐛 故障排查

### Release 工作流失败

**问题**: 构建失败
```
Error: go: command not found
```
**解决**: 检查 Go 版本是否正确（需要 1.25.5+）

**问题**: 权限错误
```
Error: Resource not accessible by integration
```
**解决**: 检查 GitHub Actions 权限设置（见上文）

### PyPI 发布失败

**问题**: Token 无效
```
Error: Invalid or non-existent authentication information
```
**解决**: 
1. 检查 Secret 名称是否正确
2. 重新生成 PyPI token
3. 更新 GitHub Secret

**问题**: 版本已存在
```
Error: File already exists
```
**解决**: PyPI 不允许重新上传相同版本，需要增加版本号

---

## 📝 最佳实践

### 发布流程

1. **开发阶段**
   - 在 `develop` 分支开发
   - 每次 push 触发 CI 测试

2. **准备发布**
   - 更新 `CHANGELOG.md`
   - 更新版本号（如果需要）
   - 合并到 `main` 分支

3. **创建发布**
   - 创建 git tag
   - 推送 tag 触发 Release 工作流
   - Release 工作流创建 GitHub Release
   - PyPI 工作流自动发布到 PyPI

4. **验证发布**
   - 检查 GitHub Release 页面
   - 测试安装: `uv tool install mcp-toolkit`
   - 验证功能

### 版本号规范

遵循语义化版本 (Semantic Versioning):
- **主版本** (Major): 不兼容的 API 变更 → `v2.0.0`
- **次版本** (Minor): 向后兼容的新功能 → `v1.1.0`
- **修订版** (Patch): 向后兼容的 Bug 修复 → `v1.0.1`
- **预发布**: `v1.0.0-alpha.1`, `v1.0.0-beta.1`, `v1.0.0-rc.1`

---

## 🔗 相关链接

- [GitHub Actions 文档](https://docs.github.com/en/actions)
- [PyPI 发布指南](https://packaging.python.org/tutorials/packaging-projects/)
- [语义化版本](https://semver.org/lang/zh-CN/)

