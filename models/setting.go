package models


/**
 * 用于管理配置信息
 */

import (
    . "happyblog/library"
)

type Setting struct{
    Id int `gorm:"primary_key" json:"id" form:"id"`
    ConfigGroup string `gorm:"column:configGroup;index:group_configKey" json:"group" form:"ConfigGroup"`
    ConfigName string `gorm:"column:configName;index:group_configKey" json:"configName" form:"ConfigName"`
    ConfigValue string `gorm:"column:configValue" json:"configValue" form:"ConfigValue"`
    ConfigType string `gorm:"column:configType" json:"type" form:"ConfigType"`
    ConfigOrder int `gorm:"column:configOrder" json:"order" form:"ConfigOrder"`
}


func (Setting) TableName() string {
    return "happyblog_settings"
}

/**
 * 获取所有配置信息,并按照组来组织配置信息
 *
 * return: []Config
 * return: error
 */
func (Setting) GetConfigs() (groupRows map[string][]Setting,err error) {
    rows := new([]Setting)
    err = ConnInstance.Where("configGroup <> configName").Find(&rows).Error
    if err != nil {
        Logf("get configs list err","","")
    }
    groupRows = make(map[string][]Setting)
    if rows != nil {

        for _,row := range *rows{
            groupRows[row.ConfigGroup] = append(groupRows[row.ConfigGroup],row)
        }
    }
    return groupRows,err
}

/**
 * 获取所有的组信息
 *
 * return: []string
 * return: error
 */
func (model Setting) GetGroups()(groups []string,err error){
    allRows,err := ConnInstance.Model(&model).Select("group").Group("configGroup").Rows()
    if allRows != nil {
        defer allRows.Close()
        for allRows.Next() {
            groupString := ""
            _ = ConnInstance.ScanRows(allRows,groupString)
            groups = append(groups,groupString)
        }
    }
    return groups,err
}

/**
 * 更新设置的方法
 *
 * param: Setting row
 * return: int64
 * return: error
 */
func (model Setting) UpdateConfigInfo(row Setting) (affectRows int64 ,err error){
    err = ConnInstance.Model(&model).Save(row).Error
    affectRows = ConnInstance.RowsAffected
    return affectRows,err
}

/**
 * 通过id获取单个配置的信息
 *
 * param: int id
 * return: Config
 * return: error
 */
func (model Setting) GetConfigInfo(id int) (row Setting,err error) {
    err = ConnInstance.
        Model(&model).
        Where("id = ?", id).
        Find(&row).
        Error
    if err != nil {
        Logf("get config info error","%v",nil)
    }
    return
}

/**
 获取一组配置文件信息
 */

func (model Setting) GetConfigsByGroup(group string) (rows []Setting,err error) {
    err = ConnInstance.
        Model(&model).
        Where("configGroup = ?",group).
        Order("configOrder ASC").
        Find(&rows).Error
    if err != nil {
        Logf("get config info error","%v",nil)
    }
    return rows,err
}

/**
 * 创建配置项的方法
 *
 * param: string albumName
 * param: int8   isPublic
 * param: int    authorId
 * return: int64
 * return: error
 */
func (model Setting) CreateConfig (row Setting) (id int,err error) {
    err = ConnInstance.Model(&model).Create(&row).Error
    if err != nil {
        Logf("create config err","%v",row)
    }
    return row.Id,err
}

/**
 * 移除设置的一行
 * param: int id
 * return: int64
 */
func (model Setting) RemoveRow(id int)(affectRows int64) {
    err := ConnInstance.Delete(Setting{},"id = ?",id).Error
    if err != nil {
        Logf("reomve row error","%v",err)
    }
    return ConnInstance.RowsAffected
}