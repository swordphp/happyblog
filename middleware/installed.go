package middleware

import (
    "net/http"

    "github.com/gin-gonic/gin"
    "happyblog/library"
)

/**
 * APIInstalledGuard 保护依赖数据库的接口:未安装时返回 409 JSON.
 */
func APIInstalledGuard() gin.HandlerFunc {
    return func(c *gin.Context) {
        if !library.IsInstalled() {
            c.AbortWithStatusJSON(http.StatusConflict, gin.H{"errNo": 409, "errMsg": "站点尚未安装"})
            return
        }
        c.Next()
    }
}

/**
 * PageInstalledGuard 保护读者端页面:未安装时重定向到安装向导.
 */
func PageInstalledGuard() gin.HandlerFunc {
    return func(c *gin.Context) {
        if !library.IsInstalled() {
            c.Redirect(http.StatusFound, "/admin/install")
            c.Abort()
            return
        }
        c.Next()
    }
}
