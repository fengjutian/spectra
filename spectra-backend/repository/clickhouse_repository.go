package repository

import (
	"context"
	"database/sql"
	"fmt"
	"spectra-backend/config"
	"spectra-backend/models"
	"strings"
	"time"

	_ "github.com/ClickHouse/clickhouse-go/v2"
	"go.uber.org/zap"
)

// ClickHouseRepository ClickHouse仓库实现
type ClickHouseRepository struct {
	DB     *sql.DB
	Logger *zap.Logger
}

// NewClickHouseRepository 创建ClickHouse仓库实例
func NewClickHouseRepository(cfg *config.Config, logger *zap.Logger) (*ClickHouseRepository, error) {
	// 记录数据库配置信息（隐藏敏感信息）
	logger.Info("Initializing ClickHouse connection",
		zap.String("host", cfg.DB.Host),
		zap.Int("port", cfg.DB.Port),
		zap.String("database", cfg.DB.Database),
		zap.String("username", cfg.DB.Username),
		zap.Bool("debug", cfg.DB.Debug))

	// 使用https协议连接远程ClickHouse云服务
	dsn := fmt.Sprintf("https://%s:%d?database=%s&username=%s&password=%s&secure=true",
		cfg.DB.Host, cfg.DB.Port, cfg.DB.Database, cfg.DB.Username, cfg.DB.Password)

	logger.Debug("ClickHouse DSN constructed", zap.String("dsn", maskPassword(dsn)))

	db, err := sql.Open("clickhouse", dsn)
	if err != nil {
		logger.Error("Failed to open database connection", zap.Error(err))
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// 设置连接池参数
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(time.Minute * 5)

	// 测试连接
	logger.Info("Testing database connection...")
	if err := db.Ping(); err != nil {
		logger.Error("Failed to ping database", zap.Error(err))
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	logger.Info("Successfully connected to ClickHouse database")
	return &ClickHouseRepository{
		DB:     db,
		Logger: logger,
	}, nil
}

// maskPassword 隐藏DSN中的密码信息
func maskPassword(dsn string) string {
	// 简单的密码掩码函数，实际使用时可以根据需要调整
	passwordStart := strings.Index(dsn, "password=")
	if passwordStart == -1 {
		return dsn
	}
	passwordStart += len("password=")
	passwordEnd := strings.Index(dsn[passwordStart:], "&")
	if passwordEnd == -1 {
		return dsn[:passwordStart] + "***"
	}
	return dsn[:passwordStart] + "***" + dsn[passwordStart+passwordEnd:]
}

// 实现 ErrorLog 相关方法
func (r *ClickHouseRepository) SaveErrorLog(ctx context.Context, log *models.ErrorLog) error {
	query := `INSERT INTO error_logs (timestamp, project_id, session_id, trace_id, user_id, url, referrer, type, name, message, extra) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	_, err := r.DB.ExecContext(ctx, query, log.Timestamp, log.ProjectID, log.SessionID, log.TraceID, log.UserID, log.URL, log.Referrer, log.Type, log.Name, log.Message, log.Extra)
	if err != nil {
		return fmt.Errorf("failed to save error log: %w", err)
	}
	return nil
}

func (r *ClickHouseRepository) GetErrorLogs(ctx context.Context, projectID string, startTime, endTime time.Time) ([]*models.ErrorLog, error) {
	query := `SELECT timestamp, project_id, session_id, trace_id, user_id, url, referrer, type, name, message, extra FROM error_logs WHERE project_id = ? AND timestamp >= ? AND timestamp <= ? ORDER BY timestamp DESC`
	rows, err := r.DB.QueryContext(ctx, query, projectID, startTime, endTime)
	if err != nil {
		return nil, fmt.Errorf("failed to query error logs: %w", err)
	}
	defer rows.Close()

	var logs []*models.ErrorLog
	for rows.Next() {
		var log models.ErrorLog
		err := rows.Scan(&log.Timestamp, &log.ProjectID, &log.SessionID, &log.TraceID, &log.UserID, &log.URL, &log.Referrer, &log.Type, &log.Name, &log.Message, &log.Extra)
		if err != nil {
			return nil, fmt.Errorf("failed to scan error log: %w", err)
		}
		logs = append(logs, &log)
	}
	return logs, nil
}

func (r *ClickHouseRepository) GetErrorLogByTraceID(ctx context.Context, traceID string) (*models.ErrorLog, error) {
	query := `SELECT timestamp, project_id, session_id, trace_id, user_id, url, referrer, type, name, message, extra FROM error_logs WHERE trace_id = ? LIMIT 1`
	var log models.ErrorLog
	err := r.DB.QueryRowContext(ctx, query, traceID).Scan(&log.Timestamp, &log.ProjectID, &log.SessionID, &log.TraceID, &log.UserID, &log.URL, &log.Referrer, &log.Type, &log.Name, &log.Message, &log.Extra)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to query error log by traceID: %w", err)
	}
	return &log, nil
}

// 实现 PerformanceMetric 相关方法
func (r *ClickHouseRepository) SavePerformanceMetric(ctx context.Context, metric *models.PerformanceMetric) error {
	query := `INSERT INTO performance_metrics (timestamp, project_id, session_id, trace_id, user_id, url, referrer, type, name, value, extra) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	_, err := r.DB.ExecContext(ctx, query, metric.Timestamp, metric.ProjectID, metric.SessionID, metric.TraceID, metric.UserID, metric.URL, metric.Referrer, metric.Type, metric.Name, metric.Value, metric.Extra)
	if err != nil {
		return fmt.Errorf("failed to save performance metric: %w", err)
	}
	return nil
}

func (r *ClickHouseRepository) GetPerformanceMetrics(ctx context.Context, projectID string, startTime, endTime time.Time) ([]*models.PerformanceMetric, error) {
	query := `SELECT timestamp, project_id, session_id, trace_id, user_id, url, referrer, type, name, value, extra FROM performance_metrics WHERE project_id = ? AND timestamp >= ? AND timestamp <= ? ORDER BY timestamp DESC`
	rows, err := r.DB.QueryContext(ctx, query, projectID, startTime, endTime)
	if err != nil {
		return nil, fmt.Errorf("failed to query performance metrics: %w", err)
	}
	defer rows.Close()

	var metrics []*models.PerformanceMetric
	for rows.Next() {
		var metric models.PerformanceMetric
		err := rows.Scan(&metric.Timestamp, &metric.ProjectID, &metric.SessionID, &metric.TraceID, &metric.UserID, &metric.URL, &metric.Referrer, &metric.Type, &metric.Name, &metric.Value, &metric.Extra)
		if err != nil {
			return nil, fmt.Errorf("failed to scan performance metric: %w", err)
		}
		metrics = append(metrics, &metric)
	}
	return metrics, nil
}

func (r *ClickHouseRepository) GetPerformanceMetricsByType(ctx context.Context, projectID string, metricType string, startTime, endTime time.Time) ([]*models.PerformanceMetric, error) {
	query := `SELECT timestamp, project_id, session_id, trace_id, user_id, url, referrer, type, name, value, extra FROM performance_metrics WHERE project_id = ? AND name = ? AND timestamp >= ? AND timestamp <= ? ORDER BY timestamp DESC`
	rows, err := r.DB.QueryContext(ctx, query, projectID, metricType, startTime, endTime)
	if err != nil {
		return nil, fmt.Errorf("failed to query performance metrics by type: %w", err)
	}
	defer rows.Close()

	var metrics []*models.PerformanceMetric
	for rows.Next() {
		var metric models.PerformanceMetric
		err := rows.Scan(&metric.Timestamp, &metric.ProjectID, &metric.SessionID, &metric.TraceID, &metric.UserID, &metric.URL, &metric.Referrer, &metric.Type, &metric.Name, &metric.Value, &metric.Extra)
		if err != nil {
			return nil, fmt.Errorf("failed to scan performance metric: %w", err)
		}
		metrics = append(metrics, &metric)
	}
	return metrics, nil
}

// 实现 UserAction 相关方法
func (r *ClickHouseRepository) SaveUserAction(ctx context.Context, action *models.UserAction) error {
	query := `INSERT INTO user_actions (timestamp, project_id, session_id, trace_id, user_id, url, referrer, type, name, message, method, status, value, extra) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	_, err := r.DB.ExecContext(ctx, query, action.Timestamp, action.ProjectID, action.SessionID, action.TraceID, action.UserID, action.URL, action.Referrer, action.Type, action.Name, action.Message, action.Method, action.Status, action.Value, action.Extra)
	if err != nil {
		return fmt.Errorf("failed to save user action: %w", err)
	}
	return nil
}

func (r *ClickHouseRepository) GetUserActions(ctx context.Context, projectID string, startTime, endTime time.Time) ([]*models.UserAction, error) {
	query := `SELECT timestamp, project_id, session_id, trace_id, user_id, url, referrer, type, name, message, method, status, value, extra FROM user_actions WHERE project_id = ? AND timestamp >= ? AND timestamp <= ? ORDER BY timestamp DESC`
	rows, err := r.DB.QueryContext(ctx, query, projectID, startTime, endTime)
	if err != nil {
		return nil, fmt.Errorf("failed to query user actions: %w", err)
	}
	defer rows.Close()

	var actions []*models.UserAction
	for rows.Next() {
		var action models.UserAction
		err := rows.Scan(&action.Timestamp, &action.ProjectID, &action.SessionID, &action.TraceID, &action.UserID, &action.URL, &action.Referrer, &action.Type, &action.Name, &action.Message, &action.Method, &action.Status, &action.Value, &action.Extra)
		if err != nil {
			return nil, fmt.Errorf("failed to scan user action: %w", err)
		}
		actions = append(actions, &action)
	}
	return actions, nil
}

func (r *ClickHouseRepository) GetUserActionsByType(ctx context.Context, projectID string, actionType string, startTime, endTime time.Time) ([]*models.UserAction, error) {
	query := `SELECT timestamp, project_id, session_id, trace_id, user_id, url, referrer, type, name, message, method, status, value, extra FROM user_actions WHERE project_id = ? AND name = ? AND timestamp >= ? AND timestamp <= ? ORDER BY timestamp DESC`
	rows, err := r.DB.QueryContext(ctx, query, projectID, actionType, startTime, endTime)
	if err != nil {
		return nil, fmt.Errorf("failed to query user actions by type: %w", err)
	}
	defer rows.Close()

	var actions []*models.UserAction
	for rows.Next() {
		var action models.UserAction
		err := rows.Scan(&action.Timestamp, &action.ProjectID, &action.SessionID, &action.TraceID, &action.UserID, &action.URL, &action.Referrer, &action.Type, &action.Name, &action.Message, &action.Method, &action.Status, &action.Value, &action.Extra)
		if err != nil {
			return nil, fmt.Errorf("failed to scan user action: %w", err)
		}
		actions = append(actions, &action)
	}
	return actions, nil
}

// 实现 CustomEvent 相关方法
func (r *ClickHouseRepository) SaveCustomEvent(ctx context.Context, event *models.CustomEvent) error {
	query := `INSERT INTO custom_events (timestamp, project_id, session_id, trace_id, user_id, url, referrer, type, name, message, extra) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	_, err := r.DB.ExecContext(ctx, query, event.Timestamp, event.ProjectID, event.SessionID, event.TraceID, event.UserID, event.URL, event.Referrer, event.Type, event.Name, event.Message, event.Extra)
	if err != nil {
		return fmt.Errorf("failed to save custom event: %w", err)
	}
	return nil
}

func (r *ClickHouseRepository) GetCustomEvents(ctx context.Context, projectID string, startTime, endTime time.Time) ([]*models.CustomEvent, error) {
	query := `SELECT timestamp, project_id, session_id, trace_id, user_id, url, referrer, type, name, message, extra FROM custom_events WHERE project_id = ? AND timestamp >= ? AND timestamp <= ? ORDER BY timestamp DESC`
	rows, err := r.DB.QueryContext(ctx, query, projectID, startTime, endTime)
	if err != nil {
		return nil, fmt.Errorf("failed to query custom events: %w", err)
	}
	defer rows.Close()

	var events []*models.CustomEvent
	for rows.Next() {
		var event models.CustomEvent
		err := rows.Scan(&event.Timestamp, &event.ProjectID, &event.SessionID, &event.TraceID, &event.UserID, &event.URL, &event.Referrer, &event.Type, &event.Name, &event.Message, &event.Extra)
		if err != nil {
			return nil, fmt.Errorf("failed to scan custom event: %w", err)
		}
		events = append(events, &event)
	}
	return events, nil
}

func (r *ClickHouseRepository) GetCustomEventsByName(ctx context.Context, projectID string, eventName string, startTime, endTime time.Time) ([]*models.CustomEvent, error) {
	query := `SELECT timestamp, project_id, session_id, trace_id, user_id, url, referrer, type, name, message, extra FROM custom_events WHERE project_id = ? AND name = ? AND timestamp >= ? AND timestamp <= ? ORDER BY timestamp DESC`
	rows, err := r.DB.QueryContext(ctx, query, projectID, eventName, startTime, endTime)
	if err != nil {
		return nil, fmt.Errorf("failed to query custom events by name: %w", err)
	}
	defer rows.Close()

	var events []*models.CustomEvent
	for rows.Next() {
		var event models.CustomEvent
		err := rows.Scan(&event.Timestamp, &event.ProjectID, &event.SessionID, &event.TraceID, &event.UserID, &event.URL, &event.Referrer, &event.Type, &event.Name, &event.Message, &event.Extra)
		if err != nil {
			return nil, fmt.Errorf("failed to scan custom event: %w", err)
		}
		events = append(events, &event)
	}
	return events, nil
}

// 实现 PageStay 相关方法
func (r *ClickHouseRepository) SavePageStay(ctx context.Context, pageStay *models.PageStay) error {
	query := `INSERT INTO page_stay (timestamp, project_id, session_id, trace_id, user_id, url, referrer, type, name, value, extra) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	_, err := r.DB.ExecContext(ctx, query, pageStay.Timestamp, pageStay.ProjectID, pageStay.SessionID, pageStay.TraceID, pageStay.UserID, pageStay.URL, pageStay.Referrer, pageStay.Type, pageStay.Name, pageStay.Value, pageStay.Extra)
	if err != nil {
		return fmt.Errorf("failed to save page stay: %w", err)
	}
	return nil
}

func (r *ClickHouseRepository) GetPageStays(ctx context.Context, projectID string, startTime, endTime time.Time) ([]*models.PageStay, error) {
	query := `SELECT timestamp, project_id, session_id, trace_id, user_id, url, referrer, type, name, value, extra FROM page_stay WHERE project_id = ? AND timestamp >= ? AND timestamp <= ? ORDER BY timestamp DESC`
	rows, err := r.DB.QueryContext(ctx, query, projectID, startTime, endTime)
	if err != nil {
		return nil, fmt.Errorf("failed to query page stays: %w", err)
	}
	defer rows.Close()

	var stays []*models.PageStay
	for rows.Next() {
		var stay models.PageStay
		err := rows.Scan(&stay.Timestamp, &stay.ProjectID, &stay.SessionID, &stay.TraceID, &stay.UserID, &stay.URL, &stay.Referrer, &stay.Type, &stay.Name, &stay.Value, &stay.Extra)
		if err != nil {
			return nil, fmt.Errorf("failed to scan page stay: %w", err)
		}
		stays = append(stays, &stay)
	}
	return stays, nil
}

func (r *ClickHouseRepository) GetAveragePageStay(ctx context.Context, projectID string, startTime, endTime time.Time) (float64, error) {
	query := `SELECT avg(value) FROM page_stay WHERE project_id = ? AND timestamp >= ? AND timestamp <= ?`
	var avg float64
	err := r.DB.QueryRowContext(ctx, query, projectID, startTime, endTime).Scan(&avg)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, nil
		}
		return 0, fmt.Errorf("failed to query average page stay: %w", err)
	}
	return avg, nil
}

// Close 关闭数据库连接
func (r *ClickHouseRepository) Close() error {
	return r.DB.Close()
}
