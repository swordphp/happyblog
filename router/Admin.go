package router

import (
    "github.com/gin-contrib/sessions"
    "github.com/gin-contrib/sessions/cookie"
    "github.com/gin-gonic/gin"
    . "happyblog/controller"
    lib "happyblog/library"
    "happyblog/middleware"
)

/**
 * AdminRouter 注册后台 JSON 接口(挂在 /admin/api 下):
 *   - 登录接口公开(仅需已安装),其余接口需要会话鉴权
 *   - 页面由前端 SPA 负责(见 SpaRouter)
 */
func AdminRouter(router *gin.Engine) {
    // 会话存储密钥:优先取配置,未安装时用随机临时密钥
    secret := lib.SessionSecret()
    store := cookie.NewStore([]byte(secret))
    sessionMW := sessions.Sessions("happyBlogSession", store)

    // 登录接口:需已安装,但不需要登录态
    pub := router.Group("/admin/api", sessionMW, middleware.APIInstalledGuard())
    {
        pub.POST("/login", CAdmin.Login)
    }

    // 其余接口:需已安装 + 登录态
    api := router.Group("/admin/api", sessionMW, middleware.APIInstalledGuard(), middleware.AdminAuth())
    api.Use(middleware.LoggerToFile())
    {
        api.GET("/me", CAdmin.Me)
        api.POST("/logout", CAdmin.LogOut)

        //文章
        api.GET("/articles", CArticle.Articles)
        api.GET("/article/:id", CArticle.ArticleEdit)
        api.POST("/article/save", CAdminApi.ArticleSave)
        api.POST("/article/remove", CAdminApi.ArticleRemove)

        //标签
        api.GET("/tags", CTags.Tags)
        api.POST("/tag/add", CAdminApi.TagAdd)
        api.POST("/tag/remove", CAdminApi.TagRemove)

        //专辑
        api.GET("/albums", CAlbums.Albums)
        api.GET("/album/:id", CAlbums.AlbumsGetInfo)
        api.POST("/album/save", CAlbums.AlbumsSaveInfo)

        //用户
        api.GET("/users", CUsers.Users)

        //设置
        api.GET("/settings", CSetting.Settings)
        api.GET("/setting", CAdminApi.GetSettingInfo)
        api.POST("/setting/add", CAdminApi.SettingAdd)
        api.POST("/setting/save", CAdminApi.SettingSave)
        api.POST("/setting/remove", CAdminApi.SettingRemove)

        //上传
        api.POST("/upload", CAdminApi.Upload)
    }
}
