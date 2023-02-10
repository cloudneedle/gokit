package web

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gocrud/kit/logx"
	"net"
	"net/http"
	"time"
)

func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method

		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token, Authorization, Token")
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
		c.Header("Access-Control-Allow-Credentials", "true")

		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}
		c.Next()
	}
}

// panic recover中间件 并记录日志
func (r *Router) Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				errCode := time.Now().Unix()
				if r.log != nil {
					r.log.Error("panic recover", err)
					r.log.WithField("path", c.Request.URL.Path).
						WithField("stack", fmt.Sprintf("%+v", err)).
						WithField("err_code", errCode).
						Error("服务异常")
				} else {
					fmt.Printf(" err_msg:%v\n code:%d\n stack:\n %+v \n", err, errCode, err)
				}
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"code": errCode,
					"msg":  "Api接口异常，请联系管理员",
				})
			}
		}()
		c.Next()
	}
}

type RouteContext struct {
	Std  gin.IRoutes
	Safe gin.IRoutes
}

type IRoute interface {
	Routes(ctx *RouteContext)
}

type Router struct {
	host   string
	g      *gin.Engine
	auth   AuthFunc
	log    *logx.Log
	isProd bool // 是否开发环境
}

type RouterOptions func(*Router)

func NewRouter(opts ...RouterOptions) *Router {
	r := gin.New()
	r.Use(Cors())
	router := new(Router)
	for _, f := range opts {
		f(router)
	}
	router.g = r
	return router
}

func WithRouterHost(host string) RouterOptions {
	return func(router *Router) {
		router.host = host
	}
}

func WithRouterIsProd(isProd bool) RouterOptions {
	return func(router *Router) {
		router.isProd = isProd
	}
}

func WithRouterAuthFunc(f AuthFunc) RouterOptions {
	return func(router *Router) {
		router.auth = f
	}
}

func WithRouterLog(l *logx.Log) RouterOptions {
	return func(router *Router) {
		router.log = l
	}
}

func (r *Router) Routes(rs ...IRoute) {
	// panic 捕获中间件
	r.g.Use(r.Recovery())
	auth := r.g.Group("", r.auth)
	ctx := &RouteContext{
		Std:  r.g,
		Safe: auth,
	}

	for _, v := range rs {
		v.Routes(ctx)
	}
}

func (r *Router) Run() error {
	if r.host == "" {
		if !r.isProd {
			port, _ := getFreePort()
			r.host = fmt.Sprintf(":%d", port)
		} else {
			r.host = ":80"
		}
	}

	return r.g.Run(r.host)
}

func getFreePort() (int, error) {
	ln, err := net.Listen("tcp", ":0")
	if err != nil {
		return 0, err
	}
	defer ln.Close()
	return ln.Addr().(*net.TCPAddr).Port, nil
}
