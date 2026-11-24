package controller

import (
    "strconv"
    "strings"

    "github.com/gin-gonic/gin"
    . "happyblog/models"
)

type ArticleController struct {
    controller
}

var CArticle ArticleController

/**
 * Articles 返回后台文章列表(分页 + 排序)。
 */
func (ctrl ArticleController) Articles(c *gin.Context) {
    ob := c.DefaultQuery("obstring", "ctime")
    obStatus := c.DefaultQuery("obstatus", "desc")
    nowPage, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
    if nowPage < 1 {
        nowPage = 1
    }

    modelArticle := new(Article)
    articleList, err := modelArticle.GetArticlesList(ctrl.makeOrderStr(ob, obStatus), nowPage)
    if err != nil {
        respFail(c, errCodeServer, err.Error())
        return
    }
    total, _ := modelArticle.GetArticlesTotal()
    respOK(c, gin.H{
        "list":     articleList,
        "total":    total,
        "page":     nowPage,
        "obstring": ob,
        "obstatus": obStatus,
    })
}

// makeOrderStr 把排序字段/方向组织成白名单内的 ORDER BY 子句
func (ctrl ArticleController) makeOrderStr(ob string, obStatus string) string {
    asc := obStatus == "asc"
    switch ob {
    case "mtime":
        if asc {
            return "updateTime ASC"
        }
        return "updateTime DESC"
    case "pub":
        if asc {
            return "pubStatus ASC"
        }
        return "pubStatus DESC"
    default: // ctime
        if asc {
            return "createTime ASC"
        }
        return "createTime DESC"
    }
}

/**
 * ArticleEdit 返回编辑文章所需的数据:文章内容、全部专辑/标签、已关联标签与专辑。
 */
func (ctrl ArticleController) ArticleEdit(c *gin.Context) {
    articleId, _ := strconv.Atoi(c.Param("id"))
    articleInfo, _ := new(Article).GetArticleInfo(articleId)

    relationAlbumId := new(RelationArticleAlbums).GetBelongAlbumByArticleId(articleId)
    albumsList, _ := new(Album).GetAlbumList(-1) //-1:全部专辑
    tags, _ := new(Tags).GetAllTags()

    relationTagsMap := ctrl.mapRelationTags(new(RelationArticleTags).GetArticleTags([]int{articleId}), tags)
    relationTagsStr, relationTagsIds := ctrl.getTagStrAndIds(relationTagsMap[articleId])

    respOK(c, gin.H{
        "articleInfo":     articleInfo,
        "albumsList":      albumsList,
        "tagsList":        tags,
        "relationTagsStr": relationTagsStr,
        "relationTagsIds": relationTagsIds,
        "relationAlbumId": relationAlbumId,
    })
}

// mapRelationTags 将文章-标签关联组织为 文章ID -> []Tags
func (ctrl ArticleController) mapRelationTags(relationTagsInfo []RelationArticleTags, allTags []Tags) map[int][]Tags {
    tagsMap := make(map[int]Tags, len(allTags))
    for _, tagInfo := range allTags {
        tagsMap[tagInfo.Id] = tagInfo
    }
    relationTags := make(map[int][]Tags)
    for _, rel := range relationTagsInfo {
        relationTags[rel.ArticleId] = append(relationTags[rel.ArticleId], tagsMap[rel.TagId])
    }
    return relationTags
}

// getTagStrAndIds 把文章关联标签转为「名称串」与「ID 串」(逗号分隔)
func (ctrl ArticleController) getTagStrAndIds(tags []Tags) (tagsStr string, tagsIds string) {
    names := make([]string, 0, len(tags))
    ids := make([]string, 0, len(tags))
    for _, t := range tags {
        names = append(names, t.TagName)
        ids = append(ids, strconv.Itoa(t.Id))
    }
    return strings.Join(names, ","), strings.Join(ids, ",")
}
