package web

import (
	"context"
	"fmt"
	"github.com/cloudneedle/gokit/errorx"
	"github.com/gin-gonic/gin"
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

// SetReqHeader set header
func (c *Context) SetReqHeader(key string, value string) {
	c.g.Request.Header.Set(key, value)
}

// GetReqHeader get header
func (c *Context) GetReqHeader(key string) string {
	return c.g.Request.Header.Get(key)
}

// SetRespHeader set header
func (c *Context) SetRespHeader(key string, value string) {
	c.g.Header(key, value)
}

// GetRespHeader get header
func (c *Context) GetRespHeader(key string) string {
	return c.g.GetHeader(key)
}

func (c *Context) BindJson(v any) error {
	err := c.g.BindJSON(v)
	return handleErr(err, v)
}

func (c *Context) BindQuery(v any) error {
	err := c.g.BindQuery(v)
	return handleErr(err, v)
}

func (c *Context) BindForm(v any) error {
	err := c.g.Bind(v)
	return handleErr(err, v)
}

func (c *Context) BindHeader(v any) error {
	err := c.g.BindHeader(v)
	return handleErr(err, v)
}

func (c *Context) BindUri(v any) error {
	err := c.g.BindUri(v)
	return handleErr(err, v)
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

type std struct {
	status int `json:"-"`
	data   any `json:"data,omitempty"`
}

func (s *std) Status() int {
	return s.status
}

func (s *std) GetData() any {
	return s.data
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

// Biz 业务数据,可以指定code和msg
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
func (c *Context) Biz(code int, msg string, data any) ICustomResp {
	return &biz{
		status: 200,
		Code:   code,
		Msg:    msg,
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

// Data http状态码为200的数据返回
//
// http status: 200
//
// example:
//
//	{
//	  "id": 1,
//	  "name": "张三"
//	}
func (c *Context) Data(data any) ICustomResp {
	return &std{
		status: 200,
		data:   data,
	}
}

// Bad http状态码为400的错误返回，可以指定code和msg
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

// BadError http状态码为400的错误返回，code=400
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

// BadCode http状态码为400的错误返回，自定义code
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

// UnAuth 未授权
//
// http status: 401
//
// example:
//
//	{
//	  "code": 401,
//	  "msg": "未授权"
//	}

func (c *Context) UnAuth() ICustomResp {
	return &biz{
		status: 401,
		Code:   401,
		Msg:    "未授权",
	}
}

// Forbidden 禁止访问
//
// http status: 403
//
// example:
//
//	{
//	  "code": 403,
//	  "msg": "禁止访问"
//	}
func (c *Context) Forbidden() ICustomResp {
	return &biz{
		status: 403,
		Code:   403,
		Msg:    "禁止访问",
	}
}
