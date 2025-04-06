package config

// TokenValidateConfig Token 验证中间件配置
type TokenValidateConfig struct {
	WhitePathList  []string `json:"whitePathList"`  // 白名单路径列表
	LoginKeyPrefix string   `json:"loginKeyPrefix"` // 登录令牌的Redis键前缀，默认为"user:login:"
}
