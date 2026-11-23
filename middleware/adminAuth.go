package middleware

import (
    "github.com/gin-contrib/sessions"
    "github.com/gin-gonic/gin"
    "github.com/patrickmn/go-cache"
    "net/http"
    "strconv"
    "time"
)

var UserInfoCache *cache.Cache

/**
 * InitCache 初始化保存登录用户信息的全局缓存.
 * 由 main / 安装流程在配置就绪后显式调用(替代原来读取配置的包级 init).
 */
func InitCache(aliveSeconds int64) {
    if aliveSeconds <= 0 {
        aliveSeconds = 7200
    }
    UserInfoCache = cache.New(time.Duration(aliveSeconds)*time.Second, 600*time.Second)
}

/**
 * AdminAuth 处理管理接口的鉴权.
 * 仅挂在需要登录的 /admin/api 分组上(登录接口不经过它).
 * 未登录或会话失效时返回 401 JSON,由前端 SPA 跳转到登录页.
 */
func AdminAuth() gin.HandlerFunc {
    return func(c *gin.Context) {
        if UserInfoCache == nil {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"errNo": 401, "errMsg": "not logged in"})
            return
        }
        session := sessions.Default(c)
        if session.Get("UserId") == nil {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"errNo": 401, "errMsg": "not logged in"})
            return
        }
        userId := session.Get("UserId").(int)
        userInfo, found := UserInfoCache.Get("userInfo" + strconv.Itoa(userId))
        if !found {
            session.Delete("UserId")
            session.Save()
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"errNo": 401, "errMsg": "session expired"})
            return
        }
        //刷新缓存存在时间
        UserInfoCache.Set("userInfo"+strconv.Itoa(userId), userInfo, 0)
        c.Next()
    }
}
