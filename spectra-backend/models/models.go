package models

import (
	"encoding/json"
	"time"
)

// BaseLog 基础日志结构，包含所有表共有的字段
type BaseLog struct {
	Timestamp time.Time          `json:"timestamp"`
	ProjectID string             `json:"project_id"`
	SessionID string             `json:"session_id"`
	TraceID   string             `json:"trace_id"`
	UserID    string             `json:"user_id"`
	URL       string             `json:"url"`
	Referrer  string             `json:"referrer"`
	Type      string             `json:"type"`
	Name      string             `json:"name"`
	Extra     json.RawMessage    `json:"extra"`
}

// ErrorLog 错误日志表对应的结构体
type ErrorLog struct {
	BaseLog
	Message string `json:"message"`
}

// PerformanceMetric 性能指标表对应的结构体
type PerformanceMetric struct {
	BaseLog
	Value float64 `json:"value"`
}

// UserAction 用户行为表对应的结构体
type UserAction struct {
	BaseLog
	Message string  `json:"message"`
	Method  string  `json:"method"`
	Status  uint16  `json:"status"`
	Value   float64 `json:"value"`
}

// CustomEvent 自定义事件表对应的结构体
type CustomEvent struct {
	BaseLog
	Message string `json:"message"`
}

// PageStay 页面停留时长表对应的结构体
type PageStay struct {
	BaseLog
	Value float64 `json:"value"`
}