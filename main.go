package main

import (
	"flag"
	"github.com/zeromicro/go-zero/core/logx"

	"gpu-container-service/internal/config"
	"gpu-container-service/internal/handler"
	"gpu-container-service/internal/svc"
	"gpu-container-service/pkg"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
)

var configFile = flag.String("f", "etc/config.yaml", "the config file")

func main() {
	flag.Parse()

	logx.Errorf("GPU Container Service Initializing...")
	var c config.Config
	conf.UseEnv()
	conf.MustLoad(*configFile, &c)
	logx.Errorf("Config loaded from environment: %+v", c)

	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()

	ctx := svc.NewGpuContainerServiceContext(c)
	logx.Errorf("Router Registration Processing")
	r := handler.NewRouter(server)
	builder, err := pkg.NewClientBuilder("")
	if err != nil {
		logx.Errorf("New Client Builder error: %v", err)
	}
	client, err := builder.Client()
	if err != nil {
		logx.Errorf("New Kubernetes Client error: %v", err)
	}
	handler.RegisterHandlers(r, ctx, client)

	logx.Errorf("Starting Server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}
