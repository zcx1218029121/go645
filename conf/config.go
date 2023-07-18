package conf

import (
	"time"

	"github.com/tarm/serial"
)

type Config struct {
	Retry int
	Com   []Com
	Mqtt  Mqtt
}
type Com struct {
	Config SerialConfig
	Size   int
	Type   uint8
	Start  int
}

type SerialConfig struct {
	Address  string
	BaudRate int
	Timeout  time.Duration
	Size     byte
	Parity   string
	StopBits byte
}
type Mqtt struct {
	Clinet    string
	Broker    string
	Username  string
	Password  string
	InfoTopic string
	AvgTopic  string

	ReceiveTopic         string
	TimeOut              uint8
	SendTime             uint8
	ControlTopic         string
	ControlCallBackTopic string
}

func (s SerialConfig) NewSerialConfig() serial.Config {
	var p serial.Parity
	switch s.Parity {
	case "N":
		p = serial.ParityNone
	case "E":
		p = serial.ParityEven
	case "O":
		p = serial.ParityOdd
	}

	return serial.Config{
		Parity:      p,
		Name:        s.Address,
		Baud:        s.BaudRate,
		ReadTimeout: s.Timeout,
		StopBits:    serial.StopBits(s.StopBits),
	}

}
