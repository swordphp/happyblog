package models

import "time"
import . "happyblog/library"

type Tags struct{
    Id int `gorm:"primary_key" json:"id"`
    TagName string `gorm:"column:tagName;unique" json:"tagName"`
    CreateTime time.Time `gorm:"column:createTime" json:"-"`
}


func (Tags) TableName() string{
    return "happyblog_tblTag"
}

/**
 * 获取数据库中的所有标签
 *
 * return: []Tags
 */
func (model Tags) GetAllTags() (res []Tags,err error) {
    err = ConnInstance.Model(&model).Find(&res).Error
    if err != nil {
        Logf("get tags err","%v",err)
    }
    return res,err
}
/**
 * 通过标签ID 获取标签信息
 *
 * param: []int ids
 * return: []Tags
 * return: error
 */
func (model Tags) GetTagsByIds(ids []int)(res []Tags,err error) {
    err = ConnInstance.Model(&model).Where("id in (?)",ids).Find(&res).Error
    if err != nil {
        Logf("get tags err","%v",err)
    }
    return res,err
}

/**
 * 通过标签名称获取标签信息
 *
 * param: string tagName
 * return: Tags
 */
func (model Tags) GetTagInfoByName(tagName string) (res Tags,err error){
    err = ConnInstance.Model(&model).Where("tagName = ?",tagName).First(&res).Error
    if err != nil {
        Logf("get tag info err","%v",err)
    }
    return res,err
}






func (model Tags) GetTagsInfo(tagsName []string)(res []Tags) {
    return
}


func (model Tags) RemoveTagByName(tagName string) (affectRows int) {
    return
}

/**
 * 创建tag的方法
 * param: string tagName
 * return: int
 * return: error
 */
func (model Tags) AddTag(tagName string) (id int,err error) {
    tags := new(Tags)
    tags.TagName = tagName
    err = ConnInstance.Model(&model).Create(&tags).Error
    if err != nil {
        Logf("create tag error ","%v",err)
    }
    return tags.Id,err
}