package logger

import (
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"runtime/debug"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/natefinch/lumberjack"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func Init() error {
	// 0. 文件
	// file, _ := os.Create("./test.log")
	// writeSyncer := zapcore.AddSync(file)

	// 1. writeSyncer循环文件
	lumberJackLogger := &lumberjack.Logger{
		Filename:   viper.GetString("log.filename"),
		LocalTime:  true,
		MaxSize:    viper.GetInt("log.maxsize"),    //M
		MaxBackups: viper.GetInt("log.max_backup"), //个
		MaxAge:     viper.GetInt("log.max_age"),    //天
		Compress:   false,
	}
	writeSyncer := zapcore.AddSync(lumberJackLogger)

	// 2. console 打印 调试
	// writeSyncer := zapcore.AddSync(os.Stdout)

	// 3. encoder JSON格式
	// encoder := zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())

	// 4. encoder Console格式
	// encoder := zapcore.NewConsoleEncoder(zap.NewProductionEncoderConfig())

	// 5. 自定义格式
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	encoder := zapcore.NewJSONEncoder(encoderConfig)

	// 6.core
	var core zapcore.Core
	if viper.GetString("app.mode") == "dev" {
		//进入到开发模式，日志输出到终端，
		consoleencoder := zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())
		core = zapcore.NewTee(
			zapcore.NewCore(encoder, writeSyncer, zapcore.InfoLevel),
			zapcore.NewCore(consoleencoder,zapcore.Lock(os.Stdout),zap.DebugLevel),
		)
	}else {
		core = zapcore.NewCore(encoder, writeSyncer, zapcore.DebugLevel)
	}
	// 7.log
	lg := zap.New(core)
	zap.ReplaceGlobals(lg)
	return nil
}

// GinLogger 接收gin框架默认的日志
func GinLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery
		c.Next()

		cost := time.Since(start)
		zap.L().Info(path,
			zap.Int("status", c.Writer.Status()),
			zap.String("method", c.Request.Method),
			zap.String("path", path),
			zap.String("query", query),
			zap.String("ip", c.ClientIP()),
			zap.String("user-agent", c.Request.UserAgent()),
			zap.String("errors", c.Errors.ByType(gin.ErrorTypePrivate).String()),
			zap.Duration("cost", cost),
		)
	}
}

// GinRecovery recover掉项目可能出现的panic
func GinRecovery(stack bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// Check for a broken connection, as it is not really a
				// condition that warrants a panic stack trace.
				var brokenPipe bool
				if ne, ok := err.(*net.OpError); ok {
					if se, ok := ne.Err.(*os.SyscallError); ok {
						if strings.Contains(strings.ToLower(se.Error()), "broken pipe") || strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
							brokenPipe = true
						}
					}
				}

				httpRequest, _ := httputil.DumpRequest(c.Request, false)
				if brokenPipe {
					zap.L().Error(c.Request.URL.Path,
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
					)
					// If the connection is dead, we can't write a status to it.
					c.Error(err.(error)) // nolint: errcheck
					c.Abort()
					return
				}

				if stack {
					zap.L().Error("[Recovery from panic]",
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
						zap.String("stack", string(debug.Stack())),
					)
				} else {
					zap.L().Error("[Recovery from panic]",
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
					)
				}
				c.AbortWithStatus(http.StatusInternalServerError)
			}
		}()
		c.Next()
	}
}
