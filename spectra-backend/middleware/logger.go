package middleware

import (
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

func InitLogger() *zap.Logger {
	// 确保 logs 目录存在
	if _, err := os.Stat("logs"); os.IsNotExist(err) {
		os.Mkdir("logs", 0755)
	}

	// 配置 lumberjack 进行日志轮转
	hook := &lumberjack.Logger{
		Filename:   "logs/app.log", // 日志文件路径
		MaxSize:    10,             // 每个日志文件最大 10MB
		MaxBackups: 7,              // 最多保留 7 个旧文件
		MaxAge:     30,             // 最多保留 30 天
		Compress:   true,           // 压缩归档
	}

	// 编码器配置（时间格式、日志级别）
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalColorLevelEncoder, // 彩色级别输出
		EncodeTime:     customTimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	consoleEncoder := zapcore.NewConsoleEncoder(encoderConfig)
	fileEncoder := zapcore.NewJSONEncoder(encoderConfig)

	// 同时输出到控制台和文件
	core := zapcore.NewTee(
		zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), zapcore.DebugLevel),
		zapcore.NewCore(fileEncoder, zapcore.AddSync(hook), zapcore.InfoLevel),
	)

	logger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
	return logger
}

// GinLogger 是 Gin 的日志中间件
func GinLogger(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()

		latency := time.Since(start)
		status := c.Writer.Status()

		logger.Info("HTTP Request",
			zap.Int("status", status),
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
			zap.String("ip", c.ClientIP()),
			zap.Duration("latency", latency),
		)
	}
}

func customTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 15:04:05"))
}
