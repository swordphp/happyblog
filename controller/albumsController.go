package controller

import (
    "strconv"
    "time"

    "github.com/gin-gonic/gin"
    . "happyblog/models"
)

type AlbumsController struct {
    controller
}

var CAlbums AlbumsController

/**
 * Albums 返回专辑列表。
 */
func (ctrl AlbumsController) Albums(c *gin.Context) {
    albums, _ := new(Album).GetAlbumList(1)
    respOK(c, albums)
}

/**
 * AlbumsSaveInfo 新建或更新专辑(id 为 0 表示新建)。
 */
func (ctrl AlbumsController) AlbumsSaveInfo(c *gin.Context) {
    albumName := c.PostForm("albumName")
    if albumName == "" {
        respFail(c, errCodeParam, "params error")
        return
    }
    publicNum, _ := strconv.ParseInt(c.PostForm("isPublic"), 10, 8)
    id, _ := strconv.Atoi(c.PostForm("id"))

    album := Album{
        Id:         id,
        AlbumName:  albumName,
        IsPublic:   int8(publicNum),
        AuthorId:   ctrl.GetCacheUinfo(c).Id,
        CreateTime: time.Time{},
        AuthorInfo: &User{},
    }
    if id == 0 {
        album.CreateAlbum(album)
    } else {
        album.UpdateAlbum(album)
    }
    respOK(c, nil)
}

/**
 * AlbumsGetInfo 返回单个专辑信息。
 */
func (ctrl AlbumsController) AlbumsGetInfo(c *gin.Context) {
    id, _ := strconv.Atoi(c.Param("id"))
    albumInfo, err := new(Album).GetAlbumInfo(id)
    if err != nil {
        respFail(c, errCodeServer, err.Error())
        return
    }
    respOK(c, albumInfo)
}
