package go2p

import (
	"net"

	"github.com/google/uuid"
)

type Peer interface {
	//func
}

type peer struct {
	id         string
	conn       net.Conn
	attributes map[string]interface{}
}

func newPeer(conn net.Conn) *peer {
	result := &peer{conn: conn}
	result.attributes = make(map[string]interface{})
	result.id = uuid.New().String()
	return result
}

type Peers struct {
	inner []Peer
}

func newPeers() *Peers {
	result := &Peers{}
	return result
}

func (self *Peers) Add(peer Peer) {
	self.inner = append(self.inner, peer)
}

func (self *Peers) Remove(peer Peer) {
	for i := range self.inner {
		if self.inner[i] == peer {
			copy(self.inner[i:], self.inner[i+1:])
			self.inner[len(self.inner)-1] = nil
			self.inner = self.inner[:len(self.inner)-1]
		}
	}
}
