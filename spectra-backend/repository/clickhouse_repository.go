package repository

import (
	"context"
	"database/sql"
	"fmt"
	"spectra-backend/config"
	"spectra-backend/models"
	"strings"
	"time"

	// 匿名导入 ClickHouse 驱动以确保驱动被正确注册
	_ "github.com/ClickHouse/clickhouse-go/v2"
	"go.uber.org/zap"
)

// Repository 接口定义了所有数据访问操作
// 这是一个接口，ClickHouseRepository 是它的具体实现

// ClickHouseRepository 是Repository接口的ClickHouse具体实现
// 负责与ClickHouse数据库进行交互，执行所有数据存取操作
type ClickHouseRepository struct {
	DB     *sql.DB     // 数据库连接对象
	Logger *zap.Logger // 日志记录器
}

// NewClickHouseRepository 创建ClickHouse仓库实例
// 参数:
//   - cfg: 应用程序配置，包含数据库连接信息
//   - logger: 日志记录器实例
// 返回:
//   - *ClickHouseRepository: 初始化成功的仓库实例
//   - error: 初始化过程中的错误信息
func NewClickHouseRepository(cfg *config.Config, logger *zap.Logger) (*ClickHouseRepository, error) {
	// 记录数据库配置信息（隐藏敏感信息）
	logger.Info("Initializing ClickHouse connection",
		zap.String("host", cfg.DB.Host),
		zap.Int("port", cfg.DB.Port),
		zap.String("database", cfg.DB.Database),
		zap.String("username", cfg.DB.Username),
		zap.Bool("debug", cfg.DB.Debug))

	// 构建DSN连接字符串，使用https协议连接远程ClickHouse云服务
	// 包含主机、端口、数据库名、用户名、密码和安全连接参数
	dsn := fmt.Sprintf("https://%s:%d?database=%s&username=%s&password=%s&secure=true",
		cfg.DB.Host, cfg.DB.Port, cfg.DB.Database, cfg.DB.Username, cfg.DB.Password)

	// 记录构建好的DSN（已隐藏密码）用于调试
	logger.Debug("ClickHouse DSN constructed", zap.String("dsn", maskPassword(dsn)))

	// 打开数据库连接
	db, err := sql.Open("clickhouse", dsn)
	if err != nil {
		logger.Error("Failed to open database connection", zap.Error(err))
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// 设置连接池参数
	db.SetMaxOpenConns(10)       // 最大打开连接数
	db.SetMaxIdleConns(5)        // 最大空闲连接数
	db.SetConnMaxLifetime(time.Minute * 5) // 连接最大生命周期

	// 测试数据库连接是否正常
	logger.Info("Testing database connection...")
	if err := db.Ping(); err != nil {
		logger.Error("Failed to ping database", zap.Error(err))
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	logger.Info("Successfully connected to ClickHouse database")
	// 返回初始化成功的仓库实例
	return &ClickHouseRepository{
		DB:     db,
		Logger: logger,
	}, nil
}

// maskPassword 隐藏DSN中的密码信息，避免敏感数据泄露到日志中
// 参数:
//   - dsn: 原始DSN字符串
// 返回:
//   - string: 密码被替换为***的安全DSN字符串
func maskPassword(dsn string) string {
	// 查找password=在DSN字符串中的位置
	passwordStart := strings.Index(dsn, "password=")
	if passwordStart == -1 {
		// 如果没有找到password=，则直接返回原字符串
		return dsn
	}
	
	// 计算密码起始位置（跳过"password="部分）
	passwordStart += len("password=")
	
	// 查找密码后的&分隔符位置
	passwordEnd := strings.Index(dsn[passwordStart:], "&")
	if passwordEnd == -1 {
		// 如果&不存在，说明密码在字符串末尾
		return dsn[:passwordStart] + "***"
	}
	
	// 替换密码部分为***并返回
	return dsn[:passwordStart] + "***" + dsn[passwordStart+passwordEnd:]
}

// SaveErrorLog 保存错误日志到数据库
// 参数:
//   - ctx: 上下文对象，用于控制请求超时和取消
//   - log: 错误日志对象，包含要保存的错误信息
// 返回:
//   - error: 保存过程中的错误信息，成功则为nil
func (r *ClickHouseRepository) SaveErrorLog(ctx context.Context, log *models.ErrorLog) error {
	// 定义SQL插入语句，包含错误日志的所有字段
	query := `INSERT INTO error_logs (timestamp, project_id, session_id, trace_id, user_id, url, referrer, type, name, message, extra) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	
	// 将Extra字段转换为字符串，如果为空则使用空JSON对象
	extraStr := string(log.Extra)
	if extraStr == "" {
		extraStr = "{}"
	}
	
	// 执行插入操作，使用ExecContext支持上下文取消和超时
	_, err := r.DB.ExecContext(ctx, query, 
		log.Timestamp, log.ProjectID, log.SessionID, log.TraceID, log.UserID,
		log.URL, log.Referrer, log.Type, log.Name, log.Message, extraStr)
	if err != nil {
		return fmt.Errorf("failed to save error log: %w", err)
	}
	return nil
}

// GetErrorLogs 获取指定项目在时间范围内的错误日志列表
// 参数:
//   - ctx: 上下文对象，用于控制请求超时和取消
//   - projectID: 项目标识符
//   - startTime: 开始时间
//   - endTime: 结束时间
// 返回:
//   - []*models.ErrorLog: 错误日志列表
//   - error: 查询过程中的错误信息，成功则为nil
func (r *ClickHouseRepository) GetErrorLogs(ctx context.Context, projectID string, startTime, endTime time.Time) ([]*models.ErrorLog, error) {
	// 定义SQL查询语句，按时间倒序排列
	query := `SELECT timestamp, project_id, session_id, trace_id, user_id, url, referrer, type, name, message, extra 
		FROM error_logs 
		WHERE project_id = ? AND timestamp >= ? AND timestamp <= ? 
		ORDER BY timestamp DESC`
	
	// 执行查询，使用QueryContext支持上下文取消和超时
	rows, err := r.DB.QueryContext(ctx, query, projectID, startTime, endTime)
	if err != nil {
		return nil, fmt.Errorf("failed to query error logs: %w", err)
	}
	defer rows.Close() // 确保查询结果集在函数返回前被关闭

	var logs []*models.ErrorLog
	// 遍历查询结果，将每一行数据扫描到ErrorLog结构体中
	for rows.Next() {
		var log models.ErrorLog
		err := rows.Scan(
			&log.Timestamp, &log.ProjectID, &log.SessionID, &log.TraceID, &log.UserID,
			&log.URL, &log.Referrer, &log.Type, &log.Name, &log.Message, &log.Extra)
		if err != nil {
			return nil, fmt.Errorf("failed to scan error log: %w", err)
		}
		logs = append(logs, &log)
	}
	return logs, nil
}

// GetErrorLogByTraceID 根据traceID获取特定的错误日志
// 参数:
//   - ctx: 上下文对象，用于控制请求超时和取消
//   - traceID: 唯一的跟踪标识符
// 返回:
//   - *models.ErrorLog: 错误日志对象，如果不存在则为nil
//   - error: 查询过程中的错误信息，成功或未找到则为nil
func (r *ClickHouseRepository) GetErrorLogByTraceID(ctx context.Context, traceID string) (*models.ErrorLog, error) {
	// 定义SQL查询语句，使用LIMIT 1确保只返回一个结果
	query := `SELECT timestamp, project_id, session_id, trace_id, user_id, url, referrer, type, name, message, extra 
		FROM error_logs 
		WHERE trace_id = ? 
		LIMIT 1`
	
	var log models.ErrorLog
	// 使用QueryRowContext执行查询并直接扫描结果
	err := r.DB.QueryRowContext(ctx, query, traceID).Scan(
		&log.Timestamp, &log.ProjectID, &log.SessionID, &log.TraceID, &log.UserID,
		&log.URL, &log.Referrer, &log.Type, &log.Name, &log.Message, &log.Extra)
	
	if err != nil {
		if err == sql.ErrNoRows {
			// 如果没有找到记录，返回nil, nil
			return nil, nil
		}
		return nil, fmt.Errorf("failed to query error log by traceID: %w", err)
	}
	return &log, nil
}

// SavePerformanceMetric 保存性能指标数据到数据库
// 参数:
//   - ctx: 上下文对象，用于控制请求超时和取消
//   - metric: 性能指标对象，包含要保存的性能数据
// 返回:
//   - error: 保存过程中的错误信息，成功则为nil
func (r *ClickHouseRepository) SavePerformanceMetric(ctx context.Context, metric *models.PerformanceMetric) error {
	// 定义SQL插入语句，包含性能指标的所有字段
	query := `INSERT INTO performance_metrics (timestamp, project_id, session_id, trace_id, user_id, url, referrer, type, name, value, extra) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	
	// 将Extra字段转换为字符串，如果为空则使用空JSON对象
	extraStr := string(metric.Extra)
	if extraStr == "" {
		extraStr = "{}"
	}
	
	// 执行插入操作
	_, err := r.DB.ExecContext(ctx, query, 
		metric.Timestamp, metric.ProjectID, metric.SessionID, metric.TraceID, metric.UserID,
		metric.URL, metric.Referrer, metric.Type, metric.Name, metric.Value, extraStr)
	if err != nil {
		return fmt.Errorf("failed to save performance metric: %w", err)
	}
	return nil
}

// GetPerformanceMetrics 获取指定项目在时间范围内的性能指标列表
// 参数:
//   - ctx: 上下文对象，用于控制请求超时和取消
//   - projectID: 项目标识符
//   - startTime: 开始时间
//   - endTime: 结束时间
// 返回:
//   - []*models.PerformanceMetric: 性能指标列表
//   - error: 查询过程中的错误信息，成功则为nil
func (r *ClickHouseRepository) GetPerformanceMetrics(ctx context.Context, projectID string, startTime, endTime time.Time) ([]*models.PerformanceMetric, error) {
	// 定义SQL查询语句，按时间倒序排列
	query := `SELECT timestamp, project_id, session_id, trace_id, user_id, url, referrer, type, name, value, extra 
		FROM performance_metrics 
		WHERE project_id = ? AND timestamp >= ? AND timestamp <= ? 
		ORDER BY timestamp DESC`
	
	// 执行查询
	rows, err := r.DB.QueryContext(ctx, query, projectID, startTime, endTime)
	if err != nil {
		return nil, fmt.Errorf("failed to query performance metrics: %w", err)
	}
	defer rows.Close()

	var metrics []*models.PerformanceMetric
	// 遍历查询结果
	for rows.Next() {
		var metric models.PerformanceMetric
		err := rows.Scan(
			&metric.Timestamp, &metric.ProjectID, &metric.SessionID, &metric.TraceID, &metric.UserID,
			&metric.URL, &metric.Referrer, &metric.Type, &metric.Name, &metric.Value, &metric.Extra)
		if err != nil {
			return nil, fmt.Errorf("failed to scan performance metric: %w", err)
		}
		metrics = append(metrics, &metric)
	}
	return metrics, nil
}

// GetPerformanceMetricsByType 获取指定项目、指定类型在时间范围内的性能指标
// 参数:
//   - ctx: 上下文对象，用于控制请求超时和取消
//   - projectID: 项目标识符
//   - metricType: 性能指标类型名称
//   - startTime: 开始时间
//   - endTime: 结束时间
// 返回:
//   - []*models.PerformanceMetric: 性能指标列表
//   - error: 查询过程中的错误信息，成功则为nil
func (r *ClickHouseRepository) GetPerformanceMetricsByType(ctx context.Context, projectID string, metricType string, startTime, endTime time.Time) ([]*models.PerformanceMetric, error) {
	// 定义SQL查询语句，按类型和时间范围筛选，时间倒序排列
	query := `SELECT timestamp, project_id, session_id, trace_id, user_id, url, referrer, type, name, value, extra 
		FROM performance_metrics 
		WHERE project_id = ? AND name = ? AND timestamp >= ? AND timestamp <= ? 
		ORDER BY timestamp DESC`
	
	// 执行查询
	rows, err := r.DB.QueryContext(ctx, query, projectID, metricType, startTime, endTime)
	if err != nil {
		return nil, fmt.Errorf("failed to query performance metrics by type: %w", err)
	}
	defer rows.Close()

	var metrics []*models.PerformanceMetric
	// 遍历查询结果
	for rows.Next() {
		var metric models.PerformanceMetric
		err := rows.Scan(
			&metric.Timestamp, &metric.ProjectID, &metric.SessionID, &metric.TraceID, &metric.UserID,
			&metric.URL, &metric.Referrer, &metric.Type, &metric.Name, &metric.Value, &metric.Extra)
		if err != nil {
			return nil, fmt.Errorf("failed to scan performance metric: %w", err)
		}
		metrics = append(metrics, &metric)
	}
	return metrics, nil
}

// SaveUserAction 保存用户行为数据到数据库
// 参数:
//   - ctx: 上下文对象，用于控制请求超时和取消
//   - action: 用户行为对象，包含要保存的用户交互数据
// 返回:
//   - error: 保存过程中的错误信息，成功则为nil
func (r *ClickHouseRepository) SaveUserAction(ctx context.Context, action *models.UserAction) error {
	// 定义SQL插入语句，包含用户行为的所有字段
	query := `INSERT INTO user_actions (timestamp, project_id, session_id, trace_id, user_id, url, referrer, type, name, message, method, status, value, extra) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	
	// 将Extra字段转换为字符串，如果为空则使用空JSON对象
	extraStr := string(action.Extra)
	if extraStr == "" {
		extraStr = "{}"
	}
	
	// 执行插入操作
	_, err := r.DB.ExecContext(ctx, query, 
		action.Timestamp, action.ProjectID, action.SessionID, action.TraceID, action.UserID,
		action.URL, action.Referrer, action.Type, action.Name, action.Message, action.Method,
		action.Status, action.Value, extraStr)
	if err != nil {
		return fmt.Errorf("failed to save user action: %w", err)
	}
	return nil
}

// GetUserActions 获取指定项目在时间范围内的用户行为列表
// 参数:
//   - ctx: 上下文对象，用于控制请求超时和取消
//   - projectID: 项目标识符
//   - startTime: 开始时间
//   - endTime: 结束时间
// 返回:
//   - []*models.UserAction: 用户行为列表
//   - error: 查询过程中的错误信息，成功则为nil
func (r *ClickHouseRepository) GetUserActions(ctx context.Context, projectID string, startTime, endTime time.Time) ([]*models.UserAction, error) {
	// 定义SQL查询语句，按时间倒序排列
	query := `SELECT timestamp, project_id, session_id, trace_id, user_id, url, referrer, type, name, message, method, status, value, extra 
		FROM user_actions 
		WHERE project_id = ? AND timestamp >= ? AND timestamp <= ? 
		ORDER BY timestamp DESC`
	
	// 执行查询
	rows, err := r.DB.QueryContext(ctx, query, projectID, startTime, endTime)
	if err != nil {
		return nil, fmt.Errorf("failed to query user actions: %w", err)
	}
	defer rows.Close()

	var actions []*models.UserAction
	// 遍历查询结果
	for rows.Next() {
		var action models.UserAction
		err := rows.Scan(
			&action.Timestamp, &action.ProjectID, &action.SessionID, &action.TraceID, &action.UserID,
			&action.URL, &action.Referrer, &action.Type, &action.Name, &action.Message, action.Method,
			action.Status, action.Value, &action.Extra)
		if err != nil {
			return nil, fmt.Errorf("failed to scan user action: %w", err)
		}
		actions = append(actions, &action)
	}
	return actions, nil
}

// GetUserActionsByType 获取指定项目、指定类型在时间范围内的用户行为
// 参数:
//   - ctx: 上下文对象，用于控制请求超时和取消
//   - projectID: 项目标识符
//   - actionType: 用户行为类型名称
//   - startTime: 开始时间
//   - endTime: 结束时间
// 返回:
//   - []*models.UserAction: 用户行为列表
//   - error: 查询过程中的错误信息，成功则为nil
func (r *ClickHouseRepository) GetUserActionsByType(ctx context.Context, projectID string, actionType string, startTime, endTime time.Time) ([]*models.UserAction, error) {
	// 定义SQL查询语句，按类型和时间范围筛选，时间倒序排列
	query := `SELECT timestamp, project_id, session_id, trace_id, user_id, url, referrer, type, name, message, method, status, value, extra 
		FROM user_actions 
		WHERE project_id = ? AND name = ? AND timestamp >= ? AND timestamp <= ? 
		ORDER BY timestamp DESC`
	
	// 执行查询
	rows, err := r.DB.QueryContext(ctx, query, projectID, actionType, startTime, endTime)
	if err != nil {
		return nil, fmt.Errorf("failed to query user actions by type: %w", err)
	}
	defer rows.Close()

	var actions []*models.UserAction
	// 遍历查询结果
	for rows.Next() {
		var action models.UserAction
		err := rows.Scan(
			&action.Timestamp, &action.ProjectID, &action.SessionID, &action.TraceID, &action.UserID,
			&action.URL, &action.Referrer, &action.Type, &action.Name, &action.Message, action.Method,
			action.Status, action.Value, &action.Extra)
		if err != nil {
			return nil, fmt.Errorf("failed to scan user action: %w", err)
		}
		actions = append(actions, &action)
	}
	return actions, nil
}

// SaveCustomEvent 保存自定义事件数据到数据库
// 参数:
//   - ctx: 上下文对象，用于控制请求超时和取消
//   - event: 自定义事件对象，包含要保存的自定义事件数据
// 返回:
//   - error: 保存过程中的错误信息，成功则为nil
func (r *ClickHouseRepository) SaveCustomEvent(ctx context.Context, event *models.CustomEvent) error {
	// 定义SQL插入语句，包含自定义事件的所有字段
	query := `INSERT INTO custom_events (timestamp, project_id, session_id, trace_id, user_id, url, referrer, type, name, message, extra) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	
	// 将Extra字段转换为字符串，如果为空则使用空JSON对象
	extraStr := string(event.Extra)
	if extraStr == "" {
		extraStr = "{}"
	}
	
	// 执行插入操作
	_, err := r.DB.ExecContext(ctx, query, 
		event.Timestamp, event.ProjectID, event.SessionID, event.TraceID, event.UserID,
		event.URL, event.Referrer, event.Type, event.Name, event.Message, extraStr)
	if err != nil {
		return fmt.Errorf("failed to save custom event: %w", err)
	}
	return nil
}

// GetCustomEvents 获取指定项目在时间范围内的自定义事件列表
// 参数:
//   - ctx: 上下文对象，用于控制请求超时和取消
//   - projectID: 项目标识符
//   - startTime: 开始时间
//   - endTime: 结束时间
// 返回:
//   - []*models.CustomEvent: 自定义事件列表
//   - error: 查询过程中的错误信息，成功则为nil
func (r *ClickHouseRepository) GetCustomEvents(ctx context.Context, projectID string, startTime, endTime time.Time) ([]*models.CustomEvent, error) {
	// 定义SQL查询语句，按时间倒序排列
	query := `SELECT timestamp, project_id, session_id, trace_id, user_id, url, referrer, type, name, message, extra 
		FROM custom_events 
		WHERE project_id = ? AND timestamp >= ? AND timestamp <= ? 
		ORDER BY timestamp DESC`
	
	// 执行查询
	rows, err := r.DB.QueryContext(ctx, query, projectID, startTime, endTime)
	if err != nil {
		return nil, fmt.Errorf("failed to query custom events: %w", err)
	}
	defer rows.Close()

	var events []*models.CustomEvent
	// 遍历查询结果
	for rows.Next() {
		var event models.CustomEvent
		err := rows.Scan(
			&event.Timestamp, &event.ProjectID, &event.SessionID, &event.TraceID, &event.UserID,
			&event.URL, &event.Referrer, &event.Type, &event.Name, &event.Message, &event.Extra)
		if err != nil {
			return nil, fmt.Errorf("failed to scan custom event: %w", err)
		}
		events = append(events, &event)
	}
	return events, nil
}

// GetCustomEventsByName 获取指定项目、指定名称在时间范围内的自定义事件
// 参数:
//   - ctx: 上下文对象，用于控制请求超时和取消
//   - projectID: 项目标识符
//   - eventName: 自定义事件名称
//   - startTime: 开始时间
//   - endTime: 结束时间
// 返回:
//   - []*models.CustomEvent: 自定义事件列表
//   - error: 查询过程中的错误信息，成功则为nil
func (r *ClickHouseRepository) GetCustomEventsByName(ctx context.Context, projectID string, eventName string, startTime, endTime time.Time) ([]*models.CustomEvent, error) {
	// 定义SQL查询语句，按名称和时间范围筛选，时间倒序排列
	query := `SELECT timestamp, project_id, session_id, trace_id, user_id, url, referrer, type, name, message, extra 
		FROM custom_events 
		WHERE project_id = ? AND name = ? AND timestamp >= ? AND timestamp <= ? 
		ORDER BY timestamp DESC`
	
	// 执行查询
	rows, err := r.DB.QueryContext(ctx, query, projectID, eventName, startTime, endTime)
	if err != nil {
		return nil, fmt.Errorf("failed to query custom events by name: %w", err)
	}
	defer rows.Close()

	var events []*models.CustomEvent
	// 遍历查询结果
	for rows.Next() {
		var event models.CustomEvent
		err := rows.Scan(
			&event.Timestamp, &event.ProjectID, &event.SessionID, &event.TraceID, &event.UserID,
			&event.URL, &event.Referrer, &event.Type, &event.Name, &event.Message, &event.Extra)
		if err != nil {
			return nil, fmt.Errorf("failed to scan custom event: %w", err)
		}
		events = append(events, &event)
	}
	return events, nil
}

// SavePageStay 保存页面停留时间数据到数据库
// 参数:
//   - ctx: 上下文对象，用于控制请求超时和取消
//   - pageStay: 页面停留对象，包含要保存的页面停留时间数据
// 返回:
//   - error: 保存过程中的错误信息，成功则为nil
func (r *ClickHouseRepository) SavePageStay(ctx context.Context, pageStay *models.PageStay) error {
	// 定义SQL插入语句，包含页面停留数据的所有字段
	query := `INSERT INTO page_stay (timestamp, project_id, session_id, trace_id, user_id, url, referrer, type, name, value, extra) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	
	// 执行插入操作
	_, err := r.DB.ExecContext(ctx, query, 
		pageStay.Timestamp, pageStay.ProjectID, pageStay.SessionID, pageStay.TraceID, pageStay.UserID,
		pageStay.URL, pageStay.Referrer, pageStay.Type, pageStay.Name, pageStay.Value, pageStay.Extra)
	if err != nil {
		return fmt.Errorf("failed to save page stay: %w", err)
	}
	return nil
}

// GetPageStays 获取指定项目在时间范围内的页面停留时间列表
// 参数:
//   - ctx: 上下文对象，用于控制请求超时和取消
//   - projectID: 项目标识符
//   - startTime: 开始时间
//   - endTime: 结束时间
// 返回:
//   - []*models.PageStay: 页面停留时间列表
//   - error: 查询过程中的错误信息，成功则为nil
func (r *ClickHouseRepository) GetPageStays(ctx context.Context, projectID string, startTime, endTime time.Time) ([]*models.PageStay, error) {
	// 定义SQL查询语句，按时间倒序排列
	query := `SELECT timestamp, project_id, session_id, trace_id, user_id, url, referrer, type, name, value, extra 
		FROM page_stay 
		WHERE project_id = ? AND timestamp >= ? AND timestamp <= ? 
		ORDER BY timestamp DESC`
	
	// 执行查询
	rows, err := r.DB.QueryContext(ctx, query, projectID, startTime, endTime)
	if err != nil {
		return nil, fmt.Errorf("failed to query page stays: %w", err)
	}
	defer rows.Close()

	var stays []*models.PageStay
	// 遍历查询结果
	for rows.Next() {
		var stay models.PageStay
		err := rows.Scan(
			&stay.Timestamp, &stay.ProjectID, &stay.SessionID, &stay.TraceID, &stay.UserID,
			&stay.URL, &stay.Referrer, &stay.Type, &stay.Name, &stay.Value, &stay.Extra)
		if err != nil {
			return nil, fmt.Errorf("failed to scan page stay: %w", err)
		}
		stays = append(stays, &stay)
	}
	return stays, nil
}

// GetAveragePageStay 获取指定项目在时间范围内的平均页面停留时间
// 参数:
//   - ctx: 上下文对象，用于控制请求超时和取消
//   - projectID: 项目标识符
//   - startTime: 开始时间
//   - endTime: 结束时间
// 返回:
//   - float64: 平均页面停留时间（秒）
//   - error: 查询过程中的错误信息，成功则为nil
func (r *ClickHouseRepository) GetAveragePageStay(ctx context.Context, projectID string, startTime, endTime time.Time) (float64, error) {
	// 使用ClickHouse的avg函数计算平均值
	query := `SELECT avg(value) FROM page_stay WHERE project_id = ? AND timestamp >= ? AND timestamp <= ?`
	var avg float64
	// 执行聚合查询
	err := r.DB.QueryRowContext(ctx, query, projectID, startTime, endTime).Scan(&avg)
	if err != nil {
		if err == sql.ErrNoRows {
			// 如果没有数据，返回0
			return 0, nil
		}
		return 0, fmt.Errorf("failed to query average page stay: %w", err)
	}
	return avg, nil
}

// Close 关闭数据库连接
// 释放所有资源，包括连接池中的连接
// 返回:
//   - error: 关闭过程中的错误信息，成功则为nil
func (r *ClickHouseRepository) Close() error {
	return r.DB.Close()
}
