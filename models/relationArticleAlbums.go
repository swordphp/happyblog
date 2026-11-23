package models

import (
    "errors"
    . "happyblog/library"
    "gorm.io/gorm"
)

//定义数据行
type RelationArticleAlbums struct{
    Id int `gorm:"column:id;primaryKey"`
    AlbumId int `gorm:"column:albumId;index:ablum_article"`
    ArticleId int `gorm:"column:articleId;index:ablum_article"`
}


func (model RelationArticleAlbums) TableName() string {
    return "happyblog_tblAlbumRe"
}

/**
 * 获取文章所属于的专辑ID
 *
 * param: int ArticleId
 * return: int
 */
func (model RelationArticleAlbums) GetBelongAlbumByArticleId (ArticleId int) (AlbumId int) {
    albumObj := new(RelationArticleAlbums)
    err := ConnInstance.Model(&model).Where("articleId = ?", ArticleId).First(&albumObj).Error
    if errors.Is(err, gorm.ErrRecordNotFound) {
        return 0
    }
    if err != nil {
        Logf(" get belong album err","%v",err)
    }
    return albumObj.AlbumId
}

/**
 * 更新文章所属的专辑
 *
 * param: int ArticleId
 * param: int AlbumId
 * return: int
 */
func (model RelationArticleAlbums) UpdateRowByArticleId (articleId int ,albumId int) (affectRows int64) {
    albumObj :=RelationArticleAlbums{
        ArticleId:articleId,
        AlbumId:albumId,
    }
    db := ConnInstance.Model(&model).Where("articleId = ?",articleId).Assign(&albumObj).FirstOrCreate(&albumObj)
    if db.Error != nil {
        Logf("update relation error!","%v",db.Error)
    }
    return db.RowsAffected
}

/**
 * 通过文章ID 移除行
 *
 * param: int ArticleId
 * return: int
 */
func (model RelationArticleAlbums) RemoveRowByArticleId (articleId int) (affectRows int64) {
    db := ConnInstance.Delete(&model,"articleId = ?",articleId)
    if db.Error != nil {
        Logf("remove relation err","%v",db.Error)
    }
    return db.RowsAffected
}

/**
 * 获取专辑的引用数量，通过文章id。
 * 目前没有限制，也没有做缓存，如果专辑数量过多，此处可能成为性能瓶颈#todo
 * param: []int ArticleIds
 * return: map[int]int64
 */
func (model RelationArticleAlbums) GetQuoteTotalByAlbumIds(AlbumIds []int) (totalCounts map[int]int64){
    totalCounts = make(map[int]int64)
    if len(AlbumIds) == 0 {
        return
    }
    type countRow struct{
        AlbumId int
        Total int64
    }
    var rows []countRow
    err := ConnInstance.Model(&model).
        Select("albumId, count(*) as total").
        Where("albumId in (?)",AlbumIds).
        Group("albumId").
        Scan(&rows).Error
    if err != nil {
        Logf("get album quote total err","%v",err)
        return
    }
    for _,r := range rows {
        totalCounts[r.AlbumId] = r.Total
    }
    return
}

/**
 * 获取某专辑下已发布文章的ID列表
 *
 * param: int albumId
 * return: []int
 */
func (model RelationArticleAlbums) GetArticleIdsByAlbumId(albumId int) (articleIds []int) {
    var relations []RelationArticleAlbums
    err := ConnInstance.Model(&model).Where("albumId = ?",albumId).Find(&relations).Error
    if err != nil {
        Logf("get article ids by album err","%v",err)
        return
    }
    for _,r := range relations {
        articleIds = append(articleIds,r.ArticleId)
    }
    return
}
