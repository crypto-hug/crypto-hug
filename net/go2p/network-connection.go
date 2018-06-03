package go2p

import (
	"net"

	"github.com/crypto-hug/crypto-hug/errors"
)

type NetworkConnection struct {
	IO       *IO
	listener *listener
	hub      *hub
	dialer   *dialer
}

func Connect(localAddress string) (*NetworkConnection, error) {
	result := &NetworkConnection{}

	result.IO = newIO()

	result.hub = newHub()

	tcpListener, err := net.Listen("tcp", localAddress)
	if err != nil {
		return nil, errors.Wrap(err, "could not create tcp listener")
	}
	result.listener = startListener(tcpListener, result.IO, result.hub)

	result.dialer = newDialer()

	return result, nil
}
