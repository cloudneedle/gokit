package main

import (
	"github.com/gocrud/kit/errorx"
	"github.com/gocrud/kit/log"
	"github.com/gocrud/kit/web"
	"github.com/pkg/errors"
	webx "go-micro.dev/v4/web"
)

func main() {
	routes := web.WithRoutes(Greeter{})
	server, _ := web.NewServer(routes, web.WithServerHost(":8080"))

	srv := webx.NewService(
		webx.Handler(server.GIN()),
		webx.Address(server.Host()),
	)
	srv.Run()
}

type Greeter struct {
	log *log.Logger
}

func (g Greeter) Routes(ctx *web.RouteContext) {
	ctx.Std.POST("/hello", ctx.Handle(g.SayHello))
}

type helloReq struct {
	Name string `json:"name" binding:"required" msg:"name is required"`
}

// SayHello say hello handler
func (g Greeter) SayHello(ctx *web.Context) any {
	return errors.WithStack(errorx.Code(UserNameErr))
	//if err := ctx.Bind(&req); err != nil {
	//	return ctx.BadError(err)
	//}
	//return ctx.BizData(req)
}
