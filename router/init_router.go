package router

import (
	"github.com/gin-gonic/gin"
	"kylin-lab/handler"
	"kylin-lab/middleware"
	config2 "kylin-lab/tools/config"
)

func InitRouter() *gin.Engine {
	r := gin.New()
	if config2.ApplicationConfig.IsHttps {
		r.Use(handler.TlsHandler())
	}
	middleware.InitMiddleware(r)
	InitSysRouter(r)
	return r
}
