package controller

import (
    "github.com/gin-gonic/gin"
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
