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

func GetVMInfo(c *gin.Context) {
	var data models.LabVirtualMachine
	var user models.LabUser
	vmId := c.Request.FormValue("vmId")
	if vmId == "" {
		app.OK(c, nil, "参数错误")
		return
	}
	data.VmId, _ = tools.StringToInt(vmId)

	virtualMachineInfo, err := data.GetVirtualMachineInfo()
	if err != nil {
		// 处理GetVirtualMachineInfo方法的错误
		app.Error(c, 500, err, "查询失败")
		return
	}
	userInfo, err := user.GetUserInfo(virtualMachineInfo.UserID)
	if err != nil {
		app.Error(c, 500, err, "查询失败")
		return
	}
	virtualMachineInfo.UserName.Username = userInfo.Username
	result := map[string]interface{}{
		"data": virtualMachineInfo,
	}
	tools.HasError(err, "查询失败", 500)
	app.OK(c, result, "查询成功")
}
