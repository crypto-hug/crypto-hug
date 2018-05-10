package go2p

import (
	"net"
	"sync"

	"github.com/pkg/errors"
)

type listernerWrapper struct {
	inner  net.Listener
	closed bool
	mutex  sync.Mutex
}

func (self *listernerWrapper) close() error {
	self.mutex.Lock()
	defer self.mutex.Unlock()
	err := self.inner.Close()
	self.closed = self.closed || err == nil
	return err
}

func (self *listernerWrapper) isClosed() bool {
	self.mutex.Lock()
	defer self.mutex.Unlock()
	return self.closed
}

func (self *listernerWrapper) accept() (net.Conn, error) {
	c, err := self.inner.Accept()

	return c, err
}

// func (self *listernerWrapper) addr() net.Addr {
// 	return self.inner.Addr()
// }

type Server struct {
	ctx *context

	listener *listernerWrapper
}

func NewServer(localAddress string) (*Server, error) {
	result := &Server{}
	result.ctx = newContext()
	listener, err := net.Listen("tcp", localAddress)
	if err != nil {
		return nil, errors.Wrap(err, "could not create tcp listener")
	}

	result.listener = &listernerWrapper{inner: listener}

	go result.accept()
	go result.run()

	return result, nil
}

func (self *Server) accept() {
	for {
		tcpCon, err := self.listener.accept()
		if self.listener.isClosed() {
			return
		} else if err == nil && tcpCon != nil {
			self.ctx.io.newCon <- newPeer(tcpCon)
		} else if tmpErr, ok := err.(net.Error); ok && tmpErr.Temporary() {
			self.ctx.io.err <- &netError{err: err, isTemp: true}
			continue
		} else if err != nil {
			self.ctx.io.err <- &netError{err: err, isTemp: false}
			return
		}
	}
}

func (self *Server) run() {
	for {
		select {

		case <-self.ctx.io.done: // done signal
			return

		case out := <-self.ctx.io.out:
			go self.send(out)
			return

		case newCon := <-self.ctx.io.newCon:
			self.ctx.peers.Add(newCon)
			go self.read(newCon)

		case endCon := <-self.ctx.io.endCon:
			_ = endCon.conn.Close()
			self.ctx.peers.Remove(endCon)

		}
	}
}

func (self *Server) read(newCon *peer) {
	for {

	}
}
