# 发布检查清单 / Release Checklist

在发布新版本之前，请确保完成以下所有步骤。

Before releasing a new version, make sure to complete all the following steps.

---

## 发布前准备 / Pre-Release Preparation

### 1. 代码质量检查 / Code Quality Check

- [ ] 所有单元测试通过 / All unit tests pass
  ```bash
  make test
  ```

- [ ] 测试覆盖率达标 / Test coverage meets requirements
  ```bash
  make test-cover
  # 检查 coverage.html
  ```

- [ ] 代码格式化 / Code formatted
  ```bash
  make fmt
  ```

- [ ] Linter 检查通过 / Linter checks pass
  ```bash
  make lint
  ```

- [ ] 没有已知的 bug / No known bugs

### 2. 文档更新 / Documentation Update

- [ ] README.md 更新 / README.md updated
  - [ ] 版本号正确 / Version number correct
  - [ ] 新功能已记录 / New features documented
  - [ ] 示例代码可运行 / Example code works

- [ ] CHANGELOG.md 更新 / CHANGELOG.md updated
  - [ ] 新增功能列表 / New features list
  - [ ] Bug 修复列表 / Bug fixes list
  - [ ] 破坏性变更说明 / Breaking changes noted

- [ ] API 文档更新 / API documentation updated

- [ ] 安装文档更新 / Installation docs updated

### 3. 版本号管理 / Version Management

- [ ] 确定版本号 (遵循语义化版本) / Determine version number (follow semantic versioning)
  - 主版本号 (Major): 不兼容的 API 变更
  - 次版本号 (Minor): 向后兼容的功能新增
  - 修订号 (Patch): 向后兼容的问题修正

- [ ] 更新 `pyproject.toml` 中的版本号 / Update version in `pyproject.toml`

- [ ] 更新 `python/mcp_toolkit_wrapper/installer.py` 中的版本号 / Update version in installer

### 4. 依赖检查 / Dependency Check

- [ ] 依赖项是最新的 / Dependencies are up to date
  ```bash
  go mod tidy
  ```

- [ ] 没有安全漏洞 / No security vulnerabilities
  ```bash
  go list -json -m all | nancy sleuth
  ```

---

## 构建和测试 / Build and Test

### 5. 本地构建测试 / Local Build Test

- [ ] 本地构建成功 / Local build succeeds
  ```bash
  make build
  ```

- [ ] 跨平台构建成功 / Cross-platform build succeeds
  ```bash
  make build-all
  ```

- [ ] 检查构建产物 / Check build artifacts
  ```bash
  ls -lh dist/
  ```

### 6. 功能测试 / Functional Testing

- [ ] Stdio 传输模式测试 / Stdio transport test
  ```bash
  ./dist/mcp-toolkit -transport stdio
  ```

- [ ] HTTP 传输模式测试 / HTTP transport test
  ```bash
  ./dist/mcp-toolkit -transport http -http-port 8080
  ```

- [ ] SSE 传输模式测试 / SSE transport test
  ```bash
  ./dist/mcp-toolkit -transport sse -sse-port 8081
  ```

- [ ] 客户端测试通过 / Client tests pass
  ```bash
  make test-client
  ```

### 7. 平台测试 / Platform Testing

- [ ] Windows (amd64) 测试 / Windows (amd64) tested
- [ ] Windows (arm64) 测试 / Windows (arm64) tested
- [ ] Linux (amd64) 测试 / Linux (amd64) tested
- [ ] Linux (arm64) 测试 / Linux (arm64) tested
- [ ] macOS (Intel) 测试 / macOS (Intel) tested
- [ ] macOS (Apple Silicon) 测试 / macOS (Apple Silicon) tested

---

## 发布流程 / Release Process

### 8. 创建 Git 标签 / Create Git Tag

```bash
# 创建带注释的标签 / Create annotated tag
git tag -a v1.0.0 -m "Release version 1.0.0"

# 查看标签 / View tag
git tag -l -n9 v1.0.0

# 推送标签 / Push tag
git push origin v1.0.0
```

- [ ] 标签已创建 / Tag created
- [ ] 标签已推送 / Tag pushed

### 9. GitHub Actions 检查 / GitHub Actions Check

- [ ] GitHub Actions 工作流触发 / GitHub Actions workflow triggered
- [ ] 所有平台构建成功 / All platforms built successfully
- [ ] Release 自动创建 / Release automatically created
- [ ] 所有文件已上传 / All files uploaded
  - [ ] Windows (amd64) .zip
  - [ ] Windows (arm64) .zip
  - [ ] Linux (amd64) .tar.gz
  - [ ] Linux (arm64) .tar.gz
  - [ ] macOS (Intel) .tar.gz
  - [ ] macOS (Apple Silicon) .tar.gz
  - [ ] checksums.txt

### 10. Release 页面完善 / Release Page Enhancement

- [ ] Release 标题清晰 / Release title clear
- [ ] Release 描述完整 / Release description complete
  - [ ] 主要变更 / Major changes
  - [ ] 新功能 / New features
  - [ ] Bug 修复 / Bug fixes
  - [ ] 已知问题 / Known issues
  - [ ] 升级说明 / Upgrade notes

- [ ] 下载链接可用 / Download links work
- [ ] 校验和正确 / Checksums correct

---

## 发布后验证 / Post-Release Verification

### 11. 安装测试 / Installation Testing

- [ ] 使用安装脚本测试 (Linux/macOS) / Test with install script (Linux/macOS)
  ```bash
  curl -fsSL https://raw.githubusercontent.com/shibingli/mcp-toolkit/main/scripts/install.sh | bash
  ```

- [ ] 使用安装脚本测试 (Windows) / Test with install script (Windows)
  ```powershell
  Invoke-WebRequest -Uri "..." -OutFile "install.ps1"
  .\install.ps1
  ```

- [ ] 手动下载安装测试 / Test manual download installation

### 12. Python 包发布 (可选) / Python Package Release (Optional)

如果要发布到 PyPI:

If publishing to PyPI:

- [ ] 构建 Python 包 / Build Python package
  ```bash
  python -m build
  ```

- [ ] 上传到 TestPyPI / Upload to TestPyPI
  ```bash
  python -m twine upload --repository testpypi dist/*
  ```

- [ ] 从 TestPyPI 安装测试 / Test install from TestPyPI
  ```bash
  uv tool install --index-url https://test.pypi.org/simple/ mcp-toolkit
  ```

- [ ] 上传到 PyPI / Upload to PyPI
  ```bash
  python -m twine upload dist/*
  ```

- [ ] 从 PyPI 安装测试 / Test install from PyPI
  ```bash
  uv tool install mcp-toolkit
  ```

### 13. 文档和公告 / Documentation and Announcements

- [ ] 更新官方文档 / Update official documentation
- [ ] 发布公告 (如果适用) / Post announcement (if applicable)
  - [ ] GitHub Discussions
  - [ ] 社交媒体 / Social media
  - [ ] 邮件列表 / Mailing list

### 14. 监控和反馈 / Monitoring and Feedback

- [ ] 监控 GitHub Issues / Monitor GitHub Issues
- [ ] 监控下载统计 / Monitor download statistics
- [ ] 收集用户反馈 / Collect user feedback

---

## 回滚计划 / Rollback Plan

如果发现严重问题需要回滚:

If serious issues are found and rollback is needed:

1. **删除有问题的 Release** / Delete problematic release
   - 在 GitHub Releases 页面删除

2. **删除 Git 标签** / Delete Git tag
   ```bash
   git tag -d v1.0.0
   git push origin :refs/tags/v1.0.0
   ```

3. **发布修复版本** / Release fix version
   - 增加修订号 (例如 v1.0.1)
   - 重新执行发布流程

---

## 版本发布时间表 / Release Schedule

建议的发布周期:

Recommended release cycle:

- **主版本 (Major)**: 每年 1-2 次 / 1-2 times per year
- **次版本 (Minor)**: 每季度 1-2 次 / 1-2 times per quarter
- **修订版本 (Patch)**: 根据需要 / As needed

---

## 联系人 / Contacts

- **发布负责人** / Release Manager: [Name]
- **技术负责人** / Tech Lead: [Name]
- **QA 负责人** / QA Lead: [Name]

---

## 附录: 快速发布命令 / Appendix: Quick Release Commands

```bash
# 1. 运行所有测试
make test test-cover

# 2. 格式化和检查代码
make fmt lint

# 3. 构建所有平台
make build-all

# 4. 创建并推送标签
git tag -a v1.0.0 -m "Release v1.0.0"
git push origin v1.0.0

# 5. 等待 GitHub Actions 完成
# 6. 验证 Release 页面
# 7. 测试安装
```

