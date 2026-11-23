.PHONY: build run test vet tidy clean

# 构建二进制
build:
	go build -o happyblog .

# 本地运行(需先准备 configs/webconfig.yaml 与 MySQL)
run:
	go run main.go

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
