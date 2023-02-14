package web

import (
	"context"
	"fmt"
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
	Std  gin.IRoutes
	Safe gin.IRoutes
}

func (r *RouteContext) Handle(fn func(ctx *Context)) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := &Context{c}
		fn(ctx)
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
}

// ServerOption Server Option type
type ServerOption func(*Server)

// WithServerIsDebug 设置debug模式
func WithServerIsDebug(isDebug bool) ServerOption {
	return func(s *Server) {
		s.isDebug = isDebug
	}
}

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
func NewServer(opts ...ServerOption) *Server {
	s := &Server{
		isDebug: true,
	}
	for _, opt := range opts {
		opt(s)
	}
	return s
}

// Run 运行Server
func (s *Server) Run() {
	if !s.isDebug {
		gin.SetMode(gin.ReleaseMode)
	}

	s.getFreeHost()
	r := gin.New()
	r.Use(Cors())
	r.Use(gin.Recovery())

	srv := &http.Server{
		Addr:    s.host,
		Handler: r,
	}

	authRoute := r.Group("", s.authMiddleware)
	// 注册路由
	routeContext := &RouteContext{
		Std:  r,
		Safe: authRoute,
	}

	for _, route := range s.routes {
		route.Routes(routeContext)
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

// getFreeHost 获取一个空闲的host,如果未指定host，开发模式下动态获取有效服务端host,生产模式下使用默认端口80
func (s *Server) getFreeHost() error {
	if s.host != "" {
		return nil
	}

	// 判读是否debug模式
	if s.isDebug {
		ln, err := net.Listen("tcp", ":0")
		if err != nil {
			return err
		}
		defer ln.Close()
		port := ln.Addr().(*net.TCPAddr).Port
		s.host = fmt.Sprintf(":%d", port)
		return nil
	}

	s.host = ":80"
	return nil
}
