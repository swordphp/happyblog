package models

import (
    "fmt"
    . "happyblog/library"
)

//定义数据行
type RelationArticleTags struct{
    Id int `gorm:"column:id;primary_key"`
    TagId int `gorm:"column:tagId;index:tags_articleid"`
    ArticleId int `gorm:"column:articleId;index:tags_articleid"`
}



func (model RelationArticleTags) TableName() string {
    return "happylblog_tblTagRe"
}


/**
 * 添加或者删除文章与tags之间的关联
 *
 * param: []Tags tags
 * param: int    articleId
 * return: int
 */
func (model RelationArticleTags) Relations(tagIds []int,articleId int) (affectRows int) {
    if len(tagIds) >0 {
        for _,tagId := range tagIds{
            querySql := "REPLACE INTO " + model.TableName() + "(`tagId`,`articleId`) VALUES (%d,%d)"
            querySql = fmt.Sprintf(querySql,tagId,articleId)
            res,err := ConnInstance.DB().Exec(querySql)
            if err != nil {
                Logf("manage relation between article and tag err","%v",err)
            } else {
                rows,_ := res.RowsAffected()
                affectRows += int(rows)
            }
        }
    } else {
        return 0
    }
    err := ConnInstance.Delete(&model,"tagId not in (?) and articleId in (?)",tagIds,articleId).Error
    if err != nil {
        Logf("manage relation between article and tag err,remove err","%v",err)
    } else {
        affectRows -= int(ConnInstance.RowsAffected)
    }
    return affectRows
}





/**
 * 通过文章ID获取文章所属的所有TAG
 *
 * param: int articleId
 * return: []RelationArticleTags
 */
func (model RelationArticleTags) GetArticleTags(articleIds []int) (relations []RelationArticleTags) {
    err := ConnInstance.Model(&model).Where("articleId in (?)" ,articleIds).Find(&relations).Error
    if err != nil {
        Logf("get article relation error","%v",err)
    }
    return relations
}