package go645

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"testing"
)

//TestRead 测试上报请求
func TestDecode(t *testing.T) {
	str := "681401003182216891083333343339333333f116"
	decodeString, err := hex.DecodeString(str)
	if err != nil {
		return
	}
	p, _ := Decode(bytes.NewBuffer(decodeString))
	if p.Address.strValue != "140100318221" {
		t.Errorf("地址解析错误")
	}
	if p.Address.GetLen() != 6 {
		t.Errorf("长度错误")
	}
	if !p.Control.IsState(IsSlave) || !p.Control.IsState(Read) {
		t.Errorf("状态解析错误")
	}
	println(p.Data.(*ReadData).GetFloat64Value())
	if p.GetLen() != 0x08 {
		t.Errorf("长度错误")
	}

	print(GetHex(p))

}

//TestRead 测试解析读请求
func TestRead(t *testing.T) {
	str := "68610100000000681104333334331416"
	data := make([]byte, 0)
	c := NewControl()
	c.SetState(Read)
	r := ReadRequest(NewAddress("610100000000", BigEndian), 0x00_01_00_00)
	bf := bytes.NewBuffer(data)
	_ = r.Encode(bf)

	decodeString, _ := hex.DecodeString(str)
	p2, _ := Decode(bytes.NewBuffer(decodeString))
	p, _ := Decode(bf)
	p.Data.(*ReadData).GetDataTypeStr()
	p.Data.(*ReadData).GetDataType()
	Assert("状态解析错误", func() bool { return p.Control.IsState(Read) }, t)
	AssertEquest("数据项解析错误", p.Data.(*ReadData).dataType, p2.Data.(*ReadData).dataType, t)
	AssertEquest("校验码解析错误", p.CS, p2.CS, t)
}

//TestSend 测试发送
func TestSend(t *testing.T) {
	str := "68610100000000681104333334331416"
	data := make([]byte, 0)
	c := NewControl()
	c.SetState(Read)
	r := ReadRequest(NewAddress("610100000000", BigEndian), 0x00_01_00_00)
	bf := bytes.NewBuffer(data)
	_ = r.Encode(bf)
	p, _ := Decode(bf)
	decodeString, _ := hex.DecodeString(str)
	p2, _ := Decode(bytes.NewBuffer(decodeString))
	print(p.Data.(*ReadData).GetValue())
	AssertEquest("地址错误", p2.Address.strValue, p.Address.strValue, t)
	AssertEquest("校验码错误", p.CS, p2.CS, t)

}
func TestLEnd(t *testing.T) {
	str := "68610100000000681104333334331416"
	data := make([]byte, 0)
	c := NewControl()
	c.SetState(Read)
	r := ReadRequest(NewAddress("610100000000", LittleEndian), 0x00_01_00_00)
	bf := bytes.NewBuffer(data)
	_ = r.Encode(bf)
	p, _ := Decode(bf)
	decodeString, _ := hex.DecodeString(str)
	p2, _ := Decode(bytes.NewBuffer(decodeString))
	AssertEquest("地址错误", p2.Address.GetStrAddress(LittleEndian), "000000000161", t)
	AssertEquest("地址错误", p2.Address.GetStrAddress(BigEndian), "610100000000", t)
	AssertEquest("校验码错误", p.CS, p2.CS, t)
}
func Assert(msg string, assert func() bool, t *testing.T) {
	if !assert() {
		t.Errorf(msg)
	}
}
func AssertEquest(msg string, exp interface{}, act interface{}, t *testing.T) {
	Assert(msg, func() bool { return exp == act }, t)
}

func AssertState(assert func() bool, t *testing.T) {
	Assert("状态解析错误", assert, t)
}
func TestControl_IsState(t *testing.T) {
	c := new(Control)
	c.SetState(SlaveErr)
	if !c.IsStates(SlaveErr) {
		t.Errorf("设置状态错误")
	}
	c.Reset()
	if c.IsStates(SlaveErr) {
		t.Errorf("复归错误")
	}
}
func TestControl(t *testing.T) {
	c := &Control{}
	if c.getLen() != 1 {
		t.Errorf("长度错误")
	}
	c.SetState(SlaveErr)
	if !c.IsState(SlaveErr) {
		t.Errorf("设置错误")
	}
	c.Reset()
	if c.IsState(SlaveErr) {
		t.Errorf("复归错误")
	}
	c.SetStates(SlaveErr, IsSlave, HasNext, Retain, Broadcast, ReadNext, ReadAddress)
	if !c.IsStates(SlaveErr, IsSlave) {
		t.Errorf("设置错误")
	}
}
func (c *Control) TestErr(buffer *bytes.Buffer) error {

	var bf *bytes.Buffer
	r := ReadRequest(NewAddress("610100000000", LittleEndian), 0x00_01_00_00)
	_ = r.Encode(bf)

	if err := binary.Write(buffer, binary.BigEndian, c.Data); err != nil {
		s := fmt.Sprintf("Control , %v", err)
		return errors.New(s)
	}
	return nil
}
func TestReadResponse(t *testing.T) {
	rp := ReadResponse(NewAddress("610100000000", LittleEndian), 0x00_01_00_00, NewControl(), "200")

	if rp.Address.GetStrAddress(LittleEndian) != "610100000000" {
		t.Errorf("地址错误")
	}
}
