package websocket

import (
	"net/http"
	"time"
)

type Upgrader struct {
	// 握手超时时间
	HandshakeTimeout time.Duration

	//读 写缓冲大小
	ReadBufferSize, WriteBufferSize int

	//写缓冲池
	WriteBufferPool BufferPool

	//协议结构
	SubProtocols []string

	//错误
	Error func(w http.ResponseWriter, r *http.Request, status int, reason error)

	//同源策略校验
	CheckOrigin func(r *http.Request) bool

	//消息压缩
	EnableCompression bool
}

// Upgrade http协议升级到WebSocket
func (u *Upgrader) Upgrade(w http.ResponseWriter, r http.Request, responseHeader http.Header) (*Conn, error) {
	const badHandshake = "websocket: the client is not using the websocket protocol"
	if !tokenListContainsValues {

	}

	return nil, nil
}
