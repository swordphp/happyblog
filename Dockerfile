# ---- frontend build stage (Vite + Vue) ----
FROM node:22-alpine AS web
WORKDIR /web
COPY web/package.json ./
RUN npm install
COPY web/ ./
RUN npm run build

# ---- go build stage ----
FROM golang:1.26-alpine AS builder
WORKDIR /src
ENV GOPROXY=https://goproxy.cn,direct
COPY go.mod go.sum ./
RUN go mod download
COPY . .
# 用前端构建产物覆盖占位 dist,供 go:embed 内嵌
COPY --from=web /web/dist ./web/dist
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/happyblog .

# ---- run stage ----
FROM alpine:3.20
WORKDIR /app
RUN apk add --no-cache tzdata ca-certificates
COPY --from=builder /app/happyblog /app/happyblog
# 配置模板(安装向导会在运行时生成实际配置到 ./configs)
COPY configs/ /app/configs/
EXPOSE 8888
CMD ["/app/happyblog"]
