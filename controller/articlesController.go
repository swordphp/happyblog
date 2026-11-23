package controller

import (
    "net/http"
    "strconv"

    "github.com/gin-gonic/gin"
    . "happyblog/models"
)

type ArticleController struct {
    controller
}

var CArticle ArticleController

/**
 * Articles 返回后台文章列表(JSON).
 *
 * param: *gin.Context c
 */
func (ctrl ArticleController) Articles(c *gin.Context) {
    ob := c.DefaultQuery("obstring", "ctime")
    obStatus := c.DefaultQuery("obstatus", "desc")
    obString := ctrl.makeOrderStr(ob, obStatus)
    nowPage, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
    if nowPage < 1 {
        nowPage = 1
    }
    ModelArticle := new(Article)
    articleList, err := ModelArticle.GetArticlesList(obString, nowPage)
    if err != nil {
        c.JSON(http.StatusOK, gin.H{"errNo": 500, "errMsg": err.Error()})
        return
    }
    totalRows, _ := ModelArticle.GetArticlesTotal()
    c.JSON(http.StatusOK, gin.H{"errNo": 0, "errMsg": "success", "data": gin.H{
        "list":      articleList,
        "total":     totalRows,
        "page":      nowPage,
        "obstring":  ob,
        "obstatus":  obStatus,
    }})
}

/**
 * makeOrderStr 组织 ORDER BY 字段
 */
func (ctrl ArticleController) makeOrderStr(ob string, obStatus string) (obString string) {
    if ob == "ctime" {
        if obStatus == "asc" {
            return "createTime ASC"
        }
        return "createTime DESC"
    }
    if ob == "mtime" {
        if obStatus == "asc" {
            return "updateTime ASC"
        }
        return "updateTime DESC"
    }
    if ob == "pub" {
        if obStatus == "asc" {
            return "pubStatus ASC"
        }
        return "pubStatus DESC"
    }
    return "createTime DESC"
}

/**
 * ArticleEdit 返回编辑文章所需的数据(文章内容、专辑列表、标签列表、已关联标签/专辑).
 *
 * param: *gin.Context c
 */
func (ctrl ArticleController) ArticleEdit(c *gin.Context) {
    articleId, _ := strconv.Atoi(c.Param("id"))
    modelArticle := new(Article)
    articleInfo, _ := modelArticle.GetArticleInfo(articleId)

    modelRelationArticleAlbum := new(RelationArticleAlbums)
    relationAlbumId := modelRelationArticleAlbum.GetBelongAlbumByArticleId(articleId)

    modelAlbums := new(Album)
    albumsList, _ := modelAlbums.GetAlbumList(-1) //获取所有专辑

    modelTags := new(Tags)
    tags, _ := modelTags.GetAllTags()
    reModelTags := new(RelationArticleTags)
    articleIdSlice := []int{articleId}
    relationTagsMap := ctrl.mapRelationTags(reModelTags.GetArticleTags(articleIdSlice), tags)
    relationTagsStr, relationTagsIds := ctrl.getTagStrAndIds(relationTagsMap[articleId])

    c.JSON(http.StatusOK, gin.H{"errNo": 0, "errMsg": "success", "data": gin.H{
        "articleInfo":     articleInfo,
        "albumsList":      albumsList,
        "tagsList":        tags,
        "relationTagsStr": relationTagsStr,
        "relationTagsIds": relationTagsIds,
        "relationAlbumId": relationAlbumId,
    }})
}

/**
 * mapRelationTags 将文章-标签关联组织为 文章ID -> []Tags
 */
func (ctrl ArticleController) mapRelationTags(relationTagsInfo []RelationArticleTags, allTags []Tags) (relationTags map[int][]Tags) {
    tagsMap := make(map[int]Tags, len(allTags))
    for _, tagInfo := range allTags {
        tagsMap[tagInfo.Id] = tagInfo
    }
    relationTags = make(map[int][]Tags)
    for _, relations := range relationTagsInfo {
        tagInfo := tagsMap[relations.TagId]
        relationTags[relations.ArticleId] = append(relationTags[relations.ArticleId], tagInfo)
    }
    return relationTags
}

/**
 * getTagStrAndIds 把文章关联标签转为 名称串 与 ID 串
 */
func (ctrl ArticleController) getTagStrAndIds(articleRelationInfo []Tags) (tagsStr string, tagsIds string) {
    for _, tagInfo := range articleRelationInfo {
        if tagsStr != "" {
            tagsStr = tagsStr + "," + tagInfo.TagName
        } else {
            tagsStr = tagInfo.TagName
        }
        if tagsIds != "" {
            tagsIds = tagsIds + "," + strconv.Itoa(tagInfo.Id)
        } else {
            tagsIds = strconv.Itoa(tagInfo.Id)
        }
    }
    return
}
