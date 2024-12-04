package config

// TokenValidateConfig Token 验证中间件配置
type TokenValidateConfig struct {
	WhitePathList []string `json:"whitePathList"` // 白名单路径列表
}
