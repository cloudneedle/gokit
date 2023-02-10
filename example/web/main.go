package main

import (
	"github.com/gocrud/kit/web"
)

func main() {
	r := web.NewRouter(web.WithRouterIsProd(true))
	r.Routes(TestHandler{})
	r.Run()
}

type TestHandler struct {
}

func (t TestHandler) Routes(ctx *web.RouteContext) {
	ctx.Std.POST("/test", web.Handle(t.Post))
}

func (TestHandler) Post(c *web.Context) web.Resp {
	//return c.BizBad(errorx.Code(UserNameErr))
	//return c.BizData("sss")
	//return c.Data("sss")
	return c.UnAuth()
}
