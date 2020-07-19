package controller

import (
    "fmt"
    "github.com/gin-gonic/gin"
    lib "happyblog/library"
    . "happyblog/models"
    "log"
    "net/http"
    "strconv"
)

type ArticleController struct{
    controller
}

var CArticle ArticleController


func init(){
    webConf ,_ := lib.ReadWebConfig()
    CArticle.webConf = *webConf
    menu,_ := lib.ReadLanguageConfig(CArticle.webConf["language"])
    CArticle.menu = *menu
}

/**
 * 文章列表页面
 *
 * param: *gin.Context c
 */
func (ctrl ArticleController) Articles(c *gin.Context){
    obString := ctrl.makeOrderStr(c.Query("obstring"),c.Query("obstatus"))
    nowPage,_ := strconv.Atoi(c.DefaultQuery("page","1"))
    ModelArticle := new(Article)
    articleList ,err := ModelArticle.GetArticlesList(obString,nowPage)
    if err != nil {
        log.Fatal(fmt.Sprintf("ARTICLE LIST GET ERROR :err:%v",err))
    }

    totalRows,_ := ModelArticle.GetArticlesTotal()
    c.HTML(http.StatusOK, "index", gin.H{
        "title": ctrl.webConf["sitename"],
        "userInfo": ctrl.GetCacheUinfo(c),
        "menu":ctrl.menu,
        "articles":1,
        "articleslist":articleList,
        "obstatus" : c.DefaultQuery("obstatus","desc"),
        "ob" : c.DefaultQuery("obstring","ctime"),
        "curPage":nowPage,
        "totalRows":totalRows,
    })
}

/**
 * 用于组织ORDER BY 字段
 *
 * param: string ob
 * param: string obstatus
 * return: string
 */
func (ctrl ArticleController) makeOrderStr(ob string,obStatus string) (obString string) {
    if ob == "ctime" {
        if obStatus == "asc" {
            return "createTime ASC"
        } else {
            return "createTime DESC"
        }
    }
    if ob == "mtime" {
        if obStatus == "asc" {
            return "updateTime ASC"
        } else {
            return "updateTime DESC"
        }
    }
    if ob == "pub" {
        if obStatus == "asc" {
            return "pubStatus ASC"
        } else {
            return "pubStatus DESC"
        }
    }
    return "createTime DESC"
}
/**
 * 编辑文章的方法
 *
 * param: *gin.Context c
 */
func (ctrl ArticleController) ArticleEdit(c *gin.Context) {
    articleId,_ := strconv.Atoi(c.Param("id"))
    //如果获取到articleId  则去读取articleInfo
    modelArticle := new(Article)
    articleInfo,err := modelArticle.GetArticleInfo(articleId)
    if err != nil {
        log.Println("get article info err")
    }
    modelRelationArticleAlbum := new(RelationArticleAlbums)
    relationAlbumId := modelRelationArticleAlbum.GetBelongAlbumByArticleId(articleId)
    if err != nil {
        log.Println("get albumInfo err")
    }
    modelAlbums := new(Album)
    albumsList,_ := modelAlbums.GetAlbumList(-1)//获取所有的专辑信息

    modelTags := new(Tags)
    tags,_ := modelTags.GetAllTags()
    reModelTags := new(RelationArticleTags)
    articleIdSlice := make([]int,1)
    articleIdSlice = append(articleIdSlice, articleId)
    reModelTags.GetArticleTags(articleIdSlice)
    relationTagsMap := ctrl.mapRelationTags(reModelTags.GetArticleTags(articleIdSlice),tags)
    relationTagsStr,relationTagsIds := ctrl.getTagStrAndIds(relationTagsMap[articleId])

    //如果获取到了articleId ,则是编辑旧文章,读取旧文章内容.
    c.HTML(http.StatusOK, "index", gin.H{
        "title": ctrl.webConf["sitename"],
        "userInfo": ctrl.GetCacheUinfo(c),
        "menu":ctrl.menu,
        "articleEdit":1,
        "articleInfo":articleInfo,
        "albumsList":albumsList,//专辑列表
        "relationTagsStr":relationTagsStr,//文章关联的相关tags
        "relationTagsIds":relationTagsIds,//文章关联的相关tags,id 串
        "tagsList":tags,
        "relationAlbumId" : relationAlbumId,//关联专辑ID
    })
}

/**
 * 将所有标签和文章的标签关联组织成文章对应的关联标签的形式.
 * param: map[int] articleIdMap
 * return: map[int][]Tags
 */
func (ctrl ArticleController) mapRelationTags(relationTagsInfo []RelationArticleTags,allTags []Tags) (relationTags map[int][]Tags) {
    //组织一个taginfo的map
    tagsMap := make(map[int]Tags,len(allTags))
    for _,tagInfo := range allTags {
        tagsMap[tagInfo.Id] = tagInfo
    }
    relationTags = make(map[int][]Tags)
    for _,relations := range relationTagsInfo{
        tagInfo := tagsMap[relations.TagId]
        if relationTags[relations.ArticleId] == nil {
            relationTags[relations.ArticleId] = *new([]Tags)
        }
        relationTags[relations.ArticleId] =  append(relationTags[relations.ArticleId],tagInfo)
    }
    return relationTags
}

/**
 * 通过文章关联的tags信息,获取需要传送到前端的字符串和id串
 *
 * param: []Tags articleRelationInfo
 * return: string
 * return: string
 */
func (ctrl ArticleController) getTagStrAndIds (articleRelationInfo []Tags) (tagsStr string,tagsIds string) {
    if len(articleRelationInfo) >0 {
        for _,tagInfo := range articleRelationInfo {
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
    }
    return
}