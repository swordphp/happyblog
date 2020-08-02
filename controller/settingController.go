package controller

import (
    "fmt"
    "github.com/gin-gonic/gin"
    lib "happyblog/library"
    . "happyblog/models"
    "net/http"
)

type SettingController struct {
    controller
}

var CSetting SettingController


func init(){
    webConf ,_ := lib.ReadWebConfig()
    CSetting.webConf = *webConf
    menu,_ := lib.ReadLanguageConfig(CSetting.webConf["language"])
    CSetting.menu = *menu
}
/**
 * 用户列表
 *
 * param: *gin.Context c
 */
func (ctrl SettingController) Settings(c *gin.Context){
    ModelSetting := new(Setting)
    configs,_ := ModelSetting.GetConfigs()
    fmt.Print(configs)
    c.HTML(http.StatusOK, "index", gin.H{
        "title": ctrl.webConf["sitename"],
        "userInfo": ctrl.GetCacheUinfo(c),
        "menu":ctrl.menu,
        "configs":configs,
        "setting":1,
    })
}