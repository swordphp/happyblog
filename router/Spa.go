package router

import (
    "embed"
    "io/fs"
    "net/http"
    "strings"

    "github.com/gin-gonic/gin"
)

/**
 * SpaRouter 服务内嵌的前端单页应用(读者端 + 管理端共用一套 SPA):
 *   - /spa/*     Vite 构建出的静态资源(assetsDir=spa)
 *   - /favicon.ico 等根级静态文件
 *   - 其余未匹配到 API/静态资源的 GET 请求一律返回 index.html,
 *     交给前端路由(history 模式)处理
 *
 * 必须在所有 API 路由注册之后调用(NoRoute 作为兜底)。
 */
func SpaRouter(router *gin.Engine, distFS embed.FS) {
    dist, _ := fs.Sub(distFS, "web/dist")

    // Vite 构建产物的静态资源目录(vite.config 中 build.assetsDir = 'spa')
    if spaSub, err := fs.Sub(dist, "spa"); err == nil {
        router.StaticFS("/spa", http.FS(spaSub))
    }

    serveIndex := func(c *gin.Context) {
        data, err := fs.ReadFile(dist, "index.html")
        if err != nil {
            c.String(http.StatusInternalServerError, "前端未构建,请先执行 `make build-web`(或在 web/ 下 `npm install && npm run build`)")
            return
        }
        c.Data(http.StatusOK, "text/html; charset=utf-8", data)
    }

    // 根级静态文件(如 favicon),存在则直接返回
    serveRootFile := func(name, contentType string) gin.HandlerFunc {
        return func(c *gin.Context) {
            data, err := fs.ReadFile(dist, name)
            if err != nil {
                c.Status(http.StatusNotFound)
                return
            }
            c.Data(http.StatusOK, contentType, data)
        }
    }
    router.GET("/favicon.ico", serveRootFile("favicon.ico", "image/x-icon"))

    // 兜底:API 路径返回 JSON 404,其余交给前端 SPA
    router.NoRoute(func(c *gin.Context) {
        p := c.Request.URL.Path
        if strings.HasPrefix(p, "/api") || strings.HasPrefix(p, "/admin/api") {
            c.JSON(http.StatusNotFound, gin.H{"errNo": 404, "errMsg": "not found"})
            return
        }
        if c.Request.Method != http.MethodGet {
            c.Status(http.StatusNotFound)
            return
        }
        serveIndex(c)
    })
}
