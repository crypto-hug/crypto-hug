package go2p

type protocol interface {
	Run(adapter *adapter, message *Message) ProtocolResult
}

type ProtocolResult struct {
	msg  *Message
	next bool
	err  NetError
}

// type protocol struct {
// 	hub     *hub
// 	io      *IO
// 	peer    *peerConn
// 	running bool

// 	in  chan *Message
// 	out chan *Message
// }

// func runProtocol(hub *hub, io *IO, peer *peerConn) *protocol {
// 	result := &protocol{
// 		hub: hub,
// 		io:  io}

// 	go result.listen()
// 	go result.send()

// 	return result
// }

// func (self *protocol) listen() {
// 	for self.running {
// 		msg, err := self.readMsg(sefl.peer)
// 		if err == nil {
// 			self.handleMessage(msg)
// 		} else if err == io.EOF {
// 			self.handleEnd()
// 		} else {
// 			self.handleErr(err)
// 		}
// 	}
// }

// func (self *protocol) handleMessage(msg *Message) {
// 	path := msg.header.FieldPath()
// 	if path == "HANDSHAKE" {
// 		self.handleMessageHandshake(msg)
// 		return
// 	}

// 	if pubKey, ok := self.peer.attributes["PUBKEY"]; !ok {
// 		self.handleMessageWithoutHandshake(msg)
// 		return
// 	}

// 	self.handleMessageCommon(msg)
// }
// func (self *protocol) handleMessageWithoutHandshake(msg *Message) {
// 	e := &netError{err: errors.Errorf("received message from peer %s without a handshake, end connection with peer", msg.header.senderAddr), isTemp: true}
// 	self.handleErr(e)
// }
// func (self *protocol) handleMessageHandshake(msg *Message) {
// 	self.peer.attributes["PUBKEY"] = msg.body
// }
// func (self *protocol) handleMessageCommon(msg *Message) {
// 	self.io.in <- msg
// }

// func (self *protocol) handleEnd() {
// 	self.running = false
// 	self.hub.remove(self.peer)
// }

// func (self *protocol) handleErr(err NetError) {
// 	self.running = false
// 	self.hub.remove(self.peer)
// 	go func(io *IO, e NetError) {
// 		io.err <- e
// 	}(self.io, err)
// }

// func (self *protocol) handshake() {
// 	panic("not implemented")
// }

// func (self *protocol) discover() {
// 	panic("not implemented")
// }

// func (self *protocol) bye() {
// 	panic("not implemented")
// }

// func (self *protocol) message(msg *Message) {
// 	panic("not implemented")
// }
