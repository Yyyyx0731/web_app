.PHONY: all build run gotool clean help

BINARY="web_app"

all: gotool build

build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ${BINARY}
	#不在句首加艾特符，在执行“make xx”命令后先打印具体命令；加了则不打印

run:
	@go run ./

gotool:
	go fmt ./ #格式化代码
	go vet ./ #做一些检查

clean: #如果当前下面有文件，就删除
	@if [ -f ${BINARY} ] ; then rm ${BINARY} ; fi

help: #输出一些提示
	@echo "make - 格式化 Go 代码, 并编译生成二进制文件"
	@echo "make build - 编译 Go 代码, 生成二进制文件"
	@echo "make run - 直接运行 Go 代码"
	@echo "make clean - 移除二进制文件和 vim swap files"
	@echo "make gotool - 运行 Go 工具 'fmt' and 'vet'"
