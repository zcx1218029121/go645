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
	p := Decode(bytes.NewBuffer(decodeString))
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
	r := ReadRequest("610100000000", 0x00_01_00_00, c)
	bf := bytes.NewBuffer(data)
	_ = r.Encode(bf)

	decodeString, _ := hex.DecodeString(str)
	p2 := Decode(bytes.NewBuffer(decodeString))
	p := Decode(bf)


	if !p.Control.IsState(Read)  {
		t.Errorf("状态解析错误")
	}

	if p.Data.dataType != p2.Data.dataType {
		t.Errorf("数据项解析错误")
	}
	if p.CS != p2.CS {
		t.Errorf("校验码解析错误")
	}

}
func TestSend(t *testing.T) {
	str := "68610100000000681104333334331416"
	data := make([]byte, 0)
	c := NewControl()
	c.SetState(Read)
	r := ReadRequest("610100000000", 0x00_01_00_00, c)
	bf := bytes.NewBuffer(data)
	_ = r.Encode(bf)
	p := Decode(bf)
	decodeString, _ := hex.DecodeString(str)
	p2 := Decode(bytes.NewBuffer(decodeString))
	if p2.Address.strValue != p.Address.strValue {
		t.Errorf("地址错误")
	}
	if p.CS != p2.CS {
		t.Errorf("校验码错误")
	}
	print(GetHex(ReadRequest("218231000988", 0x00_01_00_00, c)))
}
