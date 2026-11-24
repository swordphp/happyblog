package main

import "embed"

// 前端单页应用(Vite + Vue 3)的构建产物,开发阶段用 node 构建,产物内嵌进二进制。
// 使用 all: 以包含点文件;web/dist 必须存在(含 index.html)才能编译。
// 读者端与管理端共用这一套 SPA,由前端路由区分。
//
//go:embed all:web/dist
var distFS embed.FS
