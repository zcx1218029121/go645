package go645

type Client interface {
	ClientProvider
	//Read 发送读请求
	Read(address Address, itemCode int32) (*ReadData, bool, error)
	//ReadWithBlock  读请求使能块
	ReadWithBlock(address Address, data ReadRequestData) (*Protocol, error)
	//Broadcast 开始广播
	Broadcast(p InformationElement, control Control) error
}
