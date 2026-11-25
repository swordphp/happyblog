# 部署 / Deploy

单文件二进制可以直接跑,也可以用下面任意一种托管方式。程序监听端口默认 `:8888`,若设置了 `PORT` 环境变量则用它(PaaS 会自动注入)。

## 1. systemd(自己的 VPS / 服务器)

见 [happyblog.service](happyblog.service),文件顶部有完整安装步骤。要点:

```bash
sudo useradd --system --home /opt/happyblog --shell /usr/sbin/nologin happyblog
sudo mkdir -p /opt/happyblog && sudo cp happyblog /opt/happyblog/
sudo chown -R happyblog:happyblog /opt/happyblog
sudo cp deploy/happyblog.service /etc/systemd/system/
sudo systemctl daemon-reload && sudo systemctl enable --now happyblog
```

数据全部落在 `/opt/happyblog`。生产环境建议前面挂 Nginx/Caddy 反代 + HTTPS。

## 2. Render

仓库里有 [`render.yaml`](../render.yaml) Blueprint,点 README 的 **Deploy to Render** 按钮即可。

- 用 Dockerfile 构建,工作目录 `/data` 挂在持久盘上。
- **持久盘需要付费实例(starter+)**;free 实例没有磁盘,重新部署会丢数据,只适合试跑。
- 端口走 Render 注入的 `$PORT`,无需手动配。

## 3. Railway

仓库里有 [`railway.json`](../railway.json)(Dockerfile 构建 + 启动命令)。Railway 的一键按钮需要你**在 Railway 后台先把本仓库发布成一个 Template**(拿到形如 `https://railway.com/deploy/xxxx` 的链接后,替换掉 README 里按钮的占位 URL)。部署时:

- 新建服务 → 从本仓库部署(会自动读 `railway.json` 用 Dockerfile 构建)。
- **加一个 Volume,挂载路径填 `/data`**(SQLite 库等数据持久化在这里)。
- 端口走 Railway 注入的 `PORT`,程序已支持。

> 说明:Render 的 `render.yaml`、Railway 的 `railway.json` 我已经配好;但**发布 Railway 模板**和**选 Render 付费档开持久盘**需要你用自己的账号在网页上点。这两步没法代劳。
