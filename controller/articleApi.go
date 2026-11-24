package controller

import (
    "strconv"

    "github.com/gin-gonic/gin"
    . "happyblog/models"
)

type ArticleApi struct{}

type ResArticle struct {
    ArticleInfo Article `json:"articleInfo"`
    Tags        []Tags  `json:"tags"`
    Album       *Album  `json:"album"`
}

/**
 * ArticleInfo 返回单篇已发布文章的详情(含标签与所属专辑)。
 */
func (api ArticleApi) ArticleInfo(c *gin.Context) {
    articleId, _ := strconv.Atoi(c.Param("id"))
    //公开接口只返回已发布文章,草稿/未发布内容不可见
    articleInfo, err := new(Article).GetPublishedArticleInfo(articleId)
    if err != nil || articleInfo.Id == 0 {
        respFail(c, errCodeNotFound, "article not found")
        return
    }
    resp := ResArticle{ArticleInfo: articleInfo}
    resp.Tags = api.getTagsByArticleId([]int{articleId})
    resp.Album, _ = api.getAlbumInfoByArticleId(articleId)
    respOK(c, resp)
}

/**
 * ArticleList 返回已发布文章的分页列表。
 */
func (ArticleApi) ArticleList(c *gin.Context) {
    page, _ := strconv.Atoi(c.Param("page"))
    if page < 1 {
        page = 1
    }
    articleList, err := new(Article).GetArticleListView(page)
    if err != nil {
        respFail(c, errCodeServer, "get article list error")
        return
    }
    respOK(c, articleList)
}

// getTagsByArticleId 通过文章 ID 取其标签
func (ArticleApi) getTagsByArticleId(ids []int) []Tags {
    articleTags := new(RelationArticleTags).GetArticleTags(ids)
    if articleTags == nil {
        return nil
    }
    tagIds := make([]int, len(articleTags))
    for k, info := range articleTags {
        tagIds[k] = info.TagId
    }
    tags, _ := new(Tags).GetTagsByIds(tagIds)
    return tags
}

// getAlbumInfoByArticleId 通过文章 ID 取其所属专辑,无归属时返回占位专辑
func (ArticleApi) getAlbumInfoByArticleId(id int) (*Album, error) {
    defaultAlbum := Album{Id: 0, AlbumName: "无专辑信息"}
    relationId := new(RelationArticleAlbums).GetBelongAlbumByArticleId(id)
    if relationId == 0 {
        return &defaultAlbum, nil
    }
    info, _ := new(Album).GetAlbumInfoView(relationId)
    if info.Id == 0 {
        return &defaultAlbum, nil
    }
    return &info, nil
}
