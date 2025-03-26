package common

import (
	"strings"

	"github.com/gofrs/uuid/v5"
)

// GenerateRequestID
// @描述: 获取无横线的UUID作为请求ID
// @return string
func GenerateRequestID() string {
	// 生成新的UUID
	id := uuid.Must(uuid.NewV4())

	// 去除横线
	requestID := strings.ReplaceAll(id.String(), "-", "")

	return requestID
}
