package models

import (
    . "happyblog/library"
    "time"
)
//定义数据行
type User struct{
    Id int `gorm:"primaryKey" json:"-"`
    AccountEmail string `gorm:"column:accountEmail" json:"accountEmail"`
    AccountPassword string `gorm:"column:accountPassword" json:"-"`
    NickName string `gorm:"column:nickName" json:"nickName"`
    CreateTime time.Time `gorm:"column:createTime" json:"createTime"`
    UpdateTime time.Time `gorm:"column:updateTime" json:"-"`
    LastLogin time.Time `gorm:"column:lastLogin" json:"-"`
    HeadImageUri string `gorm:"column:headImageUri" json:"headImageUri"`
    EmailVerify int8 `gorm:"column:emailVerify" json:"-"`
}

func (User) TableName() string {
    return "happyblog_tblUser"
}


/**
 * 通过邮件地址获取单个用户的信息
 *
 * param: string email
 * return: UserRow
 * return: error
 */
func (User) GetUserInfo(email string) (row User,err error) {
    err = ConnInstance.Where("accountEmail = ?",email).Find(&row).Error
    if err != nil {
        Logf("get user info err","%v",err)
    }
    return row ,err
}

/**
 * 统计用户总数,用于首次初始化时判断是否需要创建默认管理员
 *
 * return: int64
 */
func (model User) Count() (total int64) {
    ConnInstance.Model(&model).Count(&total)
    return
}

/**
 * 创建用户
 *
 * param: User row
 * return: int
 * return: error
 */
func (model User) CreateUser(row User) (id int,err error) {
    now := time.Now()
    row.CreateTime = now
    row.UpdateTime = now
    err = ConnInstance.Model(&model).Create(&row).Error
    if err != nil {
        Logf("create user err","%v",err)
    }
    return row.Id,err
}

/**
 * 获取全部用户列表
 *
 * return: []User
 * return: error
 */
func (model User) GetUserList() (rows []User, err error) {
    err = ConnInstance.Model(&model).Order("id ASC").Find(&rows).Error
    if err != nil {
        Logf("get user list err", "%v", err)
    }
    return
}
