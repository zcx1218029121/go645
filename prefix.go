package go645

import (
	"bytes"
	"io"
)

type PrefixHandler interface {
	EncodePrefix(buffer *bytes.Buffer) error
	DecodePrefix(reader io.Reader) ([]byte, error)
}
