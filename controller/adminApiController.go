package controller

import (
    "bytes"
    "crypto/md5"
    "fmt"
    "io"
    "log"
    "regexp"
    "strconv"
    "strings"

    "github.com/gin-gonic/gin"
    lib "happyblog/library"
    . "happyblog/models"
)

// htmlTagRe 用于从富文本正文中剥离 HTML 标签,生成纯文本简介
var htmlTagRe = regexp.MustCompile(`<[^>]+>`)

type AdminApiController struct {
    controller
}

var CAdminApi AdminApiController

/**
 * GetSettingInfo 获取指定的配置项信息
 */
func (ctrl AdminApiController) GetSettingInfo(c *gin.Context) {
    settingId, _ := strconv.Atoi(c.DefaultQuery("id", "0"))
    info, _ := new(Setting).GetConfigInfo(settingId)
    respOK(c, info)
}

/**
 * SettingAdd 添加新的配置项
 */
func (ctrl AdminApiController) SettingAdd(c *gin.Context) {
    setting := new(Setting)
    if err := c.ShouldBind(setting); err != nil {
        respFail(c, errCodeParam, "参数错误: "+err.Error())
        return
    }
    rowId, err := setting.CreateConfig(*setting)
    if err != nil {
        respFail(c, errCodeServer, err.Error())
        return
    }
    respOK(c, gin.H{"rowId": rowId})
}

/**
 * SettingSave 修改已有的配置项
 */
func (ctrl AdminApiController) SettingSave(c *gin.Context) {
    setting := new(Setting)
    if err := c.ShouldBind(setting); err != nil {
        respFail(c, errCodeParam, "参数错误: "+err.Error())
        return
    }
    if _, err := setting.UpdateConfigInfo(*setting); err != nil {
        respFail(c, errCodeServer, err.Error())
        return
    }
    respOK(c, nil)
}
/**
 * Upload 上传图片,按配置选择本地磁盘或七牛云存储,返回可访问 URL。
 */
func (ctrl AdminApiController) Upload(c *gin.Context) {
    uploadFile, err := c.FormFile("file")
    if err != nil {
        respFail(c, errCodeParam, "no file uploaded")
        return
    }
    fileHandle, err := uploadFile.Open()
    if err != nil {
        respFail(c, errCodeServer, err.Error())
        return
    }
    defer fileHandle.Close()

    rawName, _ := c.GetPostForm("filename")
    hashed := fmt.Sprintf("%x", md5.Sum([]byte(rawName)))
    ext := ""
    if i := strings.LastIndex(uploadFile.Filename, "."); i >= 0 {
        ext = uploadFile.Filename[i:]
    }
    filename := hashed + ext

    buf := bytes.NewBuffer(nil)
    if _, err := io.Copy(buf, fileHandle); err != nil {
        respFail(c, errCodeServer, err.Error())
        return
    }

    //按配置选择存储后端,默认本地磁盘;qiniu 走七牛云
    driver := strings.ToLower(strings.TrimSpace(webConf["storage.driver"]))
    var viewUrl string
    var upErr error
    if driver == "qiniu" {
        var upName string
        upName, upErr = lib.QiNiuUploader.UpLoadFile(filename, buf.Bytes(), uploadFile.Size)
        if upErr == nil {
            viewUrl = webConf["fileviewurlpre"] + upName
        }
    } else {
        var upName string
        upName, upErr = lib.LocalUploader.UpLoadFile(filename, buf.Bytes())
        if upErr == nil {
            viewUrl = lib.LocalUploader.ViewUrl(upName)
        }
    }
    if upErr != nil {
        log.Println(upErr)
        respFail(c, errCodeServer, upErr.Error())
        return
    }
    respOK(c, gin.H{"viewUrl": viewUrl})
}

/**
 * ArticleRemove 删除文章,并清理其专辑关联与计数。
 */
func (ctrl AdminApiController) ArticleRemove(c *gin.Context) {
    articleId, _ := strconv.Atoi(c.DefaultPostForm("id", "0"))
    rowAffect := new(Article).RemoveRow(articleId)

    //删除专辑关联并减少计数
    relationAlbum := new(RelationArticleAlbums)
    if albumId := relationAlbum.GetBelongAlbumByArticleId(articleId); albumId != 0 {
        relationAlbum.RemoveRowByArticleId(articleId)
        new(Album).UpdateAlbumArticleTotal(albumId, -1)
    }
    respOK(c, gin.H{"affectRows": rowAffect})
}

/**
 * SettingRemove 根据 ID 删除一条配置项。
 */
func (ctrl AdminApiController) SettingRemove(c *gin.Context) {
    settingId, _ := strconv.Atoi(c.DefaultPostForm("id", "0"))
    rowAffect := new(Setting).RemoveRow(settingId)
    respOK(c, gin.H{"affectRows": rowAffect})
}

/**
 * TagAdd 新增标签。
 */
func (ctrl AdminApiController) TagAdd(c *gin.Context) {
    tagName, _ := c.GetPostForm("tagName")
    if strings.TrimSpace(tagName) == "" {
        respFail(c, errCodeParam, "tagName is empty")
        return
    }
    tagId, err := new(Tags).AddTag(tagName)
    if err != nil {
        respFail(c, errCodeServer, err.Error())
        return
    }
    respOK(c, gin.H{"tagName": tagName, "tagId": tagId})
}

/**
 * TagRemove 根据标签 ID 删除标签(同时清理文章关联)。
 */
func (ctrl AdminApiController) TagRemove(c *gin.Context) {
    tagId, _ := strconv.Atoi(c.DefaultPostForm("id", "0"))
    if tagId == 0 {
        respFail(c, errCodeParam, "id is required")
        return
    }
    affect := new(Tags).RemoveTagById(tagId)
    respOK(c, gin.H{"affectRows": affect})
}
/**
 * ArticleSave 新建或更新文章,同步维护专辑关联/计数与标签关联。
 * id 为 0 表示新建;albumId 为 -1 表示无专辑。
 */
func (ctrl AdminApiController) ArticleSave(c *gin.Context) {
    articleId, _ := strconv.Atoi(c.DefaultPostForm("id", "0"))
    pubStatus, _ := strconv.ParseInt(c.PostForm("pubStatus"), 0, 8)
    independPage, _ := strconv.ParseInt(c.PostForm("independPage"), 0, 8)

    articleModel := new(Article)
    articleModel.Content = c.PostForm("content")
    //内容格式:markdown(默认)或 html(富文本编辑器)
    contentType := strings.ToLower(strings.TrimSpace(c.DefaultPostForm("contentType", "markdown")))
    if contentType != "html" {
        contentType = "markdown"
    }
    articleModel.ContentType = contentType
    articleModel.PubStatus = int8(pubStatus)
    articleModel.IndependPage = int8(independPage)
    articleModel.Title = c.DefaultPostForm("title", "no title")
    articleModel.Keywords = c.DefaultPostForm("keywords", "")
    articleModel.Brief = c.DefaultPostForm("describe", "")
    if articleModel.Brief == "" {
        //简介为空时,截取正文前 200 字作为简介;富文本先去掉 HTML 标签
        briefSrc := articleModel.Content
        if contentType == "html" {
            briefSrc = htmlTagRe.ReplaceAllString(briefSrc, "")
        }
        articleModel.Brief = lib.MakeBrief(briefSrc, 200)
    }
    articleModel.Headimage = c.DefaultPostForm("headimage", "")
    articleModel.Uri = c.DefaultPostForm("uri", "")

    albumId, _ := strconv.Atoi(c.PostForm("albumId"))
    albumModel := new(Album)
    relationAlbum := new(RelationArticleAlbums)

    if articleId != 0 {
        //更新:先取原专辑,再按新专辑调整归属与计数
        originAlbumId := relationAlbum.GetBelongAlbumByArticleId(articleId)
        articleModel.Id = articleId
        articleModel.UpdateArticleRow(*articleModel)
        if albumId != -1 {
            _ = relationAlbum.UpdateRowByArticleId(articleId, albumId)
            if albumId != originAlbumId {
                if originAlbumId != 0 {
                    albumModel.UpdateAlbumArticleTotal(originAlbumId, -1)
                }
                albumModel.UpdateAlbumArticleTotal(albumId, 1)
            }
        } else if originAlbumId != 0 {
            //albumId == -1:移除专辑归属
            relationAlbum.RemoveRowByArticleId(articleId)
            albumModel.UpdateAlbumArticleTotal(originAlbumId, -1)
        }
    } else {
        //新建
        articleModel.AuthorId = ctrl.GetCacheUinfo(c).Id
        articleId = articleModel.CreateArticleRow(*articleModel)
        if albumId != -1 {
            _ = relationAlbum.UpdateRowByArticleId(articleId, albumId)
            albumModel.UpdateAlbumArticleTotal(albumId, 1)
        }
    }

    //处理标签关联:tags 为逗号分隔的标签 ID
    if tagsStr := c.DefaultPostForm("tags", ""); tagsStr != "" {
        tagIds := make([]int, 0)
        for _, s := range strings.Split(tagsStr, ",") {
            if id, _ := strconv.Atoi(s); id != 0 {
                tagIds = append(tagIds, id)
            }
        }
        new(RelationArticleTags).Relations(tagIds, articleId)
    }
    respOK(c, gin.H{"articleId": articleId})
}


