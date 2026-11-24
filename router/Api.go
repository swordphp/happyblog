package router

import (
    "github.com/gin-gonic/gin"
    "happyblog/controller"
    "happyblog/middleware"
)

// InstallRouter 注册安装向导接口,未安装时也可访问
func InstallRouter(router *gin.Engine) {
    router.GET("/api/install/status", controller.CInstall.Status)
    router.POST("/api/install", controller.CInstall.Install)
}

// ApiRouter 注册读者端公开接口(未安装时返回 409)
func ApiRouter(router *gin.Engine) {
    apiGroup := router.Group("/api", middleware.APIInstalledGuard())
    {
        article := apiGroup.Group("/article/")
        {
            articleApi := new(controller.ArticleApi)
            article.GET("info/:id", articleApi.ArticleInfo)
            article.GET("list/:page", articleApi.ArticleList)
        }

        publicApi := new(controller.PublicApi)
        apiGroup.GET("/siteinfo", publicApi.SiteInfo)
        apiGroup.GET("/tags", publicApi.Tags)
        apiGroup.GET("/tag/:id/articles", publicApi.TagArticles)
        apiGroup.GET("/albums/:page", publicApi.Albums)
        apiGroup.GET("/album/:id", publicApi.AlbumInfo)

        //RSS 订阅
        apiGroup.GET("/rss", controller.CRss.Rss)
    }
}
