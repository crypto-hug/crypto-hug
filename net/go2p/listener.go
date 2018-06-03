package go2p

import (
	"net"
	"sync"
)

type listener struct {
	io     *IO
	inner  net.Listener
	closed bool
	mutex  sync.Mutex

	hub *hub
}

func startListener(inner net.Listener, io *IO, hub *hub) *listener {
	result := &listener{}
	result.closed = false
	result.io = io
	result.hub = hub

	go result.start()

	return result
}

func (self *listener) close() error {
	self.mutex.Lock()
	defer self.mutex.Unlock()
	err := self.inner.Close()
	self.closed = self.closed || err == nil
	return err
}

func (self *listener) isClosed() bool {
	self.mutex.Lock()
	defer self.mutex.Unlock()
	return self.closed
}

func (self *listener) start() {
	for {
		tcpCon, err := self.inner.Accept()
		if self.isClosed() {
			return
		} else if err == nil && tcpCon != nil {
			peer := newPeerConn(tcpCon)
			self.hub.add(peer)
		} else if tmpErr, ok := err.(net.Error); ok && tmpErr.Temporary() {
			self.io.err <- &netError{err: err, isTemp: true}
			continue
		} else if err != nil {
			self.io.err <- &netError{err: err, isTemp: false}
			return
		}
	}
}
