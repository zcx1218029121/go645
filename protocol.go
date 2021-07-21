package go645

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"
)

type ByteSlice []byte

func (x ByteSlice) Len() int           { return len(x) }
func (x ByteSlice) Less(i, j int) bool { return x[i] < x[j] }
func (x ByteSlice) Swap(i, j int)      { x[i], x[j] = x[j], x[i] }

const (
	Start   = 0x68
	End     = 0x16
	HeadLen = 1 + 6 + 1
)

var _ InformationElement = (*Address)(nil)

var _ InformationElement = (*Protocol)(nil)

var _ InformationElement = (*Data)(nil)

type InformationElement interface {
	Encode(buffer *bytes.Buffer) error
	getLen() uint16
}

type Address struct {
	value    []byte
	StrValue string
}

func NewAddress(address string) *Address {
	value := Number2bcd(address)
	//反转
	for i, j := 0, len(value)-1; i < j; i, j = i+1, j-1 {
		value[i], value[j] = value[j], value[i]
	}
	return &Address{value: value, StrValue: address}
}

func (a *Address) getString() string {
	return a.StrValue
}

func (a Address) Encode(buffer *bytes.Buffer) error {
	return binary.Write(buffer, binary.BigEndian, a.value)
}

func (a Address) getLen() uint16 {
	return 6
}

type Data struct {
	//数据标识 4 个字节
	dataType [4]byte

	//原始数据
	rawValue string
}

func (data Data) GetDataType() [4]byte {
	return data.dataType
}
func (data Data) GetDataTypeStr() string {
	//需要翻转
	var a = make([]byte, 4)
	for i, j := 0, len(data.dataType)-1; i < j; i, j = i+1, j-1 {
		a[i], a[j] = data.dataType[j], data.dataType[i]
	}
	return hex.EncodeToString(a)
}

func (d *Data) GetFloat64Value() float64 {
	if d.dataType[0] == 0x00 || d.dataType[0] == 0x0c {
		value, _ := strconv.Atoi(d.rawValue)
		return float64(value) * 0.01
	} else if d.dataType[3] == 0x02 {
		value, _ := strconv.Atoi(d.rawValue)
		return float64(value) * 0.0001
	}
	return 0
}

func (d *Data) GetValue() string {
	return d.rawValue
}

func (d Data) Encode(buffer *bytes.Buffer) error {
	//写入数据域 已经反转过了
	for index, b := range d.dataType {
		d.dataType[index] = b + 0x33
	}
	_ = binary.Write(buffer, binary.LittleEndian, d.dataType)
	if d.rawValue != "" {
		//写入数据
		bcd := Number2bcd(d.rawValue)
		//翻转
		for i, j := 0, len(bcd)-1; i < j; i, j = i+1, j-1 {
			bcd[i], bcd[j] = bcd[j], bcd[i]
		}
		_ = binary.Write(buffer, binary.LittleEndian, &bcd)
	}

	return nil
}

func (d Data) getLen() uint16 {

	if d.dataType[3] == 0x00 && d.dataType[0] == 0x00 {
		return 4
	} else if d.dataType[3] == 0x00 && d.dataType[0] == 0x01 {
		return 4
	} else if d.dataType[3] == 0x00 && d.dataType[0] == 0x0c {
		return 4
	} else if d.dataType[3] == 0x01 && d.dataType[0] == 0x00 {
		return 8
	} else if d.dataType[3] == 0x01 && d.dataType[0] == 0x01 {
		return 8
	} else if d.dataType[3] == 0x01 && d.dataType[0] == 0x0c {
		return 8
	} else if d.dataType[3] == 0x02 && d.dataType[0] == 0x00 && d.dataType[2] == 0x01 {
		return 2
	} else if d.dataType[3] == 0x02 && d.dataType[0] == 0x00 && d.dataType[2] <= 0x05 {
		return 3
	}
	return 0
}

func NewData(dataType int32, value string) *Data {
	return &Data{dataType: Int2bytes(dataType), rawValue: value}
}

//ReadRequest 读数据
func ReadRequest(address string, itemCode int32, control *Control) *Protocol {
	return &Protocol{
		Start:      Start,
		Start2:     Start,
		End:        End,
		Address:    NewAddress(address),
		Control:    control,
		DataLength: 0x04,
		Data:       NewData(itemCode, ""),
	}

}
func GetHex(protocol *Protocol) string {
	bf := bytes.NewBuffer(make([]byte, 0))
	protocol.Encode(bf)
	return hex.EncodeToString(bf.Bytes())
}

func ReadResponse(address string, itemCode int32, control *Control, rawValue string) *Protocol {
	return &Protocol{
		Start:      Start,
		Start2:     Start,
		End:        End,
		Address:    NewAddress(address),
		Control:    control,
		DataLength: 0x04,
		Data:       NewData(itemCode, rawValue),
	}

}

type Protocol struct {
	//Start 645协议起始符号
	Start      byte
	Address    *Address
	Start2     byte
	Control    *Control
	DataLength byte
	Data       *Data
	CS         byte
	End        byte
}

func (p Protocol) Encode(buffer *bytes.Buffer) error {
	//计算cs 需要重写开辟字节码缓冲区
	tmp := make([]byte, 0)
	bf := bytes.NewBuffer(tmp)
	_ = binary.Write(bf, binary.LittleEndian, &p.Start)
	_ = p.Address.Encode(bf)
	_ = binary.Write(bf, binary.LittleEndian, &p.Start2)
	_ = p.Control.Encode(bf)
	_ = binary.Write(bf, binary.LittleEndian, &p.DataLength)
	_ = p.Data.Encode(bf)

	//计算Cs
	var cs = 0
	for _, b := range bf.Bytes() {
		cs += int(b)
	}
	p.CS = byte(cs)
	_ = binary.Write(bf, binary.LittleEndian, p.CS)
	_ = binary.Write(bf, binary.LittleEndian, p.End)

	//写入
	_ = binary.Write(buffer, binary.LittleEndian, bf.Bytes())

	return nil

}

func (p Protocol) getLen() uint16 {
	return HeadLen + 4 + p.Data.getLen()
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

func Decode(buffer *bytes.Buffer) *Protocol {

	p := new(Protocol)
	_ = binary.Read(buffer, binary.LittleEndian, &p.Start)
	p.Address = DecodeAddress(buffer, 6)
	_ = binary.Read(buffer, binary.LittleEndian, &p.Start2)
	p.Control = DecodeControl(buffer)
	_ = binary.Read(buffer, binary.LittleEndian, &p.DataLength)
	p.Data = DecodeData(buffer, p.DataLength)
	_ = binary.Read(buffer, binary.LittleEndian, &p.CS)
	_ = binary.Read(buffer, binary.LittleEndian, &p.End)
	return p
}
func DecodeAddress(buffer *bytes.Buffer, size int) *Address {
	a := new(Address)
	value := make([]byte, size)
	_ = binary.Read(buffer, binary.LittleEndian, &value)
	a.value = value
	a.StrValue = Bcd2Number(a.value)
	return a
}
func DecodeData(buffer *bytes.Buffer, size byte) *Data {
	data := new(Data)
	var dataType [4]byte
	dataValue := make([]byte, size-4)
	_ = binary.Read(buffer, binary.LittleEndian, &dataType)
	for index, item := range dataType {
		dataType[index] = item - 0x33
	}
	_ = binary.Read(buffer, binary.LittleEndian, &dataValue)
	for index, item := range dataValue {
		dataValue[index] = item - 0x33
	}
	//反转成小端的
	for i, j := 0, len(dataValue)-1; i < j; i, j = i+1, j-1 {
		dataValue[i], dataValue[j] = dataValue[j], dataValue[i]
	}

	for i, j := 0, len(dataType)-1; i < j; i, j = i+1, j-1 {
		dataType[i], dataType[j] = dataType[j], dataType[i]
	}
	data.rawValue = Bcd2Number(dataValue)
	data.dataType = dataType
	return data
}
func Int2bytes(data int32) [4]byte {

	var b3 [4]byte

	b3[0] = uint8(data)

	b3[1] = uint8(data >> 8)

	b3[2] = uint8(data >> 16)

	b3[3] = uint8(data >> 24)
	return b3
}
