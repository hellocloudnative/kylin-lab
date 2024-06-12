package vm

import (
	"errors"
	"github.com/gin-gonic/gin"
	log "github.com/wonderivan/logger"
	"kylin-lab/models"
	"kylin-lab/tools"
	"kylin-lab/tools/app"
	"net/http"
)

type InstancesInfo struct {
	HwArchitecture string `json:"hw_architecture"`
	ImageName      string `json:"imageName"`
	Flavors        string `json:"flavors"`
	NetworkName    string `json:"networkName"`
	TimeOfuse      int    `json:"timeOfuse"`
}

func ApplyInstances(c *gin.Context) {
	var instancesInfo InstancesInfo
	err := c.ShouldBindJSON(&instancesInfo)
	tools.HasError(err, "数据解析失败", -1)
	// 创建虚拟机
	serverInfoResponse, err := PostApplyInstances(instancesInfo.HwArchitecture, instancesInfo.ImageName, instancesInfo.Flavors, instancesInfo.NetworkName, instancesInfo.TimeOfuse)
	if err != nil {
		log.Error("申请机器失败", err)
		app.Error(c, http.StatusInternalServerError, err, "申请机器失败")
		return
	}
	log.Info("申请机器成功", serverInfoResponse.Servers[0].Id)
	data := instancesInfo
	log.Info("申请机器成功", data)
	app.OK(c, data, "申请机器成功")

}

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

func UpdateVMStatus(c *gin.Context) {
	var data models.LabVirtualMachine
	err := c.Bind(&data)
	tools.HasError(err, "参数错误", -1)
	if data.VmId == 0 {
		app.Error(c, 500, errors.New("ID不能为空"), "更新失败")
		return
	}

	if data.Status == "" {
		app.Error(c, 500, errors.New("状态不能为空"), "更新失败")
		return
	}

	duration, _ := tools.StringToInt(data.Duration)
	status, _ := tools.StringToInt(data.Status)

	if status > 1 || status < 0 {
		app.Error(c, 500, errors.New("状态错误"), "更新失败")
		return
	}
	if status == 0 && duration <= 0 {
		app.Error(c, 500, errors.New("使用状态下时长不能为0"), "更新失败")
		return
	}

	if status == 1 || duration < 0 {
		//deletevm
		log.Info("delete vm")
	}

	updatedData, err := data.Update(data.VmId, data.Duration, data.Status)
	if err != nil {
		app.Error(c, 500, err, "更新失败")
		return
	}
	app.OK(c, updatedData, "更新成功")
}

func GetALLImagesList(c *gin.Context) {
	images, err := GetImagesRequest()
	if err != nil {
		log.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve token"})
		return
	}
	c.JSON(http.StatusOK, images.Images)
}
