package models

import . "happyblog/library"

//定义数据行
type RelationArticleAlbums struct{
	Id int `gorm:"column:id;primary_key"`
	AlbumId int `gorm:"column:albumId;index:ablum_article"`
	ArticleId int `gorm:"column:articleId;index:ablum_article"`
}


func (mo RelationArticleAlbums) TableName() string {
	return "happyblog_tblAlbumRe"
}

/**
 * 获取文章所属于的专辑ID
 *
 * param: int ArticleId
 * return: int
 */
func (mo RelationArticleAlbums) GetBelongAlbumByArticleId (ArticleId int) (AlbumId int) {
	albumObj := new(RelationArticleAlbums)
	err := ConnInstance.Model(&mo).Where("articleId = ?", ArticleId).Find(&albumObj).Error
	if err != nil {
		Logf("get belong album err","%v",err)
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
func (mo RelationArticleAlbums) UpdateRowByArticleId (articleId int ,albumId int) (affectRows int64) {
	albumObj :=RelationArticleAlbums{
		ArticleId:articleId,
		AlbumId:albumId,
	}
	err := ConnInstance.Model(&mo).Where("articleId = ?",articleId).Assign(&albumObj).FirstOrCreate(&albumObj).Error
	if err != nil {
		Logf("update relation error!","%v",err)
	}
	return ConnInstance.RowsAffected
}

/**
 * 通过文章ID 移除行
 *
 * param: int ArticleId
 * return: int
 */
func (mo RelationArticleAlbums) RemoveRowByArticleId (articleId int) (affectRows int64) {
	err := ConnInstance.Model(&mo).Delete("articleId = ?",articleId).Error
	if err != nil {
		Logf("remove relation err","%v",err)
	}
	return ConnInstance.RowsAffected
}

/**
 * 获取专辑的引用数量，通过文章id。
 * 目前没有限制，也没有做缓存，如果专辑数量过多，此处可能成为性能瓶颈#todo
 * param: []int ArticleIds
 * return: map[int]int64
 */
func (mo RelationArticleAlbums) GetQuoteTotalByAlbumIds(AlbumIds []int) (totalCounts map[int]int64){
	return
}
