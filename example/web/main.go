package main

import (
	"github.com/gocrud/kit/web"
)

func main() {
	routes := web.WithRoutes(Greeter{})
	server := web.NewServer(routes)
	server.Run()
}

type Greeter struct {
}

func (g Greeter) Routes(ctx *web.RouteContext) {
	ctx.Std.POST("/hello", ctx.Handle(g.SayHello))
}

type helloReq struct {
	Name string `json:"name" binding:"required" msg:"name is required"`
}

// SayHello say hello handler
func (g Greeter) SayHello(ctx *web.Context) {
	var req helloReq
	if err := ctx.Bind(&req); err != nil {
		ctx.Bad(err)
		return
	}
	ctx.Data(req)
}
