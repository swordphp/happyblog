# HappyBlog

一个用 Go 编写的开源博客系统：Gin + GORM，默认使用 **SQLite**（本地单文件、免安装），
同时支持 **MySQL / MariaDB**。自带后台管理界面、七牛云图片上传、面向读者的服务端渲染展示站，
以及一组对外的只读 API 和 RSS 订阅。

## 技术栈

- **语言 / Web 框架**：Go 1.22、[Gin](https://github.com/gin-gonic/gin)
- **ORM**：[GORM v2](https://gorm.io)（`gorm.io/gorm`）
- **数据库**：SQLite（纯 Go 驱动 `glebarez/sqlite` + `modernc.org/sqlite`，**无需 CGO**）/ MySQL / MariaDB
- **配置**：Viper（`configs/webconfig.yaml`）
- **会话**：`gin-contrib/sessions`（cookie store）+ `patrickmn/go-cache` 缓存用户信息
- **图片存储**：本地磁盘（默认）/ 七牛云 OSS，由 `storage.driver` 切换
- **Markdown 渲染**：`gomarkdown/markdown`
- **RSS**：`gorilla/feeds`
- **日志**：logrus + file-rotatelogs（按天切割）

## 目录结构

```
happyblog/
├── main.go              # 入口，装载各路由
├── router/              # 路由：Static(读者端)/Admin(后台)/Api(公开接口)
├── controller/          # 控制器：后台页面、后台接口、公开 API、读者端页面、RSS
├── models/              # 数据模型与查询（文章/专辑/标签/用户/设置/关联表）；base.go 负责选库+自动建表
├── library/             # 公共库：配置读取、渲染、七牛上传、日志、字符串/Markdown 工具
├── middleware/          # 中间件：后台鉴权、访问日志
├── templates/           # 后台模板 + 读者端 front_*.tmpl
├── assets/              # 静态资源（css/js/favicon）
├── configs/             # 配置文件与多语言
├── createTable.sql      # MySQL 建表语句（SQLite 由程序自动建表，无需它）
├── Dockerfile / docker-compose.yml
└── Makefile
```

## 快速开始

### 方式一：SQLite（默认，最简单）

无需安装任何数据库，程序会在首次启动时自动建表，并按配置创建默认管理员。

1. **准备配置**（默认即为 SQLite，一般无需改动）：
   ```bash
   cp configs/webconfig.yaml.template configs/webconfig.yaml
   ```

2. **运行**：
   ```bash
   go run main.go
   ```
   服务默认监听 `:8888`。读者端首页 http://localhost:8888/ ，���台 http://localhost:8888/admin/login 。

数据库文件默认落在 `./happyblog.db`（可用 `database.sqlitepath` 修改）。默认管理员取自配置的
`admin.email` / `admin.password`（模板里是 `admin@happyblog.local` / `admin123`，**上线前务必修改**）。

### 方式二：MySQL / MariaDB

1. 把 `configs/webconfig.yaml` 的 `database.driver` 改为 `mysql`，填好 `hostname/port/username/password/dbname`。
2. 创建数据库（表可由程序自动迁移创建，也可手动导入 `createTable.sql`）：
   ```bash
   mysql -u root -p -e "CREATE DATABASE happyblog DEFAULT CHARSET utf8mb4;"
   ```
3. `go run main.go`。

> 用户也可不依赖自动建管理员，自行在 `happyblog_tblUser` 表插入：密码字段为 `md5(明文密码 + salt)`，
> `salt` 取自 `webconfig.yaml`。

### 方式三：Docker Compose（app + MariaDB）

```bash
cp configs/webconfig.yaml.template configs/webconfig.yaml
# 将 database.driver 改为 mysql，hostname 改为 db，username/password 改为 root/happyblog
docker-compose up --build
```
会同时启动 MariaDB（自动导入 `createTable.sql`）和应用容器。

## 常用命令（Makefile）

```bash
make run     # 本地运行
make build   # 编译二进制
make test    # 运行单元测试
make vet     # 静态检查
make tidy    # 整理依赖
```

## 对外 API

统一返回体：`{ "errNo": 0, "errMsg": "success", "data": ... }`。只暴露已发布 / 公开内容。

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/article/list/:page` | 已发布文章分页列表（每页 10 条） |
| GET | `/api/article/info/:id` | 单篇已发布文章（含标签、所属专辑） |
| GET | `/api/tags` | 全部标签 |
| GET | `/api/tag/:id/articles` | 某标签下的已发布文章 |
| GET | `/api/albums/:page` | 公开专辑分页列表 |
| GET | `/api/album/:id` | 单个公开专辑及其文章 |
| GET | `/api/rss` | RSS 2.0 订阅（最新已发布文章） |

读者端页面：`/`（文章列表）、`/article/:id`（详情，Markdown 渲染）、`/tag/:id`、`/album/:id`。

## 图片存储

后台上传接口 `POST /admin/api/upload`，存储后端由 `storage.driver` 决定：

- `local`（默认）：保存到 `storage.localpath`（缺省 `./assets/uploads`），通过 `/assets/uploads/` 静态路由访问。无需任何外部服务，和 SQLite 一样开箱即用。
- `qiniu`：上传到七牛云，需填好 `qiniuOSS` 的 `ak/sk/bucket/zone`，访问 URL 用 `fileviewurlpre` 拼接。

## 功能进度

- [x] 文章 CRUD（支持 Markdown）
- [x] 专辑管理及文章-专辑关联、计数
- [x] 标签：打标签、选择已有标签、**标签列表/删除管理**
- [x] 图片上传与内嵌：**本地磁盘（默认）/ 七牛云，可切换**
- [x] 网站设置管理
- [x] **多数据库支持：SQLite（默认）/ MySQL / MariaDB，自动建表**
- [x] **对外只读 API**
- [x] **RSS 订阅**
- [x] **面向读者的服务端渲染展示站**
- [ ] 富文本编辑器（可选，当前以 Markdown 为主）
- [ ] 缓存层（Redis/Memcache，低流量暂不需要）

## License

MIT，见 [LICENSE](LICENSE)。
