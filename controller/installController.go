package controller

import (
    "os"

    "github.com/gin-gonic/gin"
    lib "happyblog/library"
)

type InstallController struct {
    controller
}

var CInstall InstallController

/**
 * Status 返回安装状态,供 SPA 判断是否进入安装向导。
 * 未安装时也可访问(不经过已安装校验)。
 */
func (ctrl InstallController) Status(c *gin.Context) {
    respOK(c, gin.H{"installed": lib.IsInstalled()})
}

/**
 * Install 执行安装:写配置、建目录、连库建表、创建管理员。
 * 成功后无需重启即可进入正常模式。
 */
func (ctrl InstallController) Install(c *gin.Context) {
    if lib.IsInstalled() {
        respFail(c, errCodeNotFound, "站点已安装")
        return
    }

    var form lib.InstallForm
    if err := c.ShouldBindJSON(&form); err != nil {
        respFail(c, errCodeParam, "请求参数错误: "+err.Error())
        return
    }
    if err := form.Validate(); err != nil {
        respFail(c, errCodeParam, err.Error())
        return
    }

    if err := lib.WriteWebConfig(form); err != nil {
        respFail(c, errCodeServer, "写入配置失败: "+err.Error())
        return
    }

    conf, err := lib.LoadConfig()
    if err != nil {
        _ = os.Remove(lib.ConfigFilePath())
        respFail(c, errCodeServer, "加载配置失败: "+err.Error())
        return
    }

    if err := Bootstrap(conf); err != nil {
        //安装未完成,删除配置文件以便用户修正后重试
        _ = os.Remove(lib.ConfigFilePath())
        respFail(c, errCodeServer, "初始化失败(请检查数据库配置): "+err.Error())
        return
    }

    respOK(c, nil)
}
