/**
  @author: 35840
  @date: 2024/5/15
  @desc: 捕捉panic错误
**/

package middleware

import (
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

type ErrorHandlerFunc func(c *gin.Context, err any)

func Recovery(handle ErrorHandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				var brokenPipe bool
				var err error

				if e, ok := r.(error); ok {
					err = e
					if ne, ok := err.(*net.OpError); ok {
						var se *os.SyscallError
						if errors.As(ne, &se) {
							seStr := strings.ToLower(se.Error())
							if strings.Contains(seStr, "broken pipe") || strings.Contains(
								seStr, "connection reset by peer",
							) {
								brokenPipe = true
							}
						}
					}
				} else {
					err = errors.New(fmt.Sprint(r))
				}

				headersToStr := sanitizeRequestHeaders(c.Request)
				ctx := c.Request.Context()

				slog.ErrorContext(
					ctx, "服务器内部错误！",
					slog.Any("err", err),
					slog.String("method", c.Request.Method),
					slog.String("url", c.Request.URL.String()),
					slog.String("headers", headersToStr),
				)

				if brokenPipe {
					c.Error(err)
					c.Abort()
					return
				}

				// 调用传入的错误处理函数
				handle(c, r)
			}
		}()
		c.Next()
	}
}

func sanitizeRequestHeaders(r *http.Request) string {
	httpRequest, _ := httputil.DumpRequest(r, false)
	headers := strings.Split(string(httpRequest), "\r\n")
	for idx, header := range headers {
		current := strings.SplitN(header, ":", 2)
		if len(current) == 2 && strings.EqualFold(strings.TrimSpace(current[0]), "Authorization") {
			headers[idx] = current[0] + ": *"
		}
	}
	return strings.Join(headers, "\r\n")
}
