package common

import (
	"strconv"
	"strings"
	"time"
)

// ParseDuration
// @描述: 解析时间
func ParseDuration(d string) (time.Duration, error) {
	// 去除字符串两端的空白字符
	d = strings.TrimSpace(d)

	// 尝试使用 time.ParseDuration 解析时间段
	dr, err := time.ParseDuration(d)
	if err == nil {
		// 如果成功解析，直接返回解析结果和 nil 错误
		return dr, nil
	}

	// 如果字符串包含 "d"，表示可能是自定义的天数格式
	if strings.Contains(d, "d") {
		// 找到第一个 "d" 的位置
		index := strings.Index(d, "d")

		// 解析 "d" 前面的部分作为小时数
		hour, _ := strconv.Atoi(d[:index])

		// 计算总的小时数对应的时间段（天数）
		dr = time.Hour * 24 * time.Duration(hour)

		// 继续解析 "d" 后面的部分作为剩余的时间段
		ndr, err := time.ParseDuration(d[index+1:])
		if err != nil {
			// 如果剩余部分无法解析为时间段，返回之前解析的天数时间段和 nil 错误
			return dr, nil
		}

		// 将天数时间段和剩余时间段相加得到最终的时间段
		return dr + ndr, nil
	}

	// 如果以上方式都无法解析，尝试将字符串解析为整数
	dv, err := strconv.ParseInt(d, 10, 64)
	if err != nil {
		// 解析整数失败，直接返回 0 时间段和解析错误
		return 0, err
	}

	// 将整数转换为时间段并返回
	return time.Duration(dv), nil
}
