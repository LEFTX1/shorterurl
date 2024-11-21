// internal/errorx/errcode.go
package errorx

/*
 * 错误码格式说明：
 *
 * 第一位：错误级别
 *   A: 客户端错误（用户端错误）
 *   B: 服务端错误（系统执行错误）
 *   C: 第三方服务错误
 *
 * 第2-3位：模块编号
 *   00: 系统通用
 *   01: 用户模块
 *   02: 短链接模块
 *   ...
 *
 * 第4-6位：具体错误编号
 */

// 一级宏观错误码（客户端错误）
const (
	ClientError = "A000001" // 用户端错误
)

// 二级宏观错误码（用户注册相关）
const (
	UserRegisterError             = "A000100" // 用户注册错误
	UserNameVerifyError           = "A000110" // 用户名校验失败
	UserNameExistError            = "A000111" // 用户名已存在
	UserNameSensitiveError        = "A000112" // 用户名包含敏感词
	UserNameSpecialCharacterError = "A000113" // 用户名包含特殊字符
	UserNotExistError             = "A000114" // 用户不存在
	PasswordVerifyError           = "A000120" // 密码校验失败
	PasswordShortError            = "A000121" // 密码长度不够
	PhoneVerifyError              = "A000151" // 手机格式校验失败

)

// 幂等性相关错误码
const (
	IdempotentTokenNullError   = "A000200" // 幂等Token为空
	IdempotentTokenDeleteError = "A000201" // 幂等Token已被使用或失效
)

// 限流相关错误码
const (
	FlowLimitError = "A000300" // 当前系统繁忙，请稍后再试
)

// 一级宏观错误码（系统错误）
const (
	ServiceError        = "B000001" // 系统执行出错
	ServiceTimeoutError = "B000100" // 系统执行超时
)

// 一级宏观错误码（第三方服务错误）
const (
	RemoteError = "C000001" // 调用第三方服务出错
)

// 错误码消息映射
var messageMap = map[string]string{
	// 客户端错误
	ClientError:                   "用户端错误",
	UserRegisterError:             "用户注册错误",
	UserNameVerifyError:           "用户名校验失败",
	UserNameExistError:            "用户名已存在",
	UserNameSensitiveError:        "用户名包含敏感词",
	UserNameSpecialCharacterError: "用户名包含特殊字符",
	PasswordVerifyError:           "密码校验失败",
	PasswordShortError:            "密码长度不够",
	PhoneVerifyError:              "手机格式校验失败",

	// 幂等性错误
	IdempotentTokenNullError:   "幂等Token为空",
	IdempotentTokenDeleteError: "幂等Token已被使用或失效",

	// 限流错误
	FlowLimitError: "当前系统繁忙，请稍后再试",

	// 系统错误
	ServiceError:        "系统执行出错",
	ServiceTimeoutError: "系统执行超时",

	// 第三方服务错误
	RemoteError: "调用第三方服务出错",
}
