BUILD_TIME=`date +%FT%T%z`
GITTAG=1.0
OUTPUT=./out/zoom
#运行信息
LDFLAGS=-ldflags "-X main.GitTag=${GITTAG} -X main.BuildTime=${BUILD_TIME} "
ARM=arm
MAC= mac
WIN =win
all: GOOS=linux GOARCH=arm go build  main.go

${ARM}:
	@echo "============ [${ARM}] 开始编译  =========="
	GOOS=linux GOARCH=arm go build -o ${OUTPUT}_${arm}_amd64
	@echo "============[${ARM}] 编译结束============"
${MAC}:
	@echo "============ [${MAC}] 开始编译  ============"
	GOOS=darwin GOARCH=amd64 go build -o ${OUTPUT}_${MAC}_amd64 ${LDFLAGS} main.go
	@echo "============ [${MAC}] 编译结束============"
${WIN}:
	@echo "============ [${WIN}] 开始编译  ============"
	GOOS=windows GOARCH=amd64 go build -o ${OUTPUT}_${WIN}_amd64 ${LDFLAGS} main.go
	@echo "============ [${WIN}] 编译结束  ============"
wasm:
	# 因为串口通讯用的是Termios所以打包不会成功的
	echo "编译到浏览器"
	GOARCH=wasm GOOS=js go build -o ${OUTPUT}.wasm ./cmd/main.go