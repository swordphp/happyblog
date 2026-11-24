.PHONY: all build build-web run test vet tidy clean

# 一键:先构建前端(Vite 产物内嵌),再编译二进制
all: build-web build

# 构建前端 SPA(读者端 + 管理端),产物落在 web/dist 供 Go 内嵌。
# 需要本机安装 Node(建议 18+)。
build-web:
	cd web && npm install && npm run build

# 编译二进制(假定 web/dist 已经构建好)
build:
	go build -o happyblog .

# 本地运行
run:
	go run .

# 运行单元测试
test:
	go test ./...

# 静态检查
vet:
	go vet ./...

# 整理依赖
tidy:
	go mod tidy

clean:
	rm -f happyblog happyblog.exe
