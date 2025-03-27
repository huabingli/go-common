/**
  @author: 35840
  @date: 2024/4/1
  @desc: 设置X-Request-ID中间件
**/

package middleware

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/huabingli/go-common"
)

// NewRequestIDMiddleware 创建带自定义 header key 的 RequestID 中间件
func NewRequestIDMiddleware(headerKeys ...string) gin.HandlerFunc {

	headerKey := "X-Request-ID"
	if len(headerKeys) > 0 && headerKeys[0] != "" {
		headerKey = headerKeys[0]
	}

	return func(c *gin.Context) {
		// 启动计时器
		common.GetStartTime(c, StartTimeKey)

		// 从 Header 获取 request ID
		requestID := c.GetHeader(headerKey)
		if requestID == "" {
			requestID = common.GenerateRequestID()
			c.Request.Header.Set(headerKey, requestID)
		}

		c.Header(headerKey, requestID)

		// 写入标准 context.Context
		ctx := context.WithValue(c.Request.Context(), headerKey, requestID)
		c.Request = c.Request.WithContext(ctx)

		// 如果你有需要，也可以设置到 gin.Context：
		c.Set(headerKey, requestID)

		c.Next()
	}
}
