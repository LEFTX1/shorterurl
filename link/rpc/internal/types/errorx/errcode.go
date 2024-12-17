package errorx

import "fmt"

const (
	// 用户错误
	ErrUserNotFound             = "A000114" // 用户不存在
	ErrUserNameExists           = "A000111" // 用户名已存在
	ErrPasswordTooShort         = "A000121" // 密码长度不够
	ErrInvalidPhoneFormat       = "A000151" // 手机号格式无效
	ErrIdempotentTokenInvalid   = "A000200" // 幂等Token无效
	ErrTooManyRequests          = "A000300" // 系统繁忙，请稍后再试
	ErrInternalServer           = "B000001" // 系统执行出错
	ErrThirdPartyServiceFailure = "C000001" // 第三方服务错误
	ErrGroupLimit               = "A000112" // 已超出最大分组数限制
	ErrGroupNotFound            = "A000113" // 分组不存在
	ErrBloomFilterCheck         = "B000002" // 布隆过滤器检查失败
	ErrDistributedLock          = "B000003" // 分布式锁操作失败
	ErrDatabaseOperation        = "B000004" // 数据库操作失败
	ErrRedisOperation           = "B000005" // Redis操作失败
	ErrShortLinkEmpty           = "A000400" // 短链接为空
	ErrShortLinkInvalid         = "A000401" // 短链接无效
	ErrShortLinkExists          = "A000402" // 短链接已存在

	// 布隆过滤器错误码
	ErrBloomFilterInit     = "B000006" // 布隆过滤器初始化失败
	ErrBloomFilterAdd      = "B000007" // 布隆过滤器添加失败
	ErrBloomFilterReset    = "B000008" // 布隆过滤器重置失败
	ErrBloomFilterNotFound = "B000009" // 布隆过滤器不存在
	ErrRedisClientNil      = "B000010" // Redis客户端为空
)

// 错误消息映射
var messages = map[string]string{
	ErrUserNotFound:             "用户不存在",
	ErrUserNameExists:           "用户名已存在",
	ErrPasswordTooShort:         "密码长度不够",
	ErrInvalidPhoneFormat:       "手机号格式无效",
	ErrIdempotentTokenInvalid:   "幂等Token无效",
	ErrTooManyRequests:          "系统繁忙，请稍后再试",
	ErrInternalServer:           "系统执行出错",
	ErrThirdPartyServiceFailure: "第三方服务调用失败",
	ErrGroupLimit:               "已超出最大分组数限制",
	ErrGroupNotFound:            "分组不存在",
	ErrBloomFilterCheck:         "布隆过滤器检查失败",
	ErrDistributedLock:          "分布式锁操作失败",
	ErrDatabaseOperation:        "数据库操作失败",
	ErrRedisOperation:           "Redis操作失败",
	ErrShortLinkEmpty:           "短链接不能为空",
	ErrShortLinkInvalid:         "短链接无效",
	ErrShortLinkExists:          "短链接已存在",

	// 布隆过滤器错误消息
	ErrBloomFilterInit:     "布隆过滤器初始化失败",
	ErrBloomFilterAdd:      "布隆过滤器添加失败",
	ErrBloomFilterReset:    "布隆过滤器重置失败",
	ErrBloomFilterNotFound: "布隆过滤器不存在",
	ErrRedisClientNil:      "Redis客户端未初始化",
}

// 布隆过滤器相关错误消息
const (
	MsgShortLinkEmpty      = "短链接不能为空"
	MsgBloomFilterInit     = "布隆过滤器初始化失败"
	MsgBloomFilterAdd      = "添加到布隆过滤器失败"
	MsgBloomFilterCheck    = "检查布隆过滤器失败"
	MsgBloomFilterReset    = "重置布隆过滤器失败"
	MsgRedisNotFound       = "Redis连接不存在"
	MsgBloomFilterNotFound = "布隆过滤器不存在"
)

type CodeError struct {
	Code    string `json:"code"`
	Msg     string `json:"msg"`
	ErrCode string `json:"err_code"`
}

func (e *CodeError) Error() string {
	return e.Msg
}

// NewCodeError 创建新的错误，支持格式化消息
func NewCodeError(code string, errCode string, msgFormat string, args ...interface{}) *CodeError {
	msg := msgFormat
	if len(args) > 0 {
		msg = fmt.Sprintf(msgFormat, args...)
	}
	return &CodeError{
		Code:    code,
		Msg:     msg,
		ErrCode: errCode,
	}
}

// NewDefaultError 创建默认错误
func NewDefaultError(msg string) *CodeError {
	return NewCodeError(ErrInternalServer, "", msg)
}

// NewValidationError 创建验证错误
func NewValidationError(msg string) *CodeError {
	return NewCodeError(ErrInvalidPhoneFormat, "", msg)
}

// Message 获取错误码对应的消息
func Message(code string) string {
	if msg, ok := messages[code]; ok {
		return msg
	}
	return "未知错误"
}
