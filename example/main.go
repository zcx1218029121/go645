package main

import (
	"github.com/goburrow/serial"
	"github.com/zcx1218029121/go645"
	"log"
	"time"
)

func main() {
	c := go645.NewClient(go645.NewRTUClientProvider(go645.WithEnableLogger(), go645.WithSerialConfig(serial.Config{
		Address:  "/dev/ttyUSB3",
		BaudRate: 19200,
		DataBits: 8,
		StopBits: 1,
		Parity:   "E",
		Timeout:  time.Second * 8,
	})))

	for {
		time.Sleep(time.Second)
		pr, ok, err := c.Read(go645.NewAddress("3a2107000481", go645.LittleEndian), 0x00_01_00_00)
		if err != nil {
			log.Print(err.Error())
		} else {
			if ok {
				println("有后续")
			} else {
				log.Print(pr.GetValue())
			}

		}

	}

}
