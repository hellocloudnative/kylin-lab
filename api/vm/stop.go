package vm

import (
	"github.com/gin-gonic/gin"
	"kylin-lab/tools/app"
	"net/http"
)

func StopInstances(c *gin.Context) {
	// 获取请求参数
	instanceID := c.Param("instance_id")

	// 调用停止机器的逻辑
	yesOrNo, err := PostStopInstances(instanceID)
	if err != nil {
		app.Error(c, http.StatusBadRequest, err, "停止机器失败")
		return
	}
	if yesOrNo {
		app.OK(c, http.StatusOK, "停止机器成功")
		return
	} else {
		app.Error(c, http.StatusBadRequest, err, "停止机器失败")
		return
	}
}
