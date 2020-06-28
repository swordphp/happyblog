package controller

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	lib "happyblog/library"
	"happyblog/middleware"
	. "happyblog/models"
	"net/http"
	"strconv"
)


type TagsController struct {
	controller
}

var CTags TagsController

func init(){
	webConf ,_ := lib.ReadWebConfig()
	CTags.webConf = *webConf
	menu,_ := lib.ReadLanguageConfig(CTags.webConf["language"])
	CTags.menu = *menu
}
/**
 * 文章列表页面
 *
 * param: *gin.Context c
 */
func (ctrl TagsController) Tags(c *gin.Context){
	session := sessions.Default(c)
	cacheUserInfo := User{}
	if session.Get("UserId") != nil {
		userId := session.Get("UserId").(int)
		tmpInfo,_ := middleware.UserInfoCache.Get("userInfo" + strconv.Itoa(userId))
		cacheUserInfo = tmpInfo.(User)
	}
	c.HTML(http.StatusOK, "index", gin.H{
		"title": ctrl.webConf["sitename"],
		"userInfo": cacheUserInfo,
		"menu":ctrl.menu,
		"tags":1,
	})
}