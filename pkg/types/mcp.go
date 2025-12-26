package types

// MCPRequest MCP JSON-RPC请求 / MCP JSON-RPC request
type MCPRequest struct {
	// JSONRPC JSON-RPC版本 / JSON-RPC version
	JSONRPC string `json:"jsonrpc"`

	// ID 请求ID / Request ID
	ID interface{} `json:"id,omitempty"`

	// Method 方法名 / Method name
	Method string `json:"method"`

	// Params 参数 / Parameters
	Params interface{} `json:"params,omitempty"`
}

// MCPResponse MCP JSON-RPC响应 / MCP JSON-RPC response
type MCPResponse struct {
	// JSONRPC JSON-RPC版本 / JSON-RPC version
	JSONRPC string `json:"jsonrpc"`

	// ID 请求ID / Request ID
	ID interface{} `json:"id,omitempty"`

	// Result 结果 / Result
	Result interface{} `json:"result,omitempty"`

	// Error 错误 / Error
	Error *MCPError `json:"error,omitempty"`
}

// MCPError MCP错误 / MCP error
type MCPError struct {
	// Code 错误码 / Error code
	Code int `json:"code"`

	// Message 错误消息 / Error message
	Message string `json:"message"`

	// Data 错误数据 / Error data
	Data interface{} `json:"data,omitempty"`
}

// MCP错误码常量 / MCP error code constants
const (
	// MCPErrorCodeParseError 解析错误 / Parse error
	MCPErrorCodeParseError = -32700

	// MCPErrorCodeInvalidRequest 无效请求 / Invalid request
	MCPErrorCodeInvalidRequest = -32600

	// MCPErrorCodeMethodNotFound 方法未找到 / Method not found
	MCPErrorCodeMethodNotFound = -32601

	// MCPErrorCodeInvalidParams 无效参数 / Invalid params
	MCPErrorCodeInvalidParams = -32602

	// MCPErrorCodeInternalError 内部错误 / Internal error
	MCPErrorCodeInternalError = -32603
)

// MCP错误消息常量 / MCP error message constants
const (
	// MCPErrorMsgParseError 解析错误消息 / Parse error message
	MCPErrorMsgParseError = "Parse error"

	// MCPErrorMsgInvalidRequest 无效请求消息 / Invalid request message
	MCPErrorMsgInvalidRequest = "Invalid request"

	// MCPErrorMsgMethodNotFound 方法未找到消息 / Method not found message
	MCPErrorMsgMethodNotFound = "Method not found"

	// MCPErrorMsgInvalidParams 无效参数消息 / Invalid params message
	MCPErrorMsgInvalidParams = "Invalid params"

	// MCPErrorMsgInternalError 内部错误消息 / Internal error message
	MCPErrorMsgInternalError = "Internal error"
)

// NewMCPError 创建MCP错误 / Create MCP error
func NewMCPError(code int, message string, data interface{}) *MCPError {
	return &MCPError{
		Code:    code,
		Message: message,
		Data:    data,
	}
}

// NewMCPResponse 创建MCP响应 / Create MCP response
func NewMCPResponse(id interface{}, result interface{}) *MCPResponse {
	return &MCPResponse{
		JSONRPC: "2.0",
		ID:      id,
		Result:  result,
	}
}

// NewMCPErrorResponse 创建MCP错误响应 / Create MCP error response
func NewMCPErrorResponse(id interface{}, err *MCPError) *MCPResponse {
	return &MCPResponse{
		JSONRPC: "2.0",
		ID:      id,
		Error:   err,
	}
}

// MCPToolsListRequest 工具列表请求参数 / Tools list request params
type MCPToolsListRequest struct {
	// Cursor 分页游标 / Pagination cursor
	Cursor string `json:"cursor,omitempty"`
}

// MCPToolsCallRequest 工具调用请求参数 / Tools call request params
type MCPToolsCallRequest struct {
	// Name 工具名称 / Tool name
	Name string `json:"name"`

	// Arguments 工具参数 / Tool arguments
	Arguments interface{} `json:"arguments,omitempty"`
}

// MCPInitializeRequest 初始化请求参数 / Initialize request params
type MCPInitializeRequest struct {
	// ProtocolVersion 协议版本 / Protocol version
	ProtocolVersion string `json:"protocolVersion"`

	// Capabilities 客户端能力 / Client capabilities
	Capabilities interface{} `json:"capabilities,omitempty"`

	// ClientInfo 客户端信息 / Client info
	ClientInfo interface{} `json:"clientInfo,omitempty"`
}

// SSEMessage SSE消息 / SSE message
type SSEMessage struct {
	// Event 事件类型 / Event type
	Event string `json:"event,omitempty"`

	// Data 消息数据 / Message data
	Data string `json:"data"`

	// ID 消息ID / Message ID
	ID string `json:"id,omitempty"`

	// Retry 重试时间(毫秒) / Retry time in milliseconds
	Retry int `json:"retry,omitempty"`
}
