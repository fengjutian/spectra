package repository

import (
	"context"
	"spectra-backend/models"
	"time"
)

// LogRepository 日志存储接口
type LogRepository interface {
	// ErrorLog 相关方法
	SaveErrorLog(ctx context.Context, log *models.ErrorLog) error
	GetErrorLogs(ctx context.Context, projectID string, startTime, endTime time.Time) ([]*models.ErrorLog, error)
	GetErrorLogByTraceID(ctx context.Context, traceID string) (*models.ErrorLog, error)

	// PerformanceMetric 相关方法
	SavePerformanceMetric(ctx context.Context, metric *models.PerformanceMetric) error
	GetPerformanceMetrics(ctx context.Context, projectID string, startTime, endTime time.Time) ([]*models.PerformanceMetric, error)
	GetPerformanceMetricsByType(ctx context.Context, projectID string, metricType string, startTime, endTime time.Time) ([]*models.PerformanceMetric, error)

	// UserAction 相关方法
	SaveUserAction(ctx context.Context, action *models.UserAction) error
	GetUserActions(ctx context.Context, projectID string, startTime, endTime time.Time) ([]*models.UserAction, error)
	GetUserActionsByType(ctx context.Context, projectID string, actionType string, startTime, endTime time.Time) ([]*models.UserAction, error)

	// CustomEvent 相关方法
	SaveCustomEvent(ctx context.Context, event *models.CustomEvent) error
	GetCustomEvents(ctx context.Context, projectID string, startTime, endTime time.Time) ([]*models.CustomEvent, error)
	GetCustomEventsByName(ctx context.Context, projectID string, eventName string, startTime, endTime time.Time) ([]*models.CustomEvent, error)

	// PageStay 相关方法
	SavePageStay(ctx context.Context, pageStay *models.PageStay) error
	GetPageStays(ctx context.Context, projectID string, startTime, endTime time.Time) ([]*models.PageStay, error)
	GetAveragePageStay(ctx context.Context, projectID string, startTime, endTime time.Time) (float64, error)

	// 通用方法
	Close() error
}