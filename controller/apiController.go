package controller

import (
    "strconv"

    lib "happyblog/library"
    "happyblog/middleware"
    "happyblog/models"
)

// webConf 保存全站配置,安装完成或启动加载配置后由 Setup 写入.
// 所有控制器共享它(替代原来每个控制器各自读取配置的包级 init).
var webConf = map[string]string{}

// Setup 注入全站配置供各控制器使用.
func Setup(conf map[string]string) {
    if conf != nil {
        webConf = conf
    }
}

/**
 * Bootstrap 在配置就绪后完成运行期初始化:
 * 连接数据库、初始化存储、创建运行目录、初始化会话缓存、注入配置.
 * main 在已安装时调用它,安装向导在写入配置后也调用它(无需重启即可进入正常模式).
 */
func Bootstrap(conf map[string]string) error {
    if err := models.Connect(conf); err != nil {
        return err
    }
    lib.ConfigureStorage(conf)
    if err := lib.EnsureRuntimeDirs(conf); err != nil {
        return err
    }
    alive, _ := strconv.ParseInt(conf["sessionalive"], 10, 64)
    middleware.InitCache(alive)
    Setup(conf)
    return nil
}
