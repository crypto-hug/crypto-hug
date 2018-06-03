package go2p

import (
	"net"
)

type peerConn struct {
	addr       string
	conn       net.Conn
	attributes map[string]interface{}

	protocols chan protocol
}

func newPeerConn(conn net.Conn) *peerConn {
	result := &peerConn{conn: conn}
	result.attributes = make(map[string]interface{})
	//result.id = uuid.New().String()
	return result
}

func (self *peerConn) execProtocol(p protocol) {
	self.protocols <- p
}

type protocol interface {
}
