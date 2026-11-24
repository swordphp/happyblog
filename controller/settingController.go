package controller

import (
    "strings"

    "github.com/gin-gonic/gin"
    lib "happyblog/library"
    . "happyblog/models"
)

type SettingController struct {
    controller
}

var CSetting SettingController

/**
 * Settings 返回按分组组织的配置项。
 */
func (ctrl SettingController) Settings(c *gin.Context) {
    configs, _ := new(Setting).GetConfigs()
    respOK(c, configs)
}

/**
 * SiteConfigGet 返回站点基础信息(名称/简介/地址),供后台「设置」展示。
 */
func (ctrl SettingController) SiteConfigGet(c *gin.Context) {
    respOK(c, gin.H{
        "sitename": webConf["sitename"],
        "sitedesc": webConf["sitedesc"],
        "siteurl":  webConf["siteurl"],
    })
}

/**
 * SiteConfigSave 修改站点基础信息:写回配置文件并刷新内存,无需重启即生效。
 */
func (ctrl SettingController) SiteConfigSave(c *gin.Context) {
    sitename := strings.TrimSpace(c.PostForm("sitename"))
    if sitename == "" {
        respFail(c, errCodeParam, "网站名称不能为空")
        return
    }
    if err := lib.UpdateSiteConfig(sitename, c.PostForm("siteurl"), c.PostForm("sitedesc")); err != nil {
        respFail(c, errCodeServer, "保存失败: "+err.Error())
        return
    }
    //重新加载配置并刷新各控制器共享的 webConf
    if conf, err := lib.LoadConfig(); err == nil {
        Setup(conf)
    }
    respOK(c, nil)
}
