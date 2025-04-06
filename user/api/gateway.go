package main

import (
	"flag"
	"fmt"

	"shorterurl/user/api/internal/config"
	"shorterurl/user/api/internal/handler"
	"shorterurl/user/api/internal/svc"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
)

var configFile = flag.String("f", "etc/gateway.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()

	ctx := svc.NewServiceContext(c)

	// 不再需要全局注册中间件，改为在API文件中通过注解声明
	// 中间件已经在 API 文件中声明，并在 ServiceContext 中实例化

	handler.RegisterHandlers(server, ctx)

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}
