package go645

import (
	"bytes"
	"encoding/hex"
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
	if !p.Control.IsState(IsSlave) || !p.Control.IsState(Read) {
		t.Errorf("状态解析错误")
	}
	println(p.Data.GetFloat64Value())

	if p.Data.dataType != [4]byte{0, 1, 0, 0} {
		t.Errorf("数据项解析错误")
	}

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
	Assert("状态解析错误", func() bool { return p.Control.IsState(Read) }, t)
	AssertEquest("数据项解析错误", p.Data.dataType, p2.Data.dataType, t)
	AssertEquest("校验码解析错误", p.CS, p2.CS, t)
}
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
