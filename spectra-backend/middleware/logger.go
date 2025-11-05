package middleware

import (
    "os"
    "spectra-backend/config"
    "time"

    "github.com/gin-gonic/gin"
    "go.uber.org/zap"
    "go.uber.org/zap/zapcore"
    "gopkg.in/natefinch/lumberjack.v2"
)

func InitLogger() *zap.Logger {

	cfg, err := config.LoadConfig()
	if err != nil {
		// 如果加载配置失败，使用默认值
		cfg = &config.Config{
			Log: config.LogConfig{
				Level:    "info",
				Path:     "./logs/app.log",
				MaxSize:  500,
				MaxAge:   30,
				Compress: true,
			},
		}
	}

	// 设置日志级别
	level := zap.InfoLevel
	switch cfg.Log.Level {
	case "debug":
		level = zap.DebugLevel
	case "info":
		level = zap.InfoLevel
	case "warn":
		level = zap.WarnLevel
	case "error":
		level = zap.ErrorLevel
	default:
		level = zap.InfoLevel
	}

    // 文件输出（JSON编码）
    fileWS := getLogWriter(cfg.Log.Path, cfg.Log.MaxSize, cfg.Log.MaxAge, 10)
    fileEncoder := getJSONEncoder()
    fileCore := zapcore.NewCore(fileEncoder, fileWS, level)

    // 控制台输出（开发环境启用，Console编码更易读）
    var core zapcore.Core
    if cfg.App.Environment == "development" {
        consoleWS := zapcore.AddSync(os.Stdout)
        consoleEncoder := getConsoleEncoder()
        consoleCore := zapcore.NewCore(consoleEncoder, consoleWS, level)
        core = zapcore.NewTee(fileCore, consoleCore)
    } else {
        core = fileCore
    }

    logger := zap.New(core, zap.AddCaller())

	return logger
}

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

func getLogWriter(filePath string, maxSize, maxAge, maxBackups int) zapcore.WriteSyncer {
    lumberJackLogger := &lumberjack.Logger{
        Filename:   filePath,
        MaxSize:    maxSize,
        MaxAge:     maxAge,
        MaxBackups: maxBackups,
        Compress:   true,
    }
    return zapcore.AddSync(lumberJackLogger)
}

// 文件 JSON 编码器
func getJSONEncoder() zapcore.Encoder {
    encoderConfig := zap.NewProductionEncoderConfig()
    encoderConfig.EncodeTime = customTimeEncoder
    encoderConfig.TimeKey = "time"
    encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
    encoderConfig.EncodeDuration = zapcore.SecondsDurationEncoder
    encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
    return zapcore.NewJSONEncoder(encoderConfig)
}

// 控制台编码器（更易读）
func getConsoleEncoder() zapcore.Encoder {
    encoderConfig := zap.NewDevelopmentEncoderConfig()
    encoderConfig.EncodeTime = customTimeEncoder
    encoderConfig.TimeKey = "time"
    encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
    encoderConfig.EncodeDuration = zapcore.SecondsDurationEncoder
    encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
    return zapcore.NewConsoleEncoder(encoderConfig)
}

func customTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 15:04:05"))
}
