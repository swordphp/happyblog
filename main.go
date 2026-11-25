package main

import (
    "log"
    "os"

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

    r := gin.Default()

    router.StaticRouter(r)        //读者上传的文件(本地磁盘 ./assets)
    router.InstallRouter(r)       //安装向导接口(未安装也可用)
    router.AdminRouter(r)         //后台 JSON 接口 /admin/api
    router.ApiRouter(r)           //读者端公开接口 /api
    router.SpaRouter(r, distFS)   //前端 SPA(读者端+管理端),放在最后兜底

    //监听端口:默认 8888;PaaS(Railway/Render/Fly 等)会注入 PORT 环境变量
    addr := ":8888"
    if p := os.Getenv("PORT"); p != "" {
        addr = ":" + p
    }
    r.Run(addr)
}
