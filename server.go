package websocket

import (
	"net/http"
	"net/url"
	"strings"
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

type HandshakeError struct {
	message string
}

func (e HandshakeError) Error() string {
	return e.message
}

// Upgrade http协议升级到WebSocket
func (u *Upgrader) Upgrade(w http.ResponseWriter, r *http.Request, responseHeader http.Header) (*Conn, error) {
	const badHandshake = "websocket: the client is not using the websocket protocol"

	if !tokenListContainsValues(r.Header, "Connection", "upgrade") {
		return u.returnError(w, r, http.StatusBadRequest, badHandshake+" ,'upgrade': token not found in 'Connection' header")
	}

	if !tokenListContainsValues(r.Header, "Upgrade", "websocket") {
		w.Header().Set("Upgrade", "websocket")
		return u.returnError(w, r, http.StatusBadRequest, badHandshake+" ,'websocket': token not found in 'Upgrade' header")
	}

	if r.Method != http.MethodGet {
		return u.returnError(w, r, http.StatusMethodNotAllowed, badHandshake+" "+r.Method+" not allowed")
	}

	if !tokenListContainsValues(r.Header, "Sec-WebSocket-Version", "13") {
		w.Header().Set("Upgrade", "websocket")
		return u.returnError(w, r, http.StatusBadRequest, badHandshake+" ,'websocket': unsupported version :13 not found in 'Sec-WebSocket-Version' header")
	}

	// 握手阶段,响应头
	if _, ok := responseHeader["Sec-WebSocket-Extensions"]; !ok {
		return u.returnError(w, r, http.StatusInternalServerError, badHandshake+" websocket: application specific 'Sec-WebSocket-Extensions' header are not unsupported")
	}

	checkOrigin := u.CheckOrigin
	if checkOrigin == nil {
		checkOrigin = checkSameOrigin
	}

	if !checkOrigin(r) {
		return u.returnError(w, r, http.StatusForbidden, "websocket: request origin not allowed by Upgrader.CheckOrigin")
	}

	challengeKey := r.Header.Get("Sec-WebSocket-Key")
	if !isValidChallengeKey(challengeKey) {
		return u.returnError(w, r, http.StatusBadRequest, "websocket: not a websocket handshake :'Sec-WebSocket-Key' header is invalid")
	}

	// 匹配子协议
	subprotocol := u.selectSubprotocol(r, responseHeader)

	//协商压缩格式

	return nil, nil
}

func (u *Upgrader) returnError(w http.ResponseWriter, r *http.Request, status int, reason string) (*Conn, error) {
	err := &HandshakeError{reason}
	if u.Error != nil {
		u.Error(w, r, status, err)
	} else {
		w.Header().Set("Sec-WebSocket-Version", "13")
		http.Error(w, http.StatusText(status), status)
	}

	return nil, err
}

func checkSameOrigin(r *http.Request) bool {
	origin := r.Header["Origin"]
	if len(origin) == 0 {
		return true
	}

	u, err := url.Parse(origin[0])
	if err != nil {
		return false
	}

	return equalASCIIFold(u.Host, r.Host)
}

func Subprotocols(r *http.Request) []string {
	h := strings.TrimSpace(r.Header.Get("Sec-WebSocket-Protocol"))
	if h == "" {
		return nil
	}

	procotols := strings.Split(h, ",")
	for i := range procotols {
		procotols[i] = strings.TrimSpace(procotols[i])
	}
	return procotols
}
func (u *Upgrader) selectSubprotocol(r *http.Request, responseHeader http.Header) string {
	if u.SubProtocols != nil {
		clientProtocols := Subprotocols(r)
		for _, clientProcotol := range clientProtocols {
			for _, serverProtocol := range u.SubProtocols {
				if clientProcotol == serverProtocol {
					return clientProcotol
				}
			}
		}
	} else if u.SubProtocols == nil {
		return responseHeader.Get("Sec-WebSocket-Protocol")
	}

	return ""
}
