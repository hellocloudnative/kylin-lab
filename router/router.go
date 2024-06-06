package router

import (
	"github.com/gin-gonic/gin"
	"kylin-lab/api/vm"
)

func InitSysRouter(r *gin.Engine) *gin.RouterGroup {
	g := r.Group("")
	// 无需认证
	//g.GET("/token", api.GetToken)
	kylinlabNoCheckRoleRouter(g)
	return g
}

func kylinlabNoCheckRoleRouter(r *gin.RouterGroup) {
	v1 := r.Group("/api/v1")
	//v1.Use(middleware.JWTMiddleware())
	registerUserRouter(v1)
}

func registerUserRouter(v1 *gin.RouterGroup) {
	router := v1.Group("/")
	router.GET("/getALLVMlist", vm.GetALLVMList)
	router.GET("/getVMInfo", vm.GetVMInfo)

}
