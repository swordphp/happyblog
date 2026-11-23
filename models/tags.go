package models

import "time"
import . "happyblog/library"

type Tags struct{
    Id int `gorm:"primaryKey" json:"id"`
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

/**
 * 通过标签名称批量获取标签信息
 *
 * param: []string tagsName
 * return: []Tags
 * return: error
 */
func (model Tags) GetTagsInfo(tagsName []string)(res []Tags,err error) {
    if len(tagsName) == 0 {
        return
    }
    err = ConnInstance.Model(&model).Where("tagName in (?)",tagsName).Find(&res).Error
    if err != nil {
        Logf("get tags info err","%v",err)
    }
    return
}


/**
 * 通过标签ID删除标签,同时清理文章-标签关联
 *
 * param: int id
 * return: int64
 */
func (model Tags) RemoveTagById(id int) (affectRows int64) {
    db := ConnInstance.Delete(Tags{},"id = ?",id)
    if db.Error != nil {
        Logf("remove tag err","%v",db.Error)
        return 0
    }
    // 清理该标签在文章关联表里的记录
    reModel := new(RelationArticleTags)
    _ = ConnInstance.Delete(reModel,"tagId = ?",id).Error
    return db.RowsAffected
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
    tags.CreateTime = time.Now()
    err = ConnInstance.Model(&model).Create(&tags).Error
    if err != nil {
        Logf("create tag error ","%v",err)
    }
    return tags.Id,err
}