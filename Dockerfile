# ---- build stage ----
FROM golang:1.22-alpine AS builder
WORKDIR /src
ENV GOPROXY=https://goproxy.cn,direct
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/happyblog .

# ---- run stage ----
FROM alpine:3.20
WORKDIR /app
RUN apk add --no-cache tzdata ca-certificates
COPY --from=builder /app/happyblog /app/happyblog
# 需要挂载/拷贝以下目录(配置、模板、静态资源)
COPY templates/ /app/templates/
COPY assets/ /app/assets/
COPY configs/ /app/configs/
EXPOSE 8888
CMD ["/app/happyblog"]
