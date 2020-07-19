package main

import (
    "github.com/gin-gonic/gin"
    "happyblog/router"
)

func main(){
    r := gin.Default()

    router.StaticRouter(r) //载入静态资源路由
    router.AdminRouter(r) //载入后台用的路由
    router.ApiRouter(r) //载入接口用的路由

    r.Run(":8888")
}