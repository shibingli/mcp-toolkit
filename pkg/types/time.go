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

// Package types 时间相关类型定义 / Time related type definitions
package types

// GetCurrentTimeRequest 获取当前时间请求 / Get current time request
type GetCurrentTimeRequest struct {
	TimeZone string `json:"timezone,omitempty"` // 时区（可选，如 "Asia/Shanghai"、"America/New_York"），为空则使用系统本地时区 / Time zone (optional, e.g. "Asia/Shanghai", "America/New_York"), empty means use system local timezone
}

// GetCurrentTimeResponse 获取当前时间响应 / Get current time response
type GetCurrentTimeResponse = GetTimeResponse

// GetTimeResponse 获取时间响应 / Get time response
type GetTimeResponse struct {
	DateTime       string `json:"datetime"`         // 格式化的日期时间字符串 / Formatted datetime string
	Date           string `json:"date"`             // 日期 (YYYY-MM-DD) / Date (YYYY-MM-DD)
	Time           string `json:"time"`             // 时间 (HH:MM:SS) / Time (HH:MM:SS)
	TimeZone       string `json:"timezone"`         // 时区名称 / Time zone name
	TimeZoneOffset string `json:"timezone_offset"`  // 时区偏移 (如 +08:00) / Time zone offset (e.g. +08:00)
	Unix           int64  `json:"unix"`             // Unix时间戳（秒）/ Unix timestamp (seconds)
	UnixMilli      int64  `json:"unix_milli"`       // Unix时间戳（毫秒）/ Unix timestamp (milliseconds)
	Weekday        string `json:"weekday"`          // 星期几 / Day of week
	IsDST          bool   `json:"is_dst,omitempty"` // 是否夏令时 / Is daylight saving time
}
