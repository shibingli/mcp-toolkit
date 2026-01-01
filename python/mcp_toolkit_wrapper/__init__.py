"""
MCP Toolkit Python Wrapper

这是一个 Python 包装器，用于通过 pip/pipx/uv 安装和管理 MCP Toolkit。

支持两种模式：
1. 嵌入式模式：二进制文件打包在 wheel 中（通过 PyPI 安装时）
2. 下载模式：运行时从 GitHub 下载二进制文件（开发/源码安装时）

This is a Python wrapper for installing and managing MCP Toolkit via pip/pipx/uv.

Supports two modes:
1. Embedded mode: Binary bundled in wheel (when installed via PyPI)
2. Download mode: Download binary from GitHub at runtime (for dev/source install)
"""

import importlib.metadata

try:
    # Try to get version from package metadata (when installed)
    __version__ = importlib.metadata.version("mcp-sandbox-toolkit")
except importlib.metadata.PackageNotFoundError:
    # Fallback version for development
    __version__ = "0.0.0-dev"

__author__ = "MCP Toolkit Authors"
__license__ = "Apache-2.0"

from .installer import (
    install_binary,
    uninstall_binary,
    get_binary_path,
    get_embedded_binary_path,
    main,
    uninstall_main,
)

__all__ = [
    "install_binary",
    "uninstall_binary",
    "get_binary_path",
    "get_embedded_binary_path",
    "main",
    "uninstall_main",
    "__version__",
]

