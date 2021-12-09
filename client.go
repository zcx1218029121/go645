package go645

import (
	"bytes"
)

var (
	// check implements Client interface.
	_ Client = (*client)(nil)
)

type client struct {
	ClientProvider
}

func (c client) ReadWithBlock(address Address, data ReadRequestData) (*Protocol, error) {
	resp, err := c.ClientProvider.Send(ReadRequestWithBlock(address, data))
	if err != nil {
		return nil, err
	}
	return Decode(bytes.NewBuffer(resp))
}

func (c client) Read(address Address, itemCode int32) (*ReadData, error) {
	resp, err := c.ClientProvider.Send(ReadRequest(address, itemCode))
	if err != nil {
		return nil, err
	}
	decode, err := Decode(bytes.NewBuffer(resp))
	if err != nil {
		return nil, err
	} else {
		return decode.Data.(*ReadData), nil
	}
}

// Option custom option
type Option func(c *client)

func NewClient(p ClientProvider, opts ...Option) Client {
	c := &client{p}
	for _, opt := range opts {
		opt(c)
	}
	return c
}
