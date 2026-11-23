package router

import (
    "github.com/gin-gonic/gin"
    "happyblog/middleware"
    . "happyblog/controller"
)

// StaticRouter 注册读者端静态资源(本地磁盘),上传文件也落在 ./assets/uploads 下
func StaticRouter(router *gin.Engine) {
    router.Static("/assets", "./assets/")
    router.StaticFile("/favicon.ico", "./assets/favicon.ico")
}

// FrontRouter 注册面向读者的服务端渲染页面路由(未安装时跳转到安装向导)
func FrontRouter(router *gin.Engine) {
    g := router.Group("/", middleware.PageInstalledGuard())
    g.GET("/", CFront.Index)
    g.GET("/article/:id", CFront.Article)
    g.GET("/tag/:id", CFront.Tag)
    g.GET("/album/:id", CFront.Album)
}
