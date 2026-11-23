package models

import (
    . "happyblog/library"
    "gorm.io/gorm/clause"
)

//定义数据行
type RelationArticleTags struct{
    Id int `gorm:"column:id;primaryKey"`
    TagId int `gorm:"column:tagId;uniqueIndex:tags_articleid"`
    ArticleId int `gorm:"column:articleId;uniqueIndex:tags_articleid"`
}



func (model RelationArticleTags) TableName() string {
    return "happylblog_tblTagRe"
}


/**
 * 添加或者删除文章与tags之间的关联
 * 跨数据库实现:用 upsert(冲突即忽略)新增关联,再删除不在列表中的旧关联
 *
 * param: []int tagIds
 * param: int   articleId
 * return: int
 */
func (model RelationArticleTags) Relations(tagIds []int,articleId int) (affectRows int) {
    if len(tagIds) == 0 {
        return 0
    }
    for _,tagId := range tagIds {
        row := RelationArticleTags{TagId: tagId, ArticleId: articleId}
        db := ConnInstance.Clauses(clause.OnConflict{DoNothing: true}).Create(&row)
        if db.Error != nil {
            Logf("manage relation between article and tag err","%v",db.Error)
        } else {
            affectRows += int(db.RowsAffected)
        }
    }
    db := ConnInstance.Where("tagId not in (?) and articleId = ?",tagIds,articleId).Delete(&RelationArticleTags{})
    if db.Error != nil {
        Logf("manage relation between article and tag err,remove err","%v",db.Error)
    } else {
        affectRows -= int(db.RowsAffected)
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

/**
 * 通过标签ID获取该标签下的文章ID列表
 *
 * param: int tagId
 * return: []int
 */
func (model RelationArticleTags) GetArticleIdsByTagId(tagId int) (articleIds []int) {
    var relations []RelationArticleTags
    err := ConnInstance.Model(&model).Where("tagId = ?",tagId).Find(&relations).Error
    if err != nil {
        Logf("get article ids by tag error","%v",err)
        return
    }
    for _,r := range relations {
        articleIds = append(articleIds,r.ArticleId)
    }
    return
}