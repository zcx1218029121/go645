# 快速开始

## 安装

```shell
    go get github.com/zcx1218029121/go645
```

### 创建串口客户端

```go
    p := go645.NewRTUClientProvider(
go645.WithSerialConfig(serial.Config{
Address:  "COM2",
BaudRate: b,
DataBits: 8,
StopBits: 1,
Parity:   "E",
Timeout:  time.Second * 30}))
```

### 创建串口客户端并打印日志

```go
p := go645.NewRTUClientProvider(
go645.WithSerialConfig(serial.Config{
Address:  "COM2",
BaudRate: b,
DataBits: 8,
StopBits: 1,
Parity:   "E",
Timeout:  time.Second * 30}),
go645.WithEnableLogger())
```

### 连接串口

```go
    err := c.Connect()
if err != nil {
panic(err)
}
```

### 自定义引导词

```go
    var _ go645.PrefixHandler = (*Handler)(nil)

    type Handler struct {
    }

    func (h Handler) EncodePrefix(buffer *bytes.Buffer) error {
    // 写入引导词
        buffer.Write([]byte{0xfe, 0xfe, 0xfe, 0xfe})
        return nil
    }
    // 写入引导解码 一般来说 在这里的引导词都会被丢弃
    func (h Handler) DecodePrefix(reader io.Reader) ([]byte, error) {
        fe := make([]byte, 4)
        _, err := io.ReadAtLeast(reader, fe, 4)
        if err != nil {
            return nil, err
        }
        return fe, err
    }
    p := go645.NewRTUClientProvider(
        go645.WithSerialConfig(serial.Config{
            Address:  "COM2",
            BaudRate: b,
            DataBits: 8,
            StopBits: 1,
            Parity:   "E",
        Timeout:  time.Second * 30}),
        go645.WithEnableLogger(),go645.WithPrefixHandler(&Handler{}))
```
### 发送读请求
1. hasNext 为是否有后续的Frame 标识
2. read 为读请求响应
```go
    read, hasNext, err := c.Read(go645.NewAddress("000703200136", go645.LittleEndian), 0x00_01_00_00)
    if err != nil{
    	log.Printf("rec %s \n", read.GetValue())
    	if hasNext{
            log.Printf("有后续Frame  \n", read.GetValue())	
        }
        
    }
    
```

### 发送广播请求
> 标准的广播请求不会有后续响应如果为特殊响应请手动调用 frame, err := c.ReadRawFrame()
#### 正常广播(无需返回值)
```go
    err:=c.Broadcast(go645.NullData{}, *go645.NewControlValue(0x0a))
    if err != nil{
    	log.Printf("广播发送失败")
    }
```
#### 特殊广播(需返回值)
```go
    err:=c.Broadcast(go645.NullData{}, *go645.NewControlValue(0x0a))
    if err != nil{
    	log.Printf("广播发送失败")
    }
    frame, err := c.ReadRawFrame()
    if err != nil {
    log.Printf(err.Error())
    return
    }
```