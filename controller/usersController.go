package controller

import (
	"github.com/gin-gonic/gin"
	lib "happyblog/library"
	"net/http"
)

type UserController struct {
	controller
}

var CUsers UserController


func init(){
	webConf ,_ := lib.ReadWebConfig()
	CUsers.webConf = *webConf
	menu,_ := lib.ReadLanguageConfig(CUsers.webConf["language"])
	CUsers.menu = *menu
}
/**
 * 用户列表
 *
 * param: *gin.Context c
 */
func (ctrl UserController) Users(c *gin.Context){
	c.HTML(http.StatusOK, "index", gin.H{
		"title": ctrl.webConf["sitename"],
		"userInfo": ctrl.GetCacheUinfo(c),
		"menu":ctrl.menu,
		"users":1,
	})
}