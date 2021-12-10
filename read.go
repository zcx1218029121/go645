package go645

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"strconv"
)

type (
	//ReadData 数据域
	ReadData struct {
		//数据标识 4 个字节
		dataType [4]byte
		//原始数据
		rawValue string
	}
	ReadRequestData struct {
		dataType  [4]byte
		recordNum byte
		min       byte
		hours     byte
		day       byte
		month     byte
		year      byte
		withTime  bool
	}
)

func (d ReadData) GetDataType() [4]byte {
	return d.dataType
}
func (d ReadData) GetDataTypeStr() string {
	//需要翻转
	var a = make([]byte, 4)
	for i, j := 0, len(d.dataType)-1; i < j; i, j = i+1, j-1 {
		a[i], a[j] = d.dataType[i], d.dataType[j]
	}
	return hex.EncodeToString(a)
}

func (d *ReadData) GetFloat64Value() float64 {
	if d.dataType[0] == 0x00 || d.dataType[0] == 0x0c {
		value, _ := strconv.Atoi(d.rawValue)
		return float64(value) * 0.01
	} else if d.dataType[3] == 0x02 {
		value, _ := strconv.Atoi(d.rawValue)
		return float64(value) * 0.0001
	}
	return 0
}

func (d ReadData) GetValue() string {
	return d.rawValue
}

func (d ReadData) Encode(buffer *bytes.Buffer) error {
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

func (d ReadData) GetLen() byte {
	if d.rawValue == "" {
		return 4
	}
	return 4 + byte(len(Number2bcd(d.rawValue)))
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
