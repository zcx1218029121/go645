package go645

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"
)

type (
	//Order 地址大小端
	Order              bool
	InformationElement interface {
		Encode(buffer *bytes.Buffer) error
		getLen() byte
	}
	//Address 表计地址
	Address struct {
		value    []byte
		strValue string
	}
	//Data 数据域
	Data struct {
		//数据标识 4 个字节
		dataType [4]byte
		//原始数据
		rawValue string
	}
	//Protocol 协议
	Protocol struct {
		//Start 645协议起始符号
		Start byte
		//设备地址 6个字节的BCD
		Address *Address
		//Start  645协议起始符号 标志报文头结束
		Start2 byte
		//Control 控制域
		Control *Control
		//Control 数据长度
		DataLength byte
		//Control 数据抽象
		Data *Data
		//CS 校验和
		CS byte
		//End 0x16
		End byte
	}
)

const (
	LittleEndian Order = false
	BigEndian    Order = true
	Start              = 0x68
	End                = 0x16
	HeadLen            = 1 + 6 + 1
)

var (
	_ InformationElement = (*Address)(nil)

	_ InformationElement = (*Protocol)(nil)

	_ InformationElement = (*Data)(nil)
)

// NewAddress ，构建设备地址
// 参数：
//      address ： 设备地址
//      order ： 大小端表示
// 返回值：
//      *Address 设备地址
func NewAddress(address string, order Order) *Address {
	value := Number2bcd(address)
	if !order {
		for i, j := 0, len(value)-1; i < j; i, j = i+1, j-1 {
			value[i], value[j] = value[j], value[i]
		}
	}

	return &Address{value: value, strValue: address}
}

func NewData(dataType int32, value string) *Data {
	return &Data{dataType: Int2bytes(dataType), rawValue: value}
}

func NewProtocol(address *Address, data *Data, control *Control) *Protocol {
	return &Protocol{
		Start:      Start,
		Start2:     Start,
		End:        End,
		Address:    address,
		Control:    control,
		DataLength: data.getLen(),
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

func (a Address) getLen() byte {
	return 6
}

func (d Data) GetDataType() [4]byte {
	return d.dataType
}
func (d Data) GetDataTypeStr() string {
	//需要翻转
	var a = make([]byte, 4)
	for i, j := 0, len(d.dataType)-1; i < j; i, j = i+1, j-1 {
		a[i], a[j] = d.dataType[i], d.dataType[j]
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
	if err := binary.Write(buffer, binary.LittleEndian, d.dataType); err != nil {
		return err
	}
	if d.rawValue != "" {
		//写入数据
		bcd := Number2bcd(d.rawValue)
		//翻转
		for i, j := 0, len(bcd)-1; i < j; i, j = i+1, j-1 {
			bcd[i], bcd[j] = bcd[j], bcd[i]
		}
		if err := binary.Write(buffer, binary.LittleEndian, &bcd); err != nil {
			return err
		}
	}

	return nil
}

func (d Data) getLen() byte {
	if d.dataType[3] == 0x00 && d.dataType[0] == 0x00 {
		return 4
	}
	return 4
}

//ReadRequest 读数据
func ReadRequest(address *Address, itemCode int32) *Protocol {
	c := NewControl()
	c.SetState(Read)
	d := NewData(itemCode, "")
	return NewProtocol(address, d, c)

}

//GetHex 返回16进制string
func GetHex(protocol *Protocol) (string, error) {
	bf := bytes.NewBuffer(make([]byte, 0))
	if err := protocol.Encode(bf); err != nil {
		return "", err
	}
	return hex.EncodeToString(bf.Bytes()), nil
}

//ReadResponse 创建读响应
func ReadResponse(address *Address, itemCode int32, control *Control, rawValue string) *Protocol {
	return &Protocol{
		Start:      Start,
		Start2:     Start,
		End:        End,
		Address:    address,
		Control:    control,
		DataLength: 0x04,
		Data:       NewData(itemCode, rawValue),
	}

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

func (p Protocol) getLen() byte {
	if p.DataLength != 0 {
		return p.DataLength
	}
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

func Decode(buffer *bytes.Buffer) (*Protocol, error) {
	var err error
	read := func(data interface{}) {
		if err != nil {
			return
		}
		err = binary.Read(buffer, binary.LittleEndian, data)
	}
	p := new(Protocol)
	read(&p.Start)
	p.Address, err = DecodeAddress(buffer, 6)
	read(&p.Start2)
	p.Control, err = DecodeControl(buffer)
	read(&p.DataLength)
	p.Data, err = DecodeData(buffer, p.DataLength)
	read(&p.CS)
	read(&p.End)
	return p, nil
}
func DecodeAddress(buffer *bytes.Buffer, size int) (*Address, error) {
	a := new(Address)
	value := make([]byte, size)
	if err := binary.Read(buffer, binary.LittleEndian, &value); err != nil {
		return nil, err
	}
	{
		a.value = value
		a.strValue = Bcd2Number(a.value)
	}
	return a, nil
}
func DecodeData(buffer *bytes.Buffer, size byte) (*Data, error) {

	var err error
	read := func(data interface{}) {
		if err != nil {
			return
		}
		err = binary.Read(buffer, binary.LittleEndian, data)
	}
	data := new(Data)
	var dataType [4]byte
	dataValue := make([]byte, size-4)
	read(&dataType)
	for index, item := range dataType {
		dataType[index] = item - 0x33
	}
	read(&dataValue)
	for index, item := range dataValue {
		dataValue[index] = item - 0x33
	}
	for i, j := 0, len(dataValue)-1; i < j; i, j = i+1, j-1 {
		dataValue[i], dataValue[j] = dataValue[j], dataValue[i]
	}
	for i, j := 0, len(dataType)-1; i < j; i, j = i+1, j-1 {
		dataType[i], dataType[j] = dataType[j], dataType[i]
	}
	data.rawValue = Bcd2Number(dataValue)
	data.dataType = dataType
	return data, nil
}
func Int2bytes(data int32) [4]byte {
	var b3 [4]byte
	b3[0] = uint8(data)
	b3[1] = uint8(data >> 8)
	b3[2] = uint8(data >> 16)
	b3[3] = uint8(data >> 24)
	return b3
}
