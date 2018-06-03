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

func (self *Peers) byId(id string) *peer {
	for _, p := range self.inner {
		peer := p.(*peer)
		if peer.id == id {
			return peer
		}
	}

	return nil
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
func (self *Peers) Close(peer2Close Peer) {
	p := peer2Close.(*peer)
	p.conn.Close()
}

func (self *Peers) CloseAll() {
	for _, p := range self.inner {
		p2 := p.(*peer)
		p2.conn.Close()
	}
}
