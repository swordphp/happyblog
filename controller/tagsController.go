package controller

import (
    "net/http"

    "github.com/gin-gonic/gin"
    . "happyblog/models"
)

type TagsController struct {
    controller
}

var CTags TagsController

/**
 * Tags 返回全部标签(JSON).后台标签管理页使用.
 *
 * param: *gin.Context c
 */
func (ctrl TagsController) Tags(c *gin.Context) {
    modelTags := new(Tags)
    tags, err := modelTags.GetAllTags()
    if err != nil {
        c.JSON(http.StatusOK, gin.H{"errNo": 500, "errMsg": err.Error(), "data": nil})
        return
    }
    c.JSON(http.StatusOK, gin.H{"errNo": 0, "errMsg": "success", "data": tags})
}
