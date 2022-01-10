package go645

import (
	"bytes"
	"encoding/hex"
	"log"
)

type GaeaUdpClient struct {
	EHandler
	Client
	Rev      chan *Protocol
	SendChan chan *Protocol
}

func NewGaeaClient(client Client, opts ...Option) *GaeaUdpClient {
	return &GaeaUdpClient{
		Client:   client,
		Rev:      make(chan *Protocol),
		SendChan: make(chan *Protocol),
	}
}
func (c *GaeaUdpClient) OnRecv(p *Protocol) {

	var err error
	if p.Control.Data == 0x8a {
		err = c.EHandler.ForceOnlineResp(c, p.Address, p.Data)
	}
	if p.Control.IsState(Read) {
		print(p.Data.(*ReadData))
		err = c.EHandler.ReadDataResp(c, p.Address, p.Data.(*ReadData))
	}
	if err != nil {
		log.Print(err.Error())
	}
}
func (c *GaeaUdpClient) Start(handler EHandler) {
	c.EHandler = handler
	err := c.Connect()
	if err != nil {
		panic(err)
	}
	go func() {
		for value := range c.SendChan {
			err := c.Send(value)
			if err != nil {
				log.Print(err.Error())
			}
		}
	}()
	go func() {
		for value := range c.Rev {
			c.OnRecv(value)
		}
	}()
	c.Loop()
}

//SendWithChan 通道写
func (c *GaeaUdpClient) SendWithChan(protocol *Protocol) {
	c.SendChan <- protocol

}
func (c *GaeaUdpClient) Loop() {
	log.Print(c == nil)
	for {
		frame, err := c.ReadRawFrame()
		if err != nil {
			return
		}
		p, err := Decode(bytes.NewBuffer(frame))
		if err != nil {
			log.Print(err.Error())
		}
		log.Printf("rec hex[%s]", hex.EncodeToString(frame))
		c.Rev <- p
	}

}

//ForceOnline 强制联机 佳和电表
func (c *GaeaUdpClient) ForceOnline() {
	c.SendWithChan(NewProtocol(NewAddress(BroadcastAddress, BigEndian), NullData{}, NewControlValue(0x0a)))
}

func (c *GaeaUdpClient) ReadAsy(address Address, itemCode int32) {
	c.SendWithChan(ReadRequest(address, itemCode))
}
