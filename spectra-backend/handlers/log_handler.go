package handlers

import (
	"encoding/json"
	"net/http"
	"spectra-backend/models"
	"spectra-backend/services"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// LogHandler 日志处理器
type LogHandler struct {
	logService services.LogService
	logger     *zap.Logger
}

// NewLogHandler 创建日志处理器实例
func NewLogHandler(logService services.LogService, logger *zap.Logger) *LogHandler {
	return &LogHandler{
		logService: logService,
		logger:     logger,
	}
}

// RecordErrorLog 记录错误日志
func (h *LogHandler) RecordErrorLog(c *gin.Context) {
	h.logger.Debug("RecordErrorLog called", zap.String("method", c.Request.Method), zap.String("content_type", c.GetHeader("Content-Type")))

	var log models.ErrorLog
	h.logger.Debug("Binding JSON request body")
	if err := c.ShouldBindJSON(&log); err != nil {
		h.logger.Error("Failed to bind error log", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	h.logger.Debug("Successfully bound request body",
		zap.String("project_id", log.ProjectID),
		zap.String("session_id", log.SessionID),
		zap.String("trace_id", log.TraceID),
		zap.String("user_id", log.UserID),
		zap.String("url", log.URL),
		zap.String("type", log.Type),
		zap.String("name", log.Name),
		zap.String("message", log.Message),
		zap.String("extra", string(log.Extra)))

	if err := h.logService.RecordErrorLog(c.Request.Context(), &log); err != nil {
		h.logger.Error("Failed to record error log 111", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to record error log"})
		return
	}

	h.logger.Debug("Error log recorded successfully")
	c.JSON(http.StatusCreated, gin.H{"message": "Error log recorded successfully"})
}

// GetErrorLogs 获取错误日志列表
func (h *LogHandler) GetErrorLogs(c *gin.Context) {
	projectID := c.Query("project_id")
	if projectID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "project_id is required"})
		return
	}

	// 调试：打印查询的项目ID
	h.logger.Debug("GetErrorLogs called", zap.String("project_id", projectID))

	startTime, endTime, err := parseTimeRange(c)
	if err != nil {
		h.logger.Error("Invalid time range for GetErrorLogs",
			zap.String("project_id", projectID),
			zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	logs, err := h.logService.GetErrorLogs(c.Request.Context(), projectID, startTime, endTime)
	if err != nil {
		h.logger.Error("Failed to get error logs",
			zap.String("project_id", projectID),
			zap.Time("start_time", startTime),
			zap.Time("end_time", endTime),
			zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get error logs"})
		return
	}

	// 调试：输出查询结果条数
	h.logger.Debug("GetErrorLogs succeeded",
		zap.String("project_id", projectID),
		zap.Time("start_time", startTime),
		zap.Time("end_time", endTime),
		zap.Int("count", len(logs)))

	c.JSON(http.StatusOK, logs)
}

// RecordPerformanceMetric 记录性能指标
func (h *LogHandler) RecordPerformanceMetric(c *gin.Context) {
	var metric models.PerformanceMetric
	if err := c.ShouldBindJSON(&metric); err != nil {
		h.logger.Error("Failed to bind performance metric", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if err := h.logService.RecordPerformanceMetric(c.Request.Context(), &metric); err != nil {
		h.logger.Error("Failed to record performance metric", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to record performance metric"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Performance metric recorded successfully"})
}

// GetPerformanceMetrics 获取性能指标列表
func (h *LogHandler) GetPerformanceMetrics(c *gin.Context) {
    projectID := c.Query("project_id")
    if projectID == "" {
        h.logger.Error("Missing project_id for performance metrics request")
        c.JSON(http.StatusBadRequest, gin.H{"error": "project_id is required"})
        return
    }

    // Debug: raw inputs
    h.logger.Debug(
        "GetPerformanceMetrics request",
        zap.String("project_id", projectID),
        zap.String("start_time_raw", c.Query("start_time")),
        zap.String("end_time_raw", c.Query("end_time")),
    )

    startTime, endTime, err := parseTimeRange(c)
    if err != nil {
        h.logger.Error(
            "Failed to parse time range for performance metrics",
            zap.String("project_id", projectID),
            zap.String("start_time_raw", c.Query("start_time")),
            zap.String("end_time_raw", c.Query("end_time")),
            zap.Error(err),
        )
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    h.logger.Debug(
        "Parsed time range for performance metrics",
        zap.String("project_id", projectID),
        zap.Time("start_time", startTime),
        zap.Time("end_time", endTime),
    )

    metrics, err := h.logService.GetPerformanceMetrics(c.Request.Context(), projectID, startTime, endTime)
    if err != nil {
        h.logger.Error(
            "Failed to get performance metrics",
            zap.String("project_id", projectID),
            zap.Time("start_time", startTime),
            zap.Time("end_time", endTime),
            zap.Error(err),
        )
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get performance metrics"})
        return
    }

    h.logger.Debug(
        "Fetched performance metrics",
        zap.String("project_id", projectID),
        zap.Int("count", len(metrics)),
        zap.Time("start_time", startTime),
        zap.Time("end_time", endTime),
    )

    c.JSON(http.StatusOK, metrics)
}

// RecordUserAction 记录用户行为
func (h *LogHandler) RecordUserAction(c *gin.Context) {
	var action models.UserAction
	if err := c.ShouldBindJSON(&action); err != nil {
		h.logger.Error("Failed to bind user action", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if err := h.logService.RecordUserAction(c.Request.Context(), &action); err != nil {
		h.logger.Error("Failed to record user action", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to record user action"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User action recorded successfully"})
}

// GetUserActions 获取用户行为列表
func (h *LogHandler) GetUserActions(c *gin.Context) {
    projectID := c.Query("project_id")
    if projectID == "" {
        h.logger.Error("Missing project_id for user actions request")
        c.JSON(http.StatusBadRequest, gin.H{"error": "project_id is required"})
        return
    }

    // Debug: raw inputs
    h.logger.Debug(
        "GetUserActions request",
        zap.String("project_id", projectID),
        zap.String("start_time_raw", c.Query("start_time")),
        zap.String("end_time_raw", c.Query("end_time")),
    )

    startTime, endTime, err := parseTimeRange(c)
    if err != nil {
        h.logger.Error(
            "Failed to parse time range for user actions",
            zap.String("project_id", projectID),
            zap.String("start_time_raw", c.Query("start_time")),
            zap.String("end_time_raw", c.Query("end_time")),
            zap.Error(err),
        )
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    h.logger.Debug(
        "Parsed time range for user actions",
        zap.String("project_id", projectID),
        zap.Time("start_time", startTime),
        zap.Time("end_time", endTime),
    )

    actions, err := h.logService.GetUserActions(c.Request.Context(), projectID, startTime, endTime)
    if err != nil {
        h.logger.Error(
            "Failed to get user actions",
            zap.String("project_id", projectID),
            zap.Time("start_time", startTime),
            zap.Time("end_time", endTime),
            zap.Error(err),
        )
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user actions"})
        return
    }

    h.logger.Debug(
        "Fetched user actions",
        zap.String("project_id", projectID),
        zap.Int("count", len(actions)),
        zap.Time("start_time", startTime),
        zap.Time("end_time", endTime),
    )

    c.JSON(http.StatusOK, actions)
}

// RecordCustomEvent 记录自定义事件
func (h *LogHandler) RecordCustomEvent(c *gin.Context) {
	var event models.CustomEvent
	if err := c.ShouldBindJSON(&event); err != nil {
		h.logger.Error("Failed to bind custom event", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// 处理自定义extra字段
	if c.Query("parse_extra") == "true" {
		// 这里可以根据需要解析额外的查询参数到extra字段
		extra := make(map[string]interface{})
		for key, values := range c.Request.URL.Query() {
			if key != "parse_extra" && len(values) > 0 {
				extra[key] = values[0]
			}
		}
		if len(extra) > 0 {
			extraData, _ := json.Marshal(extra)
			event.Extra = extraData
		}
	}

	if err := h.logService.RecordCustomEvent(c.Request.Context(), &event); err != nil {
		h.logger.Error("Failed to record custom event", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to record custom event"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Custom event recorded successfully"})
}

// GetCustomEvents 获取自定义事件列表
func (h *LogHandler) GetCustomEvents(c *gin.Context) {
    projectID := c.Query("project_id")
    if projectID == "" {
        h.logger.Error("Missing project_id for custom events request")
        c.JSON(http.StatusBadRequest, gin.H{"error": "project_id is required"})
        return
    }

    // Debug: raw inputs
    h.logger.Debug(
        "GetCustomEvents request",
        zap.String("project_id", projectID),
        zap.String("start_time_raw", c.Query("start_time")),
        zap.String("end_time_raw", c.Query("end_time")),
    )

    startTime, endTime, err := parseTimeRange(c)
    if err != nil {
        h.logger.Error(
            "Failed to parse time range for custom events",
            zap.String("project_id", projectID),
            zap.String("start_time_raw", c.Query("start_time")),
            zap.String("end_time_raw", c.Query("end_time")),
            zap.Error(err),
        )
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    h.logger.Debug(
        "Parsed time range for custom events",
        zap.String("project_id", projectID),
        zap.Time("start_time", startTime),
        zap.Time("end_time", endTime),
    )

    events, err := h.logService.GetCustomEvents(c.Request.Context(), projectID, startTime, endTime)
    if err != nil {
        h.logger.Error(
            "Failed to get custom events",
            zap.String("project_id", projectID),
            zap.Time("start_time", startTime),
            zap.Time("end_time", endTime),
            zap.Error(err),
        )
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get custom events"})
        return
    }

    h.logger.Debug(
        "Fetched custom events",
        zap.String("project_id", projectID),
        zap.Int("count", len(events)),
        zap.Time("start_time", startTime),
        zap.Time("end_time", endTime),
    )

    c.JSON(http.StatusOK, events)
}

// RecordPageStay 记录页面停留时长
func (h *LogHandler) RecordPageStay(c *gin.Context) {
	var pageStay models.PageStay
	if err := c.ShouldBindJSON(&pageStay); err != nil {
		h.logger.Error("Failed to bind page stay", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if err := h.logService.RecordPageStay(c.Request.Context(), &pageStay); err != nil {
		h.logger.Error("Failed to record page stay", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to record page stay"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Page stay recorded successfully"})
}

// GetAveragePageStay 获取平均页面停留时长
func (h *LogHandler) GetAveragePageStay(c *gin.Context) {
	projectID := c.Query("project_id")
	if projectID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "project_id is required"})
		return
	}

	startTime, endTime, err := parseTimeRange(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	average, err := h.logService.GetAveragePageStay(c.Request.Context(), projectID, startTime, endTime)
	if err != nil {
		h.logger.Error("Failed to get average page stay", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get average page stay"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"average_page_stay": average})
}

// parseTimeRange 解析时间范围参数
func parseTimeRange(c *gin.Context) (time.Time, time.Time, error) {
	startTimeStr := c.DefaultQuery("start_time", time.Now().AddDate(0, 0, -1).Format(time.RFC3339))
	endTimeStr := c.DefaultQuery("end_time", time.Now().Format(time.RFC3339))

	startTime, err := time.Parse(time.RFC3339, startTimeStr)
	if err != nil {
		return time.Time{}, time.Time{}, err
	}

	endTime, err := time.Parse(time.RFC3339, endTimeStr)
	if err != nil {
		return time.Time{}, time.Time{}, err
	}

	return startTime, endTime, nil
}
