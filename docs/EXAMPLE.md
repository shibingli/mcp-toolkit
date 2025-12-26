# MCP 文件系统服务使用示例 / MCP Filesystem Service Usage Examples

## 基本使用 / Basic Usage

### 启动服务器 / Start Server

```bash
# 使用默认沙箱目录 / Use default sandbox directory
./mcp-toolkit

# 指定自定义沙箱目录 / Specify custom sandbox directory
./mcp-toolkit -sandbox /path/to/sandbox
```

## MCP 工具调用示例 / MCP Tool Call Examples

### 1. 创建文件 / Create File

```json
{
  "method": "tools/call",
  "params": {
    "name": "create_file",
    "arguments": {
      "path": "hello.txt",
      "content": "Hello, World!"
    }
  }
}
```

**响应 / Response:**
```json
{
  "success": true,
  "message": "文件创建成功 / File created successfully"
}
```

### 2. 读取文件 / Read File

```json
{
  "method": "tools/call",
  "params": {
    "name": "read_file",
    "arguments": {
      "path": "hello.txt"
    }
  }
}
```

**响应 / Response:**
```json
{
  "content": "Hello, World!"
}
```

### 3. 创建目录 / Create Directory

```json
{
  "method": "tools/call",
  "params": {
    "name": "create_directory",
    "arguments": {
      "path": "mydir"
    }
  }
}
```

### 4. 列出目录内容 / List Directory

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

**响应 / Response:**
```json
{
  "files": [
    {
      "name": "hello.txt",
      "path": "hello.txt",
      "size": 13,
      "is_dir": false,
      "mode": "-rw-r--r--",
      "mod_time": "2025-01-01T12:00:00Z"
    },
    {
      "name": "mydir",
      "path": "mydir",
      "size": 0,
      "is_dir": true,
      "mode": "drwxr-xr-x",
      "mod_time": "2025-01-01T12:00:00Z"
    }
  ]
}
```

### 5. 搜索文件 / Search Files

```json
{
  "method": "tools/call",
  "params": {
    "name": "search_files",
    "arguments": {
      "path": ".",
      "pattern": "*.txt"
    }
  }
}
```

### 6. 复制文件 / Copy File

```json
{
  "method": "tools/call",
  "params": {
    "name": "copy",
    "arguments": {
      "source": "hello.txt",
      "destination": "hello_copy.txt"
    }
  }
}
```

### 7. 移动文件 / Move File

```json
{
  "method": "tools/call",
  "params": {
    "name": "move",
    "arguments": {
      "source": "hello.txt",
      "destination": "mydir/hello.txt"
    }
  }
}
```

### 8. 删除文件 / Delete File

```json
{
  "method": "tools/call",
  "params": {
    "name": "delete",
    "arguments": {
      "path": "hello_copy.txt"
    }
  }
}
```

### 9. 批量删除 / Batch Delete

```json
{
  "method": "tools/call",
  "params": {
    "name": "batch_delete",
    "arguments": {
      "paths": ["file1.txt", "file2.txt", "file3.txt"]
    }
  }
}
```

### 10. 获取文件信息 / Get File Info

```json
{
  "method": "tools/call",
  "params": {
    "name": "file_stat",
    "arguments": {
      "path": "hello.txt"
    }
  }
}
```

### 11. 检查文件是否存在 / Check File Exists

```json
{
  "method": "tools/call",
  "params": {
    "name": "file_exists",
    "arguments": {
      "path": "hello.txt"
    }
  }
}
```

**响应 / Response:**
```json
{
  "exists": true
}
```

### 12. 获取当前时间 / Get Current Time

```json
{
  "method": "tools/call",
  "params": {
    "name": "get_current_time",
    "arguments": {}
  }
}
```

**响应 / Response:**
```json
{
  "time": "2025-01-01T12:00:00Z",
  "time_zone": "UTC",
  "unix": 1704110400
}
```

## 安全特性 / Security Features

### 沙箱限制 / Sandbox Restriction

所有文件操作都限制在指定的沙箱目录内，任何试图访问沙箱外部文件的操作都会被拒绝。

All file operations are restricted to the specified sandbox directory. Any attempt to access files outside the sandbox will be rejected.

### 路径遍历保护 / Path Traversal Protection

系统会自动检测和阻止路径遍历攻击（如使用 `../` 访问上级目录）。

The system automatically detects and blocks path traversal attacks (such as using `../` to access parent directories).

## 错误处理 / Error Handling

所有工具调用在发生错误时都会返回适当的错误信息：

All tool calls return appropriate error messages when errors occur:

```json
{
  "error": {
    "code": -32000,
    "message": "文件不存在 / File not found"
  }
}
```

