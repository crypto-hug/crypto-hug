package go2p

import (
	"context"
	"io"
)

type peerConn struct {
	addr       string
	adapter    *adapter
	attributes map[string]interface{}

	ctx       context.Context
	ctxCancel context.CancelFunc

	hub    *hub
	outbox chan *Message
	inbox  chan *Message

	//protocol *protocol
}

func newPeer(adapter *adapter, hub *hub) *peerConn {
	result := &peerConn{}
	result.adapter = adapter
	result.attributes = make(map[string]interface{})
	result.hub = hub

	return result
}

func (self *peerConn) start(parentCtx context.Context) {
	if self.ctx != nil {
		panic("connection already started!")
	}

	myCtx, myCancelFunc := context.WithCancel(parentCtx)
	self.ctx = myCtx
	self.ctxCancel = myCancelFunc

	self.adapter.start(self.ctx)

	go self.processReceive()

	go self.processSend()

	go self.processError()

}

func (self *peerConn) cleanup() {
	if self.ctx == nil {
		panic("connection already stopped!")
	}

	close(self.outbox)
	close(self.inbox)

	self.ctx = nil
}

func (self *peerConn) processReceive() {
	for {
		select {
		case <-self.ctx.Done():
			self.cleanup()
			return
		case rawMsg := <-self.adapter.in:
			msg, err := self.processIncommingMessage(rawMsg)
			if err != nil {
				self.hub.err <- err
			} else if msg != nil {
				self.inbox <- msg
			}
		}
	}
}
func (self *peerConn) processSend() {
	for {
		select {
		case <-self.ctx.Done():
			self.cleanup()
			return
		case msg := <-self.outbox:
			if msg != nil {
				rawMsg, err := self.processOutgoingMessage(msg)
				if err != nil {
					self.hub.err <- err
				} else if rawMsg != nil {
					self.adapter.out <- rawMsg
				}
			}
		}
	}
}

func (self *peerConn) processError() {
	for {
		select {
		case <-self.ctx.Done():
			self.cleanup()
			return
		case err := <-self.adapter.err:
			if err == io.EOF {
				self.hub.remove <- self.adapter
			} else {
				e := &netError{err: err, isTemp: true}
				self.hub.err <- e
			}
		}
	}
}

func (self *peerConn) processIncommingMessage(msg *Message) (*Message, NetError) {
	panic("not impl")
}
func (self *peerConn) processOutgoingMessage(msg *Message) (*Message, NetError) {
	panic("not impl")
}
