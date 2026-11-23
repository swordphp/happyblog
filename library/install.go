package library

import (
    "crypto/rand"
    "encoding/hex"
    "errors"
    "os"
    "strings"

    yaml "gopkg.in/yaml.v2"
)

// InstallForm 安装向导提交的表单
type InstallForm struct {
    SiteName string `json:"sitename"`
    SiteURL  string `json:"siteurl"`
    SiteDesc string `json:"sitedesc"`

    DBDriver   string `json:"dbDriver"` // sqlite | mysql
    SQLitePath string `json:"sqlitePath"`
    DBHost     string `json:"dbHost"`
    DBPort     string `json:"dbPort"`
    DBName     string `json:"dbName"`
    DBUser     string `json:"dbUser"`
    DBPassword string `json:"dbPassword"`

    AdminEmail    string `json:"adminEmail"`
    AdminPassword string `json:"adminPassword"`
    AdminNickname string `json:"adminNickname"`
}

// Validate 校验安装表单的必填项
func (f InstallForm) Validate() error {
    if strings.TrimSpace(f.SiteName) == "" {
        return errors.New("站点名称不能为空")
    }
    if strings.TrimSpace(f.AdminEmail) == "" || strings.TrimSpace(f.AdminPassword) == "" {
        return errors.New("管理员邮箱与密码不能为空")
    }
    driver := strings.ToLower(strings.TrimSpace(f.DBDriver))
    if driver == "mysql" || driver == "mariadb" {
        if f.DBHost == "" || f.DBName == "" || f.DBUser == "" {
            return errors.New("MySQL 连接信息不完整")
        }
    }
    return nil
}

func randomHex(n int) string {
    b := make([]byte, n)
    if _, err := rand.Read(b); err != nil {
        return "happyblog-fallback-secret"
    }
    return hex.EncodeToString(b)
}

// RandomHex 导出的随机十六进制字符串生成器
func RandomHex(n int) string {
    return randomHex(n)
}

/**
 * SessionSecret 返回会话 cookie 的签名密钥.
 * 优先用配置里的 sessionsecret,其次 salt,未安装时回退到随机临时密钥.
 */
func SessionSecret() string {
    conf, _ := ReadWebConfig()
    c := *conf
    s := c["sessionsecret"]
    if s == "" {
        s = c["salt"]
    }
    if s == "" {
        s = randomHex(32)
    }
    return s
}

/**
 * WriteWebConfig 根据安装表单生成 configs/webconfig.yaml.
 * salt 与 sessionsecret 自动生成随机值,无需用户填写.
 */
func WriteWebConfig(f InstallForm) error {
    driver := strings.ToLower(strings.TrimSpace(f.DBDriver))
    if driver == "" {
        driver = "sqlite"
    }
    sqlitePath := f.SQLitePath
    if sqlitePath == "" {
        sqlitePath = "./happyblog.db"
    }
    port := f.DBPort
    if port == "" {
        port = "3306"
    }
    nickname := f.AdminNickname
    if nickname == "" {
        nickname = "admin"
    }

    cfg := map[string]interface{}{
        "sitename":      f.SiteName,
        "siteurl":       f.SiteURL,
        "sitedesc":      f.SiteDesc,
        "language":      "cn",
        "salt":          randomHex(16),
        "sessionsecret": randomHex(32),
        "sessionalive":  "7200",
        "database": map[string]interface{}{
            "driver":     driver,
            "sqlitepath": sqlitePath,
            "port":       port,
            "hostname":   f.DBHost,
            "dbname":     f.DBName,
            "username":   f.DBUser,
            "password":   f.DBPassword,
            "maxidel":    10,
            "maxopen":    20,
        },
        "admin": map[string]interface{}{
            "email":    f.AdminEmail,
            "password": f.AdminPassword,
            "nickname": nickname,
        },
        "timezone": "Asia/Shanghai",
        "accesslog": map[string]interface{}{
            "logfile": "./log",
            "logname": "access.log",
        },
        "servicelog": map[string]interface{}{
            "logfile": "./log",
            "logname": "service.log",
        },
        "storage": map[string]interface{}{
            "driver":      "local",
            "localpath":   "./assets/uploads",
            "localurlpre": "/assets/uploads/",
        },
        "fileviewurlpre": "",
        "qiniuoss": map[string]interface{}{
            "bucket":       "happyblog",
            "ak":           "",
            "sk":           "",
            "https":        "1",
            "zone":         "ZoneHuanan",
            "usecdndomain": "0",
            "pathto":       "",
        },
    }

    data, err := yaml.Marshal(cfg)
    if err != nil {
        return err
    }
    if err := os.MkdirAll(configDir, 0755); err != nil {
        return err
    }
    return os.WriteFile(ConfigFilePath(), data, 0644)
}

/**
 * ConfigureStorage 用配置初始化本地与七牛云上传器.
 * 在配置加载或安装完成后调用.
 */
func ConfigureStorage(conf map[string]string) {
    LocalUploader.saveDir = conf["storage.localpath"]
    if LocalUploader.saveDir == "" {
        LocalUploader.saveDir = "./assets/uploads"
    }
    LocalUploader.urlPre = conf["storage.localurlpre"]
    if LocalUploader.urlPre == "" {
        LocalUploader.urlPre = "/assets/uploads/"
    }
    if !strings.HasSuffix(LocalUploader.urlPre, "/") {
        LocalUploader.urlPre += "/"
    }

    QiNiuUploader.ak = conf["qiniuoss.ak"]
    QiNiuUploader.sk = conf["qiniuoss.sk"]
    QiNiuUploader.bucket = conf["qiniuoss.bucket"]
    if zoneInfo, err := getZone(conf["qiniuoss.zone"]); err == nil {
        QiNiuUploader.zone = zoneInfo
    }
    QiNiuUploader.https = conf["qiniuoss.https"] == "1"
    QiNiuUploader.useCdnDomain = conf["qiniuoss.usecdndomain"] == "1"
    QiNiuUploader.pathTo = conf["qiniuoss.pathto"]
}

/**
 * EnsureRuntimeDirs 创建运行期需要的目录(日志目录、本地上传目录).
 */
func EnsureRuntimeDirs(conf map[string]string) error {
    dirs := []string{}
    if lf := conf["servicelog.logfile"]; lf != "" {
        dirs = append(dirs, lf)
    }
    if lf := conf["accesslog.logfile"]; lf != "" {
        dirs = append(dirs, lf)
    }
    uploadDir := conf["storage.localpath"]
    if uploadDir == "" {
        uploadDir = "./assets/uploads"
    }
    dirs = append(dirs, uploadDir)

    for _, d := range dirs {
        if err := os.MkdirAll(d, 0755); err != nil {
            return err
        }
    }
    return nil
}
