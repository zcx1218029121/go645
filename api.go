package go645

type Client interface {
	ClientProvider
	//Read 发送读请求
	Read(address Address, itemCode int32) (*ReadData, bool, error)
	//ReadWithBlock  读请求使能块
	ReadWithBlock(address Address, data ReadRequestData) (*Protocol, error)
	//广播信息
	Broadcast(p *Protocol) error
}
