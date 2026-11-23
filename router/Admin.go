package router

import (
    "embed"
    "io/fs"
    "net/http"
    "strings"

    "github.com/gin-contrib/sessions"
    "github.com/gin-contrib/sessions/cookie"
    "github.com/gin-gonic/gin"
    . "happyblog/controller"
    lib "happyblog/library"
    "happyblog/middleware"
)

/**
 * AdminRouter 注册后台路由:
 *   - 后台 SPA(Vite+Vue 构建产物,内嵌于二进制)挂在 /admin 下
 *   - 后台 JSON 接口挂在 /admin/api 下,登录接口公开,其余需要会话鉴权
 *   - SPA 的前端路由通过 NoRoute 兜底返回 index.html
 */
func AdminRouter(router *gin.Engine, distFS embed.FS) {
    dist, _ := fs.Sub(distFS, "web/dist")

    // Vite 构建后的静态资源位于 web/dist/assets
    if assetsSub, err := fs.Sub(dist, "assets"); err == nil {
        router.StaticFS("/admin/assets", http.FS(assetsSub))
    }

    serveIndex := func(c *gin.Context) {
        data, err := fs.ReadFile(dist, "index.html")
        if err != nil {
            c.String(http.StatusInternalServerError, "admin UI 未构建,请先执行 `make build-web` 或 `npm run build`")
            return
        }
        c.Data(http.StatusOK, "text/html; charset=utf-8", data)
    }
    router.GET("/admin", serveIndex)

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

    // SPA 前端路由兜底
    router.NoRoute(func(c *gin.Context) {
        p := c.Request.URL.Path
        if strings.HasPrefix(p, "/admin/api") {
            c.JSON(http.StatusNotFound, gin.H{"errNo": 404, "errMsg": "not found"})
            return
        }
        if strings.HasPrefix(p, "/admin") {
            serveIndex(c)
            return
        }
        c.String(http.StatusNotFound, "404 not found")
    })
}
