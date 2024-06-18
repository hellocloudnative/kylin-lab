package models

import (
	"errors"
	"fmt"
	"kylin-lab/global/orm"
	"strings"
	"time"
)

func (LabVirtualMachine) TableName() string {
	return "lab_virtualMachine"
}

type LabVirtualMachinePage struct {
	LabVirtualMachineId
	LabVirtualMachineInfo
}

type LabVirtualMachineView struct {
	UserName
	LabVirtualMachineId
	LabVirtualMachineInfo
}

type LabVirtualMachine struct {
	LabVirtualMachineId
	LabVirtualMachineInfo
}

type LabVirtualMachineId struct {
	VmId int `gorm:"primary_key;AUTO_INCREMENT"  json:"vmId"` // 编码
}

type LabVirtualMachineInfo struct {
	UserID          int       `gorm:"index" json:"userId"`
	UUID            string    `gorm:"type:varchar(255)" json:"uuid"`
	CPUArchitecture string    `gorm:"type:varchar(64)" json:"cpuArchitecture"`
	OSImage         string    `gorm:"type:varchar(255)"  json:"osImage"`
	Flavors         string    `gorm:"type:varchar(255)"json:"flavors"`
	VNCAddress      string    `gorm:"type:varchar(255)" json:"vncAddress"`
	IPAddress       string    `gorm:"type:varchar(255)" json:"ipAddress"`
	NetworkName     string    `gorm:"type:varchar(255)" json:"networkName"`
	Duration        string    `gorm:"type:int(2)" json:"duration"`
	TimeOfuse       string    `gorm:"type:varchar(255)" json:"timeOfuse"`
	ApplyStatus     string    `gorm:"type:int(1)" json:"applyStatus"`
	ApplyTime       time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"applyTime"` // 默认值
	Status          string    `gorm:"type:int(1)" json:"status"`
	VmLog           string    `gorm:"type:varchar(255)" json:"vmlog"`
	CreatedAt       time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"createdAt"` // 默认值设置为当前时间戳
	UpdatedAt       time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updatedAt"` // 默认值设置为当前时间戳
}

func (e *LabVirtualMachine) GetPage(pageSize int, pageIndex int) ([]LabVirtualMachinePage, int, error) {
	var doc []LabVirtualMachinePage
	table := orm.Eloquent.Select("lab_virtualMachine.*").Table(e.TableName())

	if e.VmId != 0 {
		table = table.Where("vm_id = ?", e.VmId)
	}
	if e.UserID != 0 {
		table = table.Where("user_id = ?", e.UserID)
	}

	if e.CPUArchitecture != "" {
		table = table.Where("cpu_architecture = ?", e.CPUArchitecture)
	}
	if e.Status != "" {
		table = table.Where("status = ?", e.Status)
	}

	if e.ApplyStatus != "" {
		table = table.Where("apply_status = ?", e.ApplyStatus)
	}

	var count int

	if err := table.Offset((pageIndex - 1) * pageSize).Limit(pageSize).Find(&doc).Error; err != nil {
		return nil, 0, err
	}
	table.Count(&count)
	return doc, count, nil
}

func (e *LabVirtualMachine) GetVirtualMachineInfo() (LabVirtualMachineView LabVirtualMachineView, err error) {

	table := orm.Eloquent.Table(e.TableName()).Select([]string{"lab_virtualMachine.*"})
	if e.VmId != 0 {
		table = table.Where("vm_id = ?", e.VmId)
	}

	if err = table.First(&LabVirtualMachineView).Error; err != nil {
		return
	}
	return
}

func (e *LabVirtualMachine) Create() (LabVirtualMachine, error) {
	var doc LabVirtualMachine
	result := orm.Eloquent.Table(e.TableName()).Create(&e)
	if result.Error != nil {
		err := result.Error
		return doc, err
	}
	doc = *e
	return doc, nil
}

func (e *LabVirtualMachine) Update(uuid, vmlog, status string) (update LabVirtualMachine, err error) {
	// 检查id是否有效
	if uuid == "" {
		return update, errors.New("invalid uuid")
	}

	// 根据id获取要更新的记录
	if err = orm.Eloquent.Table(e.TableName()).Where("uuid = ?", uuid).First(&update).Error; err != nil {
		return update, err // 返回错误
	}

	// 这里假设e结构体中包含了你想要更新的字段，并且这些字段已经被修改
	// 例如，e.FieldName = newValue
	e.VmLog = vmlog
	e.Status = status

	// 更新记录
	if err = orm.Eloquent.Table(e.TableName()).Model(&update).Updates(e).Error; err != nil {
		return update, err // 返回错误
	}

	// 返回更新后的记录
	return update, nil
}

func (e *LabVirtualMachine) GetUserVirtualMachineStatus0Count(userId int) (count int, err error) {
	table := orm.Eloquent.Table(e.TableName()).Select([]string{"lab_virtualMachine.*"})
	if userId != 0 {
		table = table.Where("user_id = ? ", userId)
	}
	if err = table.Where("status = ?", "0").Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (e *LabVirtualMachine) CheckAllVirtualMachineStatusAndUpdate(status string) ([]LabVirtualMachine, error) {
	var machines []LabVirtualMachine
	table := orm.Eloquent.Table(e.TableName()).Select([]string{"lab_virtualMachine.*"})
	if status != "" {
		table = table.Where("status = ? ", status)
	}
	err := table.Scan(&machines).Error
	if err != nil {
		return nil, err
	}
	for _, machine := range machines {
		if machine.Status == "0" {
			timestamps := strings.Split(machine.TimeOfuse, "-")
			endTimestampParts := timestamps[len(timestamps)-3:]
			end_timestamp := strings.Join(endTimestampParts, "-")
			// 解析TimeOfuse字段，该字段格式为"开始时间-结束时间"
			endTime, err := time.ParseInLocation("2006-01-02 15:04:05", end_timestamp, time.FixedZone("CST", 8*60*60))
			if err != nil {
				// 处理错误
				return nil, err
			}
			// 检查结束时间是否小于当前时间
			if endTime.Before(time.Now()) {
				// 更新状态和vmlog
				machine.Status = "1" // 假设1代表到期
				machine.VmLog = fmt.Sprintf("系统自动消息: 审批通过。借用的机器时间已到。")
				err = orm.Eloquent.Table(e.TableName()).Where("uuid = ?", machine.UUID).Updates(map[string]interface{}{
					"status":     machine.Status,
					"vm_log":     machine.VmLog,
					"updated_at": time.Now(),
				}).Error
				if err != nil {
					return nil, err
				}
			} else {
				return nil, nil
			}
		}
	}
	return machines, nil
}
