package errorx

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
	ErrInvalidUsername          = "A000152" // 用户名无效
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
	ErrInvalidUsername:          "用户名只能包含ASCII字符，不能使用中文",
}

// Message 获取错误码对应的消息
func Message(code string) string {
	if msg, ok := messages[code]; ok {
		return msg
	}
	return "未知错误"
}
