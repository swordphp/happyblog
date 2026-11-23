package main

import (
    "html/template"
    "io/fs"
    "log"
    "os"
    "path/filepath"

    "github.com/gin-gonic/gin"
    "happyblog/controller"
    lib "happyblog/library"
    "happyblog/router"
)

func main() {
    // 已安装则加载配置并完成运行期初始化(连库/存储/缓存)
    if lib.IsInstalled() {
        conf, err := lib.LoadConfig()
        if err != nil {
            log.Fatalf("加载配置失败: %v", err)
        }
        if err := controller.Bootstrap(conf); err != nil {
            log.Fatalf("初始化失败: %v", err)
        }
    }

    // 首次运行时把内置读者端静态资源释放到磁盘
    extractReaderAssets()

    r := gin.Default()

    // 读者端模板从内嵌文件系统加载,无需依赖磁盘 templates 目录
    tmpl, err := template.ParseFS(frontTemplatesFS, "templates/front_*.tmpl")
    if err != nil {
        log.Fatalf("加载读者端模板失败: %v", err)
    }
    r.SetHTMLTemplate(tmpl)

    router.StaticRouter(r)          //读者端静态资源
    router.InstallRouter(r)         //安装向导接口(未安装也可用)
    router.AdminRouter(r, distFS)   //后台 SPA + 接口
    router.ApiRouter(r)             //读者端公开接口
    router.FrontRouter(r)           //读者端页面

    r.Run(":8888")
}

/**
 * extractReaderAssets 把内嵌的读者端静态文件释放到 ./assets(仅当目标不存在时).
 * 让单个二进制分发时也能自带样式与脚本。
 */
func extractReaderAssets() {
    _ = os.MkdirAll("./assets", 0755)
    entries := []string{"style.css", "script.js", "favicon.ico"}
    for _, name := range entries {
        target := filepath.Join("./assets", name)
        if _, err := os.Stat(target); err == nil {
            continue
        }
        data, err := fs.ReadFile(readerAssetsFS, "assets/"+name)
        if err != nil {
            continue
        }
        _ = os.WriteFile(target, data, 0644)
    }
}
