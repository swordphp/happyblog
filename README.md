# HappyBlog

一个用 Go 编写的开源博客系统：Gin + GORM，默认使用 **SQLite**（本地单文件、免安装），
同时支持 **MySQL / MariaDB**。前端为 **Vite + Vue 3** 单页应用（读者端 + Element Plus 后台管理，
内嵌进二进制），带图形化安装向导、七牛云图片上传，以及一组对外的只读 API 和 RSS 订阅。
最终产物是**一个自带前端的二进制文件**：下载后运行，浏览器里完成安装即可使用。

## 技术栈

- **语言 / Web 框架**：Go 1.25、[Gin](https://github.com/gin-gonic/gin)
- **前端**：[Vite](https://vitejs.dev) + [Vue 3](https://vuejs.org) + [Element Plus](https://element-plus.org)（后台）/ Vue Router / Pinia；构建产物经 `go:embed` 内嵌
- **ORM**：[GORM v2](https://gorm.io)（`gorm.io/gorm`）
- **数据库**：SQLite（纯 Go 驱动 `glebarez/sqlite` + `modernc.org/sqlite`，**无需 CGO**）/ MySQL / MariaDB
- **配置**：Viper（`configs/webconfig.yaml`，由安装向导生成）
- **会话**：`gin-contrib/sessions`（cookie store）+ `patrickmn/go-cache` 缓存用户信息
- **图片存储**：本地磁盘（默认）/ 七牛云 OSS，由 `storage.driver` 切换
- **Markdown 渲染**：前端 `marked`（读者端展示与后台编辑预览）
- **RSS**：`gorilla/feeds`
- **日志**：logrus + file-rotatelogs（按天切割）

## 目录结构

```
happyblog/
├── main.go              # 入口，装载各路由
├── embed.go             # 内嵌前端构建产物 web/dist
├── router/              # 路由：Static(上传文件)/Admin(后台接口)/Api(公开接口)/Spa(前端兜底)
├── controller/          # 控制器：后台接口、公开 API、安装向导、RSS
├── models/              # 数据模型与查询（文章/专辑/标签/用户/设置/关联表）；base.go 负责选库+自动建表
├── library/             # 公共库：配置读取、安装、渲染、七牛上传、日志、字符串工具
├── middleware/          # 中间件：后台鉴权、安装校验、访问日志
├── web/                 # 前端工程（Vite + Vue 3）
│   ├── src/             # 读者端 views/reader、后台 views/admin、安装向导、api、stores
│   └── dist/            # 构建产物（被二进制内嵌，需提交）
├── configs/             # 配置模板与多语言
├── Dockerfile / docker-compose.yml
└── Makefile
```

> 前端读者端与后台共用一套 SPA：读者端走公开 `/api`，后台走 `/admin/api`，由 Vue Router 按路径区分。

## 快速开始

### 第 0 步：构建前端（需要 Node 18+）

前端产物会被二进制内嵌，编译前先构建一次（前端代码没改动时无需重复）：

```bash
make build-web          # 等价于 cd web && npm install && npm run build
```

### 方式一：SQLite（默认，最简单）

无需安装任何数据库，安装向导会在浏览器里完成初始化并创建管理员。

1. **运行**：
   ```bash
   go run .
   ```
   服务默认监听 `:8888`。读者端首页 http://localhost:8888/ ，���台 http://localhost:8888/admin/login 。

数据库文件默认落在 `./happyblog.db`（安装向导里可改）。配置会生成到 `configs/webconfig.yaml`。

### 方式二：MySQL / MariaDB

构建前端并运行后，在安装向导中数据库类型选择 **MySQL / MariaDB**，填好 `主机/端口/库名/用户/密码` 即可。
需先创建好空库：

```bash
mysql -u root -p -e "CREATE DATABASE happyblog DEFAULT CHARSET utf8mb4;"
```

> 表由程序自动迁移创建。若想手动插入管理员：密码字段为 `md5(明文密码 + salt)`，`salt` 取自 `webconfig.yaml`。

### 方式三：Docker（前端 + 后端一次构建）

`Dockerfile` 内置 Node 构建阶段，会自动构建前端再编译二进制：

```bash
docker-compose up --build
```

## 常用命令（Makefile）

```bash
make build-web   # 构建前端 SPA(Vite + Vue),产物落在 web/dist
make build       # 编译二进制(假定 web/dist 已构建)
make all         # 构建前端 + 编译二进制
make run         # 本地运行(go run .)
make test        # 运行单元测试
make vet         # 静态检查
make tidy        # 整理依赖
```

> 本地前端开发：`cd web && npm run dev` 启动 Vite 开发服务器(:5173)，已把 `/api`、`/admin/api`
> 代理到 `:8888`；后端用 `go run .` 另开进程即可热调前端。

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
