# MCP ToolKit / MCP 工具集

基于 Model Context Protocol (MCP) 设计和实现的综合工具集，提供文件系统操作、命令执行、系统工具等多种功能，兼容 Windows、Linux 和 macOS。

A comprehensive MCP tools collection based on the Model Context Protocol (MCP), providing filesystem operations, command execution, system tools and more, compatible with Windows, Linux, and macOS.

## 项目简介 / Project Overview

MCP ToolKit 是一个功能丰富、安全可靠的 MCP 工具集合，旨在为 AI 模型提供强大的系统交互能力。项目采用模块化设计，支持灵活扩展，未来将持续集成更多实用工具。

MCP ToolKit is a feature-rich and secure MCP tools collection designed to provide powerful system interaction capabilities for AI models. The project adopts a modular design, supports flexible expansion, and will continue to integrate more practical tools in the future.

### 核心特性 / Core Features

- 🚀 **多功能集成** - 文件系统、命令执行、系统工具等多种功能
- 🔒 **安全可靠** - 沙箱隔离、黑名单机制、路径验证等多重安全保障
- ⚡ **高性能** - Sonic JSON库、结构体预热等性能优化
- 🛡️ **稳定性强** - Panic Recovery机制确保服务稳定运行
- 🔌 **灵活传输** - 支持 Stdio、HTTP、SSE 多种传输方式
- 🌍 **跨平台** - 完美支持 Windows、Linux、macOS
- 📦 **模块化设计** - 易于扩展和维护

## 功能特性 / Features

### 文件操作 / File Operations
- ✅ 创建文件 / Create files
- ✅ 读取文件 / Read files
- ✅ 写入文件 / Write files
- ✅ 删除文件 / Delete files
- ✅ 复制文件 / Copy files
- ✅ 移动文件 / Move files
- ✅ 获取文件状态 / Get file status
- ✅ 检查文件是否存在 / Check file existence

### 目录操作 / Directory Operations
- ✅ 创建目录 / Create directories
- ✅ 列出目录内容 / List directory contents
- ✅ 删除目录 / Delete directories
- ✅ 复制目录 / Copy directories
- ✅ 移动目录 / Move directories

### 批量操作 / Batch Operations
- ✅ 批量删除 / Batch delete
- ✅ 文件搜索 / File search (支持通配符 / supports wildcards)

### 命令执行 / Command Execution
- ✅ 在沙箱内执行命令 / Execute commands within sandbox
- ✅ 工作目录管理 / Working directory management
- ✅ 命令黑名单保护 / Command blacklist protection
- ✅ 目录黑名单保护 / Directory blacklist protection
- ✅ 命令超时控制 / Command timeout control
- ✅ 输出捕获(stdout/stderr) / Output capture (stdout/stderr)
- ✅ 跨平台命令支持 / Cross-platform command support
- ✅ 异步命令执行 / Asynchronous command execution
- ✅ 命令执行历史记录 / Command execution history
- ✅ 权限级别控制 / Permission level control
- ✅ 环境变量配置 / Environment variable configuration
- ✅ 审计日志 / Audit logging

### 系统工具 / System Tools
- ✅ 获取当前系统时间 / Get current system time
- ✅ 获取系统信息 / Get system information (OS, CPU, Memory, GPU, Network)

### 安全特性 / Security Features
- ✅ 沙箱目录限制 / Sandbox directory restriction
- ✅ 路径遍历保护 / Path traversal protection
- ✅ 命令黑名单机制 / Command blacklist mechanism
- ✅ 命令参数路径验证 / Command argument path validation
- ✅ 系统目录保护 / System directory protection
- ✅ 危险命令拦截 / Dangerous command interception

### 稳定性保障 / Stability Assurance
- ✅ **Panic Recovery 机制** / Panic recovery mechanism
  - 工具层 panic 恢复 / Tool-level panic recovery
  - 传输层 panic 恢复 / Transport-level panic recovery
  - 完整的堆栈跟踪记录 / Complete stack trace logging
  - 优雅的错误降级处理 / Graceful error degradation
- ✅ 多层防护确保服务稳定 / Multi-layer protection ensures service stability
- ✅ 单个工具异常不影响整体服务 / Individual tool exceptions don't affect overall service

### 性能优化 / Performance Optimization
- ✅ Sonic JSON库支持（高性能序列化/反序列化）
- ✅ 结构体预热机制（消除首次请求延迟）
- ✅ 多种JSON库可选（Sonic、go-json、jsoniter、标准库）

## 技术栈 / Tech Stack

- **语言 / Language**: Go 1.25.5+
- **MCP SDK**: github.com/modelcontextprotocol/go-sdk v1.2.0 (官方SDK / Official SDK)
- **JSON库 / JSON Library**:
  - github.com/bytedance/sonic v1.14.2 (高性能 / High performance)
  - github.com/goccy/go-json v0.10.5 (备选 / Alternative)
  - github.com/json-iterator/go v1.1.12 (备选 / Alternative)
- **日志 / Logging**: go.uber.org/zap v1.27.1
- **测试 / Testing**: github.com/stretchr/testify v1.10.0

## 安装 / Installation

### 快速安装 / Quick Installation

#### 使用 uv (推荐) / Using uv (Recommended)

```bash
# 安装 uv (如果还没有安装) / Install uv (if not already installed)
curl -LsSf https://astral.sh/uv/install.sh | sh

# 使用 uv 安装 MCP Toolkit / Install MCP Toolkit using uv
uv tool install mcp-sandbox-toolkit

# 运行程序（两个命令都可以）/ Run the program (both commands work)
mcp-sandbox-toolkit --help
mcp-toolkit --help

# 或使用 uvx 直接运行（无需安装）/ Or use uvx to run directly (no installation needed)
uvx mcp-sandbox-toolkit --help
```

**配置 PATH (如果需要) / Configure PATH (if needed)**:

如果安装后无法直接运行 `mcp-toolkit` 命令，需要将 `~/.local/bin` (Linux/macOS) 添加到 PATH：

```bash
# Bash
echo 'export PATH="$HOME/.local/bin:$PATH"' >> ~/.bashrc
source ~/.bashrc

# Zsh
echo 'export PATH="$HOME/.local/bin:$PATH"' >> ~/.zshrc
source ~/.zshrc
```

详细的安装和配置说明请参考 [INSTALLATION.md](docs/INSTALLATION.md)

#### 使用安装脚本 / Using Installation Script

**Linux/macOS:**
```bash
curl -fsSL https://raw.githubusercontent.com/shibingli/mcp-toolkit/main/scripts/install.sh | bash
```

**Windows (PowerShell):**
```powershell
Invoke-WebRequest -Uri "https://raw.githubusercontent.com/shibingli/mcp-toolkit/main/scripts/install.ps1" -OutFile "install.ps1"
.\install.ps1
```

#### 手动下载 / Manual Download

从 [Releases 页面](https://github.com/shibingli/mcp-toolkit/releases) 下载对应平台的二进制文件。

Download the binary for your platform from the [Releases page](https://github.com/shibingli/mcp-toolkit/releases).

### 从源码编译 / Build from Source

```bash
# 克隆仓库 / Clone repository
git clone https://github.com/shibingli/mcp-toolkit.git
cd mcp-toolkit

# 安装依赖 / Install dependencies
go mod download

# 编译(使用sonic JSON库以获得最佳性能) / Build (using sonic JSON library for best performance)
go build -tags="sonic" -o mcp-toolkit main.go

# 或使用 Makefile / Or use Makefile
make build
```

更多安装方式请参考 [安装指南](docs/INSTALLATION.md)。

For more installation methods, see [Installation Guide](docs/INSTALLATION.md).

## 使用方法 / Usage

### 启动服务器 / Start Server

#### Stdio 传输 (默认) / Stdio Transport (Default)

```bash
# 使用默认沙箱目录和stdio传输 / Use default sandbox directory and stdio transport
./mcp-toolkit

# 指定自定义沙箱目录 / Specify custom sandbox directory
./mcp-toolkit -sandbox /path/to/sandbox
```

#### HTTP 传输 / HTTP Transport

```bash
# 使用HTTP传输 / Use HTTP transport
./mcp-toolkit -transport http

# 自定义HTTP配置 / Customize HTTP configuration
./mcp-toolkit -transport http -http-host 0.0.0.0 -http-port 8080 -sandbox /path/to/sandbox
```

#### SSE 传输 / SSE Transport

```bash
# 使用SSE传输 / Use SSE transport
./mcp-toolkit -transport sse

# 自定义SSE配置 / Customize SSE configuration
./mcp-toolkit -transport sse -sse-host 0.0.0.0 -sse-port 8081 -sandbox /path/to/sandbox
```

详细的传输方式说明请参考：[传输方式文档](docs/TRANSPORT.md)

For detailed transport documentation, see: [Transport Documentation](docs/TRANSPORT.md)

### JSON 结构体预热 / JSON Structure Preheating

程序启动时会自动预热所有注册的结构体（仅在使用Sonic时有效），以消除首次请求的延迟。

The program automatically preheats all registered structures at startup (only effective when using Sonic) to eliminate first request delays.

查看预热日志：

Check preheating logs:

```
{"level":"INFO","msg":"preheating JSON structures","json_library":"sonic"}
{"level":"INFO","msg":"JSON structures preheated successfully"}
```

### MCP 工具列表 / MCP Tools List

#### 1. create_file
创建新文件并写入内容 / Create a new file and write content

**参数 / Parameters:**
- `path` (必填 / required): 文件路径(相对于沙箱目录) / File path (relative to sandbox directory)
- `content` (可选 / optional): 文件内容 / File content

#### 2. create_directory
创建新目录 / Create a new directory

**参数 / Parameters:**
- `path` (必填 / required): 目录路径(相对于沙箱目录) / Directory path (relative to sandbox directory)

#### 3. read_file
读取文件内容 / Read file content

**参数 / Parameters:**
- `path` (必填 / required): 文件路径(相对于沙箱目录) / File path (relative to sandbox directory)

#### 4. write_file
写入或覆盖文件内容 / Write or overwrite file content

**参数 / Parameters:**
- `path` (必填 / required): 文件路径(相对于沙箱目录) / File path (relative to sandbox directory)
- `content` (必填 / required): 文件内容 / File content

#### 5. delete
删除文件或目录 / Delete file or directory

**参数 / Parameters:**
- `path` (必填 / required): 文件或目录路径(相对于沙箱目录) / File or directory path (relative to sandbox directory)

#### 6. copy
复制文件或目录 / Copy file or directory

**参数 / Parameters:**
- `source` (必填 / required): 源路径(相对于沙箱目录) / Source path (relative to sandbox directory)
- `destination` (必填 / required): 目标路径(相对于沙箱目录) / Destination path (relative to sandbox directory)

#### 7. move
移动或重命名文件或目录 / Move or rename file or directory

**参数 / Parameters:**
- `source` (必填 / required): 源路径(相对于沙箱目录) / Source path (relative to sandbox directory)
- `destination` (必填 / required): 目标路径(相对于沙箱目录) / Destination path (relative to sandbox directory)

#### 8. list_directory
列出目录中的文件和子目录 / List files and subdirectories in a directory

**参数 / Parameters:**
- `path` (必填 / required): 目录路径(相对于沙箱目录) / Directory path (relative to sandbox directory)

#### 9. search_files
根据文件名模式搜索文件 / Search files by filename pattern

**参数 / Parameters:**
- `path` (必填 / required): 搜索起始路径(相对于沙箱目录) / Search starting path (relative to sandbox directory)
- `pattern` (必填 / required): 文件名匹配模式(支持通配符*和?) / Filename pattern (supports wildcards * and ?)

#### 10. batch_delete
批量删除多个文件或目录 / Batch delete multiple files or directories

**参数 / Parameters:**
- `paths` (必填 / required): 要删除的文件或目录路径列表(相对于沙箱目录) / List of file or directory paths to delete (relative to sandbox directory)

#### 11. file_stat
获取文件或目录的详细信息 / Get detailed information about a file or directory

**参数 / Parameters:**
- `path` (必填 / required): 文件或目录路径(相对于沙箱目录) / File or directory path (relative to sandbox directory)

#### 12. file_exists
检查文件或目录是否存在 / Check if a file or directory exists

**参数 / Parameters:**
- `path` (必填 / required): 文件或目录路径(相对于沙箱目录) / File or directory path (relative to sandbox directory)

#### 13. get_current_time
获取当前系统时间 / Get current system time

**参数 / Parameters:** 无 / None

#### 14. execute_command
在沙箱目录内执行命令行命令 / Execute command line command within sandbox directory

**参数 / Parameters:**
- `command` (必填 / required): 要执行的命令 / Command to execute
- `args` (可选 / optional): 命令参数列表 / Command arguments list
- `work_dir` (可选 / optional): 工作目录(相对于沙箱根目录) / Working directory (relative to sandbox root)
- `timeout` (可选 / optional): 超时时间(秒),0表示使用默认值 / Timeout in seconds, 0 for default

#### 15. get_command_blacklist
获取命令和目录黑名单配置 / Get command and directory blacklist configuration

**参数 / Parameters:** 无 / None

#### 16. update_command_blacklist
更新命令和目录黑名单 / Update command and directory blacklist

**参数 / Parameters:**
- `commands` (可选 / optional): 要添加的黑名单命令列表 / Commands to add to blacklist
- `directories` (可选 / optional): 要添加的黑名单目录列表 / Directories to add to blacklist

#### 17. get_working_directory
获取当前工作目录 / Get current working directory

**参数 / Parameters:** 无 / None

#### 18. change_directory
切换当前工作目录 / Change current working directory

**参数 / Parameters:**
- `path` (必填 / required): 目标目录路径(相对于沙箱根目录) / Target directory path (relative to sandbox root)

#### 19. execute_command_async
异步执行命令,返回任务ID / Execute command asynchronously, returns task ID

**参数 / Parameters:**
- `command` (必填 / required): 要执行的命令 / Command to execute
- `args` (可选 / optional): 命令参数列表 / Command arguments list
- `work_dir` (可选 / optional): 工作目录 / Working directory
- `timeout` (可选 / optional): 超时时间(秒) / Timeout in seconds
- `environment` (可选 / optional): 环境变量 / Environment variables
- `permission_level` (可选 / optional): 权限级别 / Permission level
- `user` (可选 / optional): 执行用户 / Executing user

#### 20. get_command_task
获取异步命令任务状态 / Get async command task status

**参数 / Parameters:**
- `task_id` (必填 / required): 任务ID / Task ID

#### 21. cancel_command_task
取消正在执行的命令任务 / Cancel running command task

**参数 / Parameters:**
- `task_id` (必填 / required): 任务ID / Task ID

#### 22. get_command_history
获取命令执行历史记录 / Get command execution history

**参数 / Parameters:**
- `limit` (可选 / optional): 返回记录数量限制 / Limit of returned records
- `offset` (可选 / optional): 偏移量 / Offset
- `user` (可选 / optional): 按用户过滤 / Filter by user

#### 23. clear_command_history
清空命令执行历史记录 / Clear command execution history

**参数 / Parameters:** 无 / None

#### 24. set_permission_level
设置命令执行权限级别 / Set command execution permission level

**参数 / Parameters:**
- `level` (必填 / required): 权限级别(0-3) / Permission level (0-3)

#### 25. get_permission_level
获取当前权限级别 / Get current permission level

**参数 / Parameters:** 无 / None

#### 26. get_system_info
获取系统信息 / Get system information

获取完整的系统信息，包括操作系统、CPU、内存、GPU、网络接口等详细信息。
Get complete system information including OS, CPU, memory, GPU, network interfaces and more.

**参数 / Parameters:** 无 / None

**返回 / Returns:**
- `os`: 操作系统信息 / OS information (platform, architecture, hostname, uptime, etc.)
- `cpu`: CPU信息 / CPU information (model, cores, frequency, usage, etc.)
- `memory`: 内存信息 / Memory information (total, available, used, swap, etc.)
- `gpus`: GPU信息列表 / GPU information list (name, memory, temperature, utilization, etc.)
- `networks`: 网络接口信息列表 / Network interface list (name, MAC, IPs, speed, etc.)

## 文档 / Documentation

- [命令执行使用指南](docs/COMMAND_EXECUTION.md) - 详细的命令执行功能说明
- [命令执行高级功能](docs/COMMAND_ADVANCED_FEATURES.md) - 异步执行、历史记录、权限控制等
- [命令路径验证](docs/COMMAND_PATH_VALIDATION.md) - 路径参数验证机制
- [Recovery 功能文档](docs/RECOVERY.md) - Panic 恢复机制和稳定性保障

## 测试 / Testing

```bash
# 运行所有测试 / Run all tests
go test -v ./...

# 运行测试并生成覆盖率报告 / Run tests with coverage report
go test -v ./... -cover

# 生成详细的覆盖率报告 / Generate detailed coverage report
go test -v ./... -coverprofile=coverage.out
go tool cover -html=coverage.out
```

当前测试覆盖率 / Current test coverage:
- sandbox: **53.0%**
- client: **78.0%**
- transport: **72.0%**
- json: **86.1%**
- recovery: **100.0%**

## 项目结构 / Project Structure

```
mcp-toolkit/
├── main.go                              # 主程序入口 / Main entry point
├── go.mod                               # Go 模块定义 / Go module definition
├── go.sum                               # 依赖校验和 / Dependency checksums
├── README.md                            # 项目文档 / Project documentation
├── pkg/
│   ├── types/                           # 类型定义 / Type definitions
│   │   ├── common.go                    # 通用类型 / Common types
│   │   ├── file.go                      # 文件操作类型 / File operation types
│   │   ├── command.go                   # 命令执行类型 / Command execution types
│   │   ├── time.go                      # 时间类型 / Time types
│   │   ├── sysinfo.go                   # 系统信息类型 / System info types
│   │   ├── schema.go                    # JSON Schema 定义 / JSON Schema definitions
│   │   └── constants.go                 # 常量定义 / Constants
│   └── utils/
│       └── json/                        # JSON 工具 / JSON utilities
│           ├── json.go                  # JSON 编解码 / JSON encoding/decoding
│           └── pretouch.go              # 结构体预热 / Struct pretouch
└── internal/
    └── services/
        └── sandbox/                     # 沙箱服务 / Sandbox service
            ├── service.go               # 核心服务实现 / Core service implementation
            ├── service_test.go          # 服务测试 / Service tests
            ├── sysinfo.go               # 系统信息获取 / System info retrieval
            ├── sysinfo_test.go          # 系统信息测试 / System info tests
            ├── mcp_tools.go             # MCP 工具注册 / MCP tools registration
            ├── mcp_tools_test.go        # 工具注册测试 / Tools registration tests
            ├── mcp_handlers.go          # MCP 处理器 / MCP handlers
            └── mcp_handlers_test.go     # 处理器测试 / Handlers tests
```

## 测试 / Testing

### 完整功能测试 / Complete Functionality Test

项目提供了完整的客户端测试工具,可以自动测试所有26个MCP工具。

The project provides a complete client testing tool that automatically tests all 26 MCP tools.

#### 运行测试 / Run Tests

**Linux/macOS:**
```bash
# 编译服务器和客户端 / Build server and client
go build -tags="sonic" -o mcp-toolkit .
go build -tags="sonic" -o mcp-toolkit-client ./cmd/client
```

**Windows:**
```bash
# 编译服务器和客户端 / Build server and client
go build -tags="sonic" -o mcp-toolkit.exe .
go build -tags="sonic" -o mcp-toolkit-client.exe ./cmd/client
```

#### 手动测试 / Manual Testing

```bash
# 1. 启动服务器 / Start server
./mcp-toolkit -transport http -http-port 8080

# 2. 在新终端运行客户端测试 / Run client tests in new terminal
./mcp-toolkit-client

# 3. 使用详细日志 / Use verbose logging
./mcp-toolkit-client -verbose
```

#### 测试覆盖 / Test Coverage

✅ **26个MCP工具** / 26 MCP Tools
- 文件操作 (11个) / File Operations (11)
- 目录操作 (2个) / Directory Operations (2)
- 命令执行 (3个) / Command Execution (3)
- 异步操作 (3个) / Async Operations (3)
- 命令历史 (2个) / Command History (2)
- 权限管理 (2个) / Permission Management (2)
- 系统工具 (3个) / System Tools (3)


### 单元测试 / Unit Tests

```bash
# 运行所有单元测试 / Run all unit tests
go test -v ./...

# 运行特定包的测试 / Run tests for specific package
go test -v ./internal/services/sandbox/

# 查看测试覆盖率 / View test coverage
go test -cover ./...
```

## 许可证 / License

本项目采用 Apache License 2.0 许可证。详情请参阅 [LICENSE](LICENSE) 文件。

This project is licensed under the Apache License 2.0. See the [LICENSE](LICENSE) file for details.

```
Copyright 2024 MCP Toolkit Authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
```

### 第三方依赖 / Third-Party Dependencies

本项目使用了多个开源库，详情请参阅 [NOTICE](NOTICE) 文件。

This project uses several open-source libraries. See the [NOTICE](NOTICE) file for details.
