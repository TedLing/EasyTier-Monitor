package Tools

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
)

// 日志级别
const (
	LogLevelDebug = "debug"
	LogLevelInfo  = "info"
	LogLevelWarn  = "warn"
	LogLevelError = "error"
)

// 日志记录器接口
type Logger interface {
	Debug(format string, args ...interface{})
	Info(format string, args ...interface{})
	Warn(format string, args ...interface{})
	Error(format string, args ...interface{})
}

// 简单日志记录器实现
type simpleLogger struct {
	level  string
	output *os.File
}

// 创建新的日志记录器
func NewLogger() Logger {
	logConf := AppConfig.Log

	// 打开日志文件（如果配置了）
	var output *os.File
	if logConf.File != "" {
		// 创建日志目录
		dir := filepath.Dir(logConf.File)
		if err := os.MkdirAll(dir, 0755); err != nil {
			fmt.Printf("警告: 无法创建日志目录 %s: %v\n", dir, err)
			output = os.Stdout
		} else {
			// 打开日志文件
			f, err := os.OpenFile(logConf.File, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			if err != nil {
				fmt.Printf("警告: 无法打开日志文件 %s: %v\n", logConf.File, err)
				output = os.Stdout
			} else {
				output = f
			}
		}
	} else {
		output = os.Stdout
	}

	return &simpleLogger{
		level:  logConf.Level,
		output: output,
	}
}

// 全局日志实例
var AppLogger Logger

// 初始化日志
func InitLogger() {
	AppLogger = NewLogger()
	AppLogger.Info("日志系统初始化完成，级别: %s", AppConfig.Log.Level)
}

// 判断是否需要记录指定级别的日志
func (l *simpleLogger) shouldLog(level string) bool {
	levelOrder := map[string]int{
		LogLevelDebug: 0,
		LogLevelInfo:  1,
		LogLevelWarn:  2,
		LogLevelError: 3,
	}

	currentLevel, exists := levelOrder[l.level]
	if !exists {
		currentLevel = levelOrder[LogLevelInfo] // 默认info级别
	}

	logLevel, exists := levelOrder[level]
	if !exists {
		return false
	}

	return logLevel >= currentLevel
}

// 输出日志
func (l *simpleLogger) log(level string, format string, args ...interface{}) {
	if !l.shouldLog(level) {
		return
	}

	// 格式化日志消息
	msg := fmt.Sprintf(format, args...)
	// 格式化时间
	now := time.Now().Format("2006-01-02 15:04:05")
	// 构建完整日志行
	logLine := fmt.Sprintf("%s [%s] %s\n", now, level, msg)

	// 写入日志
	l.output.WriteString(logLine)
}

// Debug 记录调试日志
func (l *simpleLogger) Debug(format string, args ...interface{}) {
	l.log(LogLevelDebug, format, args...)
}

// Info 记录信息日志
func (l *simpleLogger) Info(format string, args ...interface{}) {
	l.log(LogLevelInfo, format, args...)
}

// Warn 记录警告日志
func (l *simpleLogger) Warn(format string, args ...interface{}) {
	l.log(LogLevelWarn, format, args...)
}

// Error 记录错误日志
func (l *simpleLogger) Error(format string, args ...interface{}) {
	l.log(LogLevelError, format, args...)
}

// Gin中间件：记录API请求日志
func GinLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 开始时间
		startTime := time.Now()

		// 处理请求
		c.Next()

		// 结束时间
		endTime := time.Now()
		// 执行时间
		latencyTime := endTime.Sub(startTime)
		// 请求方式
		reqMethod := c.Request.Method
		// 请求路由
		reqUri := c.Request.RequestURI
		// 状态码
		statusCode := c.Writer.Status()
		// 请求IP
		clientIP := c.ClientIP()

		// 根据状态码选择日志级别
		if statusCode >= 500 {
			AppLogger.Error("%s | %3d | %13v | %15s | %s", reqMethod, statusCode, latencyTime, clientIP, reqUri)
		} else if statusCode >= 400 {
			AppLogger.Warn("%s | %3d | %13v | %15s | %s", reqMethod, statusCode, latencyTime, clientIP, reqUri)
		} else {
			AppLogger.Info("%s | %3d | %13v | %15s | %s", reqMethod, statusCode, latencyTime, clientIP, reqUri)
		}
	}
}
