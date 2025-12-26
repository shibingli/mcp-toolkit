# 命令路径参数验证机制 / Command Path Argument Validation Mechanism

## 概述 / Overview

为了在保证安全的前提下允许使用 `rm`、`rmdir` 等文件操作命令，系统实现了智能的路径参数验证机制。该机制会自动解析命令参数中的路径，并确保所有操作都限制在沙箱目录内。

To allow the use of file operation commands like `rm` and `rmdir` while ensuring security, the system implements an intelligent path argument validation mechanism. This mechanism automatically parses paths in command arguments and ensures all operations are restricted within the sandbox directory.

## 验证的命令 / Validated Commands

以下命令会进行路径参数验证：

The following commands undergo path argument validation:

| 命令 / Command | 平台 / Platform | 说明 / Description |
|---------------|----------------|-------------------|
| `rm`          | Linux/macOS    | 删除文件 / Delete files |
| `rmdir`       | Linux/macOS    | 删除目录 / Delete directories |
| `del`         | Windows        | 删除文件 / Delete files |
| `erase`       | Windows        | 删除文件 / Delete files |
| `rd`          | Windows        | 删除目录 / Delete directories |
| `remove`      | 通用 / Generic  | 删除操作 / Remove operation |

## 验证流程 / Validation Process

### 1. 命令识别 / Command Identification

系统首先识别命令是否需要路径验证：

```go
// 提取命令基本名称
cmdName := filepath.Base(command)
cmdName = strings.ToLower(strings.TrimSuffix(cmdName, filepath.Ext(cmdName)))

// 检查是否在路径敏感命令列表中
if pathSensitiveCommands[cmdName] {
    // 进行路径验证
}
```

### 2. 参数解析 / Argument Parsing

遍历所有命令参数，识别路径参数：

```go
for _, arg := range args {
    // 跳过选项参数（以 - 开头）
    if strings.HasPrefix(arg, "-") {
        continue
    }
    
    // 解析路径参数
    // ...
}
```

### 3. 路径解析 / Path Resolution

将相对路径转换为绝对路径：

```go
var targetPath string
if filepath.IsAbs(arg) {
    targetPath = arg
} else {
    targetPath = filepath.Join(workDir, arg)
}
targetPath = filepath.Clean(targetPath)
```

### 4. 沙箱验证 / Sandbox Validation

验证路径是否在沙箱目录内：

```go
sandboxDir := filepath.Clean(s.sandboxDir)
if !strings.HasPrefix(targetPath, sandboxDir) {
    return error // 路径在沙箱外
}
```

### 5. 黑名单检查 / Blacklist Check

验证路径是否在黑名单目录中：

```go
if s.isDirectoryBlacklisted(targetPath) {
    return error // 路径在黑名单中
}
```

## 使用示例 / Usage Examples

### ✅ 允许的操作 / Allowed Operations

```bash
# 删除沙箱内的文件
rm test.txt
rm ./subdir/file.txt
rm subdir/file1.txt subdir/file2.txt

# 删除沙箱内的目录
rmdir empty_dir
rm -rf temp_directory

# 带选项的删除操作
rm -f file.txt
rm -rf directory/
```

### ❌ 拒绝的操作 / Rejected Operations

```bash
# 尝试删除沙箱外的文件（相对路径）
rm ../../etc/passwd
# 错误: path is outside sandbox directory

# 尝试删除系统文件（绝对路径）
rm /etc/passwd
# 错误: path is outside sandbox directory

# 尝试删除系统目录
rm -rf /bin
# 错误: path is in blacklisted directory

# Windows - 尝试删除系统文件
del C:\Windows\System32\notepad.exe
# 错误: path is outside sandbox directory
```

## 跨平台支持 / Cross-platform Support

### Linux/macOS

```bash
# 删除文件
rm file.txt
rm -f file.txt
rm -rf directory/

# 删除目录
rmdir empty_dir
```

### Windows

```bash
# 删除文件
del file.txt
erase file.txt

# 删除目录
rd empty_dir
rd /s /q directory
```

## 安全保证 / Security Guarantees

1. **沙箱隔离** / **Sandbox Isolation**: 所有路径必须在沙箱目录内
2. **黑名单保护** / **Blacklist Protection**: 系统关键目录受到保护
3. **路径规范化** / **Path Normalization**: 自动处理 `..` 和 `.` 等路径符号
4. **跨平台一致性** / **Cross-platform Consistency**: 在不同操作系统上保持一致的安全策略

## 错误处理 / Error Handling

### 沙箱违规 / Sandbox Violation

```json
{
  "error": "path is outside sandbox directory: path '../../etc/passwd' is outside sandbox directory"
}
```

### 黑名单目录 / Blacklisted Directory

```json
{
  "error": "directory is blacklisted: path '/bin' is in blacklisted directory"
}
```

## 最佳实践 / Best Practices

1. **使用相对路径** / **Use Relative Paths**: 优先使用相对于工作目录的路径
2. **避免绝对路径** / **Avoid Absolute Paths**: 除非确定路径在沙箱内
3. **检查错误信息** / **Check Error Messages**: 仔细阅读错误信息以了解拒绝原因
4. **测试路径** / **Test Paths**: 在执行删除操作前，先使用 `file_exists` 验证路径

## 与黑名单机制的区别 / Difference from Blacklist Mechanism

| 特性 / Feature | 黑名单机制 / Blacklist | 路径验证 / Path Validation |
|---------------|---------------------|--------------------------|
| 验证对象 / Target | 命令本身 / Command itself | 命令参数 / Command arguments |
| 验证时机 / Timing | 命令执行前 / Before execution | 命令执行前 / Before execution |
| 灵活性 / Flexibility | 完全禁止 / Complete prohibition | 条件允许 / Conditional allowance |
| 适用场景 / Use Case | 危险命令 / Dangerous commands | 文件操作命令 / File operation commands |

## 总结 / Summary

路径参数验证机制提供了一种更加灵活和智能的安全控制方式，既保证了系统安全，又允许用户在沙箱内自由地使用文件操作命令。这种机制是对传统黑名单机制的重要补充。

The path argument validation mechanism provides a more flexible and intelligent security control method, ensuring system security while allowing users to freely use file operation commands within the sandbox. This mechanism is an important complement to the traditional blacklist mechanism.

