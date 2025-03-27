/**
  @author: 35840
  @date: 2024/4/10
  @desc: gin 日志中间件
**/

package middleware

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/huabingli/go-common"
)

const StartTimeKey = "startTime"

// SkipLogFunc 定义类型：用于判断是否跳过日志记录的函数
type SkipLogFunc func(c *gin.Context) bool

func GSlog(skipFns ...SkipLogFunc) gin.HandlerFunc {

	skipFn := func(c *gin.Context) bool { return false }
	if len(skipFns) > 0 && skipFns[0] != nil {
		skipFn = skipFns[0]
	}

	return func(c *gin.Context) {
		// 开始计时，记录请求开始时间
		start := common.GetStartTime(c, StartTimeKey)
		path := c.Request.URL.Path // 获取请求的 URL 路径
		method := c.Request.Method // 获取请求的方法（GET、POST 等）

		if skipFn(c) {
			c.Next()
			return
		}

		raw := c.Request.URL.RawQuery // 获取请求的原始查询参数

		c.Next() // 执行下一个中间件或最终的处理器函数

		// 计算请求处理耗时
		duration := time.Since(start)

		status := c.Writer.Status()

		attrs := buildRequestLogAttrs(c, status, method, path, raw, duration)

		// 自动设置日志等级
		level := slog.LevelInfo
		if status >= 500 {
			level = slog.LevelError
		} else if status >= 400 {
			level = slog.LevelWarn
		}
		//
		summary := fmt.Sprintf("%3d %v %s %s %s", status, duration, c.ClientIP(), method, constructPath(path, raw))
		// 将请求信息记录到日志
		slog.LogAttrs(
			c.Request.Context(),
			level, // 设置日志级别为 Debug
			"HTTP request",
			slog.String("summary", summary),
			slog.Any("attrs", attrs),
		)
	}
}

func buildRequestLogAttrs(
	c *gin.Context,
	status int,
	method, path, raw string,
	duration time.Duration,
) []slog.Attr {
	attrs := []slog.Attr{
		slog.Int("status", status),
		slog.String("duration", duration.String()),
		slog.String("ip", c.ClientIP()),
		slog.String("method", method),
		slog.String("path", path),
		slog.String("query", raw),
		slog.String("fullPath", constructPath(path, raw)),
		slog.Group(
			"requestDuration",
			slog.Int64("millis", duration.Milliseconds()),
			slog.Float64("seconds", duration.Seconds()),
		),
	}

	if errMsg := c.Errors.ByType(gin.ErrorTypePrivate).String(); errMsg != "" {
		attrs = append(attrs, slog.String("errorMessage", errMsg))
	}

	if endpoint := c.FullPath(); endpoint != "" {
		attrs = append(attrs, slog.String("endpoint", endpoint))
	}

	return attrs
}

// constructPath 函数组合路径和查询参数
// 如果查询参数不为空，则返回 "path?raw"，否则仅返回 path。
func constructPath(path, raw string) string {
	if raw != "" {
		return path + "?" + raw
	}
	return path
}
