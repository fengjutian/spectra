package config

import (
	"github.com/spf13/viper"
)

// Config 应用程序配置结构
type Config struct {
	App    AppConfig    `mapstructure:"app"`
	Server ServerConfig `mapstructure:"server"`
	Log    LogConfig    `mapstructure:"log"`
	DB     DBConfig     `mapstructure:"db"`
}

// AppConfig 应用基本配置
type AppConfig struct {
	Name        string `mapstructure:"name"`
	Version     string `mapstructure:"version"`
	Environment string `mapstructure:"environment"` // development, production, test
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Port         int    `mapstructure:"port"`
	Host         string `mapstructure:"host"`
	ReadTimeout  int    `mapstructure:"read_timeout"`
	WriteTimeout int    `mapstructure:"write_timeout"`
}

// LogConfig 日志配置
type LogConfig struct {
	Level    string `mapstructure:"level"` // debug, info, warn, error
	Path     string `mapstructure:"path"`
	MaxSize  int    `mapstructure:"max_size"` // MB
	MaxAge   int    `mapstructure:"max_age"`  // days
	Compress bool   `mapstructure:"compress"`
}

// LoadConfig 加载配置文件
func LoadConfig() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config/")
	viper.AddConfigPath("./")
	viper.AutomaticEnv()

	// 设置默认值
	setDefaultConfig()

	// 读取配置文件
	if err := viper.ReadInConfig(); err != nil {
		// 如果找不到配置文件，使用默认值
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// 配置文件不存在，继续使用默认值
		} else {
			// 配置文件存在但有错误
			return nil, err
		}
	}

	// 解析配置到结构体
	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	return &config, nil
}

// DBConfig 数据库配置
type DBConfig struct {
	Driver   string `mapstructure:"driver"`
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Database string `mapstructure:"database"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	Debug    bool   `mapstructure:"debug"`
}

// setDefaultConfig 设置默认配置
func setDefaultConfig() {
	// App 默认配置
	viper.SetDefault("app.name", "spectra-backend")
	viper.SetDefault("app.version", "1.0.0")
	viper.SetDefault("app.environment", "development")

	// Server 默认配置
	viper.SetDefault("server.port", 8080)
	viper.SetDefault("server.host", "0.0.0.0")
	viper.SetDefault("server.read_timeout", 15)
	viper.SetDefault("server.write_timeout", 15)

	// Log 默认配置
	viper.SetDefault("log.level", "info")
	viper.SetDefault("log.path", "./logs/app.log")
	viper.SetDefault("log.max_size", 500)
	viper.SetDefault("log.max_age", 30)
	viper.SetDefault("log.compress", true)

	// DB 默认配置
	viper.SetDefault("db.driver", "clickhouse")
	viper.SetDefault("db.host", "https://ci5eaxwoe9.asia-southeast1.gcp.clickhouse.cloud:8443")
	viper.SetDefault("db.port", 9000)
	viper.SetDefault("db.database", "spectra_log_data")
	viper.SetDefault("db.username", "default")
	viper.SetDefault("db.password", "QhH_vObgVEGw6")
	viper.SetDefault("db.debug", false)
}
