package go645

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
)

type ControlType byte

const (
	IsSlave  ControlType = 1 << 7
	SlaveErr ControlType = 1 << 6
	HasNext  ControlType = 1 << 5
	//Retain 保留
	Retain ControlType = 0b00000
	//Broadcast 广播校时
	Broadcast ControlType = 0b01000
	// ReadNext 读后续10001
	ReadNext ControlType = 0b10010
	//ReadAddress 读通讯地址
	ReadAddress ControlType = 0b10011
	//Write 写数据
	Write ControlType = 0b10100
	//WriteAddress 读通讯地址
	WriteAddress ControlType = 0b10101
	//ToChangeCommunicationRate 更改通讯速率
	ToChangeCommunicationRate ControlType = 0b10111
	Freeze                    ControlType = 0b10110
	//PassWord 修改密码
	PassWord       ControlType = 0b11000
	ResetMaxDemand ControlType = 0b11001
	//ResetEM 电表清零
	ResetEM    ControlType = 0b11010
	ResetEvent ControlType = 0b11011
	//Read 读
	Read ControlType = 0b10001
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

//SetStates 批量设置状态
func (c *Control) SetStates(state ...ControlType) {
	for _, s := range state {
		c.Data = c.Data | s
	}
}
func (c *Control) IsState(state ControlType) bool {
	return (c.Data & state) == state
}

//IsStates 判断控制域
func (c *Control) IsStates(state ...ControlType) bool {
	for _, s := range state {
		if !c.IsState(s) {
			return false
		}
	}
	return true
}

func (c *Control) Reset() {
	c.Data = 0
}
func (c *Control) getLen() uint16 {
	return 1
}

func (c *Control) Encode(buffer *bytes.Buffer) error {
	if err := binary.Write(buffer, binary.BigEndian, c.Data); err != nil {
		s := fmt.Sprintf("Control , %v", err)
		return errors.New(s)
	}
	return nil
}
