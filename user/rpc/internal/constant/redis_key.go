package constant

const (
	// 用户相关
	USER_LOGIN_KEY     = "user:login:"         // 用户登录key
	Lock_User_Register = "lock:user:register:" // 用户注册锁

	// 分组相关
	LOCK_GROUP_CREATE_KEY = "lock:group:create:" // 创建分组锁
	LOCK_GROUP_UPDATE_KEY = "lock:group:update:" // 更新分组锁
	LOCK_GROUP_DELETE_KEY = "lock:group:delete:" // 删除分组锁
	LOCK_GROUP_SORT_KEY   = "lock:group:sort:"   // 排序分组锁
)
