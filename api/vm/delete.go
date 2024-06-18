package vm

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"kylin-lab/models"
	"kylin-lab/tools/app"
	"net/http"
)

func DeleteInstances(c *gin.Context) {
	// 获取请求参数
	instanceID := c.Param("instance_id")

	// 调用删除实例的逻辑
	yesOrNo, err := PostDeleteInstances(instanceID)
	if err != nil {
		app.Error(c, http.StatusBadRequest, err, "删除实例失败")
		return
	}
	if yesOrNo {
		app.OK(c, http.StatusOK, "删除实例成功")
		return
	} else {
		app.Error(c, http.StatusBadRequest, err, "删除实例失败")
		return
	}
}

func DeleteRecycleInstances(c *gin.Context) {
	// 获取请求参数
	instanceID := c.Param("instance_id")

	// 调用删除实例的逻辑
	yesOrNo, err := PostDeleteRecycleInstances(instanceID)
	if err != nil {
		app.Error(c, http.StatusBadRequest, err, "删除回收站实例失败")
		return
	}
	if yesOrNo {
		app.OK(c, http.StatusOK, "删除回收站实例成功")
		return
	} else {
		app.Error(c, http.StatusBadRequest, err, "删除回收站实例失败")
		return
	}
}

func ReturnInstances(c *gin.Context) {
	var data models.LabVirtualMachine
	// 获取请求参数
	instanceID := c.Param("instance_id")

	_, err := PostDeleteInstances(instanceID)
	if err != nil {
		log.Error("机器删除失败", err)
	}
	log.Info("机器删除成功")

	_, err = PostDeleteRecycleInstances(instanceID)
	if err != nil {
		log.Error("机器归还失败", err)
		app.Error(c, http.StatusBadRequest, err, "机器归还失败")
		return
	}
	data.VmLog = "系统自动消息: 审批通过。借用的机器时间已到。"
	data.Status = "1"
	_, err = data.Update(instanceID, data.VmLog, data.Status)

	if err != nil {
		log.Error("归还机器数据更新失败", err)
		app.Error(c, http.StatusBadRequest, err, "归还机器数据更新失败")
		return
	}
	log.Info("归还机器数据更新成功")

	log.Info("机器归还成功")
	app.OK(c, http.StatusOK, "机器归还成功")
	return
}

func ParmReturnInstances(c *gin.Context, id string) {
	var data models.LabVirtualMachine
	// 获取请求参数
	instanceID := id

	_, err := PostDeleteInstances(instanceID)
	if err != nil {
		log.Error("机器删除失败", err)
	}
	log.Info("机器删除成功")

	_, err = PostDeleteRecycleInstances(instanceID)
	if err != nil {
		log.Error("机器归还失败", err)
		app.Error(c, http.StatusBadRequest, err, "机器归还失败")
		return
	}

	data.VmLog = "系统自动消息: 审批通过。借用的机器时间已到。"
	data.Status = "1"
	_, err = data.Update(instanceID, data.VmLog, data.Status)

	if err != nil {
		log.Error("归还机器数据更新失败", err)
		app.Error(c, http.StatusBadRequest, err, "归还机器数据更新失败")
		return
	}
	log.Info("归还机器数据更新成功")

	log.Info("机器归还成功")
	app.OK(c, http.StatusOK, "机器归还成功")

	return
}
