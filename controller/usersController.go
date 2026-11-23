package controller

import (
    "net/http"

    "github.com/gin-gonic/gin"
    . "happyblog/models"
)

type UserController struct {
    controller
}

var CUsers UserController

/**
 * Users 返回用户列表(JSON).
 *
 * param: *gin.Context c
 */
func (ctrl UserController) Users(c *gin.Context) {
    modelUser := new(User)
    users, err := modelUser.GetUserList()
    if err != nil {
        c.JSON(http.StatusOK, gin.H{"errNo": 500, "errMsg": err.Error(), "data": nil})
        return
    }
    c.JSON(http.StatusOK, gin.H{"errNo": 0, "errMsg": "success", "data": users})
}
