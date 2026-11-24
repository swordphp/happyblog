package router

import (
    "github.com/gin-gonic/gin"
)

// StaticRouter 注册读者上传的文件(本地磁盘),上传文件落在 ./assets/uploads 下。
// 前端 SPA 自身的静态资源由 SpaRouter 挂在 /spa 下,两者路径不冲突。
func StaticRouter(router *gin.Engine) {
    router.Static("/assets", "./assets/")
}
