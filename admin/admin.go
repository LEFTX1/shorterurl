package main

import (
	"flag"
	"fmt"

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
	handler.RegisterHandlers(server, ctx)

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}
