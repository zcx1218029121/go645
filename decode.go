package go645

import "bytes"

type Decoder func(buffer *bytes.Buffer) (*InformationElement, error)

func Handler(control *Control, buffer *bytes.Buffer, size byte) InformationElement {
	if control.IsState(Read) {
		return DecodeRead(buffer, int(size))
	}
	panic("没有定义的数据类型")
	return nil
}
