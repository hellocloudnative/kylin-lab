package vm

import (
	"github.com/gin-gonic/gin"
	"kylin-lab/tools/app"
	"net/http"
)

func StartInstances(c *gin.Context) {
	// 获取请求参数
	instanceID := c.Param("instance_id")

	// 调用启动机器的逻辑
	yesOrNo, err := PostStartInstances(instanceID)
	if err != nil {
		app.Error(c, http.StatusBadRequest, err, "启动机器失败")
		return
	}
	if yesOrNo {
		app.OK(c, http.StatusOK, "启动机器成功")
		return
	} else {
		app.Error(c, http.StatusBadRequest, err, "启动机器失败")
		return
	}
}
