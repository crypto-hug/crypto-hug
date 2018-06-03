package go2p_old

import (
	"encoding/binary"
	"io"
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
	ctx      *context
	listener *listernerWrapper
	peerWg   sync.WaitGroup
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
			self.listener.close()
			self.ctx.peers.CloseAll()
			self.peerWg.Wait()
			return

		case out := <-self.ctx.io.out:
			go self.processOutMsg(out)

		case newCon := <-self.ctx.io.newCon:
			self.ctx.peers.Add(newCon)
			go self.read(newCon)

		case endCon := <-self.ctx.io.endCon:
			self.ctx.peers.Close(endCon)
			self.ctx.peers.Remove(endCon)

		}
	}
}

func (self *Server) read(conn *peer) {
	self.peerWg.Add(1)
	for {
		msg, err := self.readMsg(conn)
		if err == nil {
			self.ctx.io.in <- msg
		} else if err == io.EOF {
			self.ctx.io.endCon <- conn
			self.peerWg.Done()
			return
		} else {
			self.ctx.io.err <- &netError{err: err, isTemp: true}
			self.ctx.io.endCon <- conn
			self.peerWg.Done()
			return
		}
	}
}

func (self *Server) writeMsg(msg *Message, p *peer) error {
	headerSBuffer := make([]byte, 4)
	bodySBuffer := make([]byte, 4)

	_, err := p.conn.Read(headerSBuffer)
	if err != nil {
		return errors.Wrap(err, "failed read header size")
	}
	_, err = p.conn.Read(bodySBuffer)
	if err != nil {
		return errors.Wrap(err, "failed read header size")
	}

	headerLength := binary.BigEndian.Uint32(headerSBuffer)
	bodyLength := binary.BigEndian.Uint32(bodySBuffer)

	fullMsgLength := headerLength + bodyLength
	msgBuffer := make([]byte, fullMsgLength)
	var readed uint32 = 0
	for readedLen < fullMsgLength {
		currentReaded, err := p.conn.Read(msgBuffer[readed:])
		if err != nil {
			return msgBuffer, err
		}

		readed += uint32(currentReaded)
	}

	headerBuffer := make([]byte, headerLength)
	bodyBuffer := make([]byte, bodyLength)

	var readedLen uint32 = 0
	for readedLen < msgLen {
		currentReaded, err := self.connection.Read(msgBuffer[readedLen:])
		if err != nil {
			return msgBuffer, err
		}

		readedLen += uint32(currentReaded)
	}

	panic("not impl")
}
func (self *Server) readMsg(conn *peer) (*Message, error) {
	panic("not impl")
}

func (self *Server) processOutMsg(msg *Message) {
	p := self.ctx.peers.byId(msg.receiverId)
	if p == nil {
		// not longer connected to peer, ignore
		// todo: mybe pub a msg in chan
		return
	}

	err := self.writeMsg(msg, p)
	if err != nil {
		self.ctx.io.err <- &netError{err: err, isTemp: true}
		return
	}
}
