package main

import (
	"bytes"
	"flag"
	"github.com/goburrow/serial"
	"github.com/zcx1218029121/go645"
	"log"
	"sync"
	"time"
)

var device sync.Map
var mu sync.Mutex

//特殊电表解析 同步
func main() {
	var size int
	flag.IntVar(&size, "s", 1, "电表数量")
	flag.Parse()
	p := go645.NewRTUClientProvider(go645.WithSerialConfig(serial.Config{Address: "/dev/ttyS1", BaudRate: 19200, DataBits: 8, StopBits: 1, Parity: "E", Timeout: time.Second * 40}), go645.WithEnableLogger())
	c := go645.NewClient(p)
	c.Connect()
	defer c.Close()

	forceOnline(c)

	go func() {
		time.Sleep(1 * time.Minute)
		forceOnline(c)
	}()

	for {
		//如果对扫描速度不是要求很高 需要重新扫描
		time.Sleep(50 * time.Millisecond)
		device.Range(func(key, value interface{}) bool {
			read, _, err := c.Read(go645.NewAddress(key.(string), go645.LittleEndian), 0x00_01_00_00)
			if err == nil {
				log.Printf("rec %s", read.GetValue())
			}
			return false
		})
	}

}
func forceOnline(c go645.Client) {
	c.Broadcast(go645.NullData{}, *go645.NewControlValue(0x0a))
	for {
		frame, err := c.ReadRawFrame()
		if err != nil {
			log.Printf(err.Error())
			return
		}
		if frame != nil && len(frame) > 10 {
			p, err := go645.Decode(bytes.NewBuffer(frame))
			if err != nil {
				log.Printf(err.Error())
				return
			}
			device.Store(p.Address.GetStrAddress(go645.LittleEndian), nil)
		}
	}
}
