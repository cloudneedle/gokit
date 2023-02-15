package gw

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gocrud/kit/errorx"
	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
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

// BindJsonpb 绑定数据到protobuf struct
func (c *Context) BindJsonpb(v proto.Message) error {
	return jsonpb.Unmarshal(c.g.Request.Body, v)
}

type ICustomResp interface {
	Status() int
	GetData() any
}

type biz struct {
	status int    `json:"-"`
	Code   int    `json:"code"`
	Msg    string `json:"msg,omitempty"`
	Detail string `json:"detail,omitempty"`
	Data   any    `json:"data,omitempty"`
}

func (b *biz) Status() int {
	return b.status
}

func (b *biz) GetData() any {
	return b
}

// BizData 业务数据
//
// http status: 200
//
// example:
//
//	{
//	  "code": 0,
//	  "msg": "success",
//	  "data": {
//	    "id": 1,
//	    "name": "张三"
//	  }
//	}
func (c *Context) BizData(data any) ICustomResp {
	return &biz{
		status: 200,
		Data:   data,
	}
}

// BizBadCode 业务错误，自定义错误码
//
// http status: 200
//
// example:
//
//	{
//	  "code": 1001,
//	  "msg": "用户名或密码错误"
//	}
func (c *Context) BizBadCode(err errorx.ErrorCode) ICustomResp {
	return &biz{
		status: 200,
		Code:   err.Int(),
		Msg:    fmt.Sprintf("%v", err),
	}
}

// BizBad 业务错误,可以指定code和msg
//
// http status: 200
//
// example:
//
//	{
//	  "code": 1001,
//	  "msg": "用户名或密码错误"
//	}
func (c *Context) BizBad(code int, msg string) ICustomResp {
	return &biz{
		status: 200,
		Code:   code,
		Msg:    msg,
	}
}

// BizBadError 业务错误，code=400
//
// http status: 200
//
// example:
//
//	{
//	  "code": 400,
//	  "msg": "用户名或密码错误"
//	}
func (c *Context) BizBadError(err error) ICustomResp {
	return &biz{
		status: 200,
		Code:   400,
		Msg:    fmt.Sprintf("%v", err),
	}
}

// Bad 业务错误,可以指定code和msg
//
// http status: 400
//
// example:
//
//	{
//	  "code": 1001,
//	  "msg": "用户名或密码错误"
//	}
func (c *Context) Bad(code int, msg string) ICustomResp {
	return &biz{
		status: 400,
		Code:   code,
		Msg:    msg,
	}
}

// BadError 业务错误，code=400
//
// http status: 400
//
// example:
//
//	{
//	  "code": 400,
//	  "msg": "用户名或密码错误"
//	}
func (c *Context) BadError(err error) ICustomResp {
	return &biz{
		status: 400,
		Code:   400,
		Msg:    fmt.Sprintf("%v", err),
	}
}

// BadCode 业务错误，自定义错误码
//
// http status: 400
//
// example:
//
//	{
//	  "code": 1001,
//	  "msg": "用户名或密码错误"
//	}
func (c *Context) BadCode(err errorx.ErrorCode) ICustomResp {
	return &biz{
		status: 400,
		Code:   err.Int(),
		Msg:    fmt.Sprintf("%v", err),
	}
}
