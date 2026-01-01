"""
MCP Toolkit Python 包装器
Python wrapper for MCP Toolkit

支持两种模式 / Supports two modes:
1. 嵌入式模式：二进制文件打包在 wheel 中（推荐，通过 PyPI 安装）
   Embedded mode: Binary bundled in wheel (recommended, install via PyPI)
2. 下载模式：运行时从 GitHub 下载二进制文件（开发/源码安装时使用）
   Download mode: Download binary from GitHub at runtime (for dev/source install)
"""

import os
import sys
import platform
import subprocess
import urllib.request
import tarfile
import zipfile
import json
import shutil
import sysconfig
import importlib.metadata
from pathlib import Path
from typing import Optional, List

# 配置 / Configuration
REPO = "shibingli/mcp-toolkit"
PACKAGE_NAME = "mcp-sandbox-toolkit"

# 动态获取版本号 / Get version dynamically
try:
    VERSION = importlib.metadata.version(PACKAGE_NAME)
except importlib.metadata.PackageNotFoundError:
    VERSION = "0.0.0-dev"


def is_dev_version(version: str) -> bool:
    """
    检查是否为开发版本 / Check if version is a development version

    Args:
        version: 版本号 / Version number

    Returns:
        bool: 是否为开发版本 / Whether it's a development version
    """
    if not version:
        return True
    version_lower = version.lower()
    return (
        "dev" in version_lower or
        version.startswith("0.0.0") or
        version == "0.0.0-dev"
    )


def get_package_bin_dir() -> Path:
    """
    获取包内的 bin 目录路径
    Get the bin directory path within the package

    Returns:
        Path: 包内 bin 目录的路径
    """
    return Path(__file__).parent / "bin"


def get_embedded_binary_path() -> Optional[Path]:
    """
    获取嵌入式二进制文件路径（打包在 wheel 中的二进制文件）
    Get embedded binary path (binary bundled in wheel)

    Returns:
        Optional[Path]: 嵌入式二进制文件路径，如果不存在则返回 None
    """
    bin_dir = get_package_bin_dir()
    binary_name = "mcp-toolkit.exe" if sys.platform == "win32" else "mcp-toolkit"
    binary_path = bin_dir / binary_name

    if binary_path.exists():
        return binary_path
    return None


def get_platform_info() -> tuple[str, str]:
    """
    获取平台信息 / Get platform information

    Returns:
        tuple: (os_name, architecture)
    """
    os_name = platform.system().lower()
    machine = platform.machine().lower()

    # 标准化操作系统名称 / Normalize OS name
    if os_name == "windows":
        os_name = "windows"
    elif os_name == "darwin":
        os_name = "darwin"
    elif os_name == "linux":
        os_name = "linux"
    else:
        raise RuntimeError(f"Unsupported operating system: {os_name}")

    # 标准化架构名称 / Normalize architecture
    if machine in ("x86_64", "amd64"):
        arch = "amd64"
    elif machine in ("aarch64", "arm64"):
        arch = "arm64"
    else:
        raise RuntimeError(f"Unsupported architecture: {machine}")

    return os_name, arch


def get_binary_name() -> str:
    """获取二进制文件名 / Get binary name"""
    os_name, _ = get_platform_info()
    return "mcp-toolkit.exe" if os_name == "windows" else "mcp-toolkit"


def get_scripts_dir() -> Optional[Path]:
    """
    获取当前 Python 环境的 scripts 目录
    Get the scripts directory of the current Python environment

    Returns:
        Optional[Path]: scripts 目录路径，如果无法获取则返回 None
    """
    try:
        scripts_path = sysconfig.get_path('scripts')
        if scripts_path:
            scripts_dir = Path(scripts_path)
            if scripts_dir.exists() and os.access(scripts_dir, os.W_OK):
                return scripts_dir
            try:
                scripts_dir.mkdir(parents=True, exist_ok=True)
                return scripts_dir
            except (PermissionError, OSError):
                pass
    except Exception:
        pass
    return None


def get_user_install_dir() -> Path:
    """
    获取用户级安装目录（回退方案）
    Get user-level installation directory (fallback)

    Returns:
        Path: 用户级安装目录
    """
    if sys.platform == "win32":
        base_dir = Path(os.environ.get("LOCALAPPDATA", Path.home() / "AppData" / "Local"))
        install_dir = base_dir / "Programs" / "mcp-toolkit"
    else:
        install_dir = Path.home() / ".local" / "bin"
    install_dir.mkdir(parents=True, exist_ok=True)
    return install_dir


def get_cache_dir() -> Path:
    """
    获取缓存目录
    Get cache directory

    Returns:
        Path: 缓存目录路径
    """
    return Path.home() / ".cache" / "mcp-toolkit"


def get_install_dir() -> Path:
    """
    获取安装目录 / Get installation directory

    优先使用 Python 环境的 scripts 目录，这样 pip uninstall 时会自动清理
    Prefer Python environment's scripts directory so pip uninstall cleans up automatically

    Returns:
        Path: 安装目录路径
    """
    # 优先使用 Python 环境的 scripts 目录
    scripts_dir = get_scripts_dir()
    if scripts_dir:
        return scripts_dir

    # 回退到用户级目录
    return get_user_install_dir()


def get_all_possible_binary_paths() -> List[Path]:
    """
    获取所有可能的二进制文件路径（用于卸载时清理）
    Get all possible binary paths (for cleanup during uninstall)

    Returns:
        List[Path]: 所有可能的二进制文件路径列表
    """
    binary_name = get_binary_name()
    paths = []

    # 嵌入式二进制文件（包内）
    embedded = get_embedded_binary_path()
    if embedded:
        paths.append(embedded)

    # Python 环境的 scripts 目录
    scripts_dir = get_scripts_dir()
    if scripts_dir:
        paths.append(scripts_dir / binary_name)

    # 用户级目录
    user_dir = get_user_install_dir()
    paths.append(user_dir / binary_name)

    return paths


def get_binary_path() -> Path:
    """
    获取二进制文件路径 / Get binary path

    优先级：
    1. 嵌入式二进制文件（打包在 wheel 中）
    2. Python scripts 目录中的二进制文件
    3. 用户级目录中的二进制文件
    4. 默认安装目录

    Returns:
        Path: 二进制文件路径
    """
    # 优先使用嵌入式二进制文件
    embedded = get_embedded_binary_path()
    if embedded:
        return embedded

    binary_name = get_binary_name()

    # 检查其他可能的位置
    scripts_dir = get_scripts_dir()
    if scripts_dir:
        scripts_binary = scripts_dir / binary_name
        if scripts_binary.exists():
            return scripts_binary

    user_dir = get_user_install_dir()
    user_binary = user_dir / binary_name
    if user_binary.exists():
        return user_binary

    # 返回默认安装目录（用于下载模式）
    return get_install_dir() / binary_name


def get_latest_version() -> str:
    """
    从 GitHub API 获取最新版本 / Get latest version from GitHub API

    Returns:
        str: 版本号 / Version number

    Raises:
        RuntimeError: 如果无法获取最新版本 / If failed to get latest version
    """
    try:
        url = f"https://api.github.com/repos/{REPO}/releases/latest"
        req = urllib.request.Request(url, headers={"User-Agent": "mcp-toolkit-installer"})
        with urllib.request.urlopen(req, timeout=30) as response:
            data = json.loads(response.read().decode())
            return data["tag_name"]
    except Exception as e:
        raise RuntimeError(f"Failed to get latest version from GitHub: {e}")


def download_binary(version: Optional[str] = None) -> Path:
    """
    下载二进制文件 / Download binary
    
    Args:
        version: 版本号，如果为 None 则使用最新版本 / Version number, use latest if None
    
    Returns:
        Path: 下载的文件路径 / Downloaded file path
    """
    if version is None:
        version = get_latest_version()
    
    os_name, arch = get_platform_info()
    
    # 构建下载 URL / Build download URL
    if os_name == "windows":
        filename = f"mcp-toolkit-{version}-{os_name}-{arch}.zip"
    else:
        filename = f"mcp-toolkit-{version}-{os_name}-{arch}.tar.gz"
    
    url = f"https://github.com/{REPO}/releases/download/{version}/{filename}"
    
    print(f"Downloading MCP Toolkit {version} for {os_name}-{arch}...")
    print(f"URL: {url}")
    
    # 下载到临时目录 / Download to temp directory
    temp_dir = Path.home() / ".cache" / "mcp-toolkit"
    temp_dir.mkdir(parents=True, exist_ok=True)
    temp_file = temp_dir / filename
    
    try:
        urllib.request.urlretrieve(url, temp_file)
        print(f"Downloaded to: {temp_file}")
        return temp_file
    except Exception as e:
        raise RuntimeError(f"Failed to download binary: {e}")


def extract_binary(archive_path: Path) -> Path:
    """
    解压二进制文件 / Extract binary
    
    Args:
        archive_path: 压缩包路径 / Archive path
    
    Returns:
        Path: 解压后的二进制文件路径 / Extracted binary path
    """
    print(f"Extracting {archive_path}...")
    
    extract_dir = archive_path.parent / "extracted"
    extract_dir.mkdir(exist_ok=True)
    
    # 解压 / Extract
    if archive_path.suffix == ".zip":
        with zipfile.ZipFile(archive_path, 'r') as zip_ref:
            zip_ref.extractall(extract_dir)
    else:  # .tar.gz
        with tarfile.open(archive_path, 'r:gz') as tar_ref:
            tar_ref.extractall(extract_dir)
    
    # 查找二进制文件 / Find binary
    binary_name = get_binary_name()
    for root, dirs, files in os.walk(extract_dir):
        if binary_name in files:
            binary_path = Path(root) / binary_name
            print(f"Found binary: {binary_path}")
            return binary_path
    
    raise RuntimeError(f"Binary {binary_name} not found in archive")


def get_installed_version() -> Optional[str]:
    """
    获取已安装的二进制文件版本 / Get installed binary version

    Returns:
        Optional[str]: 版本号，如果未安装则返回 None / Version number, None if not installed
    """
    binary_path = get_binary_path()
    if not binary_path.exists():
        return None

    try:
        result = subprocess.run(
            [str(binary_path), "-version"],
            capture_output=True,
            text=True,
            timeout=10
        )
        if result.returncode == 0:
            # 解析版本号，输出格式为 "mcp-toolkit version v1.5.2"
            # 提取第一行并解析版本号
            output = result.stdout.strip()
            first_line = output.split('\n')[0] if output else ""
            # 查找 "version" 后面的版本号
            parts = first_line.split()
            for i, part in enumerate(parts):
                if part == "version" and i + 1 < len(parts):
                    version = parts[i + 1].lstrip("v")
                    return version
            # 备用解析：查找 vX.X.X 格式
            for part in parts:
                if part.startswith("v") and '.' in part:
                    return part.lstrip("v")
            return None
    except (subprocess.TimeoutExpired, subprocess.SubprocessError, OSError, IndexError):
        pass

    return None


def install_binary(version: Optional[str] = None, force: bool = False) -> Path:
    """
    安装二进制文件 / Install binary

    Args:
        version: 版本号，如果为 None 则使用最新版本 / Version number, use latest if None
        force: 是否强制重新安装 / Whether to force reinstall

    Returns:
        Path: 安装后的二进制文件路径 / Installed binary path
    """
    binary_path = get_binary_path()

    # 获取目标版本 / Get target version
    target_version = version if version else get_latest_version()
    target_version_clean = target_version.lstrip("v")

    # 检查是否已安装且版本匹配 / Check if already installed with matching version
    if binary_path.exists() and not force:
        installed_version = get_installed_version()
        if installed_version:
            if installed_version == target_version_clean:
                # 版本匹配，静默返回 / Version matches, return silently
                return binary_path
            else:
                print(f"Upgrading MCP Toolkit from {installed_version} to {target_version_clean}...")
        else:
            print(f"Cannot determine installed version, reinstalling...")
    elif not binary_path.exists():
        print(f"Installing MCP Toolkit {target_version_clean}...")

    # 下载 / Download
    archive_path = download_binary(target_version)

    # 解压 / Extract
    extracted_binary = extract_binary(archive_path)

    # 安装 / Install
    install_dir = get_install_dir()
    final_binary = install_dir / get_binary_name()

    print(f"Installing to: {final_binary}")

    # 复制文件 / Copy file
    shutil.copy2(extracted_binary, final_binary)

    # 设置执行权限 (Unix-like systems) / Set execute permission
    if sys.platform != "win32":
        os.chmod(final_binary, 0o755)

    # 清理临时文件 / Clean up temp files
    shutil.rmtree(archive_path.parent / "extracted", ignore_errors=True)
    archive_path.unlink(missing_ok=True)

    print(f"✓ MCP Toolkit installed successfully!")
    print(f"  Binary: {final_binary}")

    # 检查 PATH 并提供使用说明 / Check PATH and provide usage instructions
    path_env = os.environ.get("PATH", "")
    if str(install_dir) not in path_env:
        print(f"\n⚠ Warning: {install_dir} is not in your PATH")
        print(f"\nTo use 'mcp-toolkit' command directly, add it to your PATH:")

        if sys.platform == "win32":
            print(f"\n  PowerShell:")
            print(f'  $env:Path += ";{install_dir}"')
            print(f"\n  Or add permanently via System Environment Variables")
        else:
            shell = os.environ.get("SHELL", "")
            if "zsh" in shell:
                config_file = "~/.zshrc"
            elif "fish" in shell:
                config_file = "~/.config/fish/config.fish"
            else:
                config_file = "~/.bashrc"

            print(f"\n  Add this line to your {config_file}:")
            print(f'  export PATH="$HOME/.local/bin:$PATH"')
            print(f"\n  Then reload your shell:")
            print(f"  source {config_file}")
    else:
        print(f"\n✓ {install_dir} is already in your PATH")
        print(f"  You can now use 'mcp-toolkit' command directly")

    return final_binary


def uninstall_binary() -> bool:
    """
    卸载二进制文件 / Uninstall binary

    清理所有可能的安装位置和缓存目录
    Clean up all possible installation locations and cache directory

    Returns:
        bool: 是否成功卸载 / Whether uninstall was successful
    """
    cache_dir = Path.home() / ".cache" / "mcp-toolkit"
    user_install_dir = get_user_install_dir()

    success = True
    found_any = False

    # 删除所有可能位置的二进制文件 / Delete binary files from all possible locations
    for binary_path in get_all_possible_binary_paths():
        if binary_path.exists():
            found_any = True
            try:
                binary_path.unlink()
                print(f"✓ Removed binary: {binary_path}")
            except Exception as e:
                print(f"✗ Failed to remove binary {binary_path}: {e}", file=sys.stderr)
                success = False

    if not found_any:
        print("No binary files found to remove.")

    # 清理缓存目录 / Clean up cache directory
    if cache_dir.exists():
        try:
            shutil.rmtree(cache_dir)
            print(f"✓ Removed cache directory: {cache_dir}")
        except Exception as e:
            print(f"✗ Failed to remove cache directory: {e}", file=sys.stderr)
            success = False

    # 尝试删除用户安装目录（如果为空）/ Try to remove user install directory if empty
    if user_install_dir.exists():
        try:
            # 只有当目录为空时才删除
            if not any(user_install_dir.iterdir()):
                user_install_dir.rmdir()
                print(f"✓ Removed empty directory: {user_install_dir}")
        except Exception:
            pass  # 忽略错误，目录可能不为空或无权限

    if success:
        print("\n✓ MCP Toolkit binary uninstalled successfully!")

    return success


def print_version_info():
    """打印版本信息 / Print version information"""
    print(f"Python wrapper version: {VERSION}")
    installed = get_installed_version()
    if installed:
        print(f"Binary version: {installed}")
    else:
        print("Binary: not installed")
    print(f"Binary path: {get_binary_path()}")


def print_help():
    """打印帮助信息 / Print help information"""
    embedded = get_embedded_binary_path()
    mode = "embedded (PyPI)" if embedded else "download (source)"

    print(f"""MCP Toolkit - Python Wrapper

Mode: {mode}

Usage:
  mcp-sandbox-toolkit [options] [binary-args...]

Wrapper Options:
  --binary-version     Show version information
  --binary-path        Show binary installation path
  --wrapper-help       Show this help message

All other arguments are passed directly to the mcp-toolkit binary.

Examples:
  mcp-sandbox-toolkit -version            # Show binary version
  mcp-sandbox-toolkit --help              # Show binary help
  mcp-sandbox-toolkit -stdio              # Run in stdio mode

Uninstall:
  pip:   pip uninstall mcp-sandbox-toolkit
  pipx:  pipx uninstall mcp-sandbox-toolkit
  uv:    uv pip uninstall mcp-sandbox-toolkit
  uvx:   (no uninstall needed, runs in temporary environment)

Note: When installed via PyPI (pip/pipx/uv), the binary is bundled in the
package and will be automatically removed when you uninstall the package.
""")


def main():
    """主函数，用于命令行调用 / Main function for CLI"""
    # 处理包装器特定的命令行参数 / Handle wrapper-specific command line arguments
    if len(sys.argv) > 1:
        arg = sys.argv[1]

        if arg == "--install-binary":
            # 强制安装二进制文件 / Force install binary
            target_version = None
            if not is_dev_version(VERSION):
                target_version = f"v{VERSION}" if not VERSION.startswith("v") else VERSION
            install_binary(version=target_version, force=True)
            sys.exit(0)

        elif arg == "--uninstall-binary":
            # 卸载二进制文件 / Uninstall binary
            success = uninstall_binary()
            sys.exit(0 if success else 1)

        elif arg == "--upgrade-binary":
            # 升级二进制文件 / Upgrade binary
            target_version = None
            if not is_dev_version(VERSION):
                target_version = f"v{VERSION}" if not VERSION.startswith("v") else VERSION
            install_binary(version=target_version, force=True)
            sys.exit(0)

        elif arg == "--binary-version":
            # 显示版本信息 / Show version info
            print_version_info()
            sys.exit(0)

        elif arg == "--binary-path":
            # 显示二进制路径 / Show binary path
            print(get_binary_path())
            sys.exit(0)

        elif arg == "--wrapper-help":
            # 显示帮助信息 / Show help
            print_help()
            sys.exit(0)

    # 优先检查嵌入式二进制文件 / Check embedded binary first
    embedded_binary = get_embedded_binary_path()
    if embedded_binary:
        # 使用嵌入式二进制文件（PyPI wheel 安装模式）
        # Use embedded binary (PyPI wheel installation mode)
        binary_path = embedded_binary
        # 确保有执行权限 / Ensure execute permission
        if sys.platform != "win32":
            try:
                os.chmod(binary_path, 0o755)
            except Exception:
                pass
    else:
        # 下载模式（开发/源码安装）/ Download mode (dev/source install)
        # 确定目标版本 / Determine target version
        if is_dev_version(VERSION):
            target_version = None  # 将使用 get_latest_version()
        else:
            target_version = f"v{VERSION}" if not VERSION.startswith("v") else VERSION

        # 安装或更新二进制文件 / Install or update binary
        try:
            binary_path = install_binary(version=target_version)
        except Exception as e:
            print(f"Error installing MCP Toolkit: {e}", file=sys.stderr)
            # 如果安装失败但二进制文件存在，尝试使用现有的
            binary_path = get_binary_path()
            if not binary_path.exists():
                sys.exit(1)
            print(f"Using existing binary: {binary_path}", file=sys.stderr)

    # 执行二进制文件 / Execute binary
    try:
        result = subprocess.run([str(binary_path)] + sys.argv[1:])
        sys.exit(result.returncode)
    except KeyboardInterrupt:
        sys.exit(0)
    except Exception as e:
        print(f"Error running MCP Toolkit: {e}", file=sys.stderr)
        sys.exit(1)


def uninstall_main():
    """
    卸载入口函数，清理缓存文件（如果有）
    Uninstall entry point for cleaning up cache files (if any)
    """
    print("MCP Toolkit Uninstaller")
    print("=" * 40)

    embedded = get_embedded_binary_path()
    if embedded:
        print("\nThis package was installed via PyPI with embedded binary.")
        print("The binary will be automatically removed when you uninstall the package.")
        print()

        # 只清理缓存目录 / Only clean cache directory
        cache_dir = get_cache_dir()
        if cache_dir.exists():
            try:
                shutil.rmtree(cache_dir)
                print(f"✓ Removed cache directory: {cache_dir}")
            except Exception as e:
                print(f"✗ Failed to remove cache: {e}", file=sys.stderr)

        print("\nTo uninstall the package, run one of the following:")
        print()
        print("  pip:   pip uninstall mcp-sandbox-toolkit")
        print("  pipx:  pipx uninstall mcp-sandbox-toolkit")
        print("  uv:    uv pip uninstall mcp-sandbox-toolkit")
        sys.exit(0)
    else:
        # 下载模式，清理所有文件 / Download mode, clean all files
        success = uninstall_binary()
        if success:
            print("\nTo complete uninstallation, run one of the following:")
            print()
            print("  pip:   pip uninstall mcp-sandbox-toolkit")
            print("  pipx:  pipx uninstall mcp-sandbox-toolkit")
            print("  uv:    uv pip uninstall mcp-sandbox-toolkit")
        sys.exit(0 if success else 1)


if __name__ == "__main__":
    main()

