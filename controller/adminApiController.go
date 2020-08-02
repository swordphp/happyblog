package controller

import (
    "bytes"
    "crypto/md5"
    "fmt"
    "github.com/gin-gonic/gin"
    lib "happyblog/library"
    . "happyblog/models"
    "io"
    "log"
    "net/http"
    "strconv"
    "strings"
)

type AdminApiController struct{
    controller
}

var CAdminApi AdminApiController

func init(){
    webConf ,_ := lib.ReadWebConfig()
    CAdminApi.webConf = *webConf
    menu,_ := lib.ReadLanguageConfig(CAdminApi.webConf["language"])
    CAdminApi.menu = *menu
}


/**
 * 获取指定的配置项信息
 *
 * param: *gin.Context c
 */
func (ctrl AdminApiController) GetSettingInfo(c *gin.Context) {
    settingId ,_ := strconv.Atoi(c.DefaultQuery("id","0"))
    settingModel := new(Setting)
    infos,_ := settingModel.GetConfigInfo(settingId)
    c.JSON(http.StatusOK,gin.H{"errNo":0,"errMsg":"success","data":infos})
}

/**
 * 添加配置组信息
 *
 * param: *gin.Context c
 */
func (ctrl AdminApiController) SettingGroupAdd(c *gin.Context) {

}

/**
 * 添加新的配置信息
 * param: *gin.Context c
 */
func (ctrl AdminApiController) SettingAdd(c *gin.Context) {

}

/**
 * 修改已有的配置信息
 *
 * param: *gin.Context c
 */
func (ctrl AdminApiController) SettingSave(c *gin.Context) {

}


/**
 * 上传文件
 *
 * param: *gin.Context c
 */
func (ctrl AdminApiController) Upload(c *gin.Context){
    uploadFile, _ := c.FormFile("file")
    fileHandel ,_:= uploadFile.Open()
    defer fileHandel.Close()
    filename ,_:= c.GetPostForm("filename")
    filename = fmt.Sprintf("%x", md5.Sum([]byte(filename)))
    extSlice := strings.Split(uploadFile.Filename,".")
    ext := strings.Join(extSlice[len(extSlice)-1:],"")
    filename = filename + "." + ext
    buf := bytes.NewBuffer(nil)
    io.Copy(buf, fileHandel)
    upName,err := lib.QiNiuUploader.UpLoadFile(filename,buf.Bytes(),uploadFile.Size)
    if err != nil {
        log.Println(err)
        c.JSON(http.StatusOK,gin.H{"errNo":500,"errMsg":err,"viewUrl":""})

    } else {
        viewUrl := ctrl.webConf["fileviewurlpre"] + upName
        c.JSON(http.StatusOK,gin.H{"errNo":0,"errMsg":"","viewUrl":viewUrl})
    }
}

/**
 * 通过接口移除一个文章
 *
 * param: *gin.Context c
 */
func (ctrl AdminApiController) ArticleRemove(c *gin.Context) {
    articleId ,_ := strconv.Atoi(c.DefaultQuery("id","0"))
    articleModel := new(Article)
    rowAffect := articleModel.RemoveRow(articleId)

    //删除专辑的关联
    relationAlbum := new(RelationArticleAlbums)
    albumId := relationAlbum.GetBelongAlbumByArticleId(articleId)

    //减少专辑计数,移除专辑关联
    if albumId != 0 {
        relationAlbum.RemoveRowByArticleId(articleId)
        albumModel := new(Album)
        albumModel.UpdateAlbumArticleTotal(albumId,-1)
    }
    c.JSON(http.StatusOK,gin.H{"errNo":0,"errMsg":"success","affectRows":rowAffect})
}
/**
 * 根据传递的ID
 * 来移除一行设置行
 *
 * param: *gin.Context c
 */
func (ctrl AdminApiController) SettingRemove(c *gin.Context) {
    settingId ,_ := strconv.Atoi(c.DefaultQuery("id","0"))
    settingModel := new(Setting)
    rowAffect := settingModel.RemoveRow(settingId)
    c.JSON(http.StatusOK,gin.H{"errNo":0,"errMsg":"success","affectRows":rowAffect})
}

/**
 * 添加标签的方法
 *
 * param: *gin.Context c
 */
func (ctrl AdminApiController) TagAdd(c *gin.Context) {
    tagName,_ :=c.GetPostForm("tagName")
    modelTags := new(Tags)
    tagId ,err := modelTags.AddTag(tagName)
    if err != nil {
        errInfo := fmt.Sprintf("%v",err)
        c.JSON(http.StatusOK,gin.H{"errNo":500,"errMsg":errInfo,"data":""})
    }
    resInfo := make(map[string]string)
    resInfo["tagName"] = tagName
    resInfo["tagId"] = strconv.Itoa(tagId)
    c.JSON(http.StatusOK,gin.H{"errNo":0,"errMsg":"success","data":resInfo})
}

/**
 * 保存文章的方法,接口调用 post
 *
 * param: *gin.Context c
 */
func (ctrl AdminApiController) ArticleSave(c *gin.Context) {
    articleId,_ := strconv.Atoi(c.DefaultPostForm("id",""))
    articleModel := new(Article)
    pubStatus,_ := strconv.ParseInt(c.PostForm("pubStatus"),0,8)
    independPage ,_ := strconv.ParseInt(c.PostForm("independPage"),0,8)
    articleModel.Content = c.PostForm("content")
    articleModel.PubStatus = int8(pubStatus)
    articleModel.IndependPage = int8(independPage)
    articleModel.Title = c.DefaultPostForm("title","no title")
    articleModel.Keywords = c.DefaultPostForm("keywords","")
    articleModel.Brief = c.DefaultPostForm("describe","")
    if articleModel.Brief == "" {
        //如果简介信息为空,将内容前200字截取为简介
        articleModel.Brief = articleModel.Content[0:200]
    }
    articleModel.Headimage = c.DefaultPostForm("headimage","")
    articleModel.Uri = c.DefaultPostForm("uri","")

    albumId,_ := strconv.Atoi(c.PostForm("albumId"))

    albumModel := new(Album)
    relationArticleAlbumsModel := new(RelationArticleAlbums)

    if articleId != 0 {

        //此处获取文章原始的专辑ID

        originAlbumId := relationArticleAlbumsModel.GetBelongAlbumByArticleId(articleId)

        //更新
        articleModel.Id = articleId
        articleModel.UpdateArticleRow(*articleModel)
        if albumId != -1 {
            _ = relationArticleAlbumsModel.UpdateRowByArticleId(articleId,albumId)
        }
        if albumId != originAlbumId {
            //专辑属性发生变化
            albumModel.UpdateAlbumArticleTotal(originAlbumId,-1)
            albumModel.UpdateAlbumArticleTotal(albumId,1)
        }

    } else {
        //创建
        cUserInfo := ctrl.GetCacheUinfo(c)
        articleModel.AuthorId = cUserInfo.Id
        insertId  := articleModel.CreateArticleRow(*articleModel)
        if albumId != -1 {
            _ = relationArticleAlbumsModel.UpdateRowByArticleId(insertId,albumId)
        } else {
            //增加专辑引用数量
            albumModel.UpdateAlbumArticleTotal(albumId,1)
        }
        articleId = insertId
    }

    //此处处理tags相关信息
    tagsStr := c.DefaultPostForm("tags","")
    if tagsStr != "" {
        tagSlice := strings.Split(tagsStr,",")
        tagIds := make([]int,0)
        for _,tagId := range tagSlice {
            intTagId ,_ := strconv.Atoi(tagId)
            if intTagId == 0 {
                continue
            }
            tagIds = append(tagIds, intTagId)
        }
        modelReTags := new(RelationArticleTags)
        modelReTags.Relations(tagIds,articleId)
    }
    c.JSON(http.StatusOK,gin.H{"errNo":0,"errMsg":"","articleId":articleId})
}