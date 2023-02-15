package main

import (
	"github.com/gocrud/kit/errorx"
	"github.com/gocrud/kit/gw"
	"github.com/gocrud/kit/log"
	"github.com/pkg/errors"
	"go-micro.dev/v4/web"
)

func main() {
	routes := gw.WithRoutes(Greeter{})
	server, _ := gw.NewServer(routes, gw.WithServerHost(":8080"))

	srv := web.NewService(
		web.Handler(server.GIN()),
		web.Address(server.Host()),
	)

	srv.Run()
}

type Greeter struct {
	log *log.Logger
}

func (g Greeter) Routes(ctx *gw.RouteContext) {
	ctx.Std.POST("/hello", ctx.Handle(g.SayHello))
}

type helloReq struct {
	Name string `json:"name" binding:"required" msg:"name is required"`
}

// SayHello say hello handler
func (g Greeter) SayHello(ctx *gw.Context) any {
	return errors.WithStack(errorx.Code(UserNameErr))
	//if err := ctx.Bind(&req); err != nil {
	//	return ctx.BadError(err)
	//}
	//return ctx.BizData(req)
}
