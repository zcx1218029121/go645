package go645

import (
	"bytes"
	"encoding/binary"
	"errors"
	"log"
)

type Decoder func(buffer *bytes.Buffer) (*InformationElement, error)

func Handler(control *Control, buffer *bytes.Buffer, size byte) (InformationElement, error) {
	//从站响应异常响应
	if control == nil {
		return nil, errors.New("未知错误")
	}
	if control.IsState(SlaveErr) {
		return nil, DecodeException(buffer)
	}
	//从站读正确响应
	if control.IsState(Read) {
		return DecodeRead(buffer, int(size)), nil
	}
	//佳和强制联机
	if control.Data == 0x8a {
		return DecodeNullData(buffer), nil
	}
	return nil, errors.New("未定义的数据类型")
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
	p.Data, err = Handler(p.Control, buffer, p.DataLength)
	read(&p.CS)
	read(&p.End)
	if err != nil {
		log.Print(err.Error())
	}
	return p, err
}
func DecodeAddress(buffer *bytes.Buffer, size int) (Address, error) {
	var a Address

	value := make([]byte, size)
	if err := binary.Read(buffer, binary.LittleEndian, &value); err != nil {
		return a, err
	}
	{
		a.value = value
		a.strValue = Bcd2Number(a.value)
	}
	return a, nil
}
func DecodeData(buffer *bytes.Buffer, size byte) (*ReadData, error) {
	var err error
	read := func(data interface{}) {
		if err != nil {
			return
		}
		err = binary.Read(buffer, binary.LittleEndian, data)
	}
	data := new(ReadData)
	var dataType []byte
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
func DecoderData(buffer *bytes.Buffer, size int) (*bytes.Buffer, error) {
	var err error
	read := func(data interface{}) {
		if err != nil {
			return
		}
		err = binary.Read(buffer, binary.LittleEndian, data)
	}
	var value = make([]byte, size)
	read(value)
	for i, j := 0, len(value)-1; i <= j; i, j = i+1, j-1 {
		value[i], value[j] = value[j]-0x33, value[i]-0x33
	}

	return bytes.NewBuffer(value), nil
}
func DecodeRead(buffer *bytes.Buffer, size int) InformationElement {
	df, _ := DecoderData(buffer, size)
	var err error
	read := func(data interface{}) {
		if err != nil {
			return
		}
		err = binary.Read(df, binary.LittleEndian, data)
	}
	data := new(ReadData)
	var dataType = make([]byte, 4)
	dataValue := make([]byte, size-4)
	read(&dataValue)
	read(&dataType)

	data.rawValue = Bcd2Number(dataValue)
	data.dataType = dataType
	return data
}
func DecodeException(buffer *bytes.Buffer) error {
	var data uint16
	err := binary.Read(buffer, binary.LittleEndian, &data)
	if err != nil {
		log.Print(err)
		return nil
	}
	return &Exception{data}
}
func DecodeNullData(*bytes.Buffer) InformationElement {
	return NullData{}
}
