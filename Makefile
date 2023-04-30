.PHONY: all build run gotool clean help

BINARY="bluebell-linux"

all: gotool build

build:
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ${BINARY}

run:
	@go run ./

gotool:
	go fmt ./ && go vet ./

clean:
	@if [ -f ${BINARY} ];then rm ${BINARY};fi

help:
	@echo "make - 格式化Go代码，并编译生成二进制文件"
	@echo "make build - 编译Go代码， 生成二进制文件"
	@echo "make run - 直接运行Go代码"
	@echo "make clean - 移除二进制文件"
	@echo "make gotool - 运行Go工具 'fmt' 和 'vet'"