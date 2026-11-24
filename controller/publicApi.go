package controller

import (
    "strconv"

    "github.com/gin-gonic/gin"
    . "happyblog/models"
)

// PublicApi 提供面向读者(无需登录)的只读接口,仅暴露已发布/公开内容
type PublicApi struct{}

/**
 * SiteInfo 返回站点基本信息,供读者端 SPA 渲染页头与标题。
 */
func (PublicApi) SiteInfo(c *gin.Context) {
    respOK(c, gin.H{
        "sitename": webConf["sitename"],
        "sitedesc": webConf["sitedesc"],
        "siteurl":  webConf["siteurl"],
    })
}

/**
 * Tags 返回全部标签。
 */
func (PublicApi) Tags(c *gin.Context) {
    tags, err := new(Tags).GetAllTags()
    if err != nil {
        respFail(c, errCodeServer, "get tags error")
        return
    }
    respOK(c, tags)
}

/**
 * TagArticles 返回某标签下的已发布文章。
 */
func (PublicApi) TagArticles(c *gin.Context) {
    tagId, _ := strconv.Atoi(c.Param("id"))
    if tagId == 0 {
        respFail(c, errCodeParam, "tag id is required")
        return
    }
    articleIds := new(RelationArticleTags).GetArticleIdsByTagId(tagId)
    articles, _ := new(Article).GetPublishedArticlesByIds(articleIds)
    respOK(c, articles)
}

/**
 * Albums 返回公开专辑列表。
 */
func (PublicApi) Albums(c *gin.Context) {
    page, _ := strconv.Atoi(c.Param("page"))
    if page < 1 {
        page = 1
    }
    albums, err := new(Album).GetPublicAlbumList(page)
    if err != nil {
        respFail(c, errCodeServer, "get albums error")
        return
    }
    respOK(c, albums)
}

/**
 * AlbumInfo 返回单个公开专辑及其下已发布文章。
 */
func (PublicApi) AlbumInfo(c *gin.Context) {
    albumId, _ := strconv.Atoi(c.Param("id"))
    album, err := new(Album).GetAlbumInfoView(albumId)
    if err != nil || album.Id == 0 {
        respFail(c, errCodeNotFound, "album not found")
        return
    }
    articleIds := new(RelationArticleAlbums).GetArticleIdsByAlbumId(albumId)
    articles, _ := new(Article).GetPublishedArticlesByIds(articleIds)
    respOK(c, gin.H{"album": album, "articles": articles})
}
