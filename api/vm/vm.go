package vm

import (
	"fmt"
	"github.com/gin-gonic/gin"
	log "github.com/wonderivan/logger"
	"kylin-lab/models"
	"kylin-lab/tools"
	"kylin-lab/tools/app"
	"net/http"
	"time"
)

func ApplyInstances(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Origin", "")
	var instancesInfo models.LabVirtualMachine
	err := c.ShouldBindJSON(&instancesInfo)
	tools.HasError(err, "数据解析失败", -1)

	instancesInfo.UserID = 1
	userNum, err := instancesInfo.GetUserVirtualMachineStatus0Count(instancesInfo.UserID)
	if err != nil {
		log.Error("获取用户虚拟机数量失败", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": -1,
			"msg":    "获取用户虚拟机数量失败"})
		return
	}
	if userNum >= 1 {
		log.Error("用户虚拟机数量超过限制", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": -1,
			"msg":    "用户虚拟机数量超过限制"})
		return
	}
	// 创建虚拟机
	serverInfoResponse, err := PostApplyInstances(instancesInfo.CPUArchitecture, instancesInfo.OSImage, instancesInfo.Flavors, instancesInfo.NetworkName)
	if err != nil {
		log.Error("申请机器失败", err)
		app.Error(c, http.StatusInternalServerError, err, "申请机器失败")
		return
	}

	// 轮询虚拟机状态
	maxWaitTime := 5 * time.Minute
	pollInterval := 1 * time.Second
	timeout := time.After(maxWaitTime)
	for {
		select {
		case <-time.After(pollInterval):
			serverInfo, err := GetServersRequest(serverInfoResponse.Servers[0].Id)
			if err != nil {
				log.Error("获取服务器信息失败", err)
				app.Error(c, http.StatusInternalServerError, err, "获取服务器信息失败")
				return
			}
			if len(serverInfo.Servers) == 0 {
				log.Error("未找到服务器信息")
				continue
			}
			if serverInfo.Servers[0].Status == "ACTIVE" {
				log.Info("服务器状态为 %s", serverInfo.Servers[0].Status)
				instancesInfo.IPAddress = serverInfo.Servers[0].Addresses.Vxlan[0].Addr
				instancesInfo.ApplyStatus = "0"
				instancesInfo.Status = "0"
				duration, err := tools.StringToInt(instancesInfo.Duration)
				if err != nil {
					log.Error(err)
				}
				instancesInfo.UUID = serverInfoResponse.Servers[0].Id
				vncInfo, err := GetVNCRequest(serverInfoResponse.Servers[0].Id)
				if err != nil {
					log.Error(err)
				}
				proxyAddress := vncInfo.Url
				instancesInfo.VNCAddress = proxyAddress
				instancesInfo.ApplyTime = time.Now()
				instancesInfo.TimeOfuse = time.Now().Format("2006-01-02 15:04:05") + "-" + time.Now().Add(time.Duration(duration)*time.Minute).Format("2006-01-02 15:04:05")
				instancesInfo.VmLog = "系统自动消息: 审批通过"

				instancesInfo.CreatedAt = time.Now()
				instancesInfo.UpdatedAt = time.Now()
				data := instancesInfo

				log.Info("申请机器成功", data)

				data.Create()

				app.OK(c, data, "申请机器成功")
				// 启动定时任务，n分钟后删除虚拟机
				time.AfterFunc(time.Duration(duration)*time.Minute, func() {
					ParmReturnInstances(c, data.UUID)
				})
				return
			} else if serverInfo.Servers[0].Status == "ERROR" {
				log.Error("服务器状态为 %s", serverInfo.Servers[0].Status)
				app.Error(c, http.StatusInternalServerError, fmt.Errorf("服务器状态为ERROR"), "申请机器失败")
				return
			}
			log.Info("服务器状态为: %s, 继续等待...", serverInfo.Servers[0].Status)
		case <-timeout:
			log.Info("等待超时，未能成功获取到ACTIVE状态的服务器信息")
			app.Error(c, http.StatusInternalServerError, fmt.Errorf("等待服务器状态超时"), "申请机器失败")
			return
		}
	}
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

	tools.HasError(err, "查询失败", 500)
	app.OK(c, virtualMachineInfo, "查询成功")
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
