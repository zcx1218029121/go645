package go645

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
)

type ControlType byte

const (
	IsSlave     ControlType = 1 << 7
	SlaveOk     ControlType = 1 << 6
	hasNext     ControlType = 1 << 5
	Retain      ControlType = 0b00000
	Broadcast   ControlType = 0b01000
	ReadNext    ControlType = 0b10010
	ReadAddress ControlType = 0b10011
	Read        ControlType = 0b10001
)

type Control struct {
	Data ControlType
}

func DecodeControl(buffer *bytes.Buffer) (*Control, error) {
	c := new(Control)
	if err := binary.Read(buffer, binary.LittleEndian, &c.Data); err != nil {
		return nil, err
	}
	return c, nil
}
func NewControl() *Control {
	return &Control{Data: 0}
}

func (c *Control) SetState(state ControlType) {
	c.Data = c.Data | state
}
func (c *Control) IsState(state ControlType) bool {
	return (c.Data & state) == state
}
func (c *Control) Reset() {
	c.Data = 0
}
func (c *Control) getLen() uint16 {
	return 1
}

func (c *Control) Encode(buffer *bytes.Buffer) error {
	if err := binary.Write(buffer, binary.BigEndian, c.Data); err != nil {
		s := fmt.Sprintf("Pack version error , %v", err)
		return errors.New(s)
	}
	return nil
}
