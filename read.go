package go645

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"strconv"
	"time"
)

type (
	ReadDataWithTime struct {
		ReadData
		value float64
		time  time.Time
	}
	//ReadData 数据域
	ReadData struct {
		//数据标识 4 个字节
		dataType []byte
		//bcd 码的数据
		bcdValue string
		rawValue []byte
		Negative bool
	}
	//WriteData 写数据
	WriteData struct {
		dataType    []byte
		permissions byte
		passWord    []byte
		optCode     []byte
	}
	ReadRequestData struct {
		dataType  []byte
		recordNum byte
		min       byte
		hours     byte
		day       byte
		month     byte
		year      byte
		withTime  bool
	}
)

func (d ReadData) GetDataType() []byte {
	return d.dataType
}
func (d ReadData) GetDataTypeStr() string {
	return hex.EncodeToString(d.dataType)
}
func (d *ReadData) GetFloat64ValueWithTime() *ReadDataWithTime {
	if d.dataType[3] == 0x01 {
		_, _ = strconv.Atoi(d.bcdValue[:6])
	}
	return nil
}
func (d *ReadData) GetIntValue() (int, error) {
	value, err := strconv.Atoi(d.bcdValue)
	if err != nil {
		return 0, err
	}
	return value,nil
}
func (d *ReadData) GetFloat64Value() float64 {
	var data float64
	if d.dataType[3] == 0x00 {
		value, _ := strconv.Atoi(d.bcdValue)
		data = float64(value) * 0.01
	}
	if d.dataType[3] == 0x02 && d.dataType[2] == 0x01 && d.dataType[0] == 0 {
		value, _ := strconv.Atoi(d.bcdValue)
		data = float64(value) / 10
	}
	if d.dataType[3] == 0x02 && d.dataType[2] == 0x02 && d.dataType[0] == 0 {
		value, _ := strconv.Atoi(d.bcdValue)
		data = float64(value) * 0.001
	}
	if d.dataType[3] == 0x02 && d.dataType[2] == 0x03 && d.dataType[0] == 0 {
		value, _ := strconv.Atoi(d.bcdValue)
		data = float64(value) * 0.0001
	}
	if d.dataType[3] == 0x02 && d.dataType[2] == 0x04 {
		value, _ := strconv.Atoi(d.bcdValue)
		data = float64(value) * 0.0001

	}
	if d.dataType[3] == 0x02 && d.dataType[2] == 0x05 && d.dataType[0] == 0 {
		value, _ := strconv.Atoi(d.bcdValue)
		data = float64(value) * 0.0001
	}
	if d.dataType[3] == 0x02 && d.dataType[2] == 0x06 && d.dataType[0] == 0 {
		value, _ := strconv.Atoi(d.bcdValue)
		data = float64(value) * 0.001
	}
	if d.dataType[3] == 0x02 && d.dataType[2] == 0x07 && d.dataType[0] == 0 {
		value, _ := strconv.Atoi(d.bcdValue)
		data = float64(value) * 0.01
	}
	if d.dataType[3] == 0x02 && d.dataType[2] == 0x08 && d.dataType[0] == 0 {
		value, _ := strconv.Atoi(d.bcdValue)
		data = float64(value) * 0.01
	}
	if d.dataType[3] == 0x02 && d.dataType[2] == 0x09 && d.dataType[0] == 0 {
		value, _ := strconv.Atoi(d.bcdValue)
		data = float64(value) * 0.01
	}
	if d.dataType[3] == 0x02 && d.dataType[2] == 0x0A && d.dataType[1] == 0x01 {
		value, _ := strconv.Atoi(d.bcdValue)
		data = float64(value) * 0.01
	}
	if d.dataType[3] == 0x02 && d.dataType[2] == 0x0A && d.dataType[1] == 0x02 {
		value, _ := strconv.Atoi(d.bcdValue)
		data = float64(value) * 0.01
	}
	if d.dataType[3] == 0x02 && d.dataType[2] == 0x0A && d.dataType[1] == 0x03 {
		value, _ := strconv.Atoi(d.bcdValue)
		data = float64(value) * 0.01
	}
	if d.dataType[3] == 0x02 && d.dataType[2] == 0x0B && d.dataType[1] == 0x01 {
		value, _ := strconv.Atoi(d.bcdValue)
		data = float64(value) * 0.01
	}
	if d.dataType[3] == 0x02 && d.dataType[2] == 0x0B && d.dataType[1] == 0x02 {
		value, _ := strconv.Atoi(d.bcdValue)
		data = float64(value) * 0.01
	}
	if d.dataType[3] == 0x02 && d.dataType[2] == 0x0B && d.dataType[1] == 0x03 {
		value, _ := strconv.Atoi(d.bcdValue)
		data = float64(value) * 0.01
	}
	if d.dataType[3] == 0x02 && d.dataType[2] == 0x80 && d.dataType[0] == 0x01 {
		value, _ := strconv.Atoi(d.bcdValue)
		data = float64(value) / 1000
	}
	if d.dataType[3] == 0x02 && d.dataType[2] == 0x80 && d.dataType[0] == 0x02 {
		value, _ := strconv.Atoi(d.bcdValue)
		data = float64(value) / 100
	}
	if d.dataType[3] == 0x02 && d.dataType[2] == 0x80 && d.dataType[0] == 0x03 {
		value, _ := strconv.Atoi(d.bcdValue)
		data = float64(value) / 10000
	}
	if d.dataType[3] == 0x02 && d.dataType[2] == 0x80 && d.dataType[0] == 0x04 {
		value, _ := strconv.Atoi(d.bcdValue)
		data = float64(value) / 10000
	}
	if d.dataType[3] == 0x02 && d.dataType[2] == 0x80 && d.dataType[0] == 0x05 {
		value, _ := strconv.Atoi(d.bcdValue)
		data = float64(value) / 10000
	}
	if d.dataType[3] == 0x02 && d.dataType[2] == 0x80 && d.dataType[0] == 0x06 {
		value, _ := strconv.Atoi(d.bcdValue)
		data = float64(value) / 10000
	}
	if d.dataType[3] == 0x02 && d.dataType[2] == 0x80 && d.dataType[0] == 0x07 {
		value, _ := strconv.Atoi(d.bcdValue)
		data = float64(value) / 10
	}
	if d.dataType[3] == 0x02 && d.dataType[2] == 0x80 && d.dataType[0] == 0x08 {
		value, _ := strconv.Atoi(d.bcdValue)
		data = float64(value) / 100
	}
	if d.dataType[3] == 0x02 && d.dataType[2] == 0x80 && d.dataType[0] == 0x09 {
		value, _ := strconv.Atoi(d.bcdValue)
		data = float64(value) / 100
	}
	if d.dataType[3] == 0x02 && d.dataType[2] == 0x80 && d.dataType[0] == 0x0A {
		value, _ := strconv.Atoi(d.bcdValue)
		data = float64(value)
	}
	if d.dataType[3] == 0x02 && d.dataType[2] == 0x80 && d.dataType[0] == 0x0B {
		value, _ := strconv.Atoi(d.bcdValue)
		data = float64(value) / 10000
	}
	if d.Negative {
		return data * -1
	}
	return data
}
func (d *ReadData) GetFloat64ValueUnsigned() float64 {
	for i, j := 0, len(d.dataType)-1; i < j; i, j = i+1, j-1 {
		d.dataType[j], d.dataType[i] = d.dataType[i], d.dataType[j]
	}
	if d.dataType[3] == 0x00 {
		value, _ := strconv.Atoi(d.bcdValue)
		return float64(value) * 0.01
	}
	if d.dataType[3] == 0x02 && d.dataType[2] == 0x03 && d.dataType[0] == 0 {
		value, _ := strconv.Atoi(d.bcdValue)
		return float64(value) * 0.0001

	}
	if d.dataType[3] == 0x02 && d.dataType[2] == 0x04 {
		value, _ := strconv.Atoi(d.bcdValue)
		data := float64(value) * 0.0001
		if data > 80 {
			return (data - 80) * -1
		} else {
			return data
		}
	}
	if d.dataType[3] == 0x02 && d.dataType[2] == 0x05 && d.dataType[0] == 0 {
		value, _ := strconv.Atoi(d.bcdValue)
		return float64(value) * 0.0001
	}
	if d.dataType[3] == 0x02 && d.dataType[2] == 0x06 && d.dataType[0] == 0 {
		value, _ := strconv.Atoi(d.bcdValue)
		return float64(value) * 0.001
	}
	if d.dataType[3] == 0x02 && d.dataType[2] == 0x05 && d.dataType[0] == 0 {
		value, _ := strconv.Atoi(d.bcdValue)
		return float64(value) * 0.0001
	}
	if d.dataType[3] == 0x02 && d.dataType[2] == 0x06 && d.dataType[0] == 0 {
		value, _ := strconv.Atoi(d.bcdValue)
		return float64(value) * 0.001
	}
	if d.dataType[3] == 0x02 && d.dataType[2] == 0x07 && d.dataType[0] == 0 {
		value, _ := strconv.Atoi(d.bcdValue)
		return float64(value) * 0.1
	}
	if d.dataType[3] == 0x02 && d.dataType[2] == 0x08 && d.dataType[0] == 0 {
		value, _ := strconv.Atoi(d.bcdValue)
		return float64(value) * 0.01
	}
	if d.dataType[3] == 0x02 && d.dataType[2] == 0x09 && d.dataType[0] == 0 {
		value, _ := strconv.Atoi(d.bcdValue)
		return float64(value) * 0.01
	}
	if d.dataType[3] == 0x02 && d.dataType[2] == 0x0A && d.dataType[1] == 0x01 {
		value, _ := strconv.Atoi(d.bcdValue)
		return float64(value) * 0.01
	}
	if d.dataType[3] == 0x02 && d.dataType[2] == 0x0A && d.dataType[1] == 0x02 {
		value, _ := strconv.Atoi(d.bcdValue)
		return float64(value) * 0.01
	}
	if d.dataType[3] == 0x02 && d.dataType[2] == 0x0A && d.dataType[1] == 0x03 {
		value, _ := strconv.Atoi(d.bcdValue)
		return float64(value) * 0.01
	}
	if d.dataType[3] == 0x02 && d.dataType[2] == 0x0B && d.dataType[1] == 0x01 {
		value, _ := strconv.Atoi(d.bcdValue)
		return float64(value) * 0.01
	}
	if d.dataType[3] == 0x02 && d.dataType[2] == 0x0B && d.dataType[1] == 0x02 {
		value, _ := strconv.Atoi(d.bcdValue)
		return float64(value) * 0.01
	}
	if d.dataType[3] == 0x02 && d.dataType[2] == 0x0B && d.dataType[1] == 0x03 {
		value, _ := strconv.Atoi(d.bcdValue)
		return float64(value) * 0.01
	}
	if d.dataType[3] == 0x02 && d.dataType[2] == 0x80 && d.dataType[1] == 0x01 {
		value, _ := strconv.Atoi(d.bcdValue)
		return float64(value) * 0.001
	}
	if d.dataType[3] == 0x02 && d.dataType[2] == 0x80 && d.dataType[1] == 0x02 {
		value, _ := strconv.Atoi(d.bcdValue)
		return float64(value) * 0.01
	}
	if d.dataType[3] == 0x02 && d.dataType[2] == 0x80 && d.dataType[1] == 0x03 {
		value, _ := strconv.Atoi(d.bcdValue)
		return float64(value) * 0.001
	}
	if d.dataType[3] == 0x02 && d.dataType[2] == 0x80 && d.dataType[1] == 0x04 {
		value, _ := strconv.Atoi(d.bcdValue)
		return float64(value) * 0.001
	}
	if d.dataType[3] == 0x02 && d.dataType[2] == 0x80 && d.dataType[1] == 0x05 {
		value, _ := strconv.Atoi(d.bcdValue)
		return float64(value) * 0.001
	}
	if d.dataType[3] == 0x02 && d.dataType[2] == 0x80 && d.dataType[1] == 0x06 {
		value, _ := strconv.Atoi(d.bcdValue)
		return float64(value) * 0.001
	}
	if d.dataType[3] == 0x02 && d.dataType[2] == 0x80 && d.dataType[1] == 0x07 {
		value, _ := strconv.Atoi(d.bcdValue)
		return float64(value) * 0.1
	}
	if d.dataType[3] == 0x02 && d.dataType[2] == 0x80 && d.dataType[1] == 0x08 {
		value, _ := strconv.Atoi(d.bcdValue)
		return float64(value) * 0.01
	}
	if d.dataType[3] == 0x02 && d.dataType[2] == 0x80 && d.dataType[1] == 0x09 {
		value, _ := strconv.Atoi(d.bcdValue)
		return float64(value) * 0.01
	}
	if d.dataType[3] == 0x02 && d.dataType[2] == 0x80 && d.dataType[1] == 0x09 {
		value, _ := strconv.Atoi(d.bcdValue)
		return float64(value) * 0.01
	}
	if d.dataType[3] == 0x02 && d.dataType[2] == 0x80 && d.dataType[1] == 0x0A {
		value, _ := strconv.Atoi(d.bcdValue)
		return float64(value)
	}
	if d.dataType[3] == 0x02 && d.dataType[2] == 0x80 && d.dataType[1] == 0x0B {
		value, _ := strconv.Atoi(d.bcdValue)
		return float64(value) * 0.0001
	}
	return 0
}

func (d ReadData) GetValue() string {
	return d.bcdValue
}

func (d ReadData) Encode(buffer *bytes.Buffer) error {
	//写入数据域 已经反转过了
	for index, b := range d.dataType {
		d.dataType[index] = b + 0x33
	}
	if err := binary.Write(buffer, binary.LittleEndian, d.dataType); err != nil {
		return err
	}
	if d.bcdValue != "" {
		//写入数据
		bcd := Number2bcd(d.bcdValue)
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

func (d ReadData) GetLen() byte {
	if d.bcdValue == "" {
		return 4
	}
	return 4 + byte(len(Number2bcd(d.bcdValue)))
}

//ReadRequest 读数据
func ReadRequest(address Address, itemCode int32) *Protocol {
	c := NewControl()
	c.SetState(Read)
	d := NewReadData(itemCode, "")
	return NewProtocol(address, d, c)

}

//ReadRequestWithBlock 读数据
func ReadRequestWithBlock(address Address, data ReadRequestData) *Protocol {
	c := NewControl()
	c.SetState(Read)
	return NewProtocol(address, data, c)

}

//ReadResponse 创建读响应
func ReadResponse(address Address, itemCode int32, control *Control, rawValue string) *Protocol {
	return &Protocol{
		Start:      Start,
		Start2:     Start,
		End:        End,
		Address:    address,
		Control:    control,
		DataLength: 0x04,
		Data:       NewReadData(itemCode, rawValue),
	}

}

func (r ReadRequestData) Encode(buffer *bytes.Buffer) error {
	//写入数据域 已经反转过了
	var err error
	for _, b := range r.dataType {
		err = WriteWithOffSet(buffer, b)
	}
	if r.recordNum != 0 {
		err = WriteWithOffSet(buffer, r.recordNum)

	}
	if r.withTime {
		err = WriteWithOffSet(buffer, r.min)
		err = WriteWithOffSet(buffer, r.hours)
		err = WriteWithOffSet(buffer, r.day)
		err = WriteWithOffSet(buffer, r.month)
		err = WriteWithOffSet(buffer, r.year)
	}
	return err
}

func (r ReadRequestData) GetLen() byte {
	var dataLen byte = 4
	if r.withTime {
		dataLen += 5
	}
	if r.recordNum != 0 {
		dataLen += 1
	}
	return dataLen
}
