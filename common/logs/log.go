package logs

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"runtime/debug"
	"strings"
	"time"
)

// LG : 全局日志实例
var LG *zap.Logger

// LogConfig 日志配置
type LogConfig struct {
	DebugFileName string `json:"debugFileName"` // Debug 级别日志文件路径
	InfoFileName  string `json:"infoFileName"`
	WarnFileName  string `json:"warnFileName"`
	MaxSize       int    `json:"maxSize"`    // 单个日志我呢见的最大大小（单位：MB）
	MaxAge        int    `json:"maxAge"`     // 日志文件保存的最长天数
	MaxBackups    int    `json:"maxBackups"` // 日志文件的最大备份数量
}

// InitLogger 初始化 Logger
// 传入日志配置，设置日志文件输出、日志级别和格式化
func InitLogger(cfg *LogConfig) (err error) {
	// 获取不同级别的日志写入器
	writeSyncDebug := getLogWriter(cfg.DebugFileName, cfg.MaxSize, cfg.MaxBackups, cfg.MaxAge)
	writeSyncInfo := getLogWriter(cfg.InfoFileName, cfg.MaxSize, cfg.MaxBackups, cfg.MaxAge)
	writeSyncWarn := getLogWriter(cfg.WarnFileName, cfg.MaxSize, cfg.MaxBackups, cfg.MaxAge)

	// 设置日志编码
	encoder := getEncoder()

	// 创建日志核心组件
	debugCore := zapcore.NewCore(encoder, writeSyncDebug, zapcore.DebugLevel)
	infoCore := zapcore.NewCore(encoder, writeSyncInfo, zapcore.InfoLevel)
	warnCore := zapcore.NewCore(encoder, writeSyncWarn, zapcore.WarnLevel)

	// 控制台输出（开发模式下使用）
	consoleEncoder := zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())
	std := zapcore.NewCore(consoleEncoder, zapcore.Lock(os.Stdout), zapcore.DebugLevel) // debug 级别之上都输出一下

	// 合并多个日志核心
	core := zapcore.NewTee(debugCore, infoCore, warnCore, std)

	// 创建 Logger 实例，带调试信息
	LG = zap.New(core, zap.AddCaller())

	// 替换 zap 的全局 Logger，方便在其他模块中直接使用 zap.L()
	zap.ReplaceGlobals(LG)
	return nil
}

// getEncoder 返回日志编码（JSON格式）
func getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder         // 时间格式：ISO8601
	encoderConfig.TimeKey = "time"                                // 时间字段名
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder       // 日志级别大写显示
	encoderConfig.EncodeDuration = zapcore.SecondsDurationEncoder // 耗时显示为秒
	encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder       // 调用信息：短路径
	return zapcore.NewJSONEncoder(encoderConfig)                  // 返回JSON格式的编码器
}

// getLogWriter 返回日志文件的写入器
// 使用 lumberjack 实现日志轮转
func getLogWriter(filename string, maxSize int, maxBackup int, maxAge int) zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   filename,  // 日志文件路径
		MaxSize:    maxSize,   // 单个日志文件的最大大小
		MaxBackups: maxBackup, // 最大备份数量
		MaxAge:     maxAge,    // 日志保留天数
	}
	return zapcore.AddSync(lumberJackLogger) // 返回支持同步的写入器
}

// GinLogger 接收 Gin 框架的默认日志并记录详细信息
func GinLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 记录请求开始时间
		start := time.Now()
		path := c.Request.URL.Path      // 请求路径
		query := c.Request.URL.RawQuery // 查询参数

		// 继续处理请求
		c.Next()

		// 计算请求耗时
		cost := time.Since(start)
		LG.Info(path,
			zap.Int("status", c.Writer.Status()),                                 // 响应状态码
			zap.String("method", c.Request.Method),                               // 请求方法
			zap.String("path", path),                                             // 请求路径
			zap.String("query", query),                                           // 查询参数
			zap.String("ip", c.ClientIP()),                                       // 客户端IP
			zap.String("user-agent", c.Request.UserAgent()),                      // 用户代理
			zap.String("errors", c.Errors.ByType(gin.ErrorTypePrivate).String()), // Gin中的私有错误
			zap.Duration("cost", cost),                                           // 请求耗时
		)
	}
}

// GinRecovery 捕获项目中的panic，并记录错误信息和堆栈
func GinRecovery(stack bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			// 捕获panic
			if err := recover(); err != nil {
				// 检测是否时网络错误（如：断开连接）
				var brokenPipe bool
				if ne, ok := err.(*net.OpError); ok {
					if se, ok := ne.Err.(*os.SyscallError); ok {
						if strings.Contains(strings.ToLower(se.Error()), "broken pipe") ||
							strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
							brokenPipe = true
						}
					}
				}

				// 捕获请求内容
				httpRequest, _ := httputil.DumpRequest(c.Request, false)

				if brokenPipe {
					// 如果连接终端，记录简要日志后返回
					LG.Error(c.Request.URL.Path,
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
					)
					c.Error(err.(error)) // nolint: errcheck
					c.Abort()
					return
				}

				// 记录详细日志，包括堆栈信息
				if stack {
					LG.Error("[Recovery from panic]",
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
						zap.String("stack", string(debug.Stack())),
					)
				} else {
					LG.Error("[Recovery from panic]",
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
					)
				}

				// 返回500错误
				c.AbortWithStatus(http.StatusInternalServerError)
			}

		}()

		c.Next()
	}
}
