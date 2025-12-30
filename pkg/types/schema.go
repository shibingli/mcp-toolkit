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

// JSONSchema JSON Schema 结构体 / JSON Schema structure
// 注意：properties 和 required 字段不使用 omitempty，确保始终输出
// 这是为了兼容 llama.cpp 等工具的 Jinja2 模板
// 兼容 OpenAI、Claude、llama.cpp、vLLM、SGLang 等框架的 Tools 调用规范
type JSONSchema struct {
	// Type 类型，通常为 "object" / Type, usually "object"
	Type string `json:"type"`

	// Properties 属性定义 / Property definitions
	Properties map[string]Property `json:"properties"`

	// Required 必需的属性列表 / List of required properties
	Required []string `json:"required"`

	// Description 模式描述 / Schema description
	Description string `json:"description,omitempty"`

	// AdditionalProperties 是否允许额外属性 / Whether additional properties are allowed
	AdditionalProperties *bool `json:"additionalProperties,omitempty"`
}

// Property JSON Schema 属性 / JSON Schema property
// 兼容 OpenAI Function Calling、Claude Tool Use、llama.cpp、vLLM、SGLang 等框架
type Property struct {
	// Type 属性类型：string, integer, number, boolean, array, object / Property type
	Type string `json:"type"`

	// Description 属性描述，帮助大模型理解参数用途 / Property description to help LLM understand parameter usage
	Description string `json:"description,omitempty"`

	// Items 数组元素类型定义（仅当 Type 为 "array" 时使用）/ Array item type definition (only when Type is "array")
	Items *Items `json:"items,omitempty"`

	// Enum 枚举值列表，限制参数只能取这些值 / Enumeration values, restricts parameter to these values only
	Enum []string `json:"enum,omitempty"`

	// Default 默认值 / Default value
	Default any `json:"default,omitempty"`

	// Minimum 数值最小值（仅当 Type 为 "integer" 或 "number" 时使用）/ Minimum value (only for integer/number types)
	Minimum *float64 `json:"minimum,omitempty"`

	// Maximum 数值最大值（仅当 Type 为 "integer" 或 "number" 时使用）/ Maximum value (only for integer/number types)
	Maximum *float64 `json:"maximum,omitempty"`

	// MinLength 字符串最小长度（仅当 Type 为 "string" 时使用）/ Minimum string length (only for string type)
	MinLength *int `json:"minLength,omitempty"`

	// MaxLength 字符串最大长度（仅当 Type 为 "string" 时使用）/ Maximum string length (only for string type)
	MaxLength *int `json:"maxLength,omitempty"`

	// Pattern 正则表达式模式（仅当 Type 为 "string" 时使用）/ Regex pattern (only for string type)
	Pattern string `json:"pattern,omitempty"`

	// Examples 示例值列表，帮助大模型理解参数格式 / Example values to help LLM understand parameter format
	Examples []any `json:"examples,omitempty"`

	// Format 格式提示，如 "date-time", "email", "uri" 等 / Format hint like "date-time", "email", "uri"
	Format string `json:"format,omitempty"`

	// Nullable 是否可为 null / Whether the value can be null
	Nullable bool `json:"nullable,omitempty"`
}

// Items JSON Schema 数组项 / JSON Schema array items
type Items struct {
	// Type 数组元素类型 / Array element type
	Type string `json:"type"`

	// Description 数组元素描述 / Array element description
	Description string `json:"description,omitempty"`

	// Enum 数组元素枚举值 / Array element enumeration values
	Enum []string `json:"enum,omitempty"`
}

// ToolSchemas 工具输入模式定义 / Tool input schema definitions
// 兼容 OpenAI Function Calling、Claude Tool Use、llama.cpp、vLLM、SGLang 等框架
var ToolSchemas = map[string]JSONSchema{
	// ==================== 文件操作工具 / File Operation Tools ====================

	"create_file": {
		Type:        "object",
		Description: "CREATE A NEW FILE with specified content. Use this tool when you need to: 1) Create a new file from scratch, 2) Write initial content to a file, 3) Overwrite an existing file completely. The tool automatically creates parent directories if they don't exist. IMPORTANT: This is the primary tool for file creation - use it whenever you need to create or completely replace a file's content. Keywords: create, new file, write file, save file, make file, generate file.",
		Properties: map[string]Property{
			"path": {
				Type:        "string",
				Description: "The file path to create. Can be absolute or relative to the current working directory. Supports nested paths - parent directories will be created automatically. Examples: 'src/main.go' (creates src/ directory if needed), '/tmp/test.txt' (absolute path), 'docs/api/readme.md' (creates docs/api/ directories).",
				MinLength:   intPtr(1),
				Examples:    []any{"src/main.go", "config/settings.json", "README.md", "data/output.txt", "logs/app.log"},
			},
			"content": {
				Type:        "string",
				Description: "The content to write to the file. Can be any text content including source code, configuration files, JSON/YAML data, plain text, or documentation. The content will completely replace any existing file content.",
				Examples:    []any{"Hello, World!", "package main\n\nfunc main() {\n\tprintln(\"Hello\")\n}", "{\"name\": \"config\", \"version\": \"1.0\"}", "# Documentation\n\nThis is a readme file."},
			},
		},
		Required: []string{"path", "content"},
	},

	"create_directory": {
		Type:        "object",
		Description: "CREATE NEW DIRECTORIES and folder structures. Automatically creates all parent directories if they don't exist (like 'mkdir -p'). Use this tool when you need to: 1) Create new folders, 2) Set up directory structures, 3) Organize project layout, 4) Create nested directories, 5) Prepare storage locations. The tool creates the entire path if needed. Keywords: create directory, mkdir, make folder, new folder, create dir, make directory.",
		Properties: map[string]Property{
			"path": {
				Type:        "string",
				Description: "The directory path to create. Can be absolute or relative. Supports nested paths - all parent directories will be created automatically. Examples: 'src/utils' (creates src/ and src/utils/), '/tmp/mydir' (absolute path), 'docs/api/v1' (creates docs/, docs/api/, and docs/api/v1/).",
				MinLength:   intPtr(1),
				Examples:    []any{"src/utils", "docs/api", "build/output", "data/cache", "logs/2024"},
			},
		},
		Required: []string{"path"},
	},

	"read_file": {
		Type:        "object",
		Description: "READ AND RETRIEVE the complete content of a file. Use this tool when you need to: 1) View or inspect file contents, 2) Read source code before modification, 3) Check configuration files, 4) Analyze existing data, 5) Understand file structure before editing. This is the primary tool for reading files - use it whenever you need to see what's inside a file. Keywords: read, view, show, display, get content, inspect file, check file, open file.",
		Properties: map[string]Property{
			"path": {
				Type:        "string",
				Description: "The file path to read. Can be absolute or relative to the current working directory. The file must exist. Examples: 'src/main.go' (read source code), 'package.json' (read configuration), 'README.md' (read documentation), 'data/input.csv' (read data file).",
				MinLength:   intPtr(1),
				Examples:    []any{"src/main.go", "package.json", "README.md", "config.yaml", ".env", "data/users.json"},
			},
		},
		Required: []string{"path"},
	},

	"write_file": {
		Type:        "object",
		Description: "WRITE OR UPDATE content to an existing file, completely replacing its current content. Use this tool when you need to: 1) Update an existing file with new content, 2) Modify and save changes to a file, 3) Replace file content after editing. WORKFLOW: First use 'read_file' to get current content, modify it, then use 'write_file' to save. For new files, use 'create_file' instead. Keywords: write, update, modify, save, edit file, change file, replace content.",
		Properties: map[string]Property{
			"path": {
				Type:        "string",
				Description: "The file path to write to. The file should already exist (use 'create_file' for new files). Can be absolute or relative path. Examples: 'src/main.go' (update source code), 'config.yaml' (update configuration), 'index.html' (update web page).",
				MinLength:   intPtr(1),
				Examples:    []any{"src/main.go", "config.yaml", "index.html", "package.json", "README.md"},
			},
			"content": {
				Type:        "string",
				Description: "The new content to write to the file. This will completely replace the existing content. Make sure to include all content you want to keep, as the old content will be lost.",
				Examples:    []any{"Updated content", "package main\n\nfunc main() {\n\tfmt.Println(\"Updated\")\n}", "{\"version\": \"2.0\", \"updated\": true}"},
			},
		},
		Required: []string{"path", "content"},
	},

	"delete": {
		Type:        "object",
		Description: "[RECOMMENDED - PRIMARY DELETION TOOL] DELETE any file or directory automatically. This is the MAIN deletion tool - use it for ALL deletion operations. The tool intelligently detects whether the path is a file or directory and handles it appropriately. For directories, it performs recursive deletion (removes all contents). Use this tool when you need to: 1) Remove any file, 2) Remove any directory and its contents, 3) Clean up temporary files, 4) Delete build artifacts, 5) Remove old backups. Keywords: delete, remove, erase, clean, rm, unlink, destroy, eliminate.",
		Properties: map[string]Property{
			"path": {
				Type:        "string",
				Description: "The path of the file or directory to delete. Can be a single file, empty directory, or directory with contents. The tool automatically detects the type and deletes accordingly. Examples: 'temp.txt' (delete file), 'build/' (delete directory and all contents), 'old_config.json' (delete config file), 'node_modules/' (delete dependencies).",
				MinLength:   intPtr(1),
				Examples:    []any{"temp.txt", "build/", "old_config.json", "*.log", "cache/", "dist/"},
			},
		},
		Required: []string{"path"},
	},

	"delete_file": {
		Type:        "object",
		Description: "Delete a specific file. Only works on files, not directories. Use 'delete' tool if you're unsure whether the path is a file or directory.",
		Properties: map[string]Property{
			"path": {
				Type:        "string",
				Description: "The path of the file to delete. Must be a file, not a directory.",
				MinLength:   intPtr(1),
				Examples:    []any{"temp.txt", "old_backup.zip", "unused.log"},
			},
		},
		Required: []string{"path"},
	},

	"delete_directory": {
		Type:        "object",
		Description: "Delete a directory and optionally all its contents. Only works on directories, not files. Use 'delete' tool if you're unsure whether the path is a file or directory.",
		Properties: map[string]Property{
			"path": {
				Type:        "string",
				Description: "The path of the directory to delete. Must be a directory, not a file.",
				MinLength:   intPtr(1),
				Examples:    []any{"build/", "node_modules/", "temp_dir/"},
			},
			"recursive": {
				Type:        "boolean",
				Description: "Whether to recursively delete all subdirectories and files. Default is true. Set to false to only delete empty directories.",
				Default:     true,
			},
		},
		Required: []string{"path"},
	},

	"copy": {
		Type:        "object",
		Description: "[RECOMMENDED - PRIMARY COPY TOOL] COPY files or directories to a new location. This is the MAIN copy tool - use it for ALL copy operations. Automatically detects whether source is a file or directory and handles appropriately. For directories, performs recursive copy (copies all contents). Use this tool when you need to: 1) Duplicate files, 2) Create backups, 3) Copy entire directories, 4) Clone project structures, 5) Make file/folder copies. Keywords: copy, duplicate, clone, backup, cp, replicate.",
		Properties: map[string]Property{
			"source": {
				Type:        "string",
				Description: "The source path of the file or directory to copy. Can be a single file or entire directory tree. Examples: 'src/main.go' (copy file), 'config/' (copy directory and all contents), 'template.html' (copy template file).",
				MinLength:   intPtr(1),
				Examples:    []any{"src/main.go", "config/", "template.html", "docs/", "package.json"},
			},
			"destination": {
				Type:        "string",
				Description: "The destination path where the file or directory will be copied to. If source is a file, destination should be the target file path. If source is a directory, destination should be the target directory path. Examples: 'backup/main.go' (copy file to backup), 'config_backup/' (copy directory), 'index.html' (rename while copying).",
				MinLength:   intPtr(1),
				Examples:    []any{"backup/main.go", "config_backup/", "index.html", "docs_v2/", "package.backup.json"},
			},
		},
		Required: []string{"source", "destination"},
	},

	"copy_file": {
		Type:        "object",
		Description: "Copy a specific file to a new location. Only works on files, not directories. Use 'copy' tool if you're unsure whether the path is a file or directory.",
		Properties: map[string]Property{
			"source": {
				Type:        "string",
				Description: "The source file path to copy. Must be a file, not a directory.",
				MinLength:   intPtr(1),
				Examples:    []any{"config.json", "main.go", "README.md"},
			},
			"destination": {
				Type:        "string",
				Description: "The destination file path where the file will be copied to.",
				MinLength:   intPtr(1),
				Examples:    []any{"config.backup.json", "main_copy.go", "README_backup.md"},
			},
		},
		Required: []string{"source", "destination"},
	},

	"copy_directory": {
		Type:        "object",
		Description: "Copy a directory and all its contents to a new location. Only works on directories, not files. Use 'copy' tool if you're unsure whether the path is a file or directory.",
		Properties: map[string]Property{
			"source": {
				Type:        "string",
				Description: "The source directory path to copy. Must be a directory, not a file.",
				MinLength:   intPtr(1),
				Examples:    []any{"src/", "config/", "templates/"},
			},
			"destination": {
				Type:        "string",
				Description: "The destination directory path where the directory will be copied to.",
				MinLength:   intPtr(1),
				Examples:    []any{"src_backup/", "config_copy/", "templates_v2/"},
			},
		},
		Required: []string{"source", "destination"},
	},

	"move": {
		Type:        "object",
		Description: "[RECOMMENDED - PRIMARY MOVE/RENAME TOOL] MOVE or RENAME files and directories. This is the MAIN move/rename tool - use it for ALL move and rename operations. Automatically detects whether source is a file or directory and handles appropriately. Use this tool when you need to: 1) Rename files or directories, 2) Move files to different locations, 3) Reorganize project structure, 4) Relocate directories, 5) Change file/folder names. Keywords: move, rename, mv, relocate, transfer, reorg, reorganize.",
		Properties: map[string]Property{
			"source": {
				Type:        "string",
				Description: "The source path of the file or directory to move or rename. Can be a file or directory. Examples: 'old_name.txt' (file to rename), 'src/old_module/' (directory to move), 'temp.log' (file to relocate).",
				MinLength:   intPtr(1),
				Examples:    []any{"old_name.txt", "src/old_module/", "temp.log", "draft.md", "old_config/"},
			},
			"destination": {
				Type:        "string",
				Description: "The destination path where the file or directory will be moved to. This can be a new name in the same directory (rename) or a different location (move). Examples: 'new_name.txt' (rename file), 'src/new_module/' (move directory), 'logs/app.log' (move and rename file).",
				MinLength:   intPtr(1),
				Examples:    []any{"new_name.txt", "src/new_module/", "logs/app.log", "published.md", "config/"},
			},
		},
		Required: []string{"source", "destination"},
	},

	"move_file": {
		Type:        "object",
		Description: "Move or rename a specific file. Only works on files, not directories. Use 'move' tool if you're unsure whether the path is a file or directory.",
		Properties: map[string]Property{
			"source": {
				Type:        "string",
				Description: "The source file path to move. Must be a file, not a directory.",
				MinLength:   intPtr(1),
				Examples:    []any{"old_config.json", "temp.txt", "draft.md"},
			},
			"destination": {
				Type:        "string",
				Description: "The destination file path where the file will be moved to.",
				MinLength:   intPtr(1),
				Examples:    []any{"config.json", "archive/temp.txt", "published.md"},
			},
		},
		Required: []string{"source", "destination"},
	},

	"move_directory": {
		Type:        "object",
		Description: "Move or rename a directory. Only works on directories, not files. Use 'move' tool if you're unsure whether the path is a file or directory.",
		Properties: map[string]Property{
			"source": {
				Type:        "string",
				Description: "The source directory path to move. Must be a directory, not a file.",
				MinLength:   intPtr(1),
				Examples:    []any{"old_src/", "temp_build/", "draft_docs/"},
			},
			"destination": {
				Type:        "string",
				Description: "The destination directory path where the directory will be moved to.",
				MinLength:   intPtr(1),
				Examples:    []any{"src/", "build/", "docs/"},
			},
		},
		Required: []string{"source", "destination"},
	},

	"list_directory": {
		Type:        "object",
		Description: "LIST AND EXPLORE directory contents. Shows all files and subdirectories with detailed information including names, types (file/directory), sizes, and modification times. Use this tool when you need to: 1) See what files are in a directory, 2) Explore project structure, 3) Find specific files, 4) Check directory contents before operations, 5) Understand folder organization. This is the primary tool for directory exploration. Keywords: list, ls, dir, show files, directory contents, browse, explore folder, view directory.",
		Properties: map[string]Property{
			"path": {
				Type:        "string",
				Description: "The directory path to list. Use '.' for current directory, '..' for parent directory, or specify any directory path. Examples: '.' (current directory), 'src/' (source directory), '/home/user/projects' (absolute path), '..' (parent directory).",
				MinLength:   intPtr(1),
				Examples:    []any{".", "src/", "/home/user/projects", "..", "docs/", "build/"},
			},
		},
		Required: []string{"path"},
	},

	"search_files": {
		Type:        "object",
		Description: "SEARCH AND FIND files matching patterns within directories. Recursively searches through all subdirectories using powerful glob patterns. Use this tool when you need to: 1) Find files by extension (*.go, *.js), 2) Locate specific files by name pattern, 3) Search for test files, 4) Find configuration files, 5) Discover files across project. Supports wildcards: * (any characters), ? (single character), ** (any directories). Keywords: search, find, locate, grep files, filter, pattern match, glob.",
		Properties: map[string]Property{
			"path": {
				Type:        "string",
				Description: "The directory path to start searching from. The search will recursively include all subdirectories. Use '.' for current directory, or specify any directory path. Examples: '.' (search from current), 'src/' (search in source), '/home/user/projects' (absolute path).",
				MinLength:   intPtr(1),
				Examples:    []any{".", "src/", "/home/user/projects", "docs/", "tests/"},
			},
			"pattern": {
				Type:        "string",
				Description: "The search pattern using glob syntax. Wildcards: * matches any characters, ? matches single character, ** matches any directories. Examples: '*.go' (all Go files), 'test_*.py' (test files starting with test_), '**/*.md' (all Markdown files in any subdirectory), 'config.*' (config files with any extension).",
				MinLength:   intPtr(1),
				Examples:    []any{"*.go", "*.js", "test_*.py", "**/*.md", "config.*", "*.json", "**/*.test.js"},
			},
		},
		Required: []string{"path", "pattern"},
	},

	"batch_delete": {
		Type:        "object",
		Description: "Delete multiple files or directories in a single operation. Each path is processed independently, and the tool will report success/failure for each item.",
		Properties: map[string]Property{
			"paths": {
				Type:        "array",
				Description: "List of file or directory paths to delete. Each path will be processed independently.",
				Items:       &Items{Type: "string", Description: "A file or directory path to delete"},
				Examples:    []any{[]string{"temp.txt", "old_backup/", "unused.log"}},
			},
		},
		Required: []string{"paths"},
	},

	"file_stat": {
		Type:        "object",
		Description: "Get detailed information about a file or directory, including size, permissions, modification time, and type (file/directory/symlink).",
		Properties: map[string]Property{
			"path": {
				Type:        "string",
				Description: "The path of the file or directory to get information about.",
				MinLength:   intPtr(1),
				Examples:    []any{"main.go", "src/", "config.json"},
			},
		},
		Required: []string{"path"},
	},

	"file_exists": {
		Type:        "object",
		Description: "Check if a file or directory exists at the specified path. Returns true if exists, false otherwise. Useful for conditional operations.",
		Properties: map[string]Property{
			"path": {
				Type:        "string",
				Description: "The path to check for existence. Can be a file or directory path.",
				MinLength:   intPtr(1),
				Examples:    []any{"config.json", "src/", ".env"},
			},
		},
		Required: []string{"path"},
	},

	// ==================== 时间工具 / Time Tools ====================

	"get_current_time": {
		Type:        "object",
		Description: "GET CURRENT DATE AND TIME with timezone support. Returns formatted datetime string, timezone information, and Unix timestamp. Use this tool when you need to: 1) Get current time, 2) Check time in different timezones, 3) Get timestamps for logging, 4) Schedule tasks, 5) Time-based operations. Supports all IANA timezone names. Keywords: time, date, now, current time, timestamp, timezone, clock, datetime.",
		Properties: map[string]Property{
			"timezone": {
				Type:        "string",
				Description: "IANA timezone name (optional). If empty or not provided, uses the system's local timezone. Common timezones: 'Asia/Shanghai' (China), 'America/New_York' (US East), 'Europe/London' (UK), 'UTC' (Universal), 'Asia/Tokyo' (Japan), 'America/Los_Angeles' (US West), 'Europe/Paris' (France).",
				Examples:    []any{"Asia/Shanghai", "America/New_York", "Europe/London", "UTC", "Asia/Tokyo", "America/Los_Angeles"},
			},
		},
		Required: []string{},
	},

	// ==================== 命令执行工具 / Command Execution Tools ====================

	"execute_command": {
		Type:        "object",
		Description: "EXECUTE SHELL COMMANDS synchronously and get output. Use this tool to run CLI tools, build scripts, version control, package managers, and system commands. Common use cases: 1) Run git commands (status, commit, push), 2) Execute npm/pip/go commands (install, build, test), 3) Run Python/Node.js scripts, 4) Execute build tools (make, gradle, maven), 5) Run system commands (ls, cat, grep). IMPORTANT: Do NOT use for file operations (create/delete/copy/move) - use dedicated file tools instead. The command runs synchronously and waits for completion. Keywords: execute, run, command, shell, bash, terminal, cli, script, git, npm, python.",
		Properties: map[string]Property{
			"command": {
				Type:        "string",
				Description: "The command executable name (without arguments). This is the program to run. Examples: 'git' (version control), 'npm' (Node.js package manager), 'python' (Python interpreter), 'go' (Go compiler), 'ls' (list files), 'make' (build tool), 'docker' (container tool).",
				MinLength:   intPtr(1),
				Examples:    []any{"git", "npm", "python", "go", "ls", "cat", "make", "docker", "curl"},
			},
			"args": {
				Type:        "array",
				Description: "Command arguments as separate list elements. Each argument is a separate string. Examples: For 'git commit -m \"message\"' use ['commit', '-m', 'message']. For 'npm install express' use ['install', 'express']. For 'python script.py --arg value' use ['script.py', '--arg', 'value'].",
				Items:       &Items{Type: "string", Description: "A single command argument"},
				Examples:    []any{[]string{"status"}, []string{"install", "--save", "express"}, []string{"-c", "print('hello')"}, []string{"build", "--release"}},
			},
			"work_dir": {
				Type:        "string",
				Description: "The working directory where the command will be executed. Can be absolute or relative to the sandbox root. Use '.' for current directory. The command will run as if you cd'd into this directory first.",
				MinLength:   intPtr(1),
				Examples:    []any{".", "src/", "/home/user/project", "build/", "scripts/"},
			},
			"timeout": {
				Type:        "integer",
				Description: "Maximum time in seconds to wait for command completion. 0 means no timeout (wait indefinitely). Use reasonable timeouts: 30s for quick commands, 300s (5min) for builds, 0 for long-running tasks. Default is 0.",
				Minimum:     float64Ptr(0),
				Default:     0,
				Examples:    []any{30, 60, 300, 600, 0},
			},
		},
		Required: []string{"command", "work_dir"},
	},

	"get_command_blacklist": {
		Type:        "object",
		Description: "Get the current command and directory blacklist. Returns lists of blocked commands and directories that cannot be executed or accessed.",
		Properties:  map[string]Property{},
		Required:    []string{},
	},

	"update_command_blacklist": {
		Type:        "object",
		Description: "Update the command and directory blacklist. Add commands or directories that should be blocked from execution or access.",
		Properties: map[string]Property{
			"commands": {
				Type:        "array",
				Description: "List of commands to add to the blacklist. These commands will be blocked from execution.",
				Items:       &Items{Type: "string", Description: "A command name to block"},
				Examples:    []any{[]string{"rm", "sudo", "chmod"}},
			},
			"directories": {
				Type:        "array",
				Description: "List of directories to add to the blacklist. Access to these directories will be blocked.",
				Items:       &Items{Type: "string", Description: "A directory path to block"},
				Examples:    []any{[]string{"/etc", "/root", "/var/log"}},
			},
		},
		Required: []string{},
	},

	"get_working_directory": {
		Type:        "object",
		Description: "Get the current working directory path. Returns the absolute path of the current working directory.",
		Properties:  map[string]Property{},
		Required:    []string{},
	},

	"change_directory": {
		Type:        "object",
		Description: "Change the current working directory. Similar to 'cd' command. The new directory must exist and be within the sandbox.",
		Properties: map[string]Property{
			"path": {
				Type:        "string",
				Description: "The target directory path. Can be absolute or relative to the current working directory. Use '..' to go to parent directory.",
				MinLength:   intPtr(1),
				Examples:    []any{"src/", "..", "/home/user/project", ".."},
			},
		},
		Required: []string{"path"},
	},

	"execute_command_async": {
		Type:        "object",
		Description: "Execute a command asynchronously in the background. Returns a task ID immediately that can be used to check status, get output, or cancel the command. Use this for long-running commands.",
		Properties: map[string]Property{
			"command": {
				Type:        "string",
				Description: "The command to execute asynchronously.",
				MinLength:   intPtr(1),
				Examples:    []any{"npm", "python", "go", "make"},
			},
			"args": {
				Type:        "array",
				Description: "Command arguments as a list of strings.",
				Items:       &Items{Type: "string", Description: "A command argument"},
				Examples:    []any{[]string{"run", "build"}, []string{"train.py", "--epochs", "100"}},
			},
			"work_dir": {
				Type:        "string",
				Description: "The working directory for command execution.",
				MinLength:   intPtr(1),
				Examples:    []any{".", "src/", "/home/user/project"},
			},
			"timeout": {
				Type:        "integer",
				Description: "Command timeout in seconds. 0 means no timeout limit.",
				Minimum:     float64Ptr(0),
				Default:     0,
			},
			"permission_level": {
				Type:        "integer",
				Description: "Permission level for the command (0-3). Higher levels allow more privileged operations.",
				Minimum:     float64Ptr(0),
				Maximum:     float64Ptr(3),
				Default:     0,
			},
			"user": {
				Type:        "string",
				Description: "The user to execute the command as. Leave empty to use the current user.",
			},
		},
		Required: []string{"command", "work_dir"},
	},

	"get_command_task_status": {
		Type:        "object",
		Description: "Get the status of an asynchronous command task. Returns the current state (running, completed, failed), output, and other details.",
		Properties: map[string]Property{
			"task_id": {
				Type:        "string",
				Description: "The task ID returned by execute_command_async.",
				MinLength:   intPtr(1),
				Examples:    []any{"task-12345", "abc-def-ghi"},
			},
		},
		Required: []string{"task_id"},
	},

	"cancel_command_task": {
		Type:        "object",
		Description: "Cancel a running asynchronous command task. The task will be terminated if it's still running.",
		Properties: map[string]Property{
			"task_id": {
				Type:        "string",
				Description: "The task ID of the command to cancel.",
				MinLength:   intPtr(1),
				Examples:    []any{"task-12345", "abc-def-ghi"},
			},
		},
		Required: []string{"task_id"},
	},

	"list_command_tasks": {
		Type:        "object",
		Description: "List all asynchronous command tasks. Returns information about all running and completed tasks.",
		Properties:  map[string]Property{},
		Required:    []string{},
	},

	"get_command_history": {
		Type:        "object",
		Description: "Get the history of executed commands. Returns a list of previously executed commands with their results.",
		Properties: map[string]Property{
			"limit": {
				Type:        "integer",
				Description: "Maximum number of history records to return. Default returns all records.",
				Minimum:     float64Ptr(1),
				Maximum:     float64Ptr(1000),
				Examples:    []any{10, 50, 100},
			},
		},
		Required: []string{},
	},

	"set_permission_level": {
		Type:        "object",
		Description: "Set the command execution permission level. Higher levels allow more privileged operations. Level 0 is most restrictive, level 3 is least restrictive.",
		Properties: map[string]Property{
			"level": {
				Type:        "integer",
				Description: "Permission level (0-3). 0: Read-only, 1: Basic write, 2: Extended write, 3: Full access.",
				Minimum:     float64Ptr(0),
				Maximum:     float64Ptr(3),
				Enum:        []string{"0", "1", "2", "3"},
				Examples:    []any{0, 1, 2, 3},
			},
		},
		Required: []string{"level"},
	},

	"get_permission_level": {
		Type:        "object",
		Description: "Get the current command execution permission level. Returns the current level (0-3) and its description.",
		Properties:  map[string]Property{},
		Required:    []string{},
	},

	"clear_command_history": {
		Type:        "object",
		Description: "Clear all command execution history records. This action cannot be undone.",
		Properties:  map[string]Property{},
		Required:    []string{},
	},

	"get_command_task": {
		Type:        "object",
		Description: "Get detailed information about a specific asynchronous command task, including its status, output, start time, and duration.",
		Properties: map[string]Property{
			"task_id": {
				Type:        "string",
				Description: "The task ID to get information about.",
				MinLength:   intPtr(1),
				Examples:    []any{"task-12345", "abc-def-ghi"},
			},
		},
		Required: []string{"task_id"},
	},

	// ==================== System Info Tools / 系统信息工具 ====================

	"get_system_info": {
		Type:        "object",
		Description: "GET COMPREHENSIVE SYSTEM INFORMATION about the current machine. Returns detailed hardware and software information including: 1) Operating system (name, version, architecture), 2) CPU details (model, cores, speed), 3) Memory (total, available, used), 4) GPU information (if available), 5) Network interfaces and IP addresses. Use this tool when you need to: 1) Check system specifications, 2) Verify hardware capabilities, 3) Get OS information, 4) Check available resources, 5) Diagnose system configuration. Keywords: system info, hardware, specs, os info, cpu, memory, ram, gpu, network, sysinfo.",
		Properties:  map[string]Property{},
		Required:    []string{},
	},
}

// intPtr 返回 int 指针 / Returns int pointer
func intPtr(i int) *int {
	return &i
}

// float64Ptr 返回 float64 指针 / Returns float64 pointer
func float64Ptr(f float64) *float64 {
	return &f
}

// GetToolSchema 获取工具的输入模式 / Get tool input schema
func GetToolSchema(toolName string) JSONSchema {
	if schema, ok := ToolSchemas[toolName]; ok {
		return schema
	}
	// 返回空对象模式作为默认值 / Return empty object schema as default
	return JSONSchema{
		Type:       "object",
		Properties: map[string]Property{},
		Required:   []string{},
	}
}
