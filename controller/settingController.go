package controller

import (
    "net/http"

    "github.com/gin-gonic/gin"
    . "happyblog/models"
)

type SettingController struct {
    controller
}

var CSetting SettingController

/**
 * Settings 返回按分组组织的配置项(JSON).
 *
 * param: *gin.Context c
 */
func (ctrl SettingController) Settings(c *gin.Context) {
    ModelSetting := new(Setting)
    configs, _ := ModelSetting.GetConfigs()
    c.JSON(http.StatusOK, gin.H{"errNo": 0, "errMsg": "success", "data": configs})
}
