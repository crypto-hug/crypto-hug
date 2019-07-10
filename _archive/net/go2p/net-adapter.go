package go2p

import (
	"context"
	"net"
)

type adapter struct {
	conn net.Conn
	in   chan *Message
	out  chan *Message
	err  chan error
	ctx  context.Context
}

func newAdapter(conn net.Conn) *adapter {
	result := &adapter{conn: conn}
	result.in = make(chan *Message)
	result.out = make(chan *Message)
	result.err = make(chan error)
	result.ctx = nil

	return result
}

func (self *adapter) start(ctx context.Context) {
	if self.ctx != nil {
		panic("adapter already started!")
	}

	go self.processSend()

	go self.processReceive()

	self.ctx = ctx
}

func (self *adapter) cleanup() {
	if self.ctx == nil {
		panic("adapter not running")
	}

	self.conn.Close()

	close(self.in)
	close(self.out)
	close(self.err)

	self.ctx = nil
}

func (self *adapter) processSend() {
	for {
		select {
		case m := <-self.out:
			self.writeMsg(m)
		case <-self.ctx.Done():
			self.cleanup()
			return
		}
	}
}

func (self *adapter) processReceive() {
	for {
		msg, err := self.readMsg()
		if err == nil {
			self.out <- msg
		} else {
			self.err <- err
		}
		if self.ctx.Err() != nil {
			self.cleanup()
			return
		}
	}
}

func (self *adapter) readMsg() (*Message, error) {
	panic("not impl")
}
func (self *adapter) writeMsg(msg *Message) {
	panic("not impl")
}
