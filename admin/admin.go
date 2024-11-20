package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/httpx"
	"go-zero-shorterurl/admin/internal/types/errorx"
	"net/http"

	"go-zero-shorterurl/admin/internal/config"
	"go-zero-shorterurl/admin/internal/handler"
	"go-zero-shorterurl/admin/internal/svc"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
)

var configFile = flag.String("f", "etc/admin-api.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()

	// 创建服务上下文（现在会处理错误）
	ctx, err := svc.NewServiceContext(c)
	if err != nil {
		panic(fmt.Sprintf("create service context failed: %v", err))
	}
	// 添加 defer 关闭资源
	defer func() {
		if err := ctx.Close(); err != nil {
			logx.Errorf("close service context failed: %v", err)
		}
	}()
	handler.RegisterHandlers(server, ctx)
	// 自定义错误
	httpx.SetErrorHandler(func(err error) (int, interface{}) {
		var baseErr *errorx.BaseError
		var userErr *errorx.UserError
		var sysErr *errorx.SystemError
		var remoteErr *errorx.RemoteCallError

		// 使用 errors.As 进行类型断言
		switch {
		case errors.As(err, &userErr):
			return http.StatusOK, map[string]interface{}{
				"code":    userErr.Code,
				"message": userErr.Message,
			}
		case errors.As(err, &sysErr):
			return http.StatusInternalServerError, map[string]interface{}{
				"code":    sysErr.Code,
				"message": sysErr.Message,
			}
		case errors.As(err, &remoteErr):
			return http.StatusServiceUnavailable, map[string]interface{}{
				"code":    remoteErr.Code,
				"message": remoteErr.Message,
			}
		case errors.As(err, &baseErr):
			return http.StatusInternalServerError, map[string]interface{}{
				"code":    baseErr.Code,
				"message": baseErr.Message,
			}
		default:
			return http.StatusInternalServerError, map[string]interface{}{
				"code":    "UnknownError",
				"message": err.Error(),
			}
		}
	})
	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}
