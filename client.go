package go645

import (
	"bytes"
	"sync"
)

var (
	// check implements Client interface.
	_ Client = (*client)(nil)
)

type client struct {
	ClientProvider
	mu sync.Mutex
}

func (c *client) Read(address Address, itemCode int32) (*ReadData, bool, error) {
	resp, err := c.ClientProvider.Send(ReadRequest(address, itemCode))
	if err != nil {
		return nil, false, err
	}
	decode, err := Decode(bytes.NewBuffer(resp))
	if err != nil {
		return nil, false, err
	}
	return decode.Data.(*ReadData), decode.Control.IsState(HasNext), err
}

//Broadcast 设备广播
func (c *client) Broadcast(p *Protocol) error {
	p.Address = NewAddress(BroadcastAddress, LittleEndian)
	var err error
	bf := bytes.NewBuffer(make([]byte, 0))
	err = p.Encode(bf)
	if err != nil {
		return err
	}
	return nil
}

func (c *client) ReadWithBlock(address Address, data ReadRequestData) (*Protocol, error) {
	resp, err := c.ClientProvider.Send(ReadRequestWithBlock(address, data))
	if err != nil {
		return nil, err
	}
	return Decode(bytes.NewBuffer(resp))
}

// Option custom option
type Option func(c *client)

func NewClient(p ClientProvider, opts ...Option) Client {
	c := &client{ClientProvider: p}
	for _, opt := range opts {
		opt(c)
	}
	return c
}
