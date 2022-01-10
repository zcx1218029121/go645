package main

import (
	"encoding/hex"
	"github.com/goburrow/serial"
	"github.com/zcx1218029121/go645"
	"log"
	"time"
)

func main() {
	c := go645.NewClient(go645.NewRTUClientProvider(go645.WithEnableLogger(), go645.WithSerialConfig(serial.Config{
		Address:  "COM2",
		BaudRate: 19200,
		DataBits: 8,
		StopBits: 1,
		Parity:   "E",
		Timeout:  time.Second * 8,
	})))

	go func() {
		time.Sleep(1 * time.Minute)
	}()

	for {
		time.Sleep(time.Second)
		//广播校时
		c.Broadcast(go645.NewTimeS(), *go645.NewControl())
		control := *go645.NewControl()
		control.SetState(go645.Freeze)
		c.Broadcast(go645.NewTimeS(), control)
		pr, hasNext, err := c.Read(go645.NewAddress("59050008193a", go645.BigEndian), 0x00_01_00_00)
		if err != nil {
			log.Print(err.Error())
		} else {
			if hasNext {
				println(hex.EncodeToString(pr.GetDataType()))
				println(pr.GetValue())
			} else {
				println(hex.EncodeToString(pr.GetDataType()))
				println(pr.GetValue())
			}
		}

	}
}
