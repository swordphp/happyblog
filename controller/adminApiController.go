package controller

import (
	"bytes"
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
 * 上传文件
 *
 * param: *gin.Context c
 */
func (ctrl AdminApiController) Upload(c *gin.Context){
	uploadFile, _ := c.FormFile("file")
	fileHandel ,_:= uploadFile.Open()
	defer fileHandel.Close()
	filename ,_:= c.GetPostForm("filename")
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
 * 保存文章的方法,接口调用 post
 *
 * param: *gin.Context c
 */
func (ctrl AdminApiController) ArticleSave(c *gin.Context) {
	tmpId,_ := strconv.ParseInt(c.DefaultPostForm("id","0"),0,64)
	articleId := int(tmpId)
	articleRow := new(Article)
	pubStatus,_ := strconv.ParseInt(c.PostForm("pubStatus"),0,8)
	independPage ,_ := strconv.ParseInt(c.PostForm("independPage"),0,8)
	articleRow.Content = c.PostForm("content")
	articleRow.PubStatus = int8(pubStatus)
	articleRow.IndependPage = int8(independPage)
	articleRow.Title = c.DefaultPostForm("title","no title")
	albumId,_ := strconv.Atoi(c.PostForm("albumId"))
	ArticleModel := new(Article)
	AlbumsModel := new(Album)
	RelationArticleAlbumsModel := new(RelationArticleAlbums)
	if articleId != 0 {

		//此处获取文章原始的专辑ID

		originAlbumId := RelationArticleAlbumsModel.GetBelongAlbumByArticleId(articleId)

		//更新
		articleRow.Id = articleId
		ArticleModel.UpdateArticleRow(*articleRow)
		if albumId != -1 {
			_ = RelationArticleAlbumsModel.UpdateRowByArticleId(articleId,albumId)
		}
		if albumId != originAlbumId {
			//专辑属性发生变化
			AlbumsModel.UpdateAlbumArticleTotal(originAlbumId,-1)
			AlbumsModel.UpdateAlbumArticleTotal(albumId,1)
		}

	} else {
		//创建
		cUserInfo := ctrl.GetCacheUinfo(c)
		articleRow.AuthorId = cUserInfo.Id
		insertId  := ArticleModel.CreateArticleRow(*articleRow)
		if albumId != -1 {
			_ = RelationArticleAlbumsModel.UpdateRowByArticleId(insertId,albumId)
		} else {
			//增加专辑引用数量
			AlbumsModel.UpdateAlbumArticleTotal(albumId,1)
		}
		articleId = insertId
	}
	c.JSON(http.StatusOK,gin.H{"errNo":0,"errMsg":"","articleId":articleId})
}