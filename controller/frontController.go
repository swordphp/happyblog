package controller

import (
    "github.com/gin-gonic/gin"
    lib "happyblog/library"
    . "happyblog/models"
    "net/http"
    "strconv"
)

// FrontController 面向读者的服务端渲染页面(文章列表/详情/标签/专辑)
type FrontController struct{
    controller
}

var CFront FrontController

/**
 * 读者端首页:已发布文章列表
 *
 * param: *gin.Context c
 */
func (ctrl FrontController) Index(c *gin.Context) {
    page,_ := strconv.Atoi(c.DefaultQuery("page","1"))
    if page < 1 {
        page = 1
    }
    modelArticle := new(Article)
    articles,_ := modelArticle.GetArticleListView(page)
    total,_ := modelArticle.GetPublishedTotal()
    totalPage := (total + 9) / 10

    c.HTML(http.StatusOK, "front_index", gin.H{
        "sitename":  webConf["sitename"],
        "sitedesc":  webConf["sitedesc"],
        "articles":  articles,
        "page":      page,
        "prevPage":  page - 1,
        "nextPage":  page + 1,
        "hasPrev":   page > 1,
        "hasNext":   page < totalPage,
        "totalPage": totalPage,
    })
}

/**
 * 读者端文章详情:Markdown 渲染为 HTML
 *
 * param: *gin.Context c
 */
func (ctrl FrontController) Article(c *gin.Context) {
    id,_ := strconv.Atoi(c.Param("id"))
    modelArticle := new(Article)
    article,err := modelArticle.GetPublishedArticleInfo(id)
    if err != nil || article.Id == 0 {
        c.HTML(http.StatusNotFound, "front_404", gin.H{
            "sitename": webConf["sitename"],
        })
        return
    }

    tags := new(ArticleApi).getTagsByArticleId([]int{id})

    c.HTML(http.StatusOK, "front_article", gin.H{
        "sitename":   webConf["sitename"],
        "article":    article,
        "contentHTML": lib.RenderMarkdown(article.Content),
        "tags":       tags,
    })
}

/**
 * 读者端:某标签下的已发布文章
 *
 * param: *gin.Context c
 */
func (ctrl FrontController) Tag(c *gin.Context) {
    id,_ := strconv.Atoi(c.Param("id"))
    reTags := new(RelationArticleTags)
    articleIds := reTags.GetArticleIdsByTagId(id)
    modelArticle := new(Article)
    articles,_ := modelArticle.GetPublishedArticlesByIds(articleIds)

    tagName := ""
    modelTags := new(Tags)
    if infos,_ := modelTags.GetTagsByIds([]int{id}); len(infos) > 0 {
        tagName = infos[0].TagName
    }

    c.HTML(http.StatusOK, "front_list", gin.H{
        "sitename": webConf["sitename"],
        "listTitle": "标签: " + tagName,
        "articles": articles,
    })
}

/**
 * 读者端:某公开专辑及其下已发布文章
 *
 * param: *gin.Context c
 */
func (ctrl FrontController) Album(c *gin.Context) {
    id,_ := strconv.Atoi(c.Param("id"))
    modelAlbum := new(Album)
    album,err := modelAlbum.GetAlbumInfoView(id)
    if err != nil || album.Id == 0 {
        c.HTML(http.StatusNotFound, "front_404", gin.H{
            "sitename": webConf["sitename"],
        })
        return
    }
    reAlbum := new(RelationArticleAlbums)
    articleIds := reAlbum.GetArticleIdsByAlbumId(id)
    modelArticle := new(Article)
    articles,_ := modelArticle.GetPublishedArticlesByIds(articleIds)

    c.HTML(http.StatusOK, "front_list", gin.H{
        "sitename":  webConf["sitename"],
        "listTitle": "专辑: " + album.AlbumName,
        "articles":  articles,
    })
}
