package go2p

import (
	"context"
	"net"

	"github.com/crypto-hug/crypto-hug/errors"
)

type NetworkConnection struct {
	IO       *IO
	listener *listener
	hub      *hub
	dialer   *dialer
}

func Connect(localAddress string, context context.Context) (*NetworkConnection, error) {
	result := &NetworkConnection{}

	result.IO = newIO()

	result.hub = StartupNewHub(result.IO, context)

	tcpListener, err := net.Listen("tcp", localAddress)
	if err != nil {
		return nil, errors.Wrap(err, "could not create tcp listener")
	}

	result.listener = startListener(tcpListener, result.IO, result.hub)

	result.dialer = newDialer()

	return result, nil

}
