// types/errcode/errcode.go
package errcode

const (
	// 用户相关错误码 (B0002xx)
	UserNotFound  = "B000200" // 用户记录不存在
	UserNameExist = "B000201" // 用户名已存在
	UserExist     = "B000202" // 用户记录已存在
	UserSaveError = "B000203" // 用户记录新增失败
)
