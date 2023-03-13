package main

import (
	"errors"
	"github.com/cloudneedle/gokit/log"
	"github.com/cloudneedle/gokit/web"
	"go-micro.dev/v4/logger"
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
	srv, err := web.NewServer(web.WithServerHost(":80"), web.WithRoutes(Greeter{}))
	if err != nil {
		logger.Fatal(err)
	}
	srv.Run()
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
	return ctx.BadError(testErrorsWithMessage())
}

func testErrorsWithMessage() error {
	err := errors.New("test error")
	return err
}
