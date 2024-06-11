package models

import (
	"errors"
	"kylin-lab/global/orm"
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
	CPUArchitecture string    `gorm:"type:varchar(64)" json:"cpuArchitecture"`
	OSType          string    `gorm:"type:varchar(64)" json:"osType"`
	OSImage         string    `gorm:"type:varchar(255)"  json:"osImage"`
	MachineSpec     string    `gorm:"type:varchar(255)" json:"machineSpec"`
	IPAddress       string    `gorm:"type:varchar(255)" json:"ipAddress"`
	Duration        string    `gorm:"type:int(2)" json:"duration"`
	Status          string    `gorm:"type:int(1)" json:"status"`
	VmLog           string    `gorm:"type:varchar(255)" json:"vmlog"`
	VNCAddress      string    `gorm:"type:varchar(255)" json:"vncAddress"`
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
	if e.OSType != "" {
		table = table.Where("os_type = ?", e.OSType)
	}
	if e.CPUArchitecture != "" {
		table = table.Where("cpu_architecture = ?", e.CPUArchitecture)
	}
	if e.Status != "" {
		table = table.Where("status = ?", e.Status)
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

func (e *LabVirtualMachine) Update(id int, duration, status string) (update LabVirtualMachine, err error) {
	// 检查id是否有效
	if id <= 0 {
		return update, errors.New("invalid ID")
	}

	// 根据id获取要更新的记录
	if err = orm.Eloquent.Table(e.TableName()).Where("vm_id = ?", id).First(&update).Error; err != nil {
		return update, err // 返回错误
	}

	// 这里假设e结构体中包含了你想要更新的字段，并且这些字段已经被修改
	// 例如，e.FieldName = newValue
	e.Duration = duration
	e.Status = status

	// 更新记录
	if err = orm.Eloquent.Table(e.TableName()).Model(&update).Updates(e).Error; err != nil {
		return update, err // 返回错误
	}

	// 返回更新后的记录
	return update, nil
}
