package go2p

import (
	"net"
)

type dialer struct {
	hub *hub
	io  *IO
}

func newDialer() *dialer {
	result := &dialer{}

	return result
}

func (self *dialer) dial(address string) {
	go func(a string) {
		tcpCon, err := net.Dial("tcp", a)
		if err != nil {
			self.io.err <- &netError{err: err, isTemp: true}
			return
		}

		adapter := newAdapter(tcpCon)
		self.hub.add <- adapter

	}(address)

}
