package models

import (
    . "happyblog/library"
    "time"
)

type Article struct{
    Id int `gorm:"primaryKey" json:"id"`
    Title string `gorm:"column:title" json:"title"`
    Content string `gorm:"column:content" json:"content,omitempty"`
    ContentType string `gorm:"column:contentType;size:16" json:"contentType"` // markdown(默认/空) | html(富文本)
    PubStatus int8 `gorm:"column:pubStatus" json:"pubStatus,omitempty"`
    CreateTime time.Time `gorm:"column:createTime"   json:"createTime"`
    UpdateTime time.Time `gorm:"column:updateTime"  json:"updateTime"`
    AuthorInfo User `gorm:"foreignKey:AuthorId;references:Id" json:"authorInfo"`
    AuthorId int `gorm:"column:authorId" json:"-"`
    IndependPage int8 `gorm:"column:independPage" json:"independPage"`
    Brief string `gorm:"column:brief" json:"brief"`
    Keywords string `gorm:"column:keywords" json:"keywords"`
    Headimage string `gorm:"column:headimage" json:"headimage"`
    Uri string `gorm:"column:uri" json:"uri,omitempty"`
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
    err = ConnInstance.Preload("AuthorInfo").First(&row,id).Error
    return
}

/**
 * 通过文章ID 获取已发布文章的内容(供公开API/读者端使用)
 * 仅返回 pubStatus=1 的文章,草稿与未发布内容不可见
 *
 * param: int id
 * return: Article
 * return: error
 */
func (Article) GetPublishedArticleInfo(id int) (row Article,err error) {
    err = ConnInstance.Preload("AuthorInfo").
        Where("id = ? and pubStatus = ?",id,1).
        First(&row).Error
    return
}

/**
 * 返回总数据行数
 *
 * return: int
 * return: error
 */
func (model Article) GetArticlesTotal() (total int,err error) {
    var cnt int64
    err = ConnInstance.Model(&model).Count(&cnt).Error
    total = int(cnt)
    return
}

/**
 * 返回已发布(非独立页)文章的总数,供读者端分页使用
 *
 * return: int
 * return: error
 */
func (model Article) GetPublishedTotal() (total int,err error) {
    var cnt int64
    err = ConnInstance.Model(&model).
        Where("pubStatus = ? and independPage <> ?",1,1).
        Count(&cnt).Error
    total = int(cnt)
    return
}

/**
 * 获取用于展示的文章列表
 *
 * param: int page
 * return: []Article
 * return: error
 */
func (Article) GetArticleListView(page int) (rows []Article,err error) {
    limitStart := 0
    perPage := 10
    if page<= 1 {
        limitStart = 0
    } else {
        limitStart = (page-1) *perPage
    }
    err = ConnInstance.Select("id,title,createTime,updateTime,authorId,brief,keywords,headimage,uri").
        Preload("AuthorInfo").
        Where("pubStatus = ? and independPage <> ?",1,1).
        Order("updateTime DESC").
        Offset(limitStart).Limit(perPage).
        Find(&rows).Error
    if err != nil{
        return nil ,err
    }
    return rows,err
}

/**
 * 通过一组文章ID获取已发布文章(供标签页/专辑页使用),按更新时间倒序
 *
 * param: []int ids
 * return: []Article
 * return: error
 */
func (Article) GetPublishedArticlesByIds(ids []int) (rows []Article,err error) {
    if len(ids) == 0 {
        return
    }
    err = ConnInstance.Select("id,title,createTime,updateTime,authorId,brief,keywords,headimage,uri").
        Preload("AuthorInfo").
        Where("id in (?) and pubStatus = ?",ids,1).
        Order("updateTime DESC").
        Find(&rows).Error
    if err != nil {
        Logf("get published articles by ids err","%v",err)
    }
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
    db := ConnInstance.Model(&model).Updates(row)
    if db.Error != nil {
        Logf("update article err","%v",row)
        return 0
    }
    return db.RowsAffected
}

/**
 * 通过文章id删除文章
 *
 * param: int id
 * return: int64
 */
func (model Article) RemoveRow(id int)(affectRows int64) {
    db := ConnInstance.Delete(Article{},"id = ?",id)
    if db.Error != nil {
        Logf("reomve row error","%v",db.Error)
    }
    return db.RowsAffected
}