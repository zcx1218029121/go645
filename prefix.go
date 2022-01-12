package go645

import (
	"bytes"
	"io"
)

var _ PrefixHandler = (*DefaultPrefix)(nil)

type PrefixHandler interface {
	EncodePrefix(buffer *bytes.Buffer) error
	DecodePrefix(reader io.Reader) ([]byte, error)
}

type DefaultPrefix struct {
}

func (d DefaultPrefix) EncodePrefix(buffer *bytes.Buffer) error {

	return nil
}

func (d DefaultPrefix) DecodePrefix(reader io.Reader) ([]byte, error) {
	return nil, nil
}
