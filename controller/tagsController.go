package controller

import (
    "github.com/gin-gonic/gin"
    . "happyblog/models"
)

type TagsController struct {
    controller
}

var CTags TagsController

/**
 * Tags 返回全部标签,后台标签管理页使用。
 */
func (ctrl TagsController) Tags(c *gin.Context) {
    tags, err := new(Tags).GetAllTags()
    if err != nil {
        respFail(c, errCodeServer, err.Error())
        return
    }
    respOK(c, tags)
}
