"""
MCP Toolkit Python Wrapper

这是一个 Python 包装器，用于通过 uv/pipx 安装和管理 MCP Toolkit。
实际的二进制文件会在安装时自动下载。

This is a Python wrapper for installing and managing MCP Toolkit via uv/pipx.
The actual binary will be automatically downloaded during installation.
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

from .installer import install_binary, get_binary_path, main

__all__ = ["install_binary", "get_binary_path", "main", "__version__"]

