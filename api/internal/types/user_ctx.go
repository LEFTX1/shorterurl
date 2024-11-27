package types

import "context"

// ContextKey 定义 context key 类型
type ContextKey string

const (
	// UserContextKey 用户信息的 context key
	UserContextKey ContextKey = "user_info"
)

// UserInfo 用户上下文信息
type UserInfo struct {
	ID       string `json:"id"`        // 用户ID
	Username string `json:"username"`  // 用户名
	RealName string `json:"real_name"` // 真实姓名
}

// GetUserFromCtx 从 context 中获取用户信息
func GetUserFromCtx(ctx context.Context) (*UserInfo, bool) {
	userInfo, ok := ctx.Value(UserContextKey).(*UserInfo)
	return userInfo, ok
}

// MustGetUserFromCtx 从 context 中获取用户信息，如果不存在则 panic
func MustGetUserFromCtx(ctx context.Context) *UserInfo {
	userInfo, ok := GetUserFromCtx(ctx)
	if !ok {
		panic("user info not found in context")
	}
	return userInfo
}
