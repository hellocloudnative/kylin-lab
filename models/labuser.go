package models

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
	"kylin-lab/global/orm"
	"kylin-lab/tools"
	"log"
	"strings"
)

func (LabUser) TableName() string {
	return "lab_user"
}

type UserName struct {
	Username string `gorm:"type:varchar(64)" json:"username"`
}

type PassWord struct {
	// 密码
	Password string `gorm:"type:varchar(128)" json:"password"`
}

func (u *UserName) CheckUserName() error {
	if len(u.Username) == 0 {
		return errors.New("用户名不能为空")
	}
	return nil
}

func (u *PassWord) CheckPassWord() error {
	if len(u.Password) == 0 {
		return errors.New("密码不能为空")
	}
	return nil
}

type LoginM struct {
	UserName
	PassWord
}

type LabUserView struct {
	LabUserId
	LoginM
}

type LabUserId struct {
	UserId int `gorm:"primary_key;AUTO_INCREMENT"  json:"userId"` // 编码
}

type LabUser struct {
	LabUserId
	LoginM
}

type LabUserPwd struct {
	OldPassword string `json:"oldPassword"`
	NewPassword string `json:"newPassword"`
}

type LabUserPage struct {
	LabUserId
	LoginM
}

func (e *LabUser) GetList() (LabUserView []LabUserView, err error) {

	table := orm.Eloquent.Table(e.TableName()).Select([]string{"lab_user.*"})
	if e.UserId != 0 {
		table = table.Where("user_id = ?", e.UserId)
	}

	if e.Username != "" {
		table = table.Where("username = ?", e.Username)
	}

	if e.Password != "" {
		table = table.Where("password = ?", e.Password)
	}

	if err = table.Find(&LabUserView).Error; err != nil {
		return
	}
	return
}

func (e *LabUser) GetUserInfo() (LabUserView LabUserView, err error) {

	table := orm.Eloquent.Table(e.TableName()).Select([]string{"lab_user.*"})
	if e.UserId != 0 {
		table = table.Where("user_id = ?", e.UserId)
	}

	if e.Username != "" {
		table = table.Where("username = ?", e.Username)
	}

	if e.Password != "" {
		table = table.Where("password = ?", e.Password)
	}

	if err = table.First(&LabUserView).Error; err != nil {
		return
	}
	return
}

func (e *LabUser) GetPage(pageSize int, pageIndex int) ([]LabUserPage, int, error) {
	var doc []LabUserPage
	table := orm.Eloquent.Select("lab_user.*").Table(e.TableName())

	if e.Username != "" {
		table = table.Where("username = ?", e.Username)
	}
	var count int

	if err := table.Offset((pageIndex - 1) * pageSize).Limit(pageSize).Find(&doc).Error; err != nil {
		return nil, 0, err
	}
	table.Where("prince_user.deleted_at IS NULL").Count(&count)
	return doc, count, nil
}

// 获取用户数据
func (e *LabUser) Get() (LabUserView LabUserView, err error) {

	table := orm.Eloquent.Table(e.TableName()).Select([]string{"lab_user.*"})
	if e.UserId != 0 {
		table = table.Where("user_id = ?", e.UserId)
	}

	if e.Username != "" {
		table = table.Where("username = ?", e.Username)
	}

	if e.Password != "" {
		table = table.Where("password = ?", e.Password)
	}

	if err = table.First(&LabUserView).Error; err != nil {
		return
	}
	LabUserView.Password = ""
	return
}

// 加密
func (e *LabUser) Encrypt() (err error) {
	if e.Password == "" {
		return
	}

	var hash []byte
	if hash, err = bcrypt.GenerateFromPassword([]byte(e.Password), bcrypt.DefaultCost); err != nil {
		return
	} else {
		e.Password = string(hash)
		return
	}
}

// 添加
func (e LabUser) Insert() (id int, err error) {
	if err = e.Encrypt(); err != nil {
		return
	}

	// check 用户名
	var count int
	orm.Eloquent.Table(e.TableName()).Where("username = ?", e.Username).Count(&count)
	if count > 0 {
		err = errors.New("账户已存在！")
		return
	}

	//添加数据
	if err = orm.Eloquent.Table(e.TableName()).Create(&e).Error; err != nil {
		return
	}
	id = e.UserId
	return
}

// 修改
func (e *LabUser) Update(id int) (update LabUser, err error) {
	if e.Password != "" {
		if err = e.Encrypt(); err != nil {
			return
		}
	}

	if err = orm.Eloquent.Table(e.TableName()).First(&update, id).Error; err != nil {
		return
	}

	//参数1:是要修改的数据
	//参数2:是修改的数据
	if err = orm.Eloquent.Table(e.TableName()).Model(&update).Updates(&e).Error; err != nil {
		return
	}
	return
}

func (e *LabUser) BatchDelete(id []int) (Result bool, err error) {
	if err = orm.Eloquent.Unscoped().Where("user_id in (?)", id).Delete(&LabUser{}).Error; err != nil {
		return
	}
	Result = true
	return
}

func (e *LabUser) SetPwd(pwd LabUserPwd) (Result bool, err error) {
	user, err := e.GetUserInfo()
	if err != nil {
		tools.HasError(err, "获取用户数据失败(代码202)", 500)
	}
	_, err = tools.CompareHashAndPassword(user.Password, pwd.OldPassword)
	if err != nil {
		if strings.Contains(err.Error(), "hashedPassword is not the hash of the given password") {
			tools.HasError(err, "密码错误(代码202)", 500)
		}
		log.Print(err)
		return
	}
	e.Password = pwd.NewPassword
	_, err = e.Update(e.UserId)
	tools.HasError(err, "更新密码失败(代码202)", 500)
	return
}
