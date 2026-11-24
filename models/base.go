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

    //首次运行(无任何文章)时,插入两篇示例文章,让读者端首页有内容可预览
    SeedArticles(webConf)
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
 * SeedArticles 当文章表为空时插入两篇示例文章(已发布),
 * 让全新站点的读者端首页有内容可预览。已有文章则跳过。
 */
func SeedArticles(webConf map[string]string) {
    articleModel := Article{}
    if total, _ := articleModel.GetArticlesTotal(); total > 0 {
        return
    }
    authorId := 0
    if email := strings.TrimSpace(webConf["admin.email"]); email != "" {
        if u, err := (User{}).GetUserInfo(email); err == nil {
            authorId = u.Id
        }
    }

    samples := []Article{
        {
            Title:     "欢迎使用 HappyBlog",
            Uri:       "welcome",
            Brief:     "一篇自动生成的示例文章,帮助你快速了解 HappyBlog,可随时在后台修改或删除。",
            PubStatus: 1,
            AuthorId:  authorId,
            Content:   sampleWelcome,
        },
        {
            Title:     "Markdown 语法示例",
            Uri:       "markdown-demo",
            Brief:     "展示标题、列表、代码块、引用、表格等常见 Markdown 渲染效果。",
            PubStatus: 1,
            AuthorId:  authorId,
            Content:   sampleMarkdown,
        },
    }
    for _, a := range samples {
        articleModel.CreateArticleRow(a)
    }
    log.Printf("seeded %d sample articles", len(samples))
}

// SEED_CONTENT_PLACEHOLDER

const sampleWelcome = `## 👋 欢迎

这是你用 **HappyBlog** 搭建的新站点。这篇示例文章可以随时在后台修改或删除。

### 快速上手

1. 打开 ` + "`/admin`" + ` 进入后台管理
2. 用安装时设置的管理员邮箱与密码登录
3. 在「文章」里写下你的第一篇博客

> 支持 Markdown 写作、标签与专辑归类、图片上传、RSS 订阅。

祝写作愉快!
`

const sampleMarkdown = `本文演示常见的 Markdown 渲染效果。

## 文本样式

**粗体**、*斜体*、` + "`行内代码`" + `,以及 [链接示例](https://github.com/swordphp/happyblog)。

## 列表

- 无序项 A
- 无序项 B
    - 嵌套项

1. 有序项一
2. 有序项二

## 代码块

~~~go
func main() {
    fmt.Println("Hello, HappyBlog!")
}
~~~

## 引用

> 写作是最好的思考方式。

## 表格

| 功能 | 支持 |
| --- | --- |
| Markdown | ✅ |
| 标签 / 专辑 | ✅ |
| RSS 订阅 | ✅ |
`

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
