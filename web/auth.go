package web

import "github.com/gin-gonic/gin"

type AuthFunc = func(ctx *gin.Context)
