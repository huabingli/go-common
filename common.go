package common

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
	"unicode/utf8"

	"github.com/gofrs/uuid/v5"
)

// 用 sync.Once 确保项目根目录路径只被计算一次
var (
	projectPath string
	once        sync.Once
)

// GenerateRequestID 生成无横线的UUID作为请求ID
func GenerateRequestID() string {
	id := uuid.Must(uuid.NewV4())
	// 去除UUID中的横线并返回
	return strings.ReplaceAll(id.String(), "-", "")
}

// StrToUint 将字符串转换为uint类型
func StrToUint(s string) (uint, error) {
	val, err := strconv.Atoi(s)
	if err != nil {
		return 0, err // 转换失败，返回错误
	}
	// 直接返回转换后的值，根据系统架构转换为 uint
	return uint(val), nil
}

// EqualParentIDs 判断两个 *uint 类型的指针是否相等，考虑它们可能为空
func EqualParentIDs(a, b *uint) bool {
	// 如果两个指针都为nil，视为相等
	if a == nil && b == nil {
		return true
	}
	// 如果两个指针都不为nil，比较它们指向的值
	if a != nil && b != nil {
		return *a == *b
	}
	// 如果只有一个指针为nil，视为不相等
	return false
}

// UniqueUintSlice 高效去重方法
func UniqueUintSlice(nums []uint) []uint {
	// 排序切片
	sort.Slice(nums, func(i, j int) bool { return nums[i] < nums[j] })

	// 使用新切片保存去重后的结果
	var result []uint
	for i, v := range nums {
		if i == 0 || v != nums[i-1] {
			result = append(result, v)
		}
	}
	return result
}

// SafeGo 捕获并记录任何在goroutine中panic的情况
func SafeGo(fn func(), onError func(interface{})) {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				// 如果发生panic，调用onError处理
				if onError != nil {
					onError(r)
				}
			}
		}()
		// 执行传入的函数
		fn()
	}()
}

// CreateTempFileWithCleanup 创建一个临时文件并返回文件路径和清理函数
func CreateTempFileWithCleanup(ctx context.Context, suffix string) (*os.File, func(), error) {
	// 确保后缀合法性
	suffix = strings.ToLower(strings.TrimSpace(suffix))
	if !strings.HasPrefix(suffix, ".") {
		suffix = "." + suffix
	}
	if suffix == "." {
		suffix = ".txt"
	}

	// 创建临时文件
	tempFile, err := os.CreateTemp("", "devops-api*"+suffix)
	if err != nil {
		return nil, nil, err
	}

	// 定义清理函数
	cleanup := func() {
		tempFilePath := tempFile.Name()
		// 关闭并删除临时文件
		err := tempFile.Close()
		if err != nil {
			slog.ErrorContext(ctx, "关闭临时文件失败", slog.Any("err", err))
			return
		}
		err = os.Remove(tempFilePath)
		if err != nil {
			slog.ErrorContext(ctx, "删除临时文件失败", slog.Any("err", err))
		} else {
			slog.DebugContext(ctx, fmt.Sprintf("临时文件已删除 %s", tempFilePath))
		}
	}
	return tempFile, cleanup, nil
}

// CreateTempDirWithCleanup 创建一个临时目录并返回目录路径和清理函数
func CreateTempDirWithCleanup(ctx context.Context) (string, func(), error) {
	// 创建临时目录
	tempDir, err := os.MkdirTemp("", "devops-api-*")
	if err != nil {
		return "", nil, err
	}
	slog.DebugContext(ctx, fmt.Sprintf("创建临时目录 %s", tempDir))

	// 定义清理函数
	cleanup := func() {
		// 删除临时目录
		err := os.RemoveAll(tempDir)
		if err != nil {
			slog.ErrorContext(ctx, "删除临时目录失败", slog.Any("err", err))
		} else {
			slog.InfoContext(ctx, fmt.Sprintf("临时目录已删除 %s", tempDir))
		}
	}

	return tempDir, cleanup, nil
}

// TruncateString 截取字符串的前n个字符，支持Unicode字符
func TruncateString(s string, n int) string {
	if utf8.RuneCountInString(s) <= n {
		return s
	}
	// 将字符串转换为rune切片进行截取，支持Unicode字符
	runes := []rune(s)
	return string(runes[:n])
}

// ExecutionTime 计算函数执行时间
func ExecutionTime(ctx context.Context, msg string, args ...any) func(additionalArgs ...any) {
	start := time.Now()

	// 获取调用者函数名、文件名和行号信息
	funcName, file, line := getCallerInfo(2)

	// 获取项目根目录并规范化文件路径
	relativePath := filepath.Clean(file)
	if strings.HasPrefix(relativePath, GetProjectPath()) {
		relativePath = strings.TrimPrefix(relativePath, GetProjectPath())
	}
	relativePath = strings.TrimLeft(relativePath, `\/`)
	formattedSource := fmt.Sprintf("%s:%d", relativePath, line)

	// 创建日志组，包含调用者信息
	callerInfo := slog.Group(
		"caller",
		slog.String("function", funcName),
		slog.String("file", formattedSource),
		slog.Int("line", line),
	)

	args = append(args, callerInfo)

	// 记录开始执行日志
	slog.InfoContext(ctx, fmt.Sprintf("开始执行 %s", msg), args...)

	return func(additionalArgs ...any) {
		// 计算执行时间
		endTime := time.Since(start)

		// 合并日志参数并记录结束执行日志
		finalArgs := append(args, additionalArgs...)
		finalArgs = append(
			finalArgs, slog.Group(
				"executionTime",
				slog.Int64("millis", endTime.Milliseconds()), // 毫秒
				slog.Float64("seconds", endTime.Seconds()),
			),
		)

		// 记录结束执行日志
		slog.InfoContext(
			ctx,
			fmt.Sprintf("结束执行 %s 耗时：(%s)", msg, endTime),
			finalArgs...,
		)
	}
}

// getCallerInfo 获取调用者函数的信息，包括函数名、文件名和行号
func getCallerInfo(skip int) (funcName, file string, line int) {
	pc, file, line, ok := runtime.Caller(skip)
	if !ok {
		return "未知调用者", "未知文件", 0
	}
	fn := runtime.FuncForPC(pc)
	if fn == nil {
		return "未知函数", file, line
	}
	return fn.Name(), file, line
}

// GetProjectPath 获取项目的根目录路径，使用 sync.Once 确保只执行一次
func GetProjectPath() string {
	// 使用 sync.Once 确保路径初始化只执行一次
	once.Do(
		func() {
			var err error
			projectPath, err = os.Getwd()
			if err != nil {
				slog.Error("获取项目路径失败", slog.Any("err", err))
				projectPath = "" // 如果获取失败，设置为空字符串
			}
		},
	)
	return projectPath
}

// GenerateMD5Hash 生成给定字符串的MD5哈希值
func GenerateMD5Hash(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}
