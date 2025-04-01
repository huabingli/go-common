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

const StartTimeKey = "startTime"

func GetStartTime(c *gin.Context, keys ...string) time.Time {
	key := StartTimeKey
	if len(keys) > 0 {
		key = keys[0]
	}

	if val, ok := c.Get(key); ok {
		if startTime, ok := val.(time.Time); ok {
			return startTime
		}
	}
	startTime := time.Now()
	c.Set(key, startTime)
	return startTime
}
