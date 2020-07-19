package router

import (
    "github.com/gin-gonic/gin"
    "happyblog/controller"
)

func ApiRouter (router *gin.Engine) {
    apiGroup := router.Group("/api")
    {
        article := apiGroup.Group("/article/")
        {
            articleApi := new(controller.ArticleApi)
            article.GET("info/:id",articleApi.ArticleInfo)
            article.GET("list/:page",articleApi.ArticleList)
        }
    }
}
