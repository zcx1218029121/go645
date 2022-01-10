package main

//佳和电表解析
//因为 广播强制连接不知道什么时候返回 佳和电表解析是读写分离的
//但是对单一的设备还是要加锁
//不建议使用
import (
	"github.com/goburrow/serial"
	"github.com/zcx1218029121/go645"
	"log"
	"sync"
	"time"
)

type Handler struct {
	go645.EventHandler
	DeviceMap sync.Map
}

func (handler *Handler) ForceOnlineResp(client go645.Client, address go645.Address, data go645.InformationElement) error {
	handler.DeviceMap.Store(address.GetStrAddress(go645.LittleEndian), nil)
	return nil
}
func (handler *Handler) ReadDataResp(client go645.Client, address go645.Address, data *go645.ReadData) error {
	log.Printf("rec read %s", data.GetFloat64Value())
	return nil
}

//佳和645协议解析
func main() {
	handler := &Handler{}
	c := go645.NewGaeaClient(go645.NewClient(go645.NewRTUClientProvider(go645.WithEnableLogger(), go645.WithSerialConfig(serial.Config{
		Address:  "/dev/ttyS1",
		BaudRate: 19200,
		DataBits: 8,
		StopBits: 1,
		Parity:   "E",
		Timeout:  time.Second * 8,
	}))))
	go func() {
		for {
			c.ForceOnline()
			//遍历联机到的表地址 发送读请求
			handler.DeviceMap.Range(func(key, value interface{}) bool {
				c.ReadAsy(go645.NewAddress(key.(string), go645.LittleEndian), 0x00_01_00_00)
				return false
			})
		}
	}()
	c.Start(handler)
}
