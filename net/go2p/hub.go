package go2p

import "sync"

//type

type hub struct {
	peers []*peerConn
	mutex sync.Mutex
}

func newHub() *hub {
	result := &hub{}
	return result
}

func (self *hub) add(peer *peerConn) {
	self.mutex.Lock()
	defer self.mutex.Unlock()

	self.peers = append(self.peers, peer)

	go exec(peer, nil)
}

func exec(peer *peerConn, protocol *protocol) {
	peer.execProtocol(protocol)
}

func (self *hub) remove(peer *peerConn) {
	self.mutex.Lock()
	defer self.mutex.Unlock()

	panic("not implemented")
}

func (self *hub) send() {
	self.mutex.Lock()
	defer self.mutex.Unlock()

	panic("not implemented")
}
