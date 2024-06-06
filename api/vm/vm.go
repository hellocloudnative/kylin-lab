package vm

import (
	"github.com/gin-gonic/gin"
	"kylin-lab/models"
	"kylin-lab/tools"
	"kylin-lab/tools/app"
	"net/http"
)

func GetALLVMList(c *gin.Context) {
	var data models.LabVirtualMachine
	var err error
	var pageSize = 10
	var pageIndex = 1

	size := c.Request.FormValue("pageSize")
	if size != "" {
		pageSize = tools.StrToInt(err, size)
	}

	index := c.Request.FormValue("pageIndex")
	if index != "" {
		pageIndex = tools.StrToInt(err, index)
	}
	vmId := c.Request.FormValue("vmId")
	data.VmId, _ = tools.StringToInt(vmId)

	userId := c.Request.FormValue("userId")
	data.UserID, _ = tools.StringToInt(userId)

	data.OSType = c.Request.FormValue("osType")
	data.CPUArchitecture = c.Request.FormValue("cpuArchitecture")

	data.Status = c.Request.FormValue("status")

	result, count, err := data.GetPage(pageSize, pageIndex)
	tools.HasError(err, "", -1)

	var mp = make(map[string]interface{}, 3)
	mp["list"] = result
	mp["count"] = count
	mp["pageIndex"] = pageIndex
	mp["pageSize"] = pageSize

	var res app.Response
	res.Data = mp
	c.JSON(http.StatusOK, res.ReturnOK())
}
