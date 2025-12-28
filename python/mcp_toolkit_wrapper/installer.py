"""
MCP Toolkit 二进制安装器
Binary installer for MCP Toolkit
"""

import os
import sys
import platform
import subprocess
import urllib.request
import tarfile
import zipfile
import json
import importlib.metadata
from pathlib import Path
from typing import Optional

# 配置 / Configuration
REPO = "shibingli/mcp-toolkit"

# 动态获取版本号 / Get version dynamically
try:
    VERSION = importlib.metadata.version("mcp-sandbox-toolkit")
except importlib.metadata.PackageNotFoundError:
    VERSION = "0.0.0-dev"


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


def get_install_dir() -> Path:
    """获取安装目录 / Get installation directory"""
    if sys.platform == "win32":
        base_dir = Path(os.environ.get("LOCALAPPDATA", Path.home() / "AppData" / "Local"))
        install_dir = base_dir / "Programs" / "mcp-toolkit"
    else:
        install_dir = Path.home() / ".local" / "bin"
    
    install_dir.mkdir(parents=True, exist_ok=True)
    return install_dir


def get_binary_path() -> Path:
    """获取二进制文件路径 / Get binary path"""
    return get_install_dir() / get_binary_name()


def get_latest_version() -> str:
    """
    从 GitHub API 获取最新版本 / Get latest version from GitHub API
    
    Returns:
        str: 版本号 / Version number
    """
    try:
        url = f"https://api.github.com/repos/{REPO}/releases/latest"
        with urllib.request.urlopen(url) as response:
            data = json.loads(response.read().decode())
            return data["tag_name"]
    except Exception as e:
        print(f"Warning: Failed to get latest version, using default {VERSION}: {e}")
        return f"v{VERSION}"


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


def install_binary(version: Optional[str] = None) -> Path:
    """
    安装二进制文件 / Install binary
    
    Args:
        version: 版本号，如果为 None 则使用最新版本 / Version number, use latest if None
    
    Returns:
        Path: 安装后的二进制文件路径 / Installed binary path
    """
    # 检查是否已安装 / Check if already installed
    binary_path = get_binary_path()
    if binary_path.exists():
        print(f"MCP Toolkit is already installed at: {binary_path}")
        return binary_path
    
    # 下载 / Download
    archive_path = download_binary(version)
    
    # 解压 / Extract
    extracted_binary = extract_binary(archive_path)
    
    # 安装 / Install
    install_dir = get_install_dir()
    final_binary = install_dir / get_binary_name()
    
    print(f"Installing to: {final_binary}")
    
    # 复制文件 / Copy file
    import shutil
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


def main():
    """主函数，用于命令行调用 / Main function for CLI"""
    binary_path = get_binary_path()
    
    # 如果未安装，先安装 / Install if not installed
    if not binary_path.exists():
        print("MCP Toolkit not found, installing...")
        binary_path = install_binary()
    
    # 执行二进制文件 / Execute binary
    try:
        subprocess.run([str(binary_path)] + sys.argv[1:])
    except KeyboardInterrupt:
        sys.exit(0)
    except Exception as e:
        print(f"Error running MCP Toolkit: {e}", file=sys.stderr)
        sys.exit(1)


if __name__ == "__main__":
    main()

