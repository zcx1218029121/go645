# go dlt645-2007

用go语言实现的dlt645解析
```shell
    go get github.com/zcx1218029121/go645
```

1. 读请求
```go
ReadRequest(NewAddress("610100000000", BigEndian), 0x00_01_00_00)
```