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
	"github.com/huabingli/go-common/jsonutil"
)

// SkipLogFunc 定义类型：用于判断是否跳过日志记录的函数
type SkipLogFunc func(c *gin.Context) bool

func GSlog(skipFns ...SkipLogFunc) gin.HandlerFunc {

	skipFn := func(c *gin.Context) bool { return false }
	if len(skipFns) > 0 && skipFns[0] != nil {
		skipFn = skipFns[0]
	}

	return func(c *gin.Context) {
		// 开始计时，记录请求开始时间
		start := common.GetStartTime(c)

		if skipFn(c) {
			c.Next()
			return
		}

		path := c.Request.URL.Path // 获取请求的 URL 路径
		method := c.Request.Method // 获取请求的方法（GET、POST 等）

		raw := c.Request.URL.RawQuery // 获取请求的原始查询参数

		fullPath := constructPath(path, raw)

		c.Next() // 执行下一个中间件或最终的处理器函数

		// 计算请求处理耗时
		duration := time.Since(start)

		status := c.Writer.Status()

		clientIp := c.ClientIP()

		attrs := buildRequestLogAttrs(c, status, method, path, clientIp, fullPath, duration)
		// 构建请求摘要
		summary := fmt.Sprintf("%3d %v %s %s %s", status, duration, clientIp, method, fullPath)
		// 将请求信息记录到日志
		slog.LogAttrs(
			c.Request.Context(),
			levelByStatus(status), // 设置日志级别为 Debug
			"HTTP request",
			slog.String("summary", summary),
			slog.Any("attrs", attrs),
		)
	}
}

func buildRequestLogAttrs(
	c *gin.Context,
	status int,
	method, path, clientIp, fullPath string,
	duration time.Duration,
) []slog.Attr {
	query := c.Request.URL.Query()
	attrs := []slog.Attr{
		slog.Int("status", status),
		slog.String("duration", duration.String()),
		slog.String("ip", clientIp),
		slog.String("method", method),
		slog.String("path", path),
		slog.String("query", jsonutil.MustMarshalToString(&query)),
		slog.String("fullPath", fullPath),
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

func levelByStatus(status int) slog.Level {
	switch {
	case status >= 500:
		return slog.LevelError
	case status >= 400:
		return slog.LevelWarn
	default:
		return slog.LevelInfo
	}
}
