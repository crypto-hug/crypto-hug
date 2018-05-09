package go2p

import (
	"net"
	"sync"
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
	IO       *IO
	listener *listernerWrapper
}

func NewServer() *Server {
	result := &Server{}
	result.IO = NewIO()
	result.listener = &listernerWrapper{}
	result.listener.inner = net.ListenTCP()
	result.start()

	return result
}

func (self *Server) start() {

}
