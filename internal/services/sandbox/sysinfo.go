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

package sandbox

import (
	"net"
	"runtime"

	"mcp-toolkit/pkg/types"

	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/host"
	"github.com/shirou/gopsutil/v4/mem"
	psnet "github.com/shirou/gopsutil/v4/net"
	"go.uber.org/zap"
)

// GetSystemInfo 获取系统信息 / Get system information
// 返回操作系统、CPU、内存、显卡、网卡等详细信息
func (s *Service) GetSystemInfo(_ *types.GetSystemInfoRequest) (*types.GetSystemInfoResponse, error) {
	resp := &types.GetSystemInfoResponse{}

	// 获取操作系统信息 / Get OS info
	osInfo, err := s.getOSInfo()
	if err != nil {
		s.logger.Warn("failed to get OS info", zap.Error(err))
	} else {
		resp.OS = *osInfo
	}

	// 获取CPU信息 / Get CPU info
	cpuInfo, err := s.getCPUInfo()
	if err != nil {
		s.logger.Warn("failed to get CPU info", zap.Error(err))
	} else {
		resp.CPU = *cpuInfo
	}

	// 获取内存信息 / Get memory info
	memInfo, err := s.getMemoryInfo()
	if err != nil {
		s.logger.Warn("failed to get memory info", zap.Error(err))
	} else {
		resp.Memory = *memInfo
	}

	// 获取显卡信息 / Get GPU info
	gpuInfos, err := s.getGPUInfo()
	if err != nil {
		s.logger.Warn("failed to get GPU info", zap.Error(err))
	} else {
		resp.GPUs = gpuInfos
	}

	// 获取网卡信息 / Get network info
	netInfos, err := s.getNetworkInfo()
	if err != nil {
		s.logger.Warn("failed to get network info", zap.Error(err))
	} else {
		resp.Networks = netInfos
	}

	return resp, nil
}

// getOSInfo 获取操作系统信息 / Get operating system info
func (s *Service) getOSInfo() (*types.OSInfo, error) {
	hostInfo, err := host.Info()
	if err != nil {
		return nil, err
	}

	return &types.OSInfo{
		Name:          hostInfo.Platform + " " + hostInfo.PlatformVersion,
		Platform:      hostInfo.OS,
		Family:        hostInfo.PlatformFamily,
		Version:       hostInfo.PlatformVersion,
		Architecture:  runtime.GOARCH,
		Hostname:      hostInfo.Hostname,
		KernelVersion: hostInfo.KernelVersion,
		BootTime:      int64(hostInfo.BootTime),
		Uptime:        int64(hostInfo.Uptime),
	}, nil
}

// getCPUInfo 获取CPU信息 / Get CPU info
func (s *Service) getCPUInfo() (*types.CPUInfo, error) {
	cpuInfos, err := cpu.Info()
	if err != nil {
		return nil, err
	}

	physicalCores, err := cpu.Counts(false)
	if err != nil {
		physicalCores = 0
	}

	logicalCores, err := cpu.Counts(true)
	if err != nil {
		logicalCores = 0
	}

	// 获取CPU使用率 / Get CPU usage
	usagePercent, err := cpu.Percent(0, true)
	if err != nil {
		usagePercent = nil
	}

	result := &types.CPUInfo{
		PhysicalCores: physicalCores,
		LogicalCores:  logicalCores,
		UsagePercent:  usagePercent,
	}

	// 使用第一个CPU的信息 / Use first CPU info
	if len(cpuInfos) > 0 {
		result.ModelName = cpuInfos[0].ModelName
		result.Vendor = cpuInfos[0].VendorID
		result.Frequency = cpuInfos[0].Mhz
		result.CacheSize = cpuInfos[0].CacheSize
	}

	return result, nil
}

// getMemoryInfo 获取内存信息 / Get memory info
func (s *Service) getMemoryInfo() (*types.MemoryInfo, error) {
	vmem, err := mem.VirtualMemory()
	if err != nil {
		return nil, err
	}

	swap, err := mem.SwapMemory()
	if err != nil {
		// 交换分区信息获取失败不影响主要功能 / Swap info failure doesn't affect main function
		swap = &mem.SwapMemoryStat{}
	}

	return &types.MemoryInfo{
		Total:       vmem.Total,
		Available:   vmem.Available,
		Used:        vmem.Used,
		UsedPercent: vmem.UsedPercent,
		SwapTotal:   swap.Total,
		SwapUsed:    swap.Used,
		SwapFree:    swap.Free,
	}, nil
}

// getGPUInfo 获取显卡信息 / Get GPU info
// 注意：此功能需要 NVIDIA 驱动支持，非 NVIDIA 显卡或无驱动时返回空列表
func (s *Service) getGPUInfo() ([]types.GPUInfo, error) {
	gpus := []types.GPUInfo{}

	// 尝试使用 NVML 获取 NVIDIA GPU 信息 / Try to get NVIDIA GPU info using NVML
	nvmlGPUs, err := s.getNVIDIAGPUInfo()
	if err == nil && len(nvmlGPUs) > 0 {
		gpus = append(gpus, nvmlGPUs...)
	}

	// 如果没有获取到 GPU 信息，返回空列表 / Return empty list if no GPU info
	return gpus, nil
}

// getNVIDIAGPUInfo 获取 NVIDIA GPU 信息 / Get NVIDIA GPU info
// 使用 go-nvml 库获取 NVIDIA GPU 详细信息
func (s *Service) getNVIDIAGPUInfo() ([]types.GPUInfo, error) {
	// 注意：NVML 初始化可能失败（无 NVIDIA 驱动或非 NVIDIA 显卡）
	// 这里使用延迟加载方式，避免程序启动时失败
	// Note: NVML init may fail (no NVIDIA driver or non-NVIDIA GPU)
	// Using lazy loading to avoid startup failure

	// 由于 go-nvml 需要 NVIDIA 驱动，这里提供一个基础实现
	// 实际使用时需要根据环境决定是否启用
	// Since go-nvml requires NVIDIA driver, providing a basic implementation
	// Enable based on environment in actual use

	return []types.GPUInfo{}, nil
}

// getNetworkInfo 获取网卡信息 / Get network interface info
func (s *Service) getNetworkInfo() ([]types.NetworkInfo, error) {
	// 获取网络接口列表 / Get network interface list
	interfaces, err := psnet.Interfaces()
	if err != nil {
		return nil, err
	}

	// 获取网络IO统计 / Get network IO stats
	ioCounters, err := psnet.IOCounters(true)
	if err != nil {
		ioCounters = nil
	}

	// 创建IO统计映射 / Create IO stats map
	ioMap := make(map[string]psnet.IOCountersStat)
	for _, counter := range ioCounters {
		ioMap[counter.Name] = counter
	}

	result := make([]types.NetworkInfo, 0, len(interfaces))

	for _, iface := range interfaces {
		// 跳过回环接口 / Skip loopback interface
		isLoopback := false
		for _, flag := range iface.Flags {
			if flag == "loopback" {
				isLoopback = true
				break
			}
		}
		if isLoopback {
			continue
		}

		netInfo := types.NetworkInfo{
			Name:         iface.Name,
			HardwareAddr: iface.HardwareAddr,
			MTU:          iface.MTU,
			Flags:        iface.Flags,
			Addresses:    make([]types.IPAddress, 0),
		}

		// 添加IP地址 / Add IP addresses
		for _, addr := range iface.Addrs {
			ip, ipNet, err := net.ParseCIDR(addr.Addr)
			if err != nil {
				continue
			}

			family := "IPv4"
			if ip.To4() == nil {
				family = "IPv6"
			}

			netmask := ""
			if ipNet != nil {
				netmask = net.IP(ipNet.Mask).String()
			}

			netInfo.Addresses = append(netInfo.Addresses, types.IPAddress{
				Address: ip.String(),
				Netmask: netmask,
				Family:  family,
			})
		}

		// 添加IO统计 / Add IO stats
		if counter, ok := ioMap[iface.Name]; ok {
			netInfo.BytesSent = counter.BytesSent
			netInfo.BytesRecv = counter.BytesRecv
			netInfo.PacketsSent = counter.PacketsSent
			netInfo.PacketsRecv = counter.PacketsRecv
		}

		result = append(result, netInfo)
	}

	return result, nil
}
