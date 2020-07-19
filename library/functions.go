package library

import (
    viper "github.com/spf13/viper"
)

var config = viper.New()


/**
 * 获取配置中的网站整体配置
 *
 * return: map[string]string 配置信息
 * return: error
 */
func ReadWebConfig()  (conf *map[string]string, err error){
    tmpConf := make(map[string]string)

    config.AddConfigPath("./configs")     //设置读取的文件路径
    config.SetConfigName("webconfig") //设置读取的文件名
    config.SetConfigType("yaml")
    if err := config.ReadInConfig(); err != nil {
        panic(err)
    }
    for _,key  := range config.AllKeys() {
        tmpConf[key] = config.GetString(key)
    }
    return &tmpConf,err
}

/**
 * 获取配置信息中的语言
 *
 * param: string language
 * return: map[string]string 对应的key->v 翻译信息
 * return: error
 */
func ReadLanguageConfig(language string) (conf *map[string]string,err error) {
    tmpConf := make(map[string]string)
    if language == "" {
        language = "cn"
    }
    config.AddConfigPath("./configs/language")
    config.SetConfigName(language)
    config.SetConfigType("yaml")
    if err := config.ReadInConfig(); err != nil {
        panic(err)
    }
    tmpConf = config.GetStringMapString("menu")
    return &tmpConf,err
}


