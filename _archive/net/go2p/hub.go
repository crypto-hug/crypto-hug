package go2p

import (
	"context"
	"sync"
	"time"

	"github.com/jinzhu/copier"
)

//type

type hub struct {
	peers []*peerConn

	mutex sync.Mutex

	io      *IO
	context context.Context

	err    chan NetError
	add    chan *adapter
	remove chan *adapter
}

func StartupNewHub(io *IO, context context.Context) *hub {
	result := &hub{}
	result.io = io
	result.err = make(chan NetError)
	result.add = make(chan *adapter)
	result.remove = make(chan *adapter)
	result.context = context

	go result.process()

	return result
}

func (self *hub) Send(peer *peerConn, msg *Message) {
	self.mutex.Lock()
	defer self.mutex.Unlock()
	if self.context.Err() == nil {
		return
	}

	enqueMsg(peer, msg)
}
func (self *hub) Boradcast(broadcast *Message) {
	self.mutex.Lock()
	defer self.mutex.Unlock()
	if self.context.Err() == nil {
		return
	}

	for i := range self.peers {
		peer := self.peers[i]
		msg := &Message{}
		copier.Copy(&msg, broadcast)
		msg.header.receiverAddr = peer.addr
		enqueMsg(peer, msg)
	}
}

func (self *hub) process() {
	rate := time.Second / 10
	tick := time.NewTicker(rate)
	defer tick.Stop()

	for {
		select {
		case <-self.context.Done():
			self.cleanup()
			return
		case err := <-self.err:
			self.handleErr(err)
		case new := <-self.add:
			self.handleAdd(new)
		case old := <-self.remove:
			self.handleRemove(old)
		}
	}
}

// func (self *hub) connect(peer *peerConn) {
// 	if self.add(peer) {
// 		peer.protocol.handshake()
// 	}
// }
func (self *hub) cleanup() {
	self.mutex.Lock()
	defer self.mutex.Unlock()

	if self.context.Err() == nil {
		return
	}

	for len(self.peers) > 0 {
		peer := self.peers[0]
		self.removeInternal(peer.adapter)
	}

}

func (self *hub) handleAdd(adapter *adapter) bool {
	self.mutex.Lock()
	defer self.mutex.Unlock()
	if self.context.Err() == nil {
		return false
	}

	peer := newPeer(adapter, self)
	self.peers = append(self.peers, peer)

	peer.start(self.context)

	return true
}

// peer interface
func (self *hub) handleErr(err NetError) {
	self.mutex.Lock()
	defer self.mutex.Unlock()

	self.io.err <- err
}
func (self *hub) handleRemove(adapter *adapter) {
	self.mutex.Lock()
	defer self.mutex.Unlock()

	self.removeInternal(adapter)
}

// /peer interface

// helper
func (self *hub) removeInternal(adapter *adapter) {
	for i := range self.peers {
		peer := self.peers[i]
		if peer.adapter == adapter {
			copy(self.peers[i:], self.peers[i+1:])
			self.peers[len(self.peers)-1] = nil
			self.peers = self.peers[:len(self.peers)-1]
			peer.ctxCancel()
		}
	}
}

func enqueMsg(peer *peerConn, msg *Message) {
	go (func(p *peerConn, m *Message) {
		p.inbox <- msg
	})(peer, msg)
}
