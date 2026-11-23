package main

import "embed"

// 后台 SPA 的构建产物(Vite + Vue)。开发阶段用 node 构建,产物内嵌进二进制。
// 使用 all: 以包含可能的点文件;web/dist 必须存在(含 index.html)才能编译。
//
//go:embed all:web/dist
var distFS embed.FS

// 读者端服务端渲染模板
//
//go:embed templates/front_*.tmpl
var frontTemplatesFS embed.FS

// 读者端内置静态资源(首次运行时释放到 ./assets)
//
//go:embed assets/style.css assets/script.js assets/favicon.ico
var readerAssetsFS embed.FS
