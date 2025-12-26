# 快速启动指南 / Quick Start Guide

## 安装 / Installation

### 1. 克隆项目 / Clone Project
```bash
git clone <repository-url>
cd mcp-toolkit
```

### 2. 安装依赖 / Install Dependencies
```bash
go mod download
```

### 3. 编译 / Build

#### 使用Sonic (推荐，最高性能) / Using Sonic (Recommended, Best Performance)
```bash
go build -tags="sonic" -o mcp-toolkit main.go
```

#### 使用标准库 / Using Standard Library
```bash
go build -o mcp-toolkit main.go
```

---

## 运行 / Run

### Stdio模式 (默认) / Stdio Mode (Default)
```bash
# Windows
mcp-toolkit.exe

# Linux/macOS
./mcp-toolkit
```

### HTTP模式 / HTTP Mode
```bash
# 使用默认配置 (127.0.0.1:8080)
mcp-toolkit.exe -transport http

# 自定义配置
mcp-toolkit.exe -transport http -http-host 0.0.0.0 -http-port 8080
```

### SSE模式 / SSE Mode
```bash
# 使用默认配置 (127.0.0.1:8081)
mcp-toolkit.exe -transport sse

# 自定义配置
mcp-toolkit.exe -transport sse -sse-host 0.0.0.0 -sse-port 8081
```

---

## 测试 / Testing

### 运行所有测试 / Run All Tests
```bash
go test -v ./...
```

### 查看测试覆盖率 / View Test Coverage
```bash
go test -v ./... -cover
```

### 生成覆盖率报告 / Generate Coverage Report
```bash
go test -v ./... -coverprofile=coverage.out
go tool cover -html=coverage.out
```

---

## 基本使用 / Basic Usage

### 创建文件 / Create File
```json
{
  "method": "tools/call",
  "params": {
    "name": "create_file",
    "arguments": {
      "path": "test.txt",
      "content": "Hello, World!"
    }
  }
}
```

### 读取文件 / Read File
```json
{
  "method": "tools/call",
  "params": {
    "name": "read_file",
    "arguments": {
      "path": "test.txt"
    }
  }
}
```

### 列出目录 / List Directory
```json
{
  "method": "tools/call",
  "params": {
    "name": "list_directory",
    "arguments": {
      "path": "."
    }
  }
}
```

---

## 配置 / Configuration

### 命令行参数 / Command Line Arguments

| 参数 | 说明 | 默认值 |
|------|------|--------|
| `-transport` | 传输类型 (stdio/http/sse) | stdio |
| `-sandbox` | 沙箱目录路径 | ./sandbox |
| `-http-host` | HTTP监听地址 | 127.0.0.1 |
| `-http-port` | HTTP监听端口 | 8080 |
| `-sse-host` | SSE监听地址 | 127.0.0.1 |
| `-sse-port` | SSE监听端口 | 8081 |

### 示例 / Examples

```bash
# 本地开发
mcp-toolkit.exe -transport stdio -sandbox ./dev-sandbox

# 远程HTTP服务
mcp-toolkit.exe -transport http -http-host 0.0.0.0 -http-port 8080

# SSE实时推送
mcp-toolkit.exe -transport sse -sse-host 0.0.0.0 -sse-port 8081
```

---

## 健康检查 / Health Check

### HTTP模式 / HTTP Mode
```bash
curl http://127.0.0.1:8080/health
```

### SSE模式 / SSE Mode
```bash
curl http://127.0.0.1:8081/health
```

---

## 故障排查 / Troubleshooting

### 编译错误 / Build Errors

**问题**: 找不到sonic包
```
Solution: 确保使用了正确的build tag
go build -tags="sonic" -o mcp-toolkit main.go
```

**问题**: 端口被占用
```
Solution: 更改端口或停止占用端口的进程
mcp-toolkit.exe -transport http -http-port 8090
```

### 运行时错误 / Runtime Errors

**问题**: 沙箱目录权限错误
```
Solution: 确保有足够的权限创建和访问沙箱目录
chmod 755 ./sandbox  # Linux/macOS
```

**问题**: 路径遍历错误
```
Solution: 确保所有路径都在沙箱目录内
不要使用 "../" 等路径遍历符号
```

