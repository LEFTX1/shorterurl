package constant

const (
	// 用户相关
	UserLoginKey     = "user:login:"         // 用户登录key
	LockUserRegister = "lock:user:register:" // 用户注册锁

	// 分组相关
	LockGroupCreateKey = "lock:group:create:" // 创建分组锁
	LockGroupUpdateKey = "lock:group:update:" // 更新分组锁
	LockGroupDeleteKey = "lock:group:delete:" // 删除分组锁
	LockGroupSortKey   = "lock:group:sort:"   // 排序分组锁
)
