package models


/**
 * 用于管理和获取专辑相关信息
 */

import (
	"github.com/jinzhu/gorm"
	. "happyblog/library"
	"time"
)

type Album struct{
	Id int `gorm:"primary_key" json:"id"`
	AlbumName string `gorm:"column:albumName" json:"albumName"`
	IsPublic int8 `gorm:"column:isPublic" json:"isPublic"`
	AuthorId int `gorm:"column:authorId" json:"authorInfo"`
	CreateTime time.Time `gorm:"column:createTime" json:"createTime"`
	ArticleTotal int `gorm:"column:articleTotal" json:"articleTotal"`
	AuthorInfo User `gorm:"foreignkey:authorId" json:"authorInfo"`
}


func (Album) TableName() string {
	return "happyblog_tblAlbum"
}

/**
 * 获取专辑列表
 *
 * param: string order
 * param: int    page
 * return: []ArticleRow
 * return: error
 */
func (Album) GetAlbumList(page int) (rows []Album,err error) {
	limitStart , limitEnd  := 0,0
	if page<= 1 {
		limitStart = 0
		limitEnd = 10
	} else {
		limitStart = (page-1) *10
		limitEnd = page*10
	}
	err = ConnInstance.Preload("AuthorInfo").Offset(limitStart).Limit(limitEnd).Find(&rows).Error
	if err != nil {
		Logf("get album list err","%v",page)
	}
	return
}

/**
 * 获取一个专辑的相关信息
 * #todo 还需要加入文章的专辑引用数量
 */

func (mo Album) GetAlbumInfo(id int) (row Album,err error) {
	err = ConnInstance.Preload("AuthorInfo").Find(&row,id).Error
	if err != nil {
		Logf("get album info error","%v",nil)
	}
	return row,err
}
/**
 * 创建专辑的方法
 *
 * param: string albumName
 * param: int8   isPublic
 * param: int    authorId
 * return: int64
 * return: error
 */
func (mo Album) CreateAlbum (album Album) (id int,err error) {
	album.CreateTime = time.Now()
	err = ConnInstance.Model(&mo).Create(&album).Error
	if err != nil {
		Logf("create album err","%v",album)
	}
	return album.Id,err
}
/**
 * 更新专辑的方法
 *
 * param: string albumName
 * param: int8   isPublic
 * param: int    id
 * return: bool
 */
func (mo Album) UpdateAlbum (album Album) (res bool) {
	err := ConnInstance.Model(&mo).Updates(album).Error
	if err != nil {
		Logf("update album info error","%v",album)
	}
	if ConnInstance.RowsAffected > 0 {
		return true
	}
	return false
}

/**
 * 更新专辑的方法
 *
 * param: string albumName
 * param: int8   isPublic
 * param: int    id
 * return: bool
 */
func (mo Album) UpdateAlbumArticleTotal (albumId int,totalIncr int) (res bool) {
	err := ConnInstance.Model(&mo).Where("id = ?",albumId).
		Update("articleTotal",gorm.Expr("articleTotal + ?",totalIncr)).Error
	if err != nil {
		Logf("incr article count error ","%v",nil)
	}
	return true
}