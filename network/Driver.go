package network

type Driver interface {
	Name() string                                 // 驱动名
	Create(name, subnet string) (*Network, error) // 创建网络
	Delete(nw *Network) error                     // 删除网络
	Connect(nw *Network, ep *Endpoint) error      // 连接网络与容器网络端点
	Disconnect(nw *Network, ep *Endpoint) error   // 断开网络与容器网络端点
}
