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
    setting := new(Setting)
    c.ShouldBind(&setting)
    rowId,err := setting.CreateConfig(*setting)
    if err != nil {
        c.JSON(http.StatusOK,gin.H{"errNo":500,"errMsg":err,"data":""})
    } else {
        data := gin.H{"rowId":rowId}
        c.JSON(http.StatusOK,gin.H{"errNo":0,"errMsg":"success","data":data})
    }
}

/**
 * 修改已有的配置信息
 *
 * param: *gin.Context c
 */
func (ctrl AdminApiController) SettingSave(c *gin.Context) {
    setting := new(Setting)
    c.ShouldBind(&setting)
    _,err := setting.UpdateConfigInfo(*setting)
    if err != nil {
        c.JSON(http.StatusOK,gin.H{"errNo":500,"errMsg":err,"data":""})
    } else {
        c.JSON(http.StatusOK,gin.H{"errNo":0,"errMsg":"success","data":""})
    }
}


/**
 * 上传文件
 *
 * param: *gin.Context c
 */
func (ctrl AdminApiController) Upload(c *gin.Context){
    uploadFile, err := c.FormFile("file")
    if err != nil {
        c.JSON(http.StatusOK,gin.H{"errNo":10,"errMsg":"no file uploaded","viewUrl":""})
        return
    }
    fileHandel ,_:= uploadFile.Open()
    defer fileHandel.Close()
    filename ,_:= c.GetPostForm("filename")
    filename = fmt.Sprintf("%x", md5.Sum([]byte(filename)))
    extSlice := strings.Split(uploadFile.Filename,".")
    ext := strings.Join(extSlice[len(extSlice)-1:],"")
    filename = filename + "." + ext
    buf := bytes.NewBuffer(nil)
    io.Copy(buf, fileHandel)

    //按配置选择存储后端,默认本地磁盘;qiniu 走七牛云
    driver := strings.ToLower(strings.TrimSpace(webConf["storage.driver"]))
    var viewUrl string
    var upErr error
    if driver == "qiniu" {
        var upName string
        upName, upErr = lib.QiNiuUploader.UpLoadFile(filename,buf.Bytes(),uploadFile.Size)
        if upErr == nil {
            viewUrl = webConf["fileviewurlpre"] + upName
        }
    } else {
        var upName string
        upName, upErr = lib.LocalUploader.UpLoadFile(filename,buf.Bytes())
        if upErr == nil {
            viewUrl = lib.LocalUploader.ViewUrl(upName)
        }
    }

    if upErr != nil {
        log.Println(upErr)
        c.JSON(http.StatusOK,gin.H{"errNo":500,"errMsg":fmt.Sprintf("%v",upErr),"viewUrl":""})
        return
    }
    c.JSON(http.StatusOK,gin.H{"errNo":0,"errMsg":"","viewUrl":viewUrl})
}

/**
 * 通过接口移除一个文章
 *
 * param: *gin.Context c
 */
func (ctrl AdminApiController) ArticleRemove(c *gin.Context) {
    articleId ,_ := strconv.Atoi(c.DefaultPostForm("id","0"))
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
    settingId ,_ := strconv.Atoi(c.DefaultPostForm("id","0"))
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
    if strings.TrimSpace(tagName) == "" {
        c.JSON(http.StatusOK,gin.H{"errNo":10,"errMsg":"tagName is empty","data":""})
        return
    }
    modelTags := new(Tags)
    tagId ,err := modelTags.AddTag(tagName)
    if err != nil {
        errInfo := fmt.Sprintf("%v",err)
        c.JSON(http.StatusOK,gin.H{"errNo":500,"errMsg":errInfo,"data":""})
        return
    }
    resInfo := make(map[string]string)
    resInfo["tagName"] = tagName
    resInfo["tagId"] = strconv.Itoa(tagId)
    c.JSON(http.StatusOK,gin.H{"errNo":0,"errMsg":"success","data":resInfo})
}

/**
 * 获取全部标签列表
 *
 * param: *gin.Context c
 */
func (ctrl AdminApiController) TagList(c *gin.Context) {
    modelTags := new(Tags)
    tags,err := modelTags.GetAllTags()
    if err != nil {
        c.JSON(http.StatusOK,gin.H{"errNo":500,"errMsg":fmt.Sprintf("%v",err),"data":nil})
        return
    }
    c.JSON(http.StatusOK,gin.H{"errNo":0,"errMsg":"success","data":tags})
}

/**
 * 通过标签ID删除标签(同时清理文章关联)
 *
 * param: *gin.Context c
 */
func (ctrl AdminApiController) TagRemove(c *gin.Context) {
    tagId ,_ := strconv.Atoi(c.DefaultPostForm("id","0"))
    if tagId == 0 {
        c.JSON(http.StatusOK,gin.H{"errNo":10,"errMsg":"id is required","affectRows":0})
        return
    }
    modelTags := new(Tags)
    affect := modelTags.RemoveTagById(tagId)
    c.JSON(http.StatusOK,gin.H{"errNo":0,"errMsg":"success","affectRows":affect})
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
        //如果简介信息为空,将内容前200字(按字符)截取为简介
        articleModel.Brief = lib.MakeBrief(articleModel.Content,200)
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
            if albumId != originAlbumId {
                //专辑发生变化:原专辑计数-1(原来有专辑时),新专辑计数+1
                if originAlbumId != 0 {
                    albumModel.UpdateAlbumArticleTotal(originAlbumId,-1)
                }
                albumModel.UpdateAlbumArticleTotal(albumId,1)
            }
        } else {
            //albumId == -1 表示移除专辑归属
            if originAlbumId != 0 {
                relationArticleAlbumsModel.RemoveRowByArticleId(articleId)
                albumModel.UpdateAlbumArticleTotal(originAlbumId,-1)
            }
        }

    } else {
        //创建
        cUserInfo := ctrl.GetCacheUinfo(c)
        articleModel.AuthorId = cUserInfo.Id
        insertId  := articleModel.CreateArticleRow(*articleModel)
        if albumId != -1 {
            //建立专辑关联,并增加该专辑的文章计数
            _ = relationArticleAlbumsModel.UpdateRowByArticleId(insertId,albumId)
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