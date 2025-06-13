package jsonutil

import (
	jsoniter "github.com/json-iterator/go"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

// Marshal 将结构体编码为 JSON 字节
func Marshal[T any](v *T) ([]byte, error) {
	return json.Marshal(v)
}

// Unmarshal 将 JSON 字节解码为结构体
func Unmarshal[T any](data []byte, v *T) error {
	return json.Unmarshal(data, v)
}

// MarshalToString 编码为字符串（常用于日志）
func MarshalToString[T any](v *T) (string, error) {
	return json.MarshalToString(v)
}

// UnmarshalFromString 解码 JSON 字符串
func UnmarshalFromString[T any](str string, v *T) error {
	return json.UnmarshalFromString(str, v)
}

// MustMarshalToString 忽略错误返回字符串（仅建议日志使用）
func MustMarshalToString[T any](v *T) string {
	s, _ := json.MarshalToString(v)
	return s
}

// Clone 结构体深拷贝（通过 JSON 实现）
func Clone[T any](in T) (T, error) {
	var out T
	bytes, err := json.Marshal(in)
	if err != nil {
		return out, err
	}
	err = json.Unmarshal(bytes, &out)
	return out, err
}

// CompatibleWithStd 判断是否兼容标准库的行为（供调试）
func CompatibleWithStd() bool {
	return json == jsoniter.ConfigCompatibleWithStandardLibrary
}
