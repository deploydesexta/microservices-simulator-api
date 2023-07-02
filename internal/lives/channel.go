package lives

import "github.com/gorilla/websocket"

type InputEvent struct {
	Event   string      `json:"event"`
	Payload interface{} `json:"payload"`
}

type OutputEvent struct {
	UserId  int64       `json:"user_id"`
	Event   string      `json:"event"`
	Payload interface{} `json:"payload"`
}

type Channel interface {
	Send(msg []byte)
	SendEvent(msg OutputEvent)
	Receive() ([]byte, error)
	SetCloseHandler(func(code int, text string) error)
}

type WebSocketChannel struct {
	Conn *websocket.Conn
}

func NewWebSocketChannel(conn *websocket.Conn) *WebSocketChannel {
	return &WebSocketChannel{
		Conn: conn,
	}
}

func (c *WebSocketChannel) Send(msg []byte) {
	c.Conn.WriteMessage(websocket.TextMessage, msg)
}

func (c *WebSocketChannel) SendEvent(msg OutputEvent) {
	c.Conn.WriteJSON(msg)
}

func (c *WebSocketChannel) Receive() ([]byte, error) {
	_, msg, err := c.Conn.ReadMessage()
	return msg, err
}

func (c *WebSocketChannel) SetCloseHandler(handler func(code int, text string) error) {
	c.Conn.SetCloseHandler(handler)
}
