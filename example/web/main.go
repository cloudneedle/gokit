package main

import (
	"github.com/gocrud/kit/errorx"
	"github.com/gocrud/kit/log"
	"github.com/gocrud/kit/we"
	"github.com/pkg/errors"
	"go-micro.dev/v4/web"
)

func main() {
	routes := we.WithRoutes(Greeter{})
	server, _ := we.NewServer(routes, we.WithServerHost(":8080"))

	srv := web.NewService(
		web.Handler(server.GIN()),
		web.Address(server.Host()),
	)

	srv.Run()
}

type Greeter struct {
	log *log.Logger
}

func (g Greeter) Routes(ctx *we.RouteContext) {
	ctx.Std.POST("/hello", ctx.Handle(g.SayHello))
}

type helloReq struct {
	Name string `json:"name" binding:"required" msg:"name is required"`
}

// SayHello say hello handler
func (g Greeter) SayHello(ctx *we.Context) any {
	return errors.WithStack(errorx.Code(UserNameErr))
	//if err := ctx.Bind(&req); err != nil {
	//	return ctx.BadError(err)
	//}
	//return ctx.BizData(req)
}
