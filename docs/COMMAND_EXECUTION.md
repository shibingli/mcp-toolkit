# 命令执行工具使用指南 / Command Execution Tool Guide

## 概述 / Overview

命令执行工具允许在沙箱目录内安全地执行命令行命令。该工具具有以下安全特性：

The command execution tool allows safe execution of command-line commands within the sandbox directory. It has the following security features:

- **沙箱隔离** / **Sandbox Isolation**: 所有命令只能在沙箱目录内执行
- **黑名单机制** / **Blacklist Mechanism**: 危险命令和系统目录被默认禁止
- **工作目录管理** / **Working Directory Management**: 支持在沙箱内切换工作目录
- **跨平台支持** / **Cross-platform Support**: 兼容 Linux、Windows、macOS

## 可用工具 / Available Tools

### 1. execute_command - 执行命令

在沙箱目录内执行命令行命令。

Execute command line command within sandbox directory.

**请求参数 / Request Parameters:**

```json
{
  "command": "ls",           // 要执行的命令 / Command to execute
  "args": ["-la"],          // 命令参数(可选) / Command arguments (optional)
  "work_dir": ".",          // 工作目录(相对于沙箱根目录) / Working directory (relative to sandbox root)
  "timeout": 60             // 超时时间(秒),0表示使用默认值 / Timeout in seconds, 0 for default
}
```

**响应示例 / Response Example:**

```json
{
  "success": true,
  "exit_code": 0,
  "stdout": "total 8\ndrwxr-xr-x  3 user  staff   96 Jan  1 12:00 .\n",
  "stderr": "",
  "message": "command executed successfully"
}
```

**使用示例 / Usage Examples:**

```bash
# Linux/macOS - 列出文件
{
  "command": "ls",
  "args": ["-la"],
  "work_dir": "."
}

# Windows - 列出文件
{
  "command": "cmd",
  "args": ["/c", "dir"],
  "work_dir": "."
}

# 创建文件
{
  "command": "touch",
  "args": ["test.txt"],
  "work_dir": "."
}
```

### 2. get_command_blacklist - 获取黑名单

获取当前的命令和目录黑名单配置。

Get current command and directory blacklist configuration.

**请求参数 / Request Parameters:**

```json
{}
```

**响应示例 / Response Example:**

```json
{
  "commands": [
    "shutdown", "reboot", "rm", "format", "dd"
  ],
  "directories": [
    "/bin", "/sbin", "C:\\Windows"
  ],
  "system_directories": [
    "/bin", "/sbin", "/usr/bin", "/etc"
  ]
}
```

### 3. update_command_blacklist - 更新黑名单

向黑名单中添加自定义命令或目录。

Add custom commands or directories to the blacklist.

**请求参数 / Request Parameters:**

```json
{
  "commands": ["custom_cmd"],      // 要添加的命令 / Commands to add
  "directories": ["/custom/dir"]   // 要添加的目录 / Directories to add
}
```

**响应示例 / Response Example:**

```json
{
  "success": true,
  "message": "blacklist updated successfully"
}
```

### 4. get_working_directory - 获取工作目录

获取当前工作目录(相对于沙箱根目录)。

Get current working directory (relative to sandbox root).

**请求参数 / Request Parameters:**

```json
{}
```

**响应示例 / Response Example:**

```json
{
  "work_dir": "."
}
```

### 5. change_directory - 切换工作目录

切换当前工作目录(只能在沙箱内切换)。

Change current working directory (can only change within sandbox).

**请求参数 / Request Parameters:**

```json
{
  "path": "subdir"  // 目标目录路径 / Target directory path
}
```

**响应示例 / Response Example:**

```json
{
  "success": true,
  "message": "directory changed successfully"
}
```

## 安全限制 / Security Restrictions

### 路径参数验证 / Path Argument Validation

对于某些命令（如 `rm`、`rmdir`、`del` 等），系统会自动解析命令参数中的路径，并验证这些路径是否在沙箱目录内。

For certain commands (such as `rm`, `rmdir`, `del`, etc.), the system automatically parses paths in command arguments and validates whether these paths are within the sandbox directory.

**验证的命令 / Validated Commands:**
- `rm` - 删除文件 / Delete files
- `rmdir` - 删除目录 / Delete directories
- `del` - Windows删除命令 / Windows delete command
- `erase` - Windows删除命令 / Windows delete command
- `rd` - Windows删除目录 / Windows remove directory

**验证规则 / Validation Rules:**
1. 解析命令参数中的所有路径（跳过选项参数如 `-rf`）
2. 将相对路径转换为绝对路径
3. 验证路径是否在沙箱目录内
4. 验证路径是否在黑名单目录中

**示例 / Examples:**

```bash
# ✅ 允许 - 删除沙箱内文件
rm test.txt
rm ./subdir/file.txt

# ✅ 允许 - 带选项参数
rm -rf temp_dir

# ❌ 拒绝 - 尝试删除沙箱外文件
rm ../../etc/passwd
rm /etc/passwd

# ❌ 拒绝 - 尝试删除系统目录
rm -rf /bin
```

### 默认黑名单命令 / Default Blacklisted Commands

以下命令默认被禁止执行：

The following commands are prohibited by default:

- **系统管理** / **System Management**: `shutdown`, `reboot`, `halt`, `poweroff`, `init`
- **用户管理** / **User Management**: `useradd`, `userdel`, `usermod`, `passwd`
- **磁盘管理** / **Disk Management**: `fdisk`, `mkfs`, `mount`, `umount`, `dd`
- **包管理** / **Package Management**: `apt`, `yum`, `dnf`, `rpm`, `dpkg`
- **危险命令** / **Dangerous Commands**: `format`
- **网络配置** / **Network Configuration**: `ifconfig`, `ip`, `iptables`

**注意**: `rm` 和 `rmdir` 命令不在黑名单中，但会对其参数进行严格的路径验证，确保只能删除沙箱目录内的文件和目录。

**Note**: `rm` and `rmdir` commands are not blacklisted, but their arguments undergo strict path validation to ensure only files and directories within the sandbox can be deleted.

### 默认黑名单目录 / Default Blacklisted Directories

以下目录禁止作为工作目录：

The following directories are prohibited as working directories:

- **Linux/Unix**: `/bin`, `/sbin`, `/etc`, `/boot`, `/dev`, `/proc`, `/sys`, `/root`
- **Windows**: `C:\Windows`, `C:\Program Files`
- **macOS**: `/System`, `/Library`, `/Applications`

## 完整使用示例 / Complete Usage Example

### 场景：在沙箱中执行一系列命令 / Scenario: Execute a series of commands in sandbox

```json
// 1. 获取当前工作目录
{
  "method": "tools/call",
  "params": {
    "name": "get_working_directory",
    "arguments": {}
  }
}
// 响应: {"work_dir": "."}

// 2. 创建一个新目录
{
  "method": "tools/call",
  "params": {
    "name": "create_directory",
    "arguments": {
      "path": "myproject"
    }
  }
}

// 3. 切换到新目录
{
  "method": "tools/call",
  "params": {
    "name": "change_directory",
    "arguments": {
      "path": "myproject"
    }
  }
}

// 4. 在新目录中执行命令
{
  "method": "tools/call",
  "params": {
    "name": "execute_command",
    "arguments": {
      "command": "echo",
      "args": ["Hello from myproject!"],
      "work_dir": "."
    }
  }
}

// 5. 查看黑名单配置
{
  "method": "tools/call",
  "params": {
    "name": "get_command_blacklist",
    "arguments": {}
  }
}

// 6. 添加自定义黑名单
{
  "method": "tools/call",
  "params": {
    "name": "update_command_blacklist",
    "arguments": {
      "commands": ["dangerous_cmd"],
      "directories": ["/sensitive/path"]
    }
  }
}
```

## 注意事项 / Notes

1. **命令执行超时** / **Command Timeout**: 默认超时为300秒，最大3600秒
2. **工作目录限制** / **Working Directory Restriction**: 只能在沙箱目录内切换
3. **黑名单不可删除** / **Blacklist Cannot Be Removed**: 只能添加，不能删除黑名单项
4. **命令参数** / **Command Arguments**: 建议使用args数组传递参数，而不是拼接到command中
5. **跨平台差异** / **Cross-platform Differences**: 注意不同操作系统的命令差异

## 错误处理 / Error Handling

### 常见错误 / Common Errors

1. **命令在黑名单中** / **Command Blacklisted**
   ```json
   {
     "error": "command is blacklisted"
   }
   ```

2. **目录在黑名单中** / **Directory Blacklisted**
   ```json
   {
     "error": "directory is blacklisted"
   }
   ```

3. **沙箱违规** / **Sandbox Violation**
   ```json
   {
     "error": "path is outside sandbox directory"
   }
   ```

4. **命令执行超时** / **Command Timeout**
   ```json
   {
     "success": false,
     "exit_code": -1,
     "message": "command execution failed: context deadline exceeded"
   }
   ```

## 最佳实践 / Best Practices

1. **始终检查退出码** / **Always Check Exit Code**: 即使success为true，也要检查exit_code
2. **合理设置超时** / **Set Reasonable Timeout**: 根据命令特性设置合适的超时时间
3. **使用相对路径** / **Use Relative Paths**: 工作目录使用相对于沙箱根目录的路径
4. **捕获输出** / **Capture Output**: 检查stdout和stderr以获取详细信息
5. **定期更新黑名单** / **Update Blacklist Regularly**: 根据安全需求更新黑名单配置

