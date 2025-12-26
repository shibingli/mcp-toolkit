# 命令执行功能演示 / Command Execution Feature Demo

## 演示场景 / Demo Scenario

本演示展示如何使用命令执行工具在沙箱环境中安全地执行命令。

This demo shows how to use the command execution tool to safely execute commands in a sandbox environment.

## 前置条件 / Prerequisites

1. 启动MCP服务器 / Start MCP server:
   ```bash
   ./mcp_server -sandbox ./demo_sandbox
   ```

2. 确保沙箱目录已创建 / Ensure sandbox directory is created

## 演示步骤 / Demo Steps

### 步骤 1: 查看当前工作目录 / Step 1: Check Current Working Directory

**请求 / Request:**
```json
{
  "method": "tools/call",
  "params": {
    "name": "get_working_directory",
    "arguments": {}
  }
}
```

**预期响应 / Expected Response:**
```json
{
  "work_dir": "."
}
```

### 步骤 2: 创建项目目录 / Step 2: Create Project Directory

**请求 / Request:**
```json
{
  "method": "tools/call",
  "params": {
    "name": "create_directory",
    "arguments": {
      "path": "demo_project"
    }
  }
}
```

### 步骤 3: 切换到项目目录 / Step 3: Change to Project Directory

**请求 / Request:**
```json
{
  "method": "tools/call",
  "params": {
    "name": "change_directory",
    "arguments": {
      "path": "demo_project"
    }
  }
}
```

### 步骤 4: 执行命令创建文件 / Step 4: Execute Command to Create File

**Linux/macOS:**
```json
{
  "method": "tools/call",
  "params": {
    "name": "execute_command",
    "arguments": {
      "command": "touch",
      "args": ["demo.txt"],
      "work_dir": "."
    }
  }
}
```

**Windows:**
```json
{
  "method": "tools/call",
  "params": {
    "name": "execute_command",
    "arguments": {
      "command": "cmd",
      "args": ["/c", "type", "nul", ">", "demo.txt"],
      "work_dir": "."
    }
  }
}
```

### 步骤 5: 列出目录内容 / Step 5: List Directory Contents

**Linux/macOS:**
```json
{
  "method": "tools/call",
  "params": {
    "name": "execute_command",
    "arguments": {
      "command": "ls",
      "args": ["-la"],
      "work_dir": "."
    }
  }
}
```

**Windows:**
```json
{
  "method": "tools/call",
  "params": {
    "name": "execute_command",
    "arguments": {
      "command": "cmd",
      "args": ["/c", "dir"],
      "work_dir": "."
    }
  }
}
```

### 步骤 6: 查看黑名单配置 / Step 6: View Blacklist Configuration

**请求 / Request:**
```json
{
  "method": "tools/call",
  "params": {
    "name": "get_command_blacklist",
    "arguments": {}
  }
}
```

**预期响应 / Expected Response:**
```json
{
  "commands": [
    "shutdown", "reboot", "rm", "format", "dd", ...
  ],
  "directories": [
    "/bin", "/sbin", "C:\\Windows", ...
  ],
  "system_directories": [
    "/bin", "/sbin", "/usr/bin", ...
  ]
}
```

### 步骤 7: 测试黑名单保护 / Step 7: Test Blacklist Protection

**请求 / Request:**
```json
{
  "method": "tools/call",
  "params": {
    "name": "execute_command",
    "arguments": {
      "command": "shutdown",
      "args": [],
      "work_dir": "."
    }
  }
}
```

**预期响应 / Expected Response:**
```json
{
  "error": "command is blacklisted"
}
```

### 步骤 8: 添加自定义黑名单 / Step 8: Add Custom Blacklist

**请求 / Request:**
```json
{
  "method": "tools/call",
  "params": {
    "name": "update_command_blacklist",
    "arguments": {
      "commands": ["custom_dangerous_cmd"],
      "directories": []
    }
  }
}
```

## 安全测试 / Security Tests

### 测试 1: 路径遍历保护 / Test 1: Path Traversal Protection

**请求 / Request:**
```json
{
  "method": "tools/call",
  "params": {
    "name": "change_directory",
    "arguments": {
      "path": "../../etc"
    }
  }
}
```

**预期结果 / Expected Result:** 错误 - 沙箱违规

### 测试 2: 系统目录保护 / Test 2: System Directory Protection

**请求 / Request:**
```json
{
  "method": "tools/call",
  "params": {
    "name": "change_directory",
    "arguments": {
      "path": "/bin"
    }
  }
}
```

**预期结果 / Expected Result:** 错误 - 目录在黑名单中

## 总结 / Summary

本演示展示了命令执行工具的以下特性：

This demo demonstrates the following features of the command execution tool:

✅ 工作目录管理 / Working directory management
✅ 安全的命令执行 / Safe command execution
✅ 黑名单保护机制 / Blacklist protection mechanism
✅ 沙箱隔离 / Sandbox isolation
✅ 跨平台兼容性 / Cross-platform compatibility

