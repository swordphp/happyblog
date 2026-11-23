package models

import (
    "crypto/md5"
    "fmt"
    "github.com/glebarez/sqlite"
    "gorm.io/driver/mysql"
    "gorm.io/gorm"
    "gorm.io/gorm/logger"
    lib "happyblog/library"
    "log"
    "net/url"
    "strconv"
    "strings"
    "time"
)

var ConnInstance *gorm.DB //被初始化的数据库连接

/**
 * Connect 解析配置并初始化数据库连接.
 * 由 main / 安装流程在确认配置就绪后显式调用(替代原来的包级 init).
 * 根据 database.driver 选择后端:
 *   - sqlite (默认): 纯 Go 驱动(modernc.org/sqlite),数据落在本地单文件
 *   - mysql:         连接 MySQL / MariaDB
 */
func Connect(webConf map[string]string) error {
    driver := strings.ToLower(strings.TrimSpace(webConf["database.driver"]))
    if driver == "" {
        driver = "sqlite"
    }

    gormCfg := &gorm.Config{
        Logger: logger.Default.LogMode(logger.Info),
    }

    var db *gorm.DB
    var err error

    switch driver {
    case "mysql", "mariadb":
        dbHost := webConf["database.hostname"]
        dbPort,_ := strconv.ParseInt(webConf["database.port"],10,32)
        dbUserName := webConf["database.username"]
        dbPassword := webConf["database.password"]
        dbName := webConf["database.dbname"]
        if dbName == "" {
            dbName = "happyblog"
        }
        loc := url.QueryEscape(webConf["timezone"])
        dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=true&loc=%s",
            dbUserName, dbPassword, dbHost, dbPort, dbName, loc)
        db, err = gorm.Open(mysql.Open(dsn), gormCfg)
    default:
        //sqlite: database.sqlitepath 指定文件路径,缺省 ./happyblog.db
        dbPath := webConf["database.sqlitepath"]
        if dbPath == "" {
            dbPath = "./happyblog.db"
        }
        db, err = gorm.Open(sqlite.Open(dbPath), gormCfg)
    }

    if err != nil {
        return fmt.Errorf("DB CONNECT ERR ,driver:%s,err:%w", driver, err)
    }
    ConnInstance = db

    sqlDB, err := db.DB()
    if err != nil {
        return fmt.Errorf("GET SQL DB ERR: %w", err)
    }
    maxIdle,_ := strconv.Atoi(webConf["database.maxidel"])
    maxOpen,_ := strconv.Atoi(webConf["database.maxopen"])
    if maxIdle > 0 {
        sqlDB.SetMaxIdleConns(maxIdle)
    }
    if maxOpen > 0 {
        sqlDB.SetMaxOpenConns(maxOpen)
    }

    //自动建表(对 sqlite 首次运行尤其有用;对已存在的表只补缺失列,不破坏数据)
    autoMigrate(db)

    //首次运行(无任何用户)时,按配置创建一个默认管理员,方便开箱即用
    SeedAdmin(webConf)
    return nil
}

/**
 * SeedAdmin 当用户表为空时,用配置里的 admin.email / admin.password 创建默认管理员
 * 密码存储为 md5(明文 + salt),与登录逻辑保持一致
 */
func SeedAdmin(webConf map[string]string) {
    email := strings.TrimSpace(webConf["admin.email"])
    password := webConf["admin.password"]
    if email == "" || password == "" {
        return
    }
    userModel := User{}
    if userModel.Count() > 0 {
        return
    }
    salt := webConf["salt"]
    saltPassword := fmt.Sprintf("%x", md5.Sum([]byte(password+salt)))
    nick := webConf["admin.nickname"]
    if nick == "" {
        nick = "admin"
    }
    _, err := userModel.CreateUser(User{
        AccountEmail:    email,
        AccountPassword: saltPassword,
        NickName:        nick,
        LastLogin:       time.Now(),
        EmailVerify:     1,
    })
    if err != nil {
        lib.Logf("seed admin err","%v",err)
    } else {
        log.Printf("seeded default admin user: %s", email)
    }
}

/**
 * 自动迁移所有数据模型,保证各数据库下表结构存在
 */
func autoMigrate(db *gorm.DB) {
    err := db.AutoMigrate(
        &User{},
        &Article{},
        &Album{},
        &Tags{},
        &Setting{},
        &RelationArticleAlbums{},
        &RelationArticleTags{},
    )
    if err != nil {
        lib.Logf("auto migrate err","%v",err)
    }
}
