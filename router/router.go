package router

import (
	"github.com/gin-gonic/gin"
	"kylin-lab/api"
	"kylin-lab/api/vm"
	jwt2 "kylin-lab/jwt"
)

func InitSysRouter(r *gin.Engine) *gin.RouterGroup {
	g := r.Group("")
	// 无需认证
	//g.GET("/token", api.GetToken)
	kylinlabNoCheckRoleRouter(g)
	g.GET("/get-kylincloud-token", vm.GetKylinCloudToken)
	g.GET("/get-kylincloud-images", vm.GetKylinCloudImages)
	g.GET("/lab-token", api.LabToken)
	return g
}

func kylinlabNoCheckRoleRouter(r *gin.RouterGroup) {
	v1 := r.Group("/api/v1")
	v1.Use(jwt2.AuthMiddleware)
	registerUserRouter(v1)
}

func registerUserRouter(v1 *gin.RouterGroup) {
	router := v1.Group("/")
	router.GET("/getALLVMlist", vm.GetALLVMList)
	router.GET("/getVMInfo", vm.GetVMInfo)
	router.GET("/getALLImagesList", vm.GetALLImagesList)
	router.POST("/applyInstances", vm.ApplyInstances)
	router.PUT("/updateVMStatus", vm.UpdateVMStatus)

}
