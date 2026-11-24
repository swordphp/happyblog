package controller

import (
    "github.com/gin-gonic/gin"
    . "happyblog/models"
)

type UserController struct {
    controller
}

var CUsers UserController

/**
 * Users 返回用户列表。
 */
func (ctrl UserController) Users(c *gin.Context) {
    users, err := new(User).GetUserList()
    if err != nil {
        respFail(c, errCodeServer, err.Error())
        return
    }
    respOK(c, users)
}
