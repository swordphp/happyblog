package main

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	. "happyblog/controller"
	"happyblog/middleware"
)



func main(){
	r := gin.Default()
	r.Static("/assets", "./assets/")
	r.LoadHTMLGlob("templates/*")
	r.StaticFile("/favicon.ico","./assets/favicon.ico")

	adminGroup := r.Group("/admin")
	store := cookie.NewStore([]byte("happyBlog"))
	adminGroup.Use(sessions.Sessions("happyBlogSession", store))
	adminGroup.Use(middleware.AdminAuth())
	adminGroup.Use(middleware.LoggerToFile())
	{
		//登录相关路由
		adminGroup.GET("/",CAdmin.Welcome)
		adminGroup.GET("/login",CAdmin.LoginDisplay)
		adminGroup.POST("/login",CAdmin.Login)
		adminGroup.GET("/logout",CAdmin.LogOut)

		//文章相关路由
		adminGroup.GET("/articles",CArticle.Articles)
		adminGroup.GET("/article/new",CArticle.ArticleEdit)
		adminGroup.GET("/article/edit/:id",CArticle.ArticleEdit)

		//标签管理相关路由
		adminGroup.GET("/tags",CTags.Tags)

		//专辑管理相关路由
		adminGroup.GET("/albums",CAlbums.Albums)
		adminGroup.GET("/albums/info/:id",CAlbums.AlbumsGetInfo)
		adminGroup.POST("/albums/save",CAlbums.AlbumsSaveInfo)

		//用户管理相关路由
		adminGroup.GET("/users",CUsers.Users)

		//管理后台使用的一些接口
		adminGroup.POST("/api/upload",CAdminApi.Upload)
		adminGroup.POST("/api/articlesave",CAdminApi.ArticleSave)
	}

	apiGroup := r.Group("/api")
	{
		apiGroup.GET("/")
	}



	r.Run(":8888")
}