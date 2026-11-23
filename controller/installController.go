package controller

import (
    "net/http"
    "os"

    "github.com/gin-gonic/gin"
    lib "happyblog/library"
)

type InstallController struct {
    controller
}

var CInstall InstallController

/**
 * Status 返回安装状态,供 SPA 判断是否进入安装向导.
 * 该接口在未安装时也可用(不经过已安装校验).
 *
 * param: *gin.Context c
 */
func (ctrl InstallController) Status(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{"errNo": 0, "errMsg": "success", "data": gin.H{
        "installed": lib.IsInstalled(),
    }})
}

/**
 * Install 执行安装:写配置、建目录、连库、建表、建管理员.
 * 成功后无需重启即可进入正常模式.
 *
 * param: *gin.Context c
 */
func (ctrl InstallController) Install(c *gin.Context) {
    if lib.IsInstalled() {
        c.JSON(http.StatusConflict, gin.H{"errNo": 409, "errMsg": "站点已安装"})
        return
    }

    var form lib.InstallForm
    if err := c.ShouldBindJSON(&form); err != nil {
        c.JSON(http.StatusOK, gin.H{"errNo": 10, "errMsg": "请求参数错误: " + err.Error()})
        return
    }
    if err := form.Validate(); err != nil {
        c.JSON(http.StatusOK, gin.H{"errNo": 10, "errMsg": err.Error()})
        return
    }

    if err := lib.WriteWebConfig(form); err != nil {
        c.JSON(http.StatusOK, gin.H{"errNo": 500, "errMsg": "写入配置失败: " + err.Error()})
        return
    }

    conf, err := lib.LoadConfig()
    if err != nil {
        _ = os.Remove(lib.ConfigFilePath())
        c.JSON(http.StatusOK, gin.H{"errNo": 500, "errMsg": "加载配置失败: " + err.Error()})
        return
    }

    if err := Bootstrap(conf); err != nil {
        //安装未完成,删除配置文件以便用户修正后重试
        _ = os.Remove(lib.ConfigFilePath())
        c.JSON(http.StatusOK, gin.H{"errNo": 500, "errMsg": "初始化失败(请检查数据库配置): " + err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"errNo": 0, "errMsg": "success"})
}
