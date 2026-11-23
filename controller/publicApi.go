package controller

import (
    "github.com/gin-gonic/gin"
    lib "happyblog/library"
    . "happyblog/models"
    "net/http"
    "strconv"
)

// PublicApi 提供面向读者(无需登录)的只读接口,仅暴露已发布/公开内容
type PublicApi struct{}

/**
 * 获取全部标签
 *
 * param: *gin.Context c
 */
func (PublicApi) Tags(c *gin.Context) {
    render := new(lib.Render).Init()
    modelTags := new(Tags)
    tags,err := modelTags.GetAllTags()
    if err != nil {
        render.SetErr(500,"get tags error")
        c.JSON(http.StatusOK,render)
        return
    }
    render.SetData(tags)
    c.JSON(http.StatusOK,render)
}

/**
 * 获取某标签下的已发布文章
 *
 * param: *gin.Context c
 */
func (PublicApi) TagArticles(c *gin.Context) {
    render := new(lib.Render).Init()
    tagId,_ := strconv.Atoi(c.Param("id"))
    if tagId == 0 {
        render.SetErr(10,"tag id is required")
        c.JSON(http.StatusOK,render)
        return
    }
    reTags := new(RelationArticleTags)
    articleIds := reTags.GetArticleIdsByTagId(tagId)
    modelArticle := new(Article)
    articles,_ := modelArticle.GetPublishedArticlesByIds(articleIds)
    render.SetData(articles)
    c.JSON(http.StatusOK,render)
}

/**
 * 获取公开专辑列表
 *
 * param: *gin.Context c
 */
func (PublicApi) Albums(c *gin.Context) {
    render := new(lib.Render).Init()
    page,_ := strconv.Atoi(c.Param("page"))
    if page < 1 {
        page = 1
    }
    modelAlbum := new(Album)
    albums,err := modelAlbum.GetPublicAlbumList(page)
    if err != nil {
        render.SetErr(500,"get albums error")
        c.JSON(http.StatusOK,render)
        return
    }
    render.SetData(albums)
    c.JSON(http.StatusOK,render)
}

/**
 * 获取单个公开专辑及其下已发布文章
 *
 * param: *gin.Context c
 */
func (PublicApi) AlbumInfo(c *gin.Context) {
    render := new(lib.Render).Init()
    albumId,_ := strconv.Atoi(c.Param("id"))
    modelAlbum := new(Album)
    album,err := modelAlbum.GetAlbumInfoView(albumId)
    if err != nil || album.Id == 0 {
        render.SetErr(404,"album not found")
        c.JSON(http.StatusOK,render)
        return
    }
    reAlbum := new(RelationArticleAlbums)
    articleIds := reAlbum.GetArticleIdsByAlbumId(albumId)
    modelArticle := new(Article)
    articles,_ := modelArticle.GetPublishedArticlesByIds(articleIds)

    render.SetData(gin.H{
        "album":    album,
        "articles": articles,
    })
    c.JSON(http.StatusOK,render)
}
