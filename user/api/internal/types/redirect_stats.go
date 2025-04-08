package types

import (
	"net/http"
)

// RedirectStats 存储重定向统计信息
type RedirectStats struct {
	ShortUri     string
	FullShortUrl string
	Gid          string
	Ip           string
	UserAgent    string
	Browser      string
	BrowserVer   string
	Os           string
	OsVer        string
	Device       string
	DeviceModel  string
	Network      string
	Locale       string
	User         string
	UvFirstFlag  bool
	UipFirstFlag bool
}

// RedirectStatMiddleware 处理短链接跳转和统计数据收集
type RedirectStatMiddleware interface {
	Handle(next http.HandlerFunc) http.HandlerFunc
}
