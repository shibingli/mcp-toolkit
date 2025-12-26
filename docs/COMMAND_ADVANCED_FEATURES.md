# 命令执行高级功能 / Command Execution Advanced Features

## 概述 / Overview

本文档介绍命令执行工具的高级功能，包括：
1. 命令执行历史记录
2. 异步命令执行
3. 权限级别控制
4. 环境变量配置
5. 审计日志

This document introduces advanced features of the command execution tool, including:
1. Command execution history
2. Asynchronous command execution
3. Permission level control
4. Environment variable configuration
5. Audit logging

## 1. 命令执行历史记录 / Command Execution History

### 功能说明 / Feature Description

系统自动记录所有命令的执行历史，包括命令、参数、执行时间、结果等详细信息。

The system automatically records the execution history of all commands, including command, arguments, execution time, results, and other detailed information.

### 可用工具 / Available Tools

#### get_command_history - 获取命令历史

**请求参数 / Request Parameters:**
```json
{
  "limit": 100,      // 返回记录数量限制(默认100,最大100) / Limit of returned records (default 100, max 100)
  "offset": 0,       // 偏移量(默认0) / Offset (default 0)
  "user": ""         // 按用户过滤(可选) / Filter by user (optional)
}
```

**响应示例 / Response Example:**
```json
{
  "history": [
    {
      "id": "uuid-1234",
      "command": "ls",
      "args": ["-la"],
      "work_dir": ".",
      "start_time": "2024-01-01T12:00:00Z",
      "end_time": "2024-01-01T12:00:01Z",
      "duration": 1000,
      "exit_code": 0,
      "success": true,
      "user": "admin",
      "permission_level": 1
    }
  ],
  "total": 1
}
```

#### clear_command_history - 清空命令历史

**请求参数 / Request Parameters:**
```json
{}
```

**响应示例 / Response Example:**
```json
{
  "success": true,
  "message": "operation completed successfully"
}
```

### 历史记录限制 / History Limits

- 最多保留1000条历史记录
- 超过限制时自动删除最旧的记录
- 支持按用户过滤
- 支持分页查询

## 2. 异步命令执行 / Asynchronous Command Execution

### 功能说明 / Feature Description

对于长时间运行的命令，可以使用异步执行模式。系统会立即返回任务ID，可以通过任务ID查询执行状态和结果。

For long-running commands, you can use asynchronous execution mode. The system returns a task ID immediately, which can be used to query execution status and results.

### 可用工具 / Available Tools

#### execute_command_async - 异步执行命令

**请求参数 / Request Parameters:**
```json
{
  "command": "long_running_command",
  "args": ["arg1", "arg2"],
  "work_dir": ".",
  "timeout": 3600,
  "environment": {
    "VAR1": "value1"
  },
  "permission_level": 1,
  "user": "admin"
}
```

**响应示例 / Response Example:**
```json
{
  "task_id": "task-uuid-1234",
  "message": "command task created successfully"
}
```

#### get_command_task - 获取任务状态

**请求参数 / Request Parameters:**
```json
{
  "task_id": "task-uuid-1234"
}
```

**响应示例 / Response Example:**
```json
{
  "task": {
    "id": "task-uuid-1234",
    "command": "long_running_command",
    "args": ["arg1", "arg2"],
    "status": "completed",
    "start_time": "2024-01-01T12:00:00Z",
    "end_time": "2024-01-01T12:05:00Z",
    "exit_code": 0,
    "stdout": "command output",
    "stderr": "",
    "user": "admin"
  }
}
```

#### cancel_command_task - 取消任务

**请求参数 / Request Parameters:**
```json
{
  "task_id": "task-uuid-1234"
}
```

### 任务状态 / Task Status

- `pending`: 等待执行
- `running`: 正在执行
- `completed`: 执行完成
- `failed`: 执行失败
- `cancelled`: 已取消

## 3. 权限级别控制 / Permission Level Control

### 权限级别 / Permission Levels

| 级别 / Level | 名称 / Name | 说明 / Description |
|-------------|------------|-------------------|
| 0 | ReadOnly | 只读权限 - 只能执行查询类命令 |
| 1 | Standard | 标准权限 - 可以执行大部分命令(默认) |
| 2 | Elevated | 提升权限 - 可以执行所有非黑名单命令 |
| 3 | Admin | 管理员权限 - 可以执行所有命令 |

### 可用工具 / Available Tools

#### set_permission_level - 设置权限级别

**请求参数 / Request Parameters:**
```json
{
  "level": 1
}
```

#### get_permission_level - 获取权限级别

**请求参数 / Request Parameters:**
```json
{}
```

**响应示例 / Response Example:**
```json
{
  "level": 1
}
```

### 权限限制 / Permission Restrictions

**只读权限允许的命令:**
- 查询命令: `ls`, `dir`, `cat`, `type`, `echo`, `pwd`, `find`, `grep`
- 信息命令: `stat`, `file`, `which`, `where`, `whoami`

**标准权限禁止的命令:**
- 系统管理: `chmod`, `chown`, `sudo`, `su`
- 进程管理: `kill`, `killall`

## 4. 环境变量配置 / Environment Variable Configuration

### 功能说明 / Feature Description

可以为命令执行设置自定义环境变量。

Custom environment variables can be set for command execution.

### 使用示例 / Usage Example

```json
{
  "command": "printenv",
  "args": ["MY_VAR"],
  "environment": {
    "MY_VAR": "my_value",
    "PATH": "/custom/path:$PATH"
  }
}
```

## 5. 审计日志 / Audit Logging

### 功能说明 / Feature Description

所有命令执行都会记录详细的审计日志，包括：
- 命令和参数
- 执行用户
- 权限级别
- 执行时间和时长
- 执行结果

All command executions are logged with detailed audit information, including:
- Command and arguments
- Executing user
- Permission level
- Execution time and duration
- Execution result

### 日志格式 / Log Format

```json
{
  "level": "INFO",
  "timestamp": "2024-01-01T12:00:00Z",
  "logger": "audit",
  "msg": "command executed",
  "id": "uuid-1234",
  "command": "ls",
  "args": ["-la"],
  "user": "admin",
  "permission_level": 1,
  "exit_code": 0,
  "success": true,
  "duration_ms": 1000
}
```

## 最佳实践 / Best Practices

1. **使用异步执行**: 对于预计运行时间超过10秒的命令，使用异步执行
2. **设置合适的权限**: 根据实际需求设置最小权限级别
3. **定期清理历史**: 定期清理不需要的历史记录
4. **监控审计日志**: 定期检查审计日志以发现异常行为
5. **使用环境变量**: 通过环境变量传递敏感信息，而不是命令参数

