package controller

import (
	"crypto/md5"
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	lib "happyblog/library"
	"happyblog/middleware"
	. "happyblog/models"
	"net/http"
	"strconv"
)
type controller struct{
	webConf map[string]string
	menu map[string]string
	cacheUserInfo User
}

type AdminController struct {
	controller
}

var CAdmin AdminController

func init(){
	webConf ,_ := lib.ReadWebConfig()
	CAdmin.webConf = *webConf
	menu,_ := lib.ReadLanguageConfig(CAdmin.webConf["language"])
	CAdmin.menu = *menu
}

func (ctrl controller) GetCacheUinfo (c *gin.Context) (userInfo User) {
	session := sessions.Default(c)
	userInfo = User{}
	if session.Get("UserId") != nil {
		userId := session.Get("UserId").(int)
		tmpInfo,_ := middleware.UserInfoCache.Get("userInfo" + strconv.Itoa(userId))
		userInfo = tmpInfo.(User)
	}
	return
}

/**
 * 管理系统的欢迎界面
 *
 * param: *gin.Context c
 */
func (ctrl AdminController) Welcome(c *gin.Context){
	ctrl.cacheUserInfo = ctrl.GetCacheUinfo(c)
	c.HTML(http.StatusOK, "index", gin.H{
		"title": ctrl.webConf["sitename"],
		"userInfo": ctrl.cacheUserInfo,
		"menu":ctrl.menu,
		"welcome":1,
	})
}

/**
 * 显示管理系统的登录界面
 *
 * param: *gin.Context c
 */
func (ctrl AdminController) LoginDisplay(c *gin.Context){

	c.HTML(http.StatusOK, "login", gin.H{
		"title": ctrl.webConf["sitename"],
		"menu":ctrl.menu,
	})
}

/**
 * 处理用户的登出
 *
 * param: *gin.Context c
 */
func (ctrl AdminController) LogOut (c *gin.Context) {
	session := sessions.Default(c)
	if session.Get("UserId") != nil {
		userId := session.Get("UserId").(int)
		cacheKey := "userInfo" + strconv.Itoa(userId)
		middleware.UserInfoCache.Delete(cacheKey)
		session.Delete("UserId")
		session.Save()
	}
	c.Redirect(http.StatusFound,"/admin/login")
}


/**
 * 处理用户的登录
 * 用户登录的密码使用  md5(password+salt)
 * param: *gin.Context c
 */
func (ctrl AdminController) Login(c *gin.Context){
	session := sessions.Default(c)
	loginUser := c.PostForm("userEmail")
	passWord  := c.PostForm("password")
	saltPassword := fmt.Sprintf("%x", md5.Sum([]byte(passWord + ctrl.webConf["salt"])))
	//fmt.Println(saltPassword)
	ModelUser := new(User)
	userInfo,err := ModelUser.GetUserInfo(loginUser)
	if userInfo.Id == 0 || userInfo.AccountPassword != saltPassword || err != nil {
		//用户不存在,或者密码错误,或者发生了数据库错误.

		errMsg := ctrl.menu["passwderr"]
		c.HTML(http.StatusOK, "login", gin.H{
			"title": ctrl.webConf["sitename"],
			"menu":ctrl.menu,
			"errmsg":errMsg,
			"inputUsername": c.PostForm("userEmail"),
			"moreinfo":fmt.Sprintf("%v",err),
		})
	} else {
		session.Set("UserId",userInfo.Id)
		session.Save()
		cacheKey := "userInfo" + strconv.Itoa(userInfo.Id)
		middleware.UserInfoCache.Set(cacheKey,userInfo,0)
		c.Redirect(http.StatusFound,"/admin/")
	}
}

