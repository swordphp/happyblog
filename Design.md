# HappyBlog 设计文档

> 一个用 Go 编写的开源博客系统。目标:**下载一个二进制文件、浏览器里完成安装即可运行**,
> 前端(读者端 + 管理后台)内嵌进二进制,默认使用免安装的 SQLite,也支持 MySQL / MariaDB。

## 1. 项目目的

- 练手 Go 工程化(Gin + GORM)并产出一个可实际使用的开源博客。
- 面向自部署用户:**开箱即用**,不依赖外置数据库、不要求用户手写配置文件。
- 前后端分离的开发体验 + 单文件分发的部署体验。

## 2. 总体架构

```
          浏览器
            │
   ┌────────┴─────────┐
   │  单个 Go 二进制   │
   │  (Gin HTTP :8888) │
   │                   │
   │  ├─ /            读者端 SPA (Vue)        ┐
   │  ├─ /admin       管理端 SPA (Vue+ElementPlus) ├ 同一套 SPA,go:embed 内嵌 web/dist
   │  ├─ /spa/*       前端静态资源            ┘
   │  ├─ /api/*       读者端公开接口 (只读)
   │  ├─ /admin/api/* 管理端接口 (会话鉴权)
   │  ├─ /api/install 安装向导接口
   │  └─ /assets/*    本地上传的图片
   │                   │
   │  GORM ── SQLite(默认,纯 Go,无 CGO) / MySQL
   └───────────────────┘
```

- **一套 SPA 覆盖读者端与管理端**,由 Vue Router 按路径区分;构建产物 `web/dist` 通过 `go:embed` 打进二进制。
- 后端只提供 JSON 接口 + 托管内嵌的前端资源,不再做服务端模板渲染。
- 启动流程解耦:不在包级 `init()` 里连库,而是由 `main` / 安装向导在配置就绪后显式 `Bootstrap`。

## 3. 技术栈

| 层 | 选型 |
| --- | --- |
| Web 框架 | Go 1.25 + [Gin](https://github.com/gin-gonic/gin) |
| ORM | GORM v2 |
| 数据库 | SQLite(`glebarez/sqlite`,纯 Go 免 CGO,默认)/ MySQL / MariaDB |
| 配置 | Viper,配置文件 `configs/webconfig.yaml`(由安装向导生成) |
| 会话 | `gin-contrib/sessions`(cookie store)+ `patrickmn/go-cache` 缓存用户信息 |
| 前端 | Vite + Vue 3 + Vue Router + Pinia |
| 后台 UI | Element Plus |
| 国际化 | vue-i18n(中/英,localStorage 持久化) |
| Markdown | 前端 `marked`(读者端展示 + 后台编辑实时预览) |
| 图片存储 | 本地磁盘(默认)/ 七牛云 OSS,由 `storage.driver` 切换 |
| RSS | `gorilla/feeds` |
| 日志 | logrus + file-rotatelogs |

## 4. 目录结构

```
happyblog/
├── main.go            入口:装载各路由(Static/Install/Admin/Api/Spa)
├── embed.go           go:embed web/dist(内嵌前端产物)
├── router/            路由:Static(上传文件)/Admin(后台接口)/Api(公开接口)/Spa(SPA 兜底)
├── controller/        控制器 + response.go(统一返回体 respOK/respFail)
├── models/            数据模型与查询;base.go 负责选库、自动建表、seed 管理员/示例文章
├── library/           配置读取、安装(install.go)、上传、日志、字符串工具
├── middleware/        后台鉴权、安装校验、访问日志
├── web/               前端工程(Vite + Vue 3)
│   ├── src/
│   │   ├── views/reader   读者端页面
│   │   ├── views/admin    后台页面
│   │   ├── layouts        ReaderLayout / AdminLayout
│   │   ├── i18n           zh.js / en.js / index.js
│   │   ├── api            统一 axios 接口层
│   │   └── stores         Pinia(auth / site)
│   └── dist/          构建产物(被二进制内嵌,需提交)
├── configs/           配置模板与多语言(历史遗留的后台菜单翻译)
├── Dockerfile / docker-compose.yml / Makefile / .github/workflows
```

## 5. 安装与启动流程

1. 二进制启动时检查 `configs/webconfig.yaml` 是否存在(`IsInstalled`)。
2. **未安装**:直接起服务;前端 SPA 调 `/api/install/status` 得知未安装,进入**安装向导**。
3. 安装向导(`POST /api/install`)写配置(自动生成 salt/sessionsecret)→ `LoadConfig` →
   `Bootstrap`(连库、建表、建目录、初始化存储与会话缓存)→ **seed 默认管理员 + 两篇示例文章**。
   全程无需重启。
4. **已安装**:`main` 启动时即 `Bootstrap`,正常进入读者端/后台。

数据库表由 GORM `AutoMigrate` 自动创建;SQLite 数据落在单文件(默认 `./happyblog.db`)。

## 6. 数据模型

| 模型 | 表名 | 说明 |
| --- | --- | --- |
| User | `happyblog_tblUser` | 用户/管理员,密码存 `md5(明文 + salt)` |
| Article | `happyblog_tblArticle` | 文章,含发布状态、独立页面标记、简介、别名 |
| Tag | `happyblog_tblTag` | 标签 |
| Album | `happyblog_tblAlbum` | 专辑,含公开标记与文章计数 |
| Setting | `happyblog_settings` | 按分组的键值配置项 |
| RelationArticleAlbums | `happyblog_tblAlbumRe` | 文章 ↔ 专辑 关联 |
| RelationArticleTags | `happyblog_tblTagRe` | 文章 ↔ 标签 关联 |

- 文章与专辑是「一对多归属」;文章与标签是「多对多」。
- 专辑的文章计数在文章新增/修改/删除时同步维护。
- 站点基础信息(名称/简介/地址)存在**配置文件**而非 Setting 表,可在后台实时修改并热更新内存。

## 7. 接口设计

统一返回体(`controller/response.go` 的 `respOK` / `respFail`):

```json
{ "errNo": 0, "errMsg": "success", "data": {} }
```

`errNo` 约定:`0` 成功,`10` 参数错误,`401` 未登录,`404` 不存在,`500` 服务端错误。

- **安装**(未安装可用):`GET /api/install/status`、`POST /api/install`
- **读者端公开只读**(未安装返回 409):`/api/siteinfo`、`/api/article/list/:page`、
  `/api/article/info/:id`、`/api/tags`、`/api/tag/:id/articles`、`/api/albums/:page`、
  `/api/album/:id`、`/api/rss`
- **管理端**(会话鉴权,`/admin/api`):登录/登出/me、文章列表与编辑/保存/删除、
  标签增删查、专辑增改查、用户列表、配置项增删改查、**网站基础信息读写(`/siteconfig`)**、图片上传

## 8. 前端设计

- **单 SPA 双区**:读者端(轻量、纯 CSS)与后台(Element Plus)共用一个 Vite 工程;
  后台路由异步分包,Element Plus 全局注册。
- **路由守卫**:访问需登录的后台路由时,先确认已安装(否则跳安装向导)、再确认登录态。
- **状态**:Pinia 管理登录用户(auth)与站点信息(site)。
- **接口层**:`src/api` 统一封装(axios + 同源 cookie),拆包并按 `errNo` 抛错。
- **国际化**:vue-i18n,中/英全站可切换,右上角切换器,持久化到 localStorage,
  Element Plus 自带文案随之切换;仅覆盖界面文案,不翻译文章内容。
- **读者端布局**:页脚粘性贴底;内容容器小屏保持 820px,大屏随视口放宽至最高 1280px。

## 9. 构建与部署

- 开发期需要 Node(建议 18+)构建前端;最终产物是**自带前端的单个 Go 二进制**。
- 常用命令:`make build-web`(构建前端)、`make build`(编译二进制)、`make all`(两者)。
- 本地前端热调:`cd web && npm run dev`(:5173,已把 `/api`、`/admin/api` 代理到 :8888)。
- 因 SQLite 为纯 Go 驱动(`CGO_ENABLED=0`),可轻松交叉编译多平台单文件。
- `Dockerfile` 内置 Node 构建阶段 + Go 构建阶段;`.github/workflows/ci.yml` 跑前端构建 + vet/build/test。

## 10. 关键设计决策

- **单二进制 + 内嵌前端**:为「下载即跑」服务;`web/dist` 需提交,保证无 Node 环境也能 `go build`。
- **SQLite 默认 + 纯 Go 驱动**:免外置数据库、免 CGO,降低自部署门槛。
- **图形化安装向导替代手写配置**:安装信息经表单提交,自动生成密钥与目录。
- **启动解耦(Bootstrap/Connect)**:未安装也能起服务进入向导,安装后无需重启。
- **统一返回体**:所有接口走 `respOK`/`respFail`,前后端契约一致。
