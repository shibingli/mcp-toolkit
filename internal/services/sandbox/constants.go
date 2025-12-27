package sandbox

import (
	"os"
	"runtime"
)

const (
	// DefaultDirPerm 默认目录权限 / Default directory permission
	DefaultDirPerm os.FileMode = 0755

	// DefaultFilePerm 默认文件权限 / Default file permission
	DefaultFilePerm os.FileMode = 0644

	// MaxFileSize 最大文件大小(100MB) / Maximum file size (100MB)
	MaxFileSize int64 = 100 * 1024 * 1024

	// MaxPathLength 最大路径长度 / Maximum path length
	MaxPathLength = 4096

	// MaxBatchDeleteCount 最大批量删除数量 / Maximum batch delete count
	MaxBatchDeleteCount = 1000

	// DefaultCommandTimeout 默认命令超时时间(秒) / Default command timeout in seconds
	DefaultCommandTimeout = 300

	// MaxCommandTimeout 最大命令超时时间(秒) / Maximum command timeout in seconds
	MaxCommandTimeout = 3600
)

var (
	// DefaultBlacklistCommands 默认黑名单命令列表 / Default blacklist commands
	// 这些命令可能对系统造成危害,默认禁止执行 / These commands may harm the system and are prohibited by default
	DefaultBlacklistCommands = []string{
		// 系统管理命令 / System management commands
		"shutdown", "reboot", "halt", "poweroff", "init",
		// 用户管理命令 / User management commands
		"useradd", "userdel", "usermod", "passwd", "groupadd", "groupdel",
		// 磁盘管理命令 / Disk management commands
		"fdisk", "mkfs", "mount", "umount", "dd",
		// 包管理命令 / Package management commands
		"apt", "apt-get", "yum", "dnf", "rpm", "dpkg", "pacman",
		// 服务管理命令 / Service management commands
		"systemctl", "service", "chkconfig",
		// 危险命令 / Dangerous commands
		"mkfs", "format",
		// 网络配置命令 / Network configuration commands
		"ifconfig", "ip", "route", "iptables", "firewall-cmd",
		// Windows特定命令 / Windows specific commands
		"reg", "regedit", "sc", "net", "netsh", "powercfg",
	}

	// DefaultBlacklistDirectories 默认黑名单目录列表 / Default blacklist directories
	// 这些目录禁止作为工作目录 / These directories are prohibited as working directories
	DefaultBlacklistDirectories = []string{
		// Linux/Unix系统目录 / Linux/Unix system directories
		"/bin", "/sbin", "/usr/bin", "/usr/sbin", "/usr/local/bin", "/usr/local/sbin",
		"/etc", "/boot", "/dev", "/proc", "/sys", "/root",
		// Windows系统目录 / Windows system directories
		"C:\\Windows", "C:\\Windows\\System32", "C:\\Windows\\SysWOW64",
		"C:\\Program Files", "C:\\Program Files (x86)",
		// macOS系统目录 / macOS system directories
		"/System", "/Library", "/Applications",
	}

	// SystemDirectories 系统目录列表(根据操作系统) / System directories list (based on OS)
	SystemDirectories = getSystemDirectories()
)

// getSystemDirectories 获取当前操作系统的系统目录列表 / Get system directories for current OS
func getSystemDirectories() []string {
	switch runtime.GOOS {
	case "windows":
		return []string{
			"C:\\Windows",
			"C:\\Windows\\System32",
			"C:\\Windows\\SysWOW64",
			"C:\\Program Files",
			"C:\\Program Files (x86)",
			"C:\\ProgramData",
		}
	case "darwin": // macOS
		return []string{
			"/System",
			"/Library",
			"/Applications",
			"/bin",
			"/sbin",
			"/usr/bin",
			"/usr/sbin",
			"/etc",
			"/var",
			"/private",
		}
	default: // Linux and other Unix-like systems
		return []string{
			"/bin",
			"/sbin",
			"/usr/bin",
			"/usr/sbin",
			"/usr/local/bin",
			"/usr/local/sbin",
			"/etc",
			"/boot",
			"/dev",
			"/proc",
			"/sys",
			"/root",
			"/lib",
			"/lib64",
			"/var",
		}
	}
}
