// internal/errorx/base_error.go
package errorx

import "errors"

// IError 错误接口，定义了错误对象需要实现的方法
type IError interface {
	Error() string // 返回错误信息
}

// BaseError 错误基类，包含错误码、错误信息和原始错误
type BaseError struct {
	Code    string // 错误码
	Message string // 错误信息
	Cause   error  // 原始错误
}

// Error 返回错误信息
func (e *BaseError) Error() string {
	return e.Message
}

// UserError 用户端错误，继承自 BaseError
type UserError struct {
	*BaseError
}

// SystemError 系统错误，继承自 BaseError
type SystemError struct {
	*BaseError
}

// RemoteCallError 远程调用错误，继承自 BaseError
type RemoteCallError struct {
	*BaseError
}

// NewError 创建一个新的 BaseError 实例
func NewError(code string) error {
	msg := messageMap[code] // 根据错误码从 messageMap 中获取错误信息
	return &BaseError{
		Code:    code,
		Message: msg,
	}
}

// NewErrorWithMessage 创建一个带有自定义消息的 BaseError 实例
func NewErrorWithMessage(code string, message string) error {
	return &BaseError{
		Code:    code,
		Message: message,
	}
}

// NewErrorWithCause 创建一个带有原始错误的 BaseError 实例
func NewErrorWithCause(code string, cause error) error {
	msg := messageMap[code] // 根据错误码从 messageMap 中获取错误信息
	return &BaseError{
		Code:    code,
		Message: msg,
		Cause:   cause,
	}
}

// NewUserError 创建一个新的 UserError 实例
func NewUserError(code string) error {
	msg := messageMap[code] // 根据错误码从 messageMap 中获取错误信息
	return &UserError{
		BaseError: &BaseError{
			Code:    code,
			Message: msg,
		},
	}
}

// NewSystemError 创建一个新的 SystemError 实例
func NewSystemError(code string) error {
	msg := messageMap[code] // 根据错误码从 messageMap 中获取错误信息
	return &SystemError{
		BaseError: &BaseError{
			Code:    code,
			Message: msg,
		},
	}
}

// NewRemoteError 创建一个新的 RemoteCallError 实例
func NewRemoteError(code string) error {
	msg := messageMap[code] // 根据错误码从 messageMap 中获取错误信息
	return &RemoteCallError{
		BaseError: &BaseError{
			Code:    code,
			Message: msg,
		},
	}
}

// IsUserError 判断是否是用户错误
func IsUserError(err error) bool {
	if err == nil {
		return false
	}
	var userError *UserError
	ok := errors.As(err, &userError)
	return ok
}

// IsSystemError 判断是否是系统错误
func IsSystemError(err error) bool {
	if err == nil {
		return false
	}
	var systemError *SystemError
	ok := errors.As(err, &systemError)
	return ok
}
