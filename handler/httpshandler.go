package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/unrolled/secure"
	"kylin-lab/tools/config"
)

func TlsHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		secureMiddleware := secure.New(secure.Options{
			SSLRedirect: true,
			SSLHost:     config.ApplicationConfig.Domain,
		})
		err := secureMiddleware.Process(c.Writer, c.Request)
		if err != nil {
			return
		}
		c.Next()
	}
}
