package models


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
func (model RelationArticleTags) AddRelations(tags []Tags,articleId int) (affectRows int) {
    return
}