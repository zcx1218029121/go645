package go645

import (
	"bufio"
	"bytes"
	"io"
	"time"
)

var _ ClientProvider = (*RTUClientProvider)(nil)

type RTUClientProvider struct {
	serialPort
	logger
}

//ReadRawFrame 读取Frame 线程不安全
func (sf *RTUClientProvider) ReadRawFrame() (aduResponse []byte, err error) {
	return bufio.NewReader(sf.port).ReadSlice(End)
}

//SendRawFrameNoAck 广播命令不需要返回
func (sf *RTUClientProvider) SendRawFrameNoAck(aduRequest []byte) (err error) {
	if err = sf.connect(); err != nil {
		return
	}
	// Send the request
	sf.Debugf("sending [% x]", aduRequest)
	//发送数据
	_, err = sf.port.Write(aduRequest)
	return err
}

//SendRawFrameWithHandle 广播命令不需要返回 线程不安全应该
func (sf *RTUClientProvider) SendRawFrameWithHandle(aduRequest []byte, fun func(io.ReadWriteCloser)) (err error) {
	sf.mu.Lock()
	defer sf.mu.Unlock()
	if err = sf.connect(); err != nil {
		return
	}
	// Send the request
	sf.Debugf("sending [% x]", aduRequest)
	//发送数据
	_, err = sf.port.Write(aduRequest)
	fun(sf.port)
	return err
}

func (sf *RTUClientProvider) Send(p *Protocol) (aduResponse []byte, err error) {
	bf := bytes.NewBuffer(make([]byte, 0))
	err = p.Encode(bf)
	if err != nil {
		return nil, err
	}
	return sf.SendRawFrame(bf.Bytes())
}
func (sf *RTUClientProvider) SendRawFrame(aduRequest []byte) (aduResponse []byte, err error) {
	sf.mu.Lock()
	defer sf.mu.Unlock()
	if err = sf.connect(); err != nil {
		return
	}
	// Send the request
	sf.Debugf("sending [% x]", aduRequest)
	//发送数据
	_, err = sf.port.Write(aduRequest)
	if err != nil {
		sf.close()
		return
	}
	//读取数据到结束符
	time.Sleep(sf.calculateDelay(len(aduRequest)))
	return bufio.NewReader(sf.port).ReadSlice(End)
}

// NewRTUClientProvider allocates and initializes a RTUClientProvider.
// it will use default /dev/ttyS0 19200 8 1 N and timeout 1000
func NewRTUClientProvider(opts ...ClientProviderOption) *RTUClientProvider {
	p := &RTUClientProvider{
		logger: newLogger("modbusRTUMaster => "),
		//pool:   rtuPool,
	}
	for _, opt := range opts {
		opt(p)
	}
	return p
}

// calculateDelay roughly calculates time needed for the next frame.
// See MODBUS over Serial Line - Specification and Implementation Guide (page 13).
func (sf *RTUClientProvider) calculateDelay(chars int) time.Duration {
	var characterDelay, frameDelay int // us

	if sf.BaudRate <= 0 || sf.BaudRate > 19200 {
		characterDelay = 750
		frameDelay = 1750
	} else {
		characterDelay = 15000000 / sf.BaudRate
		frameDelay = 35000000 / sf.BaudRate
	}
	return time.Duration(characterDelay*chars+frameDelay) * time.Microsecond
}
