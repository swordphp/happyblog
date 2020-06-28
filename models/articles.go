package models

import (
	. "happyblog/library"
	"time"
)

type Article struct{
	Id int `gorm:"primary_key" json:"id"`
	Title string `gorm:"column:title" json:"title"`
	Content string `gorm:"column:content" json:"content"`
	PubStatus int8 `gorm:"column:pubStatus" json:"pubStatus"`
	CreateTime time.Time `gorm:"column:createTime"   json:"createTime"`
	UpdateTime time.Time `gorm:"column:updateTime"  json:"updateTime"`
	AuthorInfo User `gorm:"foreignkey:authorId" json:"authorInfo"`
	AuthorId int `gorm:"column:authorId" json:"-"`
	IndependPage int8 `gorm:"column:independPage" json:"independPage"`
}


func (Article) TableName () string{
	return "happyblog_tblArticle"
}

/**
 * 通过文章ID 获取文章的内容
 *
 * param: int id
 * return: ArticleRow
 * return: error
 */
func (Article) GetArticleInfo(id int) (row Article,err error) {
	ConnInstance.First(&row,id)
	return
}

/**
 * 返回总数据行数
 *
 * return: int
 * return: error
 */
func (model Article) GetArticlesTotal() (total int,err error) {
	ConnInstance.Model(&model).Count(&total)
	return
}

/**
 * 通过排序字段获取
 *
 * param: string order
 * param: int    page
 * return: []ArticleRow
 * return: error
 */
func (Article) GetArticlesList(order string ,page int) (rows []Article,err error) {
	allowOrder := map[string]int8 {
		"createTime ASC":1,
		"createTime DESC":1,
		"updateTime DESC":1,
		"updateTime ASC":1,
		"pubStatus ASC":1,
		"pubStatus DESC":1,
	}
	if allowOrder[order] != 1 {
		Logf("order string not allow","%v",order)
		return
	}

	limitStart := 0
	perPage := 20
	if page<= 1 {
		limitStart = 0
	} else {
		limitStart = (page-1) *perPage
	}
	ConnInstance.Preload("AuthorInfo").Order(order).Offset(limitStart).Limit(perPage).Find(&rows)
	return rows,err
}

/**
 * 创建一个文章，返回创建的ID
 *
 * param: ArticleRow row
 * return: int
 */
func (model Article) CreateArticleRow(row Article) (insertId int) {
	row.CreateTime =time.Now()
	row.UpdateTime = time.Now()
	err := ConnInstance.Model(&model).Create(&row).Error
	if err != nil {
		Logf("create article err","%v",err)
	}
	return row.Id
}

/**
 * 更新文章内容，返回受影响的行数
 *
 * param: ArticleRow row
 * return: int
 */
func (model Article) UpdateArticleRow(row Article) (affectRows int64) {
	row.UpdateTime = time.Now()
	err := ConnInstance.Model(&model).Updates(row).Error
	if err != nil {
		Logf("update article err","%v",row)
		return 0
	}
	return ConnInstance.RowsAffected
}