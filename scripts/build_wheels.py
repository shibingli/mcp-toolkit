#!/usr/bin/env python3
"""
构建平台特定的 Python wheel 包
Build platform-specific Python wheel packages

此脚本用于 GitHub Actions 中构建包含预编译二进制文件的 wheel 包
This script is used in GitHub Actions to build wheels with pre-compiled binaries
"""

import os
import sys
import shutil
import subprocess
import platform
from pathlib import Path

# 平台映射 / Platform mapping
PLATFORM_MAP = {
    # (os, arch) -> (go_os, go_arch, wheel_platform_tag)
    ("linux", "amd64"): ("linux", "amd64", "manylinux2014_x86_64"),
    ("linux", "arm64"): ("linux", "arm64", "manylinux2014_aarch64"),
    ("darwin", "amd64"): ("darwin", "amd64", "macosx_10_9_x86_64"),
    ("darwin", "arm64"): ("darwin", "arm64", "macosx_11_0_arm64"),
    ("windows", "amd64"): ("windows", "amd64", "win_amd64"),
}

# 二进制文件名 / Binary file names
BINARY_NAMES = {
    "windows": "mcp-toolkit.exe",
    "linux": "mcp-toolkit",
    "darwin": "mcp-toolkit",
}


def get_project_root() -> Path:
    """获取项目根目录 / Get project root directory"""
    return Path(__file__).parent.parent


def build_go_binary(go_os: str, go_arch: str, output_path: Path) -> bool:
    """
    构建 Go 二进制文件 / Build Go binary
    
    Args:
        go_os: 目标操作系统 / Target OS
        go_arch: 目标架构 / Target architecture
        output_path: 输出路径 / Output path
        
    Returns:
        bool: 是否成功 / Whether successful
    """
    project_root = get_project_root()
    
    env = os.environ.copy()
    env["GOOS"] = go_os
    env["GOARCH"] = go_arch
    env["CGO_ENABLED"] = "0"
    
    # Linux 下启用实验性特性 / Enable experimental features on Linux
    if go_os == "linux":
        env["GOEXPERIMENT"] = "jsonv2,greenteagc"
    
    cmd = [
        "go", "build",
        "-tags=sonic",
        "-ldflags=-s -w",
        "-o", str(output_path),
        ".",
    ]
    
    print(f"Building for {go_os}/{go_arch}...")
    print(f"Command: {' '.join(cmd)}")
    
    try:
        result = subprocess.run(
            cmd,
            cwd=project_root,
            env=env,
            capture_output=True,
            text=True,
        )
        if result.returncode != 0:
            print(f"Build failed: {result.stderr}")
            return False
        print(f"Built: {output_path}")
        return True
    except Exception as e:
        print(f"Build error: {e}")
        return False


def build_wheel(target_os: str, target_arch: str, version: str) -> Optional[Path]:
    """
    构建平台特定的 wheel 包 / Build platform-specific wheel
    
    Args:
        target_os: 目标操作系统 / Target OS
        target_arch: 目标架构 / Target architecture
        version: 版本号 / Version number
        
    Returns:
        Optional[Path]: wheel 文件路径，失败返回 None
    """
    key = (target_os, target_arch)
    if key not in PLATFORM_MAP:
        print(f"Unsupported platform: {target_os}/{target_arch}")
        return None
    
    go_os, go_arch, wheel_tag = PLATFORM_MAP[key]
    binary_name = BINARY_NAMES[target_os]
    
    project_root = get_project_root()
    bin_dir = project_root / "python" / "mcp_toolkit_wrapper" / "bin"
    dist_dir = project_root / "dist"
    
    # 清理并创建目录 / Clean and create directories
    if bin_dir.exists():
        shutil.rmtree(bin_dir)
    bin_dir.mkdir(parents=True, exist_ok=True)
    dist_dir.mkdir(parents=True, exist_ok=True)
    
    # 构建二进制文件 / Build binary
    binary_path = bin_dir / binary_name
    if not build_go_binary(go_os, go_arch, binary_path):
        return None
    
    # 更新版本号 / Update version
    pyproject_path = project_root / "pyproject.toml"
    pyproject_content = pyproject_path.read_text()
    pyproject_content = pyproject_content.replace(
        'version = "0.0.0"',
        f'version = "{version}"'
    )
    pyproject_path.write_text(pyproject_content)
    
    try:
        # 构建 wheel / Build wheel
        cmd = [sys.executable, "-m", "build", "--wheel", "--outdir", str(dist_dir)]
        result = subprocess.run(cmd, cwd=project_root, capture_output=True, text=True)
        
        if result.returncode != 0:
            print(f"Wheel build failed: {result.stderr}")
            return None
        
        # 重命名 wheel 以包含正确的平台标签 / Rename wheel with correct platform tag
        for wheel_file in dist_dir.glob("*.whl"):
            # 解析 wheel 文件名 / Parse wheel filename
            parts = wheel_file.stem.split("-")
            if len(parts) >= 5:
                # name-version-python-abi-platform
                new_name = f"{parts[0]}-{parts[1]}-py3-none-{wheel_tag}.whl"
                new_path = dist_dir / new_name
                wheel_file.rename(new_path)
                print(f"Created: {new_path}")
                return new_path
                
    finally:
        # 恢复版本号 / Restore version
        pyproject_content = pyproject_content.replace(
            f'version = "{version}"',
            'version = "0.0.0"'
        )
        pyproject_path.write_text(pyproject_content)
        
        # 清理 bin 目录 / Clean bin directory
        if bin_dir.exists():
            shutil.rmtree(bin_dir)
        bin_dir.mkdir(parents=True, exist_ok=True)
        (bin_dir / ".gitkeep").touch()
    
    return None


def main():
    """主函数 / Main function"""
    import argparse
    from typing import Optional
    
    parser = argparse.ArgumentParser(description="Build platform-specific wheels")
    parser.add_argument("--os", required=True, choices=["linux", "darwin", "windows"])
    parser.add_argument("--arch", required=True, choices=["amd64", "arm64"])
    parser.add_argument("--version", required=True, help="Version number (e.g., 1.5.4)")
    
    args = parser.parse_args()
    
    wheel_path = build_wheel(args.os, args.arch, args.version)
    if wheel_path:
        print(f"Successfully built: {wheel_path}")
        sys.exit(0)
    else:
        print("Build failed")
        sys.exit(1)


if __name__ == "__main__":
    main()

