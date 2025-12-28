// Package main 是 MCP Toolkit 的主程序入口
//
// MCP Toolkit 是一个基于 Model Context Protocol (MCP) 的跨平台文件系统操作工具，
// 提供安全的沙箱环境下的文件系统操作能力。
//
// # 主要特性
//
//   - 多种传输方式支持：Stdio、HTTP、SSE
//   - 完整的文件系统操作：创建、读取、写入、删除、复制、移动等
//   - 安全的沙箱机制：所有操作限制在指定目录内
//   - 命令执行功能：支持同步和异步命令执行
//   - 高性能 JSON 处理：支持 Sonic、go-json、jsoniter 等多种 JSON 库
//   - 完善的错误处理：panic recovery 机制保护系统稳定性
//
// # 使用方法
//
// 启动 Stdio 传输模式（默认）：
//
//	./mcp-toolkit
//	./mcp-toolkit -sandbox /path/to/sandbox
//
// 启动 HTTP 传输模式：
//
//	./mcp-toolkit -transport http
//	./mcp-toolkit -transport http -http-host 0.0.0.0 -http-port 8080
//
// 启动 SSE 传输模式：
//
//	./mcp-toolkit -transport sse
//	./mcp-toolkit -transport sse -sse-host 0.0.0.0 -sse-port 8081
//
// # 编译
//
// 使用 Sonic JSON 库（推荐，最高性能）：
//
//	go build -tags="sonic" -o mcp-toolkit main.go
//
// 使用标准 JSON 库：
//
//	go build -o mcp-toolkit main.go
//
// # 架构
//
// 项目采用分层架构设计：
//
//   - main.go：程序入口，负责初始化和启动
//   - pkg/types：类型定义和常量
//   - pkg/transport：传输层实现（HTTP、SSE、Stdio）
//   - pkg/client：MCP 客户端实现
//   - pkg/utils：工具函数（JSON、Recovery 等）
//   - internal/services/sandbox：沙箱服务实现（文件系统、命令执行、系统信息等）
//
// # 配置
//
// 命令行参数：
//
//   - -sandbox：沙箱目录路径（默认：./sandbox）
//   - -transport：传输方式（stdio/http/sse，默认：stdio）
//   - -http-host：HTTP 服务器主机（默认：127.0.0.1）
//   - -http-port：HTTP 服务器端口（默认：8080）
//   - -sse-host：SSE 服务器主机（默认：127.0.0.1）
//   - -sse-port：SSE 服务器端口（默认：8081）
//
// # 安全性
//
// 所有文件操作都在沙箱目录内进行，系统会：
//
//   - 验证路径是否在沙箱目录内
//   - 检测路径遍历攻击（..）
//   - 限制文件大小和路径长度
//   - 提供命令黑名单机制
//
// # 性能优化
//
//   - JSON 结构体预热机制（Sonic）
//   - 工具注册表缓存
//   - Panic recovery 避免服务崩溃
//   - 并发安全的工具调用
//
// # 版本信息
//
//   - 协议版本：2025-12-26
//   - 服务器名称：mcp-toolkit
//   - 服务器版本：1.0.1
//
// # 依赖
//
//   - Go 1.25.5+
//   - github.com/modelcontextprotocol/go-sdk v1.2.0
//   - github.com/bytedance/sonic v1.14.2
//   - go.uber.org/zap v1.27.1
//
// # 文档
//
// 更多文档请参考：
//
//   - README.md：项目概述和快速开始
//   - docs/GETTING_STARTED.md：快速开始指南
//   - docs/TRANSPORT.md：传输方式详细说明
//   - docs/CLIENT.md：客户端使用指南
//   - docs/COMMAND_EXECUTION.md：命令执行功能说明
//   - docs/RECOVERY.md：错误恢复机制说明
//
// # 示例
//
// 更多示例请参考 examples/ 目录。
//
// # 许可证
//
// 本项目采用 Apache License 2.0 许可证。
//
// # Copyright 2024 MCP Toolkit Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// 详情请参阅 LICENSE 和 NOTICE 文件。
package main
