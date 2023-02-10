package web

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gocrud/kit/errorx"
	"github.com/gocrud/kit/logx"
	"time"
)

type biz struct {
	Code   int64  `json:"code"`
	Msg    string `json:"msg,omitempty"`
	Detail string `json:"detail,omitempty"`
	Data   any    `json:"data,omitempty"`
}

type Response struct {
	data    any
	_status int
	_err    error
}

func (r *Response) GetError() error {
	return r._err
}

func (r *Response) GetStatus() int {
	return r._status
}

func (r *Response) GetData() any {
	return r.data
}

type RawResp struct {
	Data    any
	_err    error
	_status int
}

func (r *RawResp) GetError() error {
	return r._err
}

func (r *RawResp) GetData() any {
	return r.Data
}

func (r *RawResp) GetStatus() int {
	return r._status
}

type Resp interface {
	GetData() any
	GetStatus() int
	GetError() error
}

type Context struct {
	GIN       *gin.Context
	log       *logx.Log
	codeField string
	msgField  string
}

func (c *Context) getResp(status int, err error) Resp {
	var ex *errorx.Error
	var res Response
	if !errors.As(err, &ex) {
		errCode := time.Now().Unix()
		res = Response{
			data: biz{
				Code: errCode,
				Msg:  "系统异常，请联系管理员",
			},
			_err:    err,
			_status: 500,
		}
		if c.log != nil {
			c.log.WithField("path", c.GIN.Request.URL.Path).
				WithField("stack", fmt.Sprintf("%+v", res.GetError())).
				WithField("err_code", errCode).
				Error(res.GetError().Error())
		} else {
			fmt.Printf(" err_msg:%v\n code:%d\n stack:\n %+v \n", res.GetError(), errCode, res.GetError())
		}
	} else {
		res = Response{
			data: biz{
				Code:   int64(ex.Code().Int()),
				Msg:    fmt.Sprintf("%v", ex.Code()),
				Detail: ex.DetailString(),
			},
			_err:    err,
			_status: status,
		}
	}
	return &res
}

// BizBad http status code 200
func (c *Context) BizBad(err error) Resp {
	return c.getResp(200, err)
}

func (c *Context) UnAuth() Resp {
	return &RawResp{
		Data: gin.H{
			"code": 401,
			"msg":  "身份验证失败",
		},
		_status: 401,
	}
}

// Bad http status code 400
// example response:
//
//	{
//		"code": 业务错误码,
//		"msg": "bad request",
//		"detail": "error detail"
//	}
func (c *Context) Bad(err error) Resp {
	return c.getResp(400, err)
}

// BizData 含code，msg业务字段
// example response:
//
//	{
//		"code": 0,
//		"msg": "ok",
//		"data": {
//			"key": "value"
//		}
func (c *Context) BizData(data any) Resp {
	return &Response{
		data: biz{
			Code: 0,
			Data: data,
		},
		_status: 200,
	}
}

// Data 不含code，msg业务字段
// example response:
//
//	{
//		"key": "value"
//	}
func (c *Context) Data(data any) Resp {
	return &RawResp{
		Data:    data,
		_status: 200,
	}
}

type ContextHandle func(ctx *Context) Resp

func Handle(h ContextHandle) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		res := h(&Context{GIN: ctx})
		ctx.JSON(res.GetStatus(), res.GetData())
	}
}
