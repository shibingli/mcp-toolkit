package client

// InitializeResponse 初始化响应 / Initialize response
type InitializeResponse struct {
	// ProtocolVersion 协议版本 / Protocol version
	ProtocolVersion string `json:"protocolVersion"`

	// Capabilities 服务器能力 / Server capabilities
	Capabilities ServerCapabilities `json:"capabilities"`

	// ServerInfo 服务器信息 / Server info
	ServerInfo ServerInfo `json:"serverInfo"`
}

// ServerCapabilities 服务器能力 / Server capabilities
type ServerCapabilities struct {
	// Tools 工具能力 / Tools capability
	Tools *ToolsCapability `json:"tools,omitempty"`

	// Resources 资源能力 / Resources capability
	Resources *ResourcesCapability `json:"resources,omitempty"`

	// Prompts 提示能力 / Prompts capability
	Prompts *PromptsCapability `json:"prompts,omitempty"`
}

// ToolsCapability 工具能力 / Tools capability
type ToolsCapability struct {
	// ListChanged 工具列表是否可能改变 / Whether tool list may change
	ListChanged bool `json:"listChanged,omitempty"`
}

// ResourcesCapability 资源能力 / Resources capability
type ResourcesCapability struct {
	// Subscribe 是否支持订阅 / Whether subscription is supported
	Subscribe bool `json:"subscribe,omitempty"`

	// ListChanged 资源列表是否可能改变 / Whether resource list may change
	ListChanged bool `json:"listChanged,omitempty"`
}

// PromptsCapability 提示能力 / Prompts capability
type PromptsCapability struct {
	// ListChanged 提示列表是否可能改变 / Whether prompt list may change
	ListChanged bool `json:"listChanged,omitempty"`
}

// ServerInfo 服务器信息 / Server info
type ServerInfo struct {
	// Name 服务器名称 / Server name
	Name string `json:"name"`

	// Version 服务器版本 / Server version
	Version string `json:"version"`
}

// ListToolsResponse 工具列表响应 / List tools response
type ListToolsResponse struct {
	// Tools 工具列表 / Tools list
	Tools []Tool `json:"tools"`

	// NextCursor 下一页游标 / Next page cursor
	NextCursor string `json:"nextCursor,omitempty"`
}

// Tool 工具定义 / Tool definition
type Tool struct {
	// Name 工具名称 / Tool name
	Name string `json:"name"`

	// Description 工具描述 / Tool description
	Description string `json:"description,omitempty"`

	// InputSchema 输入模式 / Input schema
	InputSchema map[string]interface{} `json:"inputSchema,omitempty"`
}

// CallToolResponse 工具调用响应 / Call tool response
type CallToolResponse struct {
	// Content 响应内容 / Response content
	Content []Content `json:"content"`

	// IsError 是否为错误 / Whether it is an error
	IsError bool `json:"isError,omitempty"`
}

// Content 内容 / Content
type Content struct {
	// Type 内容类型 / Content type
	Type string `json:"type"`

	// Text 文本内容(当type为text时) / Text content (when type is text)
	Text string `json:"text,omitempty"`

	// Data 数据内容(当type为其他类型时) / Data content (when type is other)
	Data interface{} `json:"data,omitempty"`
}
