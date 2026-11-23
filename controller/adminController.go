package controller

import (
    "crypto/md5"
    "fmt"
    "net/http"
    "strconv"

    "github.com/gin-contrib/sessions"
    "github.com/gin-gonic/gin"
    "happyblog/middleware"
    . "happyblog/models"
)

type controller struct{}

type AdminController struct {
    controller
}

var CAdmin AdminController

// GetCacheUinfo 从会话+缓存中取出当前登录用户信息
func (ctrl controller) GetCacheUinfo(c *gin.Context) (userInfo User) {
    userInfo = User{}
    session := sessions.Default(c)
    if session.Get("UserId") != nil && middleware.UserInfoCache != nil {
        userId := session.Get("UserId").(int)
        if tmpInfo, found := middleware.UserInfoCache.Get("userInfo" + strconv.Itoa(userId)); found {
            userInfo = tmpInfo.(User)
        }
    }
    return
}

/**
 * Login 处理管理员登录,密码为 md5(明文 + salt).
 * 成功后写入会话并返回用户信息 JSON.
 *
 * param: *gin.Context c
 */
func (ctrl AdminController) Login(c *gin.Context) {
    var form struct {
        Email    string `json:"email" form:"email"`
        Password string `json:"password" form:"password"`
    }
    _ = c.ShouldBind(&form)

    saltPassword := fmt.Sprintf("%x", md5.Sum([]byte(form.Password+webConf["salt"])))
    ModelUser := new(User)
    userInfo, err := ModelUser.GetUserInfo(form.Email)
    if userInfo.Id == 0 || userInfo.AccountPassword != saltPassword || err != nil {
        c.JSON(http.StatusOK, gin.H{"errNo": 10, "errMsg": "邮箱或密码错误"})
        return
    }

    session := sessions.Default(c)
    session.Set("UserId", userInfo.Id)
    session.Save()
    middleware.UserInfoCache.Set("userInfo"+strconv.Itoa(userInfo.Id), userInfo, 0)
    c.JSON(http.StatusOK, gin.H{"errNo": 0, "errMsg": "success", "data": userInfo})
}

/**
 * LogOut 清理会话与缓存.
 *
 * param: *gin.Context c
 */
func (ctrl AdminController) LogOut(c *gin.Context) {
    session := sessions.Default(c)
    if session.Get("UserId") != nil {
        userId := session.Get("UserId").(int)
        if middleware.UserInfoCache != nil {
            middleware.UserInfoCache.Delete("userInfo" + strconv.Itoa(userId))
        }
        session.Delete("UserId")
        session.Save()
    }
    c.JSON(http.StatusOK, gin.H{"errNo": 0, "errMsg": "success"})
}

/**
 * Me 返回当前登录用户信息,供 SPA 判断登录态.
 *
 * param: *gin.Context c
 */
func (ctrl AdminController) Me(c *gin.Context) {
    userInfo := ctrl.GetCacheUinfo(c)
    c.JSON(http.StatusOK, gin.H{"errNo": 0, "errMsg": "success", "data": gin.H{
        "user":     userInfo,
        "sitename": webConf["sitename"],
    }})
}
