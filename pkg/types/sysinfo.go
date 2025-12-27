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

// Package types 系统信息相关类型定义 / System info related type definitions
package types

// GetSystemInfoRequest 获取系统信息请求 / Get system info request
type GetSystemInfoRequest struct{}

// GetSystemInfoResponse 获取系统信息响应 / Get system info response
type GetSystemInfoResponse struct {
	OS       OSInfo        `json:"os"`                 // 操作系统信息 / Operating system info
	CPU      CPUInfo       `json:"cpu"`                // CPU信息 / CPU info
	Memory   MemoryInfo    `json:"memory"`             // 内存信息 / Memory info
	GPUs     []GPUInfo     `json:"gpus,omitempty"`     // 显卡信息列表 / GPU info list
	Networks []NetworkInfo `json:"networks,omitempty"` // 网卡信息列表 / Network interface info list
}

// OSInfo 操作系统信息 / Operating system info
type OSInfo struct {
	Name          string `json:"name"`                     // 操作系统名称 / OS name (e.g. "Windows 11", "Ubuntu 22.04")
	Platform      string `json:"platform"`                 // 平台类型 / Platform type (e.g. "windows", "linux", "darwin")
	Family        string `json:"family"`                   // 系统家族 / OS family (e.g. "debian", "rhel", "standalone")
	Version       string `json:"version"`                  // 版本号 / Version number
	Architecture  string `json:"architecture"`             // 系统架构 / Architecture (e.g. "amd64", "arm64")
	Hostname      string `json:"hostname"`                 // 主机名 / Hostname
	KernelVersion string `json:"kernel_version,omitempty"` // 内核版本 / Kernel version
	BootTime      int64  `json:"boot_time"`                // 系统启动时间(Unix时间戳) / Boot time (Unix timestamp)
	Uptime        int64  `json:"uptime"`                   // 系统运行时长(秒) / Uptime in seconds
}

// CPUInfo CPU信息 / CPU info
type CPUInfo struct {
	ModelName     string    `json:"model_name"`              // CPU型号 / CPU model name
	Vendor        string    `json:"vendor"`                  // 厂商 / Vendor (e.g. "GenuineIntel", "AuthenticAMD")
	PhysicalCores int       `json:"physical_cores"`          // 物理核心数 / Physical core count
	LogicalCores  int       `json:"logical_cores"`           // 逻辑核心数 / Logical core count
	Frequency     float64   `json:"frequency"`               // 主频(MHz) / Frequency in MHz
	CacheSize     int32     `json:"cache_size"`              // 缓存大小(KB) / Cache size in KB
	UsagePercent  []float64 `json:"usage_percent,omitempty"` // 各核心使用率(%) / Per-core usage percentage
}

// MemoryInfo 内存信息 / Memory info
type MemoryInfo struct {
	Total       uint64  `json:"total"`        // 总内存(字节) / Total memory in bytes
	Available   uint64  `json:"available"`    // 可用内存(字节) / Available memory in bytes
	Used        uint64  `json:"used"`         // 已用内存(字节) / Used memory in bytes
	UsedPercent float64 `json:"used_percent"` // 内存使用率(%) / Memory usage percentage
	SwapTotal   uint64  `json:"swap_total"`   // 交换分区总量(字节) / Swap total in bytes
	SwapUsed    uint64  `json:"swap_used"`    // 交换分区已用(字节) / Swap used in bytes
	SwapFree    uint64  `json:"swap_free"`    // 交换分区空闲(字节) / Swap free in bytes
}

// GPUInfo 显卡信息 / GPU info
type GPUInfo struct {
	Index         int     `json:"index"`                    // 显卡索引 / GPU index
	Name          string  `json:"name"`                     // 显卡名称 / GPU name
	Vendor        string  `json:"vendor,omitempty"`         // 厂商 / Vendor (e.g. "NVIDIA", "AMD", "Intel")
	DriverVersion string  `json:"driver_version,omitempty"` // 驱动版本 / Driver version
	MemoryTotal   uint64  `json:"memory_total,omitempty"`   // 显存总量(字节) / Total memory in bytes
	MemoryUsed    uint64  `json:"memory_used,omitempty"`    // 显存已用(字节) / Used memory in bytes
	MemoryFree    uint64  `json:"memory_free,omitempty"`    // 显存空闲(字节) / Free memory in bytes
	Temperature   float64 `json:"temperature,omitempty"`    // 温度(摄氏度) / Temperature in Celsius
	Utilization   float64 `json:"utilization,omitempty"`    // GPU使用率(%) / GPU utilization percentage
}

// NetworkInfo 网卡信息 / Network interface info
type NetworkInfo struct {
	Name         string      `json:"name"`                    // 网卡名称 / Interface name
	HardwareAddr string      `json:"hardware_addr,omitempty"` // MAC地址 / MAC address
	MTU          int         `json:"mtu"`                     // 最大传输单元 / Maximum transmission unit
	Flags        []string    `json:"flags,omitempty"`         // 网卡标志 / Interface flags
	Addresses    []IPAddress `json:"addresses,omitempty"`     // IP地址列表 / IP address list
	BytesSent    uint64      `json:"bytes_sent"`              // 发送字节数 / Bytes sent
	BytesRecv    uint64      `json:"bytes_recv"`              // 接收字节数 / Bytes received
	PacketsSent  uint64      `json:"packets_sent"`            // 发送包数 / Packets sent
	PacketsRecv  uint64      `json:"packets_recv"`            // 接收包数 / Packets received
	Speed        uint64      `json:"speed,omitempty"`         // 网卡速度(Mbps) / Speed in Mbps
}

// IPAddress IP地址信息 / IP address info
type IPAddress struct {
	Address string `json:"address"` // IP地址 / IP address
	Netmask string `json:"netmask"` // 子网掩码 / Netmask
	Family  string `json:"family"`  // 地址族 / Address family (e.g. "IPv4", "IPv6")
}
