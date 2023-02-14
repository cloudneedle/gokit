package web

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gocrud/kit/errorx"
)

type Context struct {
	g *gin.Context
}

func (c *Context) Context() context.Context {
	return c.g.Request.Context()
}

func (c *Context) Set(key string, value any) {
	c.g.Set(key, value)
}

func (c *Context) BindJson(v any) error {
	return c.g.BindJSON(v)
}

func (c *Context) BindQuery(v any) error {
	return c.g.BindQuery(v)
}

func (c *Context) BindForm(v any) error {
	return c.g.Bind(v)
}

func (c *Context) BindHeader(v any) error {
	return c.g.BindHeader(v)
}

func (c *Context) BindUri(v any) error {
	return c.g.BindUri(v)
}

func (c *Context) Get(key string) (value any, exists bool) {
	return c.g.Get(key)
}

func (c *Context) Bind(v any) error {
	err := c.g.ShouldBind(v)
	return handleErr(err, v)
}

func (c *Context) Data(data any) {
	c.g.JSON(200, data)
}

func (c *Context) Bad(err error) {
	c.g.JSON(400, gin.H{
		"code": 400,
		"msg":  err.Error(),
	})
}

func (c *Context) Unauthorized() {
	c.g.JSON(401, gin.H{
		"code": 401,
		"msg":  "Unauthorized",
	})
}

func (c *Context) Forbidden() {
	c.g.JSON(403, gin.H{
		"code": 403,
		"msg":  "Forbidden",
	})
}

func (c *Context) InternalServerError(data any) {
	c.g.JSON(500, data)
}

type ICustomResp interface {
	Status() int
	Data() any
}

func (c *Context) Custom(resp ICustomResp) {
	c.g.JSON(resp.Status(), resp.Data())
}

type biz struct {
	Code   int    `json:"code"`
	Msg    string `json:"msg,omitempty"`
	Detail string `json:"detail,omitempty"`
	Data   any    `json:"data,omitempty"`
}

func (c *Context) BizData(data any) {
	c.g.JSON(200, &biz{
		Code: 0,
		Data: data,
	})
}

func (c *Context) BizError(err errorx.ErrorCode) {
	c.g.JSON(200, &biz{
		Msg:  fmt.Sprintf("%s", err),
		Code: err.Int(),
	})
}
