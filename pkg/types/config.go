// Copyright 2024 MCP Toolkit Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package types

// TransportType 传输类型 / Transport type
type TransportType string

const (
	// TransportStdio 标准输入输出传输 / Standard input/output transport
	TransportStdio TransportType = "stdio"

	// TransportHTTP HTTP传输 / HTTP transport
	TransportHTTP TransportType = "http"

	// TransportSSE Server-Sent Events传输 / Server-Sent Events transport
	TransportSSE TransportType = "sse"
)

// ServerConfig MCP服务器配置 / MCP server configuration
type ServerConfig struct {
	// Transport 传输类型 / Transport type
	Transport TransportType `json:"transport"`

	// SandboxDir 沙箱目录路径 / Sandbox directory path
	SandboxDir string `json:"sandbox_dir"`

	// HTTPConfig HTTP传输配置 / HTTP transport configuration
	HTTPConfig *HTTPConfig `json:"http_config,omitempty"`

	// SSEConfig SSE传输配置 / SSE transport configuration
	SSEConfig *SSEConfig `json:"sse_config,omitempty"`
}

// HTTPConfig HTTP传输配置 / HTTP transport configuration
type HTTPConfig struct {
	// Host 监听地址 / Listen address
	Host string `json:"host"`

	// Port 监听端口 / Listen port
	Port int `json:"port"`

	// Path 服务路径 / Service path
	Path string `json:"path"`

	// EnableCORS 是否启用CORS / Whether to enable CORS
	EnableCORS bool `json:"enable_cors"`

	// AllowedOrigins 允许的来源列表 / Allowed origins list
	AllowedOrigins []string `json:"allowed_origins,omitempty"`

	// ReadTimeout 读取超时(秒) / Read timeout in seconds
	ReadTimeout int `json:"read_timeout"`

	// WriteTimeout 写入超时(秒) / Write timeout in seconds
	WriteTimeout int `json:"write_timeout"`

	// MaxHeaderBytes 最大请求头字节数 / Maximum header bytes
	MaxHeaderBytes int `json:"max_header_bytes"`
}

// SSEConfig SSE传输配置 / SSE transport configuration
type SSEConfig struct {
	// Host 监听地址 / Listen address
	Host string `json:"host"`

	// Port 监听端口 / Listen port
	Port int `json:"port"`

	// Path 服务路径 / Service path
	Path string `json:"path"`

	// EnableCORS 是否启用CORS / Whether to enable CORS
	EnableCORS bool `json:"enable_cors"`

	// AllowedOrigins 允许的来源列表 / Allowed origins list
	AllowedOrigins []string `json:"allowed_origins,omitempty"`

	// HeartbeatInterval 心跳间隔(秒) / Heartbeat interval in seconds
	HeartbeatInterval int `json:"heartbeat_interval"`

	// MaxConnections 最大连接数 / Maximum connections
	MaxConnections int `json:"max_connections"`
}

// DefaultHTTPConfig 返回默认HTTP配置 / Return default HTTP configuration
func DefaultHTTPConfig() *HTTPConfig {
	return &HTTPConfig{
		Host:           "127.0.0.1",
		Port:           8080,
		Path:           "/mcp",
		EnableCORS:     true,
		AllowedOrigins: []string{"*"},
		ReadTimeout:    30,
		WriteTimeout:   30,
		MaxHeaderBytes: 1 << 20, // 1MB
	}
}

// DefaultSSEConfig 返回默认SSE配置 / Return default SSE configuration
func DefaultSSEConfig() *SSEConfig {
	return &SSEConfig{
		Host:              "127.0.0.1",
		Port:              8081,
		Path:              "/sse",
		EnableCORS:        true,
		AllowedOrigins:    []string{"*"},
		HeartbeatInterval: 30,
		MaxConnections:    100,
	}
}

// DefaultServerConfig 返回默认服务器配置 / Return default server configuration
func DefaultServerConfig() *ServerConfig {
	return &ServerConfig{
		Transport:  TransportStdio,
		SandboxDir: "./sandbox",
		HTTPConfig: nil,
		SSEConfig:  nil,
	}
}
