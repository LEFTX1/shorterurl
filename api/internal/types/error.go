package types

// GatewayErrorResult 网关错误返回信息
type GatewayErrorResult struct {
	Status  int    `json:"status"`  // HTTP 状态码
	Message string `json:"message"` // 返回信息
}
