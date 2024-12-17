package errorx

import (
	"errors"
	"fmt"
)

// ErrorType 定义错误类型
type ErrorType string

const (
	ClientError ErrorType = "ClientError" // 客户端错误
	SystemError ErrorType = "SystemError" // 系统错误
	RemoteError ErrorType = "RemoteError" // 第三方服务错误
)

// AppError 应用错误结构体
type AppError struct {
	Type    ErrorType              // 错误类型
	Code    string                 // 错误码
	Message string                 // 错误信息
	Context map[string]interface{} // 错误上下文（可选）
}

// Error 实现 Go 的内置 error 接口
func (e *AppError) Error() string {
	if len(e.Context) > 0 {
		return fmt.Sprintf("[%s-%s] %s - Context: %v", e.Type, e.Code, e.Message, e.Context)
	}
	return fmt.Sprintf("[%s-%s] %s", e.Type, e.Code, e.Message)
}

// New 创建一个新的 AppError 实例
func New(errType ErrorType, code, message string) *AppError {
	return &AppError{
		Type:    errType,
		Code:    code,
		Message: message,
		Context: nil,
	}
}

// NewWithContext 创建一个带上下文的 AppError 实例
func NewWithContext(errType ErrorType, code, message string, context map[string]interface{}) *AppError {
	return &AppError{
		Type:    errType,
		Code:    code,
		Message: message,
		Context: context,
	}
}

// Is 判断是否是指定类型的 AppError
func Is(err error, errType ErrorType) bool {
	var appErr *AppError
	ok := errors.As(err, &appErr)
	return ok && appErr.Type == errType
}
