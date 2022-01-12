package main

import (
	"bytes"
	"flag"
	"github.com/goburrow/serial"
	"github.com/zcx1218029121/go645"
	"io"
	"log"
	"time"
)

var _ go645.PrefixHandler = (*Handler)(nil)

type Handler struct {
}

func (h Handler) EncodePrefix(buffer *bytes.Buffer) error {
	// 百富电表写入的时候不需要引导词
	buffer.Write([]byte{0xfe, 0xfe, 0xfe, 0xfe})
	return nil
}

func (h Handler) DecodePrefix(reader io.Reader) ([]byte, error) {
	fe := make([]byte, 4)
	_, err := io.ReadAtLeast(reader, fe, 4)
	if err != nil {
		return nil, err
	}
	return fe, err
}

//百富电表
func main() {
	var b int
	var code int
	flag.IntVar(&b, "b", 2400, "波特率")
	flag.IntVar(&code, "c", 0x00_03_00_00, "波特率")
	flag.Parse()
	p := go645.NewRTUClientProvider(go645.WithSerialConfig(serial.Config{Address: "COM2", BaudRate: b, DataBits: 8, StopBits: 1, Parity: "E", Timeout: time.Second * 30}), go645.WithEnableLogger(), go645.WithPrefixHandler(&Handler{}))
	c := go645.NewClient(p)
	err := c.Connect()
	if err != nil {
		panic(err)
	}
	defer c.Close()
	for {
		time.Sleep(50 * time.Millisecond)
		read, _, err := c.Read(go645.NewAddress("000703200136", go645.LittleEndian), 0x00_01_00_00)
		if err == nil {
			log.Printf("rec %s", read.GetValue())
		}

	}
}
