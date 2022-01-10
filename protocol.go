package go645

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"
	"time"
)

const (
	LittleEndian     Order = false
	BigEndian        Order = true
	Start                  = 0x68
	End                    = 0x16
	HeadLen                = 1 + 6 + 1
	BroadcastAddress       = "999999999999"
)

var (
	_ InformationElement = (*Address)(nil)

	_ InformationElement = (*Protocol)(nil)

	_ InformationElement = (*ReadData)(nil)

	_ InformationElement = (*ReadRequestData)(nil)

	_ InformationElement = (*Exception)(nil)
	_ InformationElement = (*YearDateTimeS)(nil)
	_ InformationElement = (*DateTime)(nil)
	_ InformationElement = (*NullData)(nil)
)

type (
	//Order 地址大小端
	Order              bool
	InformationElement interface {
		Encode(buffer *bytes.Buffer) error
		GetLen() byte
	}
	NullData struct {
	}

	YearDateTimeS struct {
		ss    byte
		mm    byte
		hh    byte
		day   byte
		month byte
		year  byte
	}

	DateTime struct {
		mm    byte
		hh    byte
		day   byte
		month byte
	}

	//Address 表计地址
	Address struct {
		value    []byte
		strValue string
	}

	//Protocol 协议
	Protocol struct {
		//Start 645协议起始符号
		Start byte
		//设备地址 6个字节的BCD
		Address Address
		//Start  645协议起始符号 标志报文头结束
		Start2 byte
		//Control 控制域
		Control *Control
		//Control 数据长度
		DataLength byte
		//Control 数据抽象
		Data InformationElement
		//CS 校验和
		CS byte
		//End 0x16
		End byte
	}
)

func (n NullData) Encode(buffer *bytes.Buffer) error {
	return nil
}

func (n NullData) GetLen() byte {
	return 0
}

func (t DateTime) Encode(buffer *bytes.Buffer) error {
	var err error
	err = binary.Write(buffer, binary.BigEndian, t.mm)
	if err != nil {
		return err
	}
	err = binary.Write(buffer, binary.BigEndian, t.hh)
	if err != nil {
		return err
	}
	err = binary.Write(buffer, binary.BigEndian, t.day)
	if err != nil {
		return err
	}
	err = binary.Write(buffer, binary.BigEndian, t.month)
	if err != nil {
		return err
	}
	return err
}

func (t DateTime) GetLen() byte {
	panic("implement me")
}

func NewTimeS() *YearDateTimeS {
	t := time.Now()
	return &YearDateTimeS{ss: byte(t.Second()), mm: byte(t.Minute()), hh: byte(t.Hour()), day: byte(t.Day()), year: byte(t.Year()), month: byte(t.Month())}
}
func (t YearDateTimeS) Encode(buffer *bytes.Buffer) error {
	var err error
	err = binary.Write(buffer, binary.BigEndian, t.ss)
	if err != nil {
		return err
	}
	err = binary.Write(buffer, binary.BigEndian, t.mm)
	if err != nil {
		return err
	}
	err = binary.Write(buffer, binary.BigEndian, t.hh)
	if err != nil {
		return err
	}
	err = binary.Write(buffer, binary.BigEndian, t.day)
	if err != nil {
		return err
	}
	err = binary.Write(buffer, binary.BigEndian, t.month)
	if err != nil {
		return err
	}
	err = binary.Write(buffer, binary.BigEndian, t.year)
	if err != nil {
		return err
	}
	return err

}

func (t YearDateTimeS) GetLen() byte {
	return 6
}

// NewAddress ，构建设备地址
// 参数：
//      address ： 设备地址
//      order ： 大小端表示
// 返回值：
//      *Address 设备地址
func NewAddress(address string, order Order) Address {
	value := Number2bcd(address)
	if !order {
		for i, j := 0, len(value)-1; i < j; i, j = i+1, j-1 {
			value[i], value[j] = value[j], value[i]
		}
	}

	return Address{value: value, strValue: address}
}

func NewReadData(dataType int32, value string) ReadData {
	return ReadData{dataType: Int2bytes(dataType), rawValue: value}
}

func NewProtocol(address Address, data InformationElement, control *Control) *Protocol {
	return &Protocol{
		Start:      Start,
		Start2:     Start,
		End:        End,
		Address:    address,
		Control:    control,
		DataLength: data.GetLen(),
		Data:       data,
	}
}

// Encode ，协议解码
// 参数：
//      buffer ： 字节码缓冲
// 返回值：
//      error 解码异常
func (a Address) Encode(buffer *bytes.Buffer) error {
	return binary.Write(buffer, binary.BigEndian, a.value)
}

func (a Address) GetStrAddress(order Order) string {
	if !order {
		temp := make([]byte, len(a.value))
		for i, j := 0, len(a.value)-1; i < j; i, j = i+1, j-1 {
			temp[i], temp[j] = a.value[j], a.value[i]
		}
		return Bcd2Number(temp)
	}
	return a.strValue
}

func (a Address) GetLen() byte {
	return 6
}

//GetHex 返回16进制string
func GetHex(protocol *Protocol) (string, error) {
	bf := bytes.NewBuffer(make([]byte, 0))
	if err := protocol.Encode(bf); err != nil {
		return "", err
	}
	return hex.EncodeToString(bf.Bytes()), nil
}

func (p Protocol) Encode(buffer *bytes.Buffer) error {
	//计算cs 需要重写开辟字节码缓冲区
	tmp := make([]byte, 0)
	bf := bytes.NewBuffer(tmp)
	var err error
	write := func(data interface{}) {
		if err != nil {
			return
		}
		err = binary.Write(bf, binary.LittleEndian, data)
	}
	write(&p.Start)
	err = p.Address.Encode(bf)
	write(&p.Start2)
	err = p.Control.Encode(bf)
	write(&p.DataLength)
	err = p.Data.Encode(bf)
	var cs = 0
	for _, b := range bf.Bytes() {
		cs += int(b)
	}
	p.CS = byte(cs)
	write(p.CS)
	write(p.End)
	err = binary.Write(buffer, binary.LittleEndian, bf.Bytes())
	return err

}

func (p Protocol) GetLen() byte {
	if p.DataLength != 0 {
		return p.DataLength
	}
	return HeadLen + 4 + p.Data.GetLen()
}

func Bcd2Number(bcd []byte) string {
	var number string
	for _, i := range bcd {
		number += fmt.Sprintf("%02X", i)
	}
	pos := strings.LastIndex(number, "F")
	if pos == 8 {
		return "0"
	}
	return number[pos+1:]
}

func Number2bcd(number string) []byte {
	var rNumber = number
	for i := 0; i < 8-len(number); i++ {
		rNumber = "f" + rNumber
	}
	bcd := Hex2Byte(rNumber)
	return bcd
}

func Hex2Byte(str string) []byte {
	slen := len(str)
	bHex := make([]byte, len(str)/2)
	ii := 0
	for i := 0; i < len(str); i = i + 2 {
		if slen != 1 {
			ss := string(str[i]) + string(str[i+1])
			bt, _ := strconv.ParseInt(ss, 16, 32)
			bHex[ii] = byte(bt)
			ii = ii + 1
			slen = slen - 2
		}
	}
	return bHex
}

func Int2bytes(data int32) []byte {
	var b3 = make([]byte, 4)
	b3[0] = uint8(data)
	b3[1] = uint8(data >> 8)
	b3[2] = uint8(data >> 16)
	b3[3] = uint8(data >> 24)
	return b3
}

func WriteWithOffSet(buffer *bytes.Buffer, data byte) error {
	return binary.Write(buffer, binary.LittleEndian, data+0x33)
}
