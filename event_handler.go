package go645

var _ EHandler = (*EventHandler)(nil)

type EHandler interface {
	ForceOnlineResp(client Client, address Address, data InformationElement) error
	ReadDataResp(client Client, address Address, data *ReadData) error
	OnlineResp(client Client, address Address, data InformationElement) error
}
type EventHandler struct {
}

func (e EventHandler) ForceOnlineResp(client Client, address Address, data InformationElement) error {
	return nil
}

func (e EventHandler) ReadDataResp(client Client, address Address, data *ReadData) error {
	return nil
}

func (e EventHandler) OnlineResp(client Client, address Address, data InformationElement) error {
	return nil
}
