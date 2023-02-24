package web

import (
	"context"
	"errors"
	"fmt"
	"github.com/cloudneedle/gokit/errorx"
	"github.com/gin-gonic/gin"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
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

type RouteContext struct {
	gin.IRoutes
	Auth gin.IRoutes
}

func (r *RouteContext) Handle(fn func(ctx *Context) any) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := &Context{c}
		res := fn(ctx)
		// 判断是否是自定义响应
		if customResp, ok := res.(ICustomResp); ok {
			c.JSON(customResp.Status(), customResp.GetData())
			return
		}
		// 判断是否是错误
		if err, ok := res.(error); ok {
			// 判断是否是自定义错误
			var customErr *errorx.Error
			if ok := errors.As(err, &customErr); ok {
				c.JSON(400, gin.H{
					"code": customErr.Code(),
					"msg":  customErr.Error(),
				})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{
				"code": 500,
				"msg":  err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, res)
	}
}

type IRoute interface {
	Routes(ctx *RouteContext)
}

type Server struct {
	host           string
	isDebug        bool
	routes         []IRoute
	authMiddleware gin.HandlerFunc
	g              *gin.Engine
}

// ServerOption Server Option type
type ServerOption func(*Server)

// WithServerHost 设置host
func WithServerHost(host string) ServerOption {
	return func(s *Server) {
		s.host = host
	}
}

// WithRoutes 设置路由
func WithRoutes(routes ...IRoute) ServerOption {
	return func(s *Server) {
		s.routes = routes
	}
}

// WithAuthMiddleware 设置认证中间件
func WithAuthMiddleware(authMiddleware gin.HandlerFunc) ServerOption {
	return func(s *Server) {
		s.authMiddleware = authMiddleware
	}
}

// NewServer 创建一个新的Server,默认debug模式
func NewServer(opts ...ServerOption) (*Server, error) {
	s := &Server{
		isDebug: true,
	}

	for _, opt := range opts {
		opt(s)
	}

	// 设置默认host
	err := s.getFreeHost()
	if err != nil {
		return nil, err
	}

	// 设置http server
	s.setHttpServer()

	return s, nil
}

func (s *Server) setHttpServer() {
	if !s.isDebug {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()
	r.Use(Cors())
	r.Use(gin.Recovery())

	authRoute := r.Group("", s.authMiddleware)
	// 注册路由
	routeContext := &RouteContext{
		IRoutes: r,
		Auth:    authRoute,
	}

	for _, route := range s.routes {
		route.Routes(routeContext)
	}

	s.g = r
}

func (s *Server) GIN() *gin.Engine {
	return s.g
}

func (s *Server) Host() string {
	return s.host
}

// Run 运行Server
func (s *Server) Run() {
	srv := &http.Server{
		Addr:    s.host,
		Handler: s.g,
	}

	go func() {
		// 打印服务启动信息
		log.Printf("Server is running on %s", s.host)
		// 服务连接
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()
	// 等待中断信号以优雅地关闭服务器（设置 5 秒的超时时间）
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	log.Println("Server exiting")
}

// getFreeHost 获取一个空闲的host,如果未指定host，动态获取有效服务端host
func (s *Server) getFreeHost() error {
	if s.host != "" {
		return nil
	}

	ln, err := net.Listen("tcp", ":0")
	if err != nil {
		return err
	}
	defer ln.Close()
	port := ln.Addr().(*net.TCPAddr).Port
	s.host = fmt.Sprintf(":%d", port)
	return nil
}
