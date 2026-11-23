package library

import (
    "os"
    "path/filepath"
    "sync"

    viper "github.com/spf13/viper"
)

const configDir = "./configs"

var (
    confMu     sync.RWMutex
    cachedConf map[string]string
)

// ConfigFilePath 返回主配置文件的磁盘路径
func ConfigFilePath() string {
    return filepath.Join(configDir, "webconfig.yaml")
}

// IsInstalled 判断站点是否已完成安装(主配置文件是否存在)
func IsInstalled() bool {
    _, err := os.Stat(ConfigFilePath())
    return err == nil
}

/**
 * LoadConfig 从磁盘(重新)读取 webconfig.yaml 并刷新缓存.
 * 安装完成后需要调用它让新写入的配置生效.
 * 读取失败时返回空 map(非 nil),保证调用方的 *map 解引用安全.
 */
func LoadConfig() (map[string]string, error) {
    v := viper.New()
    v.AddConfigPath(configDir)
    v.AddConfigPath("../configs") // 兼容在子包目录下运行(如 go test)
    v.SetConfigName("webconfig")
    v.SetConfigType("yaml")

    conf := make(map[string]string)
    err := v.ReadInConfig()
    if err == nil {
        for _, key := range v.AllKeys() {
            conf[key] = v.GetString(key)
        }
    }

    confMu.Lock()
    cachedConf = conf
    confMu.Unlock()
    return conf, err
}

/**
 * ReadWebConfig 返回缓存的网站配置,首次调用时按需加载.
 * 永不 panic,返回的 map 始终非 nil(未安装时为空 map).
 */
func ReadWebConfig() (conf *map[string]string, err error) {
    confMu.RLock()
    c := cachedConf
    confMu.RUnlock()
    if c != nil {
        return &c, nil
    }
    loaded, err := LoadConfig()
    return &loaded, err
}

/**
 * ReadLanguageConfig 读取语言文件中的菜单翻译信息.
 * 读取失败时返回空 map(非 nil),永不 panic.
 */
func ReadLanguageConfig(language string) (conf *map[string]string, err error) {
    tmpConf := make(map[string]string)
    if language == "" {
        language = "cn"
    }
    v := viper.New()
    v.AddConfigPath(filepath.Join(configDir, "language"))
    v.AddConfigPath(filepath.Join("../configs", "language"))
    v.SetConfigName(language)
    v.SetConfigType("yaml")
    if err = v.ReadInConfig(); err == nil {
        tmpConf = v.GetStringMapString("menu")
    }
    return &tmpConf, err
}
