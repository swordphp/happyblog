package models

import (
    _ "github.com/go-sql-driver/mysql"
    . "happyblog/library"
    "time"
)
//定义数据行
type User struct{
    Id int `gorm:"primary_key" json:"-"`
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
