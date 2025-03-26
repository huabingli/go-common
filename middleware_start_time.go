/**
  @author: 35840
  @date: 2024/5/15
  @desc: 获取gin c 中的启动时间
**/

package common

import (
	"time"

	"github.com/gin-gonic/gin"
)

func GetStartTime(c *gin.Context, key string) time.Time {
	if startTime, ok := c.Get(key); ok {
		if startTime, ok := startTime.(time.Time); ok {
			return startTime
		}
	}
	startTime := time.Now()
	c.Set(key, startTime)
	return startTime
}
