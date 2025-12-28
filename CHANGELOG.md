# Changelog / 更新日志

All notable changes to this project will be documented in this file.

本文件记录项目的所有重要变更。

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

---

## [Unreleased]

### Added / 新增
- 跨平台构建脚本 (Windows, Linux, macOS) / Cross-platform build scripts
- GitHub Actions 自动发布工作流 / GitHub Actions auto-release workflow
- 安装脚本 (Linux/macOS/Windows) / Installation scripts
- Python 包装器，支持通过 uv/pipx 安装 / Python wrapper for uv/pipx installation
- Makefile 简化构建流程 / Makefile for simplified build process
- 详细的发布文档和检查清单 / Detailed release documentation and checklist

### Changed / 变更
- 无 / None

### Deprecated / 弃用
- 无 / None

### Removed / 移除
- 无 / None

### Fixed / 修复
- 无 / None

### Security / 安全
- 无 / None

---

## [1.0.0] - YYYY-MM-DD

### Added / 新增
- 初始版本发布 / Initial release
- 支持 Stdio、HTTP、SSE 三种传输方式 / Support for Stdio, HTTP, SSE transports
- 完整的文件系统操作工具 (26个) / Complete file system operation tools (26 tools)
  - 文件操作 (11个) / File operations (11)
  - 目录操作 (2个) / Directory operations (2)
  - 命令执行 (3个) / Command execution (3)
  - 异步操作 (3个) / Async operations (3)
  - 命令历史 (2个) / Command history (2)
  - 权限管理 (2个) / Permission management (2)
  - 系统工具 (3个) / System tools (3)
- 安全的沙箱机制 / Secure sandbox mechanism
- 高性能 JSON 处理 (Sonic) / High-performance JSON processing (Sonic)
- Panic recovery 机制 / Panic recovery mechanism
- 完整的单元测试 / Complete unit tests
- 详细的文档 / Comprehensive documentation

### Technical Details / 技术细节
- Go 1.25.5+ 支持 / Go 1.25.5+ support
- MCP SDK v1.2.0 集成 / MCP SDK v1.2.0 integration
- 多 JSON 库支持 (Sonic, go-json, jsoniter) / Multiple JSON library support
- Zap 日志系统 / Zap logging system
- 跨平台兼容 (Windows, Linux, macOS) / Cross-platform compatibility
- 多架构支持 (amd64, arm64) / Multi-architecture support

---

## 版本说明 / Version Notes

### 语义化版本规则 / Semantic Versioning Rules

- **主版本号 (Major)**: 不兼容的 API 变更 / Incompatible API changes
- **次版本号 (Minor)**: 向后兼容的功能新增 / Backward-compatible new features
- **修订号 (Patch)**: 向后兼容的问题修正 / Backward-compatible bug fixes

### 版本类型标签 / Version Type Labels

- `[Added]` - 新增功能 / New features
- `[Changed]` - 功能变更 / Changes in existing functionality
- `[Deprecated]` - 即将移除的功能 / Soon-to-be removed features
- `[Removed]` - 已移除的功能 / Removed features
- `[Fixed]` - Bug 修复 / Bug fixes
- `[Security]` - 安全相关修复 / Security fixes

---

## 链接 / Links

[Unreleased]: https://github.com/shibingli/mcp-toolkit/compare/v1.0.0...HEAD
[1.0.0]: https://github.com/shibingli/mcp-toolkit/releases/tag/v1.0.0

