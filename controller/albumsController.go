package controller

import (
    "net/http"
    "strconv"
    "time"

    "github.com/gin-gonic/gin"
    . "happyblog/models"
)

type AlbumsController struct {
    controller
}

var CAlbums AlbumsController

type AlbumsJsonRes struct {
    Data   Album  `json:"data"`
    ErrNo  int    `json:"errNo"`
    ErrMsg string `json:"errMsg"`
}

/**
 * Albums 返回专辑列表(JSON).
 *
 * param: *gin.Context c
 */
func (ctrl AlbumsController) Albums(c *gin.Context) {
    ModelAlbums := new(Album)
    albumsRows, _ := ModelAlbums.GetAlbumList(1)
    c.JSON(http.StatusOK, gin.H{"errNo": 0, "errMsg": "success", "data": albumsRows})
}

/**
 * AlbumsSaveInfo 新建或更新专辑.
 *
 * param: *gin.Context c
 */
func (ctrl AlbumsController) AlbumsSaveInfo(c *gin.Context) {
    cacheUserInfo := ctrl.GetCacheUinfo(c)
    createUserId := cacheUserInfo.Id
    albumName := c.PostForm("albumName")
    if albumName == "" {
        c.JSON(http.StatusOK, gin.H{"errNo": 10, "errMsg": "params error"})
        return
    }
    publicNum, _ := strconv.ParseInt(c.PostForm("isPublic"), 10, 8)
    isPublic := int8(publicNum)
    id, _ := strconv.Atoi(c.PostForm("id"))
    Albums := Album{
        Id:           id,
        AlbumName:    albumName,
        IsPublic:     isPublic,
        AuthorId:     createUserId,
        CreateTime:   time.Time{},
        ArticleTotal: 0,
        AuthorInfo:   &User{},
    }
    if id == 0 {
        Albums.CreateAlbum(Albums)
    } else {
        Albums.UpdateAlbum(Albums)
    }
    c.JSON(http.StatusOK, gin.H{"errNo": 0, "errMsg": "success"})
}

/**
 * AlbumsGetInfo 返回单个专辑信息(JSON).
 *
 * param: *gin.Context c
 */
func (ctrl AlbumsController) AlbumsGetInfo(c *gin.Context) {
    ModelAlbums := new(Album)
    id, _ := strconv.Atoi(c.Param("id"))
    albumInfo, err := ModelAlbums.GetAlbumInfo(id)
    res := AlbumsJsonRes{Data: albumInfo, ErrNo: 0, ErrMsg: "success"}
    if err != nil {
        res.ErrNo = 500
        res.ErrMsg = err.Error()
    }
    c.JSON(http.StatusOK, res)
}
