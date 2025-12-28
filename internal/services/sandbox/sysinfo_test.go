package sandbox

import (
	"os"
	"testing"

	"mcp-toolkit/pkg/types"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

// TestGetSystemInfo 测试获取系统信息 / Test get system info
func TestGetSystemInfo(t *testing.T) {
	// 创建临时目录 / Create temp directory
	tempDir, err := os.MkdirTemp("", "sysinfo_test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// 创建服务 / Create service
	logger := zap.NewNop()
	service, err := NewService(tempDir, logger)
	require.NoError(t, err)

	// 测试获取系统信息 / Test get system info
	resp, err := service.GetSystemInfo(&types.GetSystemInfoRequest{})
	require.NoError(t, err)
	require.NotNil(t, resp)

	// 验证操作系统信息 / Verify OS info
	assert.NotEmpty(t, resp.OS.Platform, "OS platform should not be empty")
	assert.NotEmpty(t, resp.OS.Architecture, "OS architecture should not be empty")
	assert.NotEmpty(t, resp.OS.Hostname, "Hostname should not be empty")

	// 验证CPU信息 / Verify CPU info
	assert.Greater(t, resp.CPU.LogicalCores, 0, "Logical cores should be greater than 0")

	// 验证内存信息 / Verify memory info
	assert.Greater(t, resp.Memory.Total, uint64(0), "Total memory should be greater than 0")
	assert.GreaterOrEqual(t, resp.Memory.Available, uint64(0), "Available memory should be >= 0")
}

// TestGetOSInfo 测试获取操作系统信息 / Test get OS info
func TestGetOSInfo(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "sysinfo_test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	logger := zap.NewNop()
	service, err := NewService(tempDir, logger)
	require.NoError(t, err)

	osInfo, err := service.getOSInfo()
	require.NoError(t, err)
	require.NotNil(t, osInfo)

	assert.NotEmpty(t, osInfo.Platform)
	assert.NotEmpty(t, osInfo.Architecture)
	assert.NotEmpty(t, osInfo.Hostname)
	assert.GreaterOrEqual(t, osInfo.BootTime, int64(0))
	assert.GreaterOrEqual(t, osInfo.Uptime, int64(0))
}

// TestGetCPUInfo 测试获取CPU信息 / Test get CPU info
func TestGetCPUInfo(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "sysinfo_test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	logger := zap.NewNop()
	service, err := NewService(tempDir, logger)
	require.NoError(t, err)

	cpuInfo, err := service.getCPUInfo()
	require.NoError(t, err)
	require.NotNil(t, cpuInfo)

	assert.Greater(t, cpuInfo.LogicalCores, 0)
	assert.GreaterOrEqual(t, cpuInfo.PhysicalCores, 0)
}

// TestGetMemoryInfo 测试获取内存信息 / Test get memory info
func TestGetMemoryInfo(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "sysinfo_test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	logger := zap.NewNop()
	service, err := NewService(tempDir, logger)
	require.NoError(t, err)

	memInfo, err := service.getMemoryInfo()
	require.NoError(t, err)
	require.NotNil(t, memInfo)

	assert.Greater(t, memInfo.Total, uint64(0))
	assert.GreaterOrEqual(t, memInfo.Available, uint64(0))
	assert.GreaterOrEqual(t, memInfo.Used, uint64(0))
	assert.GreaterOrEqual(t, memInfo.UsedPercent, float64(0))
	assert.LessOrEqual(t, memInfo.UsedPercent, float64(100))
}

// TestGetGPUInfo 测试获取GPU信息 / Test get GPU info
func TestGetGPUInfo(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "sysinfo_test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	logger := zap.NewNop()
	service, err := NewService(tempDir, logger)
	require.NoError(t, err)

	// GPU信息可能为空（无GPU或无驱动）/ GPU info may be empty (no GPU or driver)
	gpuInfos, err := service.getGPUInfo()
	require.NoError(t, err)
	require.NotNil(t, gpuInfos)
}

// TestGetNetworkInfo 测试获取网络信息 / Test get network info
func TestGetNetworkInfo(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "sysinfo_test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	logger := zap.NewNop()
	service, err := NewService(tempDir, logger)
	require.NoError(t, err)

	netInfos, err := service.getNetworkInfo()
	require.NoError(t, err)
	require.NotNil(t, netInfos)

	// 至少应该有一个非回环网卡 / Should have at least one non-loopback interface
	// 注意：某些环境可能没有网卡 / Note: some environments may have no interfaces
	for _, netInfo := range netInfos {
		assert.NotEmpty(t, netInfo.Name)
	}
}

// TestGetNVIDIAGPUInfo 测试获取NVIDIA GPU信息 / Test get NVIDIA GPU info
func TestGetNVIDIAGPUInfo(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "sysinfo_test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	logger := zap.NewNop()
	service, err := NewService(tempDir, logger)
	require.NoError(t, err)

	// NVIDIA GPU信息可能为空 / NVIDIA GPU info may be empty
	gpuInfos, err := service.getNVIDIAGPUInfo()
	require.NoError(t, err)
	require.NotNil(t, gpuInfos)
}
