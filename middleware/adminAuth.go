package middleware

import (
    "github.com/gin-contrib/sessions"
    "github.com/gin-gonic/gin"
    "github.com/patrickmn/go-cache"
    lib "happyblog/library"
    "net/http"
    "strconv"
    "time"
)
var UserInfoCache *cache.Cache

/**
 * 初始化一个全局的go-cache缓存.
 */
func init(){
    webConf ,_ := lib.ReadWebConfig()
    tmpConf := *webConf
    alive,_ := strconv.ParseInt(tmpConf["sessionalive"],10,64)
    UserInfoCache = cache.New(time.Duration(alive)*time.Second,600*time.Second)
}

/**
 * 处理管理用户的授权.
 * 将存放UserId的session放在用户的cookie中进行管理
 * 将其余信息缓存在服务器上面
 * return: gin.HandlerFunc
 */
func AdminAuth() gin.HandlerFunc{
    return func(c *gin.Context){
        session := sessions.Default(c)
        if (session.Get("UserId") == nil ) && c.FullPath() != "/admin/login"{
            //用户session不存在
            c.Redirect(http.StatusFound,"/admin/login")
        }
        if session.Get("UserId") != nil {
            userId := session.Get("UserId").(int)
            userInfo, err := UserInfoCache.Get("userInfo" + strconv.Itoa(userId))
            if !err {
                session.Delete("UserId")
                session.Save()
                c.Redirect(http.StatusFound,"/admin/login")
            } else {
                UserInfoCache.Set("userInfo" + strconv.Itoa(userId),userInfo,0)
                //刷新缓存存在时间
            }
        }
    }
}