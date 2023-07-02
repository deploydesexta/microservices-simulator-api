package lives

type ActiveClient struct {
	Id      int64
	Name    string
	channel Channel
}

func NewActiveClient(id int64, name string, channel Channel) *ActiveClient {
	return &ActiveClient{
		Id:      id,
		Name:    name,
		channel: channel,
	}
}

func (c *ActiveClient) Write(msg []byte) {
	c.channel.Send(msg)
}
