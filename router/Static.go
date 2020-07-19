package router

import "github.com/gin-gonic/gin"

func StaticRouter (router *gin.Engine) {
    router.Static("/assets", "./assets/")
    router.LoadHTMLGlob("templates/*")
    router.StaticFile("/favicon.ico","./assets/favicon.ico")
}
