# 下载工具使用指南 / Download Tool Guide

## 概述 / Overview

下载工具（`download_file`）允许你从互联网下载文件到沙箱目录中。支持 HTTP/HTTPS 协议，可以使用 GET、POST 等多种 HTTP
方法，并支持自定义请求头和请求体。

The download tool (`download_file`) allows you to download files from the internet to the sandbox directory. It supports
HTTP/HTTPS protocols, various HTTP methods like GET and POST, and custom headers and request body.

## 功能特性 / Features

- ✅ 支持 HTTP 和 HTTPS 协议 / Support HTTP and HTTPS protocols
- ✅ 支持多种 HTTP 方法（GET、POST、PUT、DELETE、HEAD、PATCH、OPTIONS）/ Support multiple HTTP methods
- ✅ 自定义请求头 / Custom request headers
- ✅ 自定义请求体（用于 POST 等方法）/ Custom request body (for POST etc.)
- ✅ 超时控制（1-300秒）/ Timeout control (1-300 seconds)
- ✅ 自动创建父目录 / Automatic parent directory creation
- ✅ 沙箱安全保护 / Sandbox security protection
- ✅ 路径遍历攻击防护 / Path traversal attack protection
- ✅ 返回文件大小和内容类型 / Return file size and content type

## 使用方法 / Usage

### 基本用法 / Basic Usage

最简单的用法是下载一个文件：

```json
{
  "url": "https://example.com/file.pdf",
  "path": "downloads/file.pdf"
}
```

### 使用 POST 方法 / Using POST Method

下载 API 响应数据：

```json
{
  "url": "https://api.example.com/data",
  "path": "data/response.json",
  "method": "POST",
  "headers": {
    "Content-Type": "application/json",
    "Authorization": "Bearer your-token-here"
  },
  "body": "{\"query\": \"search term\"}"
}
```

### 自定义超时时间 / Custom Timeout

对于大文件或慢速连接，可以增加超时时间：

```json
{
  "url": "https://example.com/large-file.zip",
  "path": "downloads/large-file.zip",
  "timeout": 120
}
```

### 自定义请求头 / Custom Headers

添加认证或其他自定义请求头：

```json
{
  "url": "https://api.example.com/protected/file",
  "path": "downloads/protected-file.dat",
  "headers": {
    "Authorization": "Bearer your-token",
    "X-Custom-Header": "custom-value"
  }
}
```

## 参数说明 / Parameters

| 参数 / Parameter    | 类型 / Type | 必填 / Required | 默认值 / Default | 说明 / Description                                                                    |
|-------------------|-----------|---------------|---------------|-------------------------------------------------------------------------------------|
| `url`             | string    | ✅             | -             | 下载URL，必须是 http:// 或 https:// 开头 / Download URL, must start with http:// or https:// |
| `path`            | string    | ✅             | -             | 保存路径（相对于沙箱目录）/ Save path (relative to sandbox directory)                            |
| `method`          | string    | ❌             | `GET`         | HTTP方法 / HTTP method                                                                |
| `headers`         | object    | ❌             | `{}`          | 自定义请求头 / Custom headers                                                             |
| `body`            | string    | ❌             | `""`          | 请求体（用于POST等）/ Request body (for POST etc.)                                          |
| `timeout`         | integer   | ❌             | `30`          | 超时时间（秒，1-300）/ Timeout in seconds (1-300)                                           |
| `skip_tls_verify` | boolean   | ❌             | `false`       | 跳过TLS证书验证（⚠️ 仅用于开发环境）/ Skip TLS verification (⚠️ dev only)                          |

## 响应格式 / Response Format

成功响应示例：

```json
{
  "success": true,
  "message": "file downloaded successfully",
  "sandbox_path": "downloads/file.pdf",
  "absolute_path": "/path/to/sandbox/downloads/file.pdf",
  "size": 1024000,
  "content_type": "application/pdf"
}
```

## 安全限制 / Security Restrictions

1. **URL 限制** / URL Restrictions
    - 只支持 HTTP 和 HTTPS 协议
    - Only HTTP and HTTPS protocols are supported

2. **路径限制** / Path Restrictions
    - 所有文件保存在沙箱目录内
    - All files are saved within the sandbox directory
    - 自动防护路径遍历攻击（如 `../../../etc/passwd`）
    - Automatic protection against path traversal attacks

3. **文件大小限制** / File Size Limits
    - 最大文件大小：1GB（1,073,741,824 字节）
    - Maximum file size: 1GB (1,073,741,824 bytes)
    - 超过限制的文件将被拒绝下载
    - Files exceeding the limit will be rejected
    - 不完整的文件会被自动删除
    - Incomplete files are automatically deleted

4. **超时限制** / Timeout Limits
    - 最小超时：1秒 / Minimum timeout: 1 second
    - 最大超时：300秒（5分钟）/ Maximum timeout: 300 seconds (5 minutes)
    - 默认超时：30秒 / Default timeout: 30 seconds

5. **HTTP 重定向限制** / HTTP Redirect Limits
    - 最大重定向次数：10次
    - Maximum redirects: 10
    - 超过限制将停止下载
    - Download stops when limit is exceeded

6. **HTTP 方法限制** / HTTP Method Restrictions
    - 支持的方法：GET、POST、PUT、DELETE、HEAD、PATCH、OPTIONS
    - Supported methods: GET, POST, PUT, DELETE, HEAD, PATCH, OPTIONS

## 错误处理 / Error Handling

常见错误及解决方法：

1. **URL 格式错误** / Invalid URL Format
   ```
   Error: invalid URL format, must start with http:// or https://
   ```
   解决：确保 URL 以 `http://` 或 `https://` 开头

2. **HTTP 错误状态码** / HTTP Error Status Code
   ```
   Error: HTTP request failed with status code: 404
   ```
   解决：检查 URL 是否正确，资源是否存在

3. **路径遍历攻击** / Path Traversal Attack
   ```
   Error: path is outside sandbox directory
   ```
   解决：使用相对路径，不要使用 `..` 等路径遍历符号

4. **文件大小超限** / File Size Exceeded
   ```
   Error: file size (xxx bytes) exceeds maximum allowed size (1073741824 bytes)
   ```
   解决：文件太大，无法下载。考虑使用其他方式获取文件

5. **重定向次数过多** / Too Many Redirects
   ```
   Error: stopped after 10 redirects
   ```
   解决：URL 可能存在重定向循环，检查 URL 是否正确

6. **超时** / Timeout
   ```
   Error: context deadline exceeded
   ```
   解决：增加 `timeout` 参数值

## 使用示例 / Examples

### 示例 1：下载图片 / Example 1: Download Image

```json
{
  "url": "https://example.com/image.png",
  "path": "images/downloaded.png"
}
```

### 示例 2：下载 JSON 数据 / Example 2: Download JSON Data

```json
{
  "url": "https://api.github.com/repos/user/repo",
  "path": "data/repo-info.json",
  "headers": {
    "Accept": "application/vnd.github.v3+json"
  }
}
```

### 示例 3：POST 请求下载 / Example 3: POST Request Download

```json
{
  "url": "https://api.example.com/export",
  "path": "exports/data.csv",
  "method": "POST",
  "headers": {
    "Content-Type": "application/json"
  },
  "body": "{\"format\": \"csv\", \"filter\": \"active\"}"
}
```

## 最佳实践 / Best Practices

1. **使用合适的超时时间** / Use Appropriate Timeout
    - 小文件：使用默认30秒
    - 大文件：根据文件大小和网络速度调整

2. **组织文件结构** / Organize File Structure
    - 使用有意义的目录结构，如 `downloads/images/`、`downloads/documents/`
    - Use meaningful directory structure

3. **错误处理** / Error Handling
    - 始终检查响应的 `success` 字段
    - Always check the `success` field in response

4. **安全性** / Security
    - 不要下载不信任来源的文件
    - Don't download files from untrusted sources
    - 验证下载的文件内容
    - Verify downloaded file content

## TLS 证书验证 / TLS Certificate Verification

### ⚠️ 重要安全警告 / Important Security Warning

`skip_tls_verify` 参数允许跳过 HTTPS 证书验证，但这会带来严重的安全风险：

The `skip_tls_verify` parameter allows skipping HTTPS certificate verification, but this introduces serious security
risks:

- ❌ **中间人攻击** / Man-in-the-Middle Attacks：攻击者可以拦截和修改通信内容
- ❌ **数据泄露** / Data Leakage：敏感数据可能被窃取
- ❌ **身份伪造** / Identity Spoofing：无法确认服务器身份

### 何时使用 / When to Use

**仅在以下情况下使用** / Only use in the following scenarios:

1. **开发环境** / Development Environment
    - 使用自签名证书的本地测试服务器
    - Local test servers with self-signed certificates

2. **内部网络** / Internal Network
    - 企业内部使用自签名证书的服务
    - Internal services with self-signed certificates

3. **临时测试** / Temporary Testing
    - 快速验证功能，但不应用于生产环境
    - Quick functionality verification, never for production

### ❌ 禁止使用场景 / Never Use In

- 生产环境 / Production environments
- 处理敏感数据 / Handling sensitive data
- 公网服务 / Public internet services
- 金融或医疗系统 / Financial or healthcare systems

### 使用示例 / Usage Example

```json
{
  "url": "https://localhost:8443/api/data",
  "path": "downloads/local-data.json",
  "skip_tls_verify": true
}
```

**注意**：使用此参数时，日志中会记录警告信息。

**Note**: When using this parameter, a warning will be logged.

### 更好的替代方案 / Better Alternatives

1. **使用有效证书** / Use Valid Certificates
    - 从受信任的 CA 获取证书
    - Obtain certificates from trusted CAs
    - 使用 Let's Encrypt 等免费证书服务
    - Use free certificate services like Let's Encrypt

2. **配置证书信任** / Configure Certificate Trust
    - 将自签名证书添加到系统信任列表
    - Add self-signed certificates to system trust store

3. **使用 HTTP（仅限内网）** / Use HTTP (Internal Only)
    - 在完全隔离的内网环境中使用 HTTP
    - Use HTTP in completely isolated internal networks

