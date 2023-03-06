package main

import (
	"github.com/cloudneedle/gokit/config"
	"github.com/cloudneedle/gokit/errorx"
	"github.com/cloudneedle/gokit/log"
	"github.com/cloudneedle/gokit/web"
	"github.com/pkg/errors"
	"go-micro.dev/v4/logger"
	webx "go-micro.dev/v4/web"
)

func newWebServer(host string) *web.Server {
	routes := web.WithRoutes(Greeter{})

	server, err := web.NewServer(routes, web.WithServerHost(host))
	if err != nil {
		logger.Fatal(err)
	}

	return server
}

func main() {
	cfgClient, err := config.New(config.WithPath("127.0.0.1:2379"))
	if err != nil {
		panic(err)
	}
	// 读取配置
	cfg, err := cfgClient.Etcd.GetPrefix("admin")
	if err != nil {
		panic(err)
	}
	//判断配置是否为空
	if len(cfg) == 0 {
		panic("配置不能为空")
	}

	var serverHost string
	var logPath string
	// 读取配置
	for k, v := range cfg {
		switch k {
		case "admin/server_host":
			serverHost = v
		case "admin/log_path":
			logPath = v
		}
	}
	logger := log.New(log.WithLogPath(logPath))
	server := newWebServer(serverHost)
	srv := webx.NewService(
		webx.Handler(server.GIN()),
		webx.Address(server.Host()),
	)
	go func() {
		cfgClient.Etcd.WatchPrefix("admin", func(m map[string]string) {
			srv.Stop()
			for k, v := range m {
				switch k {
				case "server_host":
					serverHost = v
				case "log_path":
					logPath = v
				}
			}
			s := newWebServer(serverHost)
			srv := webx.NewService(
				webx.Handler(s.GIN()),
				webx.Address(s.Host()),
			)
			srv.Run()
		})
	}()
	go func() {
		if err := srv.Run(); err != nil {
			logger.Fatal(err)
		}
	}()
	select {}
}

type Greeter struct {
	log *log.Logger
}

func (g Greeter) Routes(ctx *web.RouteContext) {
	ctx.POST("/hello", ctx.Handle(g.SayHello))
}

type helloReq struct {
	Name string `json:"name" binding:"required" msg:"name is required"`
}

// SayHello say hello handler
func (g Greeter) SayHello(ctx *web.Context) any {
	//if err := ctx.Bind(&req); err != nil {
	//	return ctx.BadError(err)
	//}
	//return ctx.BizData(req)
	return testErrorsWithMessage()
}

func testErrorsWithMessage() error {
	err := errors.WithMessage(errorx.Code(UserNameErr), "test")
	return err
}
