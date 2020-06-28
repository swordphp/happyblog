package controller

import (
	"github.com/gin-gonic/gin"
	lib "happyblog/library"
	. "happyblog/models"
	"net/http"
	"strconv"
	"time"
)

type AlbumsController struct {
	controller
}

var CAlbums AlbumsController

type AlbumsJsonRes struct{
	Data Album `json:"data"`
	ErrNo int `json:"errNo"`
	ErrMsg string `json:"errMsg"`
}

func init(){
	webConf ,_ := lib.ReadWebConfig()
	CAlbums.webConf = *webConf
	menu,_ := lib.ReadLanguageConfig(CAlbums.webConf["language"])
	CAlbums.menu = *menu
}
/**
 * 专辑列表
 *
 * param: *gin.Context c
 */
func (ctrl AlbumsController) Albums(c *gin.Context){
	ModelAlbums := new(Album)
	albumsRows,_ := ModelAlbums.GetAlbumList(1)
	c.HTML(http.StatusOK, "index", gin.H{
		"title": ctrl.webConf["sitename"],
		"userInfo": ctrl.GetCacheUinfo(c),
		"menu":ctrl.menu,
		"AlbumsList":albumsRows,
		"albums":1,
	})
}

func (ctrl AlbumsController) AlbumsSaveInfo(c *gin.Context) {
	cacheUserInfo := ctrl.GetCacheUinfo(c)
	createUserId := cacheUserInfo.Id
	albumName := c.PostForm("albumName")
	if albumName == "" {
		c.JSON(http.StatusOK,gin.H{"errNo":10,"errMsg":"params error"})
		return
	}
	publicNum,_ := strconv.ParseInt(c.PostForm("isPublic"),10,8)
	isPublic := int8(publicNum)
	id ,_:= strconv.Atoi(c.PostForm("id"))
	Albums := Album{
		Id:           id,
		AlbumName:    albumName,
		IsPublic:     isPublic,
		AuthorId:     createUserId,
		CreateTime:   time.Time{},
		ArticleTotal: 0,
		AuthorInfo:   User{},
	}
	if id == 0 {
		Albums.CreateAlbum(Albums)
	} else {
		Albums.UpdateAlbum(Albums)
	}
	c.JSON(http.StatusOK,gin.H{"errNo":0,"errMsg":"success"})
}

/**
 * 返回json格式的专辑信息
 *
 * param: *gin.Context c
 */
func (ctrl AlbumsController) AlbumsGetInfo(c *gin.Context) {
	ModelAlbums := new(Album)
	id ,_:= strconv.Atoi(c.Param("id"))
	albumInfo,err := ModelAlbums.GetAlbumInfo(id)
	res := AlbumsJsonRes{
		Data: albumInfo,
		ErrNo:      0,
		ErrMsg:     "success",
	}
	if err != nil {
		res.ErrNo = 500
		res.ErrMsg = err.Error()
		c.JSON(http.StatusOK,res)
	} else {
		c.JSON(http.StatusOK,res)
	}
}