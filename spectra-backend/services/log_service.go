package services

import (
	"context"
	"spectra-backend/models"
	"spectra-backend/repository"
	"time"
)

// LogService 日志服务接口
type LogService interface {
	// ErrorLog 相关服务
	RecordErrorLog(ctx context.Context, log *models.ErrorLog) error
	GetErrorLogs(ctx context.Context, projectID string, startTime, endTime time.Time) ([]*models.ErrorLog, error)
	GetErrorLogByTraceID(ctx context.Context, traceID string) (*models.ErrorLog, error)

	// PerformanceMetric 相关服务
	RecordPerformanceMetric(ctx context.Context, metric *models.PerformanceMetric) error
	GetPerformanceMetrics(ctx context.Context, projectID string, startTime, endTime time.Time) ([]*models.PerformanceMetric, error)
	GetPerformanceMetricsByType(ctx context.Context, projectID string, metricType string, startTime, endTime time.Time) ([]*models.PerformanceMetric, error)

	// UserAction 相关服务
	RecordUserAction(ctx context.Context, action *models.UserAction) error
	GetUserActions(ctx context.Context, projectID string, startTime, endTime time.Time) ([]*models.UserAction, error)
	GetUserActionsByType(ctx context.Context, projectID string, actionType string, startTime, endTime time.Time) ([]*models.UserAction, error)

	// CustomEvent 相关服务
	RecordCustomEvent(ctx context.Context, event *models.CustomEvent) error
	GetCustomEvents(ctx context.Context, projectID string, startTime, endTime time.Time) ([]*models.CustomEvent, error)
	GetCustomEventsByName(ctx context.Context, projectID string, eventName string, startTime, endTime time.Time) ([]*models.CustomEvent, error)

	// PageStay 相关服务
	RecordPageStay(ctx context.Context, pageStay *models.PageStay) error
	GetPageStays(ctx context.Context, projectID string, startTime, endTime time.Time) ([]*models.PageStay, error)
	GetAveragePageStay(ctx context.Context, projectID string, startTime, endTime time.Time) (float64, error)
}

// logService 日志服务实现
type logService struct {
	repo repository.LogRepository
}

// NewLogService 创建日志服务实例
func NewLogService(repo repository.LogRepository) LogService {
	return &logService{
		repo: repo,
	}
}

// 实现 ErrorLog 相关方法
func (s *logService) RecordErrorLog(ctx context.Context, log *models.ErrorLog) error {
	if log.Timestamp.IsZero() {
		log.Timestamp = time.Now()
	}
	if log.Type == "" {
		log.Type = "error"
	}
	return s.repo.SaveErrorLog(ctx, log)
}

func (s *logService) GetErrorLogs(ctx context.Context, projectID string, startTime, endTime time.Time) ([]*models.ErrorLog, error) {
	return s.repo.GetErrorLogs(ctx, projectID, startTime, endTime)
}

func (s *logService) GetErrorLogByTraceID(ctx context.Context, traceID string) (*models.ErrorLog, error) {
	return s.repo.GetErrorLogByTraceID(ctx, traceID)
}

// 实现 PerformanceMetric 相关方法
func (s *logService) RecordPerformanceMetric(ctx context.Context, metric *models.PerformanceMetric) error {
	if metric.Timestamp.IsZero() {
		metric.Timestamp = time.Now()
	}
	if metric.Type == "" {
		metric.Type = "performance"
	}
	return s.repo.SavePerformanceMetric(ctx, metric)
}

func (s *logService) GetPerformanceMetrics(ctx context.Context, projectID string, startTime, endTime time.Time) ([]*models.PerformanceMetric, error) {
	return s.repo.GetPerformanceMetrics(ctx, projectID, startTime, endTime)
}

func (s *logService) GetPerformanceMetricsByType(ctx context.Context, projectID string, metricType string, startTime, endTime time.Time) ([]*models.PerformanceMetric, error) {
	return s.repo.GetPerformanceMetricsByType(ctx, projectID, metricType, startTime, endTime)
}

// 实现 UserAction 相关方法
func (s *logService) RecordUserAction(ctx context.Context, action *models.UserAction) error {
	if action.Timestamp.IsZero() {
		action.Timestamp = time.Now()
	}
	if action.Type == "" {
		action.Type = "user"
	}
	return s.repo.SaveUserAction(ctx, action)
}

func (s *logService) GetUserActions(ctx context.Context, projectID string, startTime, endTime time.Time) ([]*models.UserAction, error) {
	return s.repo.GetUserActions(ctx, projectID, startTime, endTime)
}

func (s *logService) GetUserActionsByType(ctx context.Context, projectID string, actionType string, startTime, endTime time.Time) ([]*models.UserAction, error) {
	return s.repo.GetUserActionsByType(ctx, projectID, actionType, startTime, endTime)
}

// 实现 CustomEvent 相关方法
func (s *logService) RecordCustomEvent(ctx context.Context, event *models.CustomEvent) error {
	if event.Timestamp.IsZero() {
		event.Timestamp = time.Now()
	}
	if event.Type == "" {
		event.Type = "custom"
	}
	if event.Message == "" {
		event.Message = "custom_event"
	}
	return s.repo.SaveCustomEvent(ctx, event)
}

func (s *logService) GetCustomEvents(ctx context.Context, projectID string, startTime, endTime time.Time) ([]*models.CustomEvent, error) {
	return s.repo.GetCustomEvents(ctx, projectID, startTime, endTime)
}

func (s *logService) GetCustomEventsByName(ctx context.Context, projectID string, eventName string, startTime, endTime time.Time) ([]*models.CustomEvent, error) {
	return s.repo.GetCustomEventsByName(ctx, projectID, eventName, startTime, endTime)
}

// 实现 PageStay 相关方法
func (s *logService) RecordPageStay(ctx context.Context, pageStay *models.PageStay) error {
	if pageStay.Timestamp.IsZero() {
		pageStay.Timestamp = time.Now()
	}
	if pageStay.Type == "" {
		pageStay.Type = "page_stay"
	}
	if pageStay.Name == "" {
		pageStay.Name = "page_stay_time"
	}
	return s.repo.SavePageStay(ctx, pageStay)
}

func (s *logService) GetPageStays(ctx context.Context, projectID string, startTime, endTime time.Time) ([]*models.PageStay, error) {
	return s.repo.GetPageStays(ctx, projectID, startTime, endTime)
}

func (s *logService) GetAveragePageStay(ctx context.Context, projectID string, startTime, endTime time.Time) (float64, error) {
	return s.repo.GetAveragePageStay(ctx, projectID, startTime, endTime)
}