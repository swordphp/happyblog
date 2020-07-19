package controller

import (
    "fmt"
    "github.com/gin-gonic/gin"
    lib "happyblog/library"
    . "happyblog/models"
    "net/http"
    "strconv"
)

type ArticleApi struct{}

type ResArticle struct{
    ArticleInfo Article
    Tags []Tags `json:"tags"`
    Album *Album `json:"album"`
}

/**
 * 获取文章信息
 *
 * param: *gin.Context c
 */
func (api ArticleApi) ArticleInfo(c *gin.Context) {

    articleId,_ := strconv.Atoi(c.Param("id"))
    response := new(ResArticle)
    modelArticle := new(Article)
    articleInfo,err := modelArticle.GetArticleInfo(articleId)
    response.ArticleInfo = articleInfo


    ids := []int{articleId}
    response.Tags = api.getTagsByArticleId(ids)
    response.Album,err = api.getAlbumInfoByArticleId(articleId)
    render := new(lib.Render).Init()
    if err == nil {
        render.SetData(response)
    } else {
        render.SetErr(500,"get article info from db error  ")
    }
    c.JSON(http.StatusOK,render)
}


/**
 * 获取文章列表
 *
 * param: *gin.Context c
 */
func (ArticleApi) ArticleList(c *gin.Context) {
    page ,_ := strconv.Atoi(c.DefaultQuery("id","1"))
    modelArticle := new(Article)
    render := new(lib.Render).Init()
    articleList ,err := modelArticle.GetArticleListView(page)
    if err == nil {
        render.SetData(articleList)
    } else {
        render.SetErr(500,"get article list from db  error")
    }
    fmt.Println(articleList)
    c.JSON(http.StatusOK,render)
}

/**
 * 通过文章ID 获取文章的tag信息
 */

func (ArticleApi) getTagsByArticleId(ids []int) (tags []Tags) {
    reTags := new(RelationArticleTags)
    articleTags := reTags.GetArticleTags(ids)
    fmt.Println(articleTags)
    if articleTags != nil {
        ids := make([]int,len(articleTags))
        modelTags := new(Tags)
        for k, info := range articleTags {
            ids[k] = info.TagId
        }
        tags,_ := modelTags.GetTagsByIds(ids)
        return tags
    }
    return nil
}

/**
 * 通过文章ID 获取文章所属专辑信息
 *
 * param: int id
 * return: Album
 * return: error
 */
func (ArticleApi) getAlbumInfoByArticleId(id int) (albumInfo *Album,err error) {
    defaultAlbum := Album{
        Id:           0,
        AlbumName:    "无专辑信息",
        IsPublic:     0,
        AuthorId:     0,
    }
    reAlbum := new(RelationArticleAlbums)
    relationId := reAlbum.GetBelongAlbumByArticleId(id)
    if relationId != 0 {
        modelAlbum := new(Album)
        tmpInfo,_ := modelAlbum.GetAlbumInfoView(relationId)
        if tmpInfo.Id == 0 {
            return &defaultAlbum,nil
        }
        return &tmpInfo,nil
    }
    return &defaultAlbum,nil
}