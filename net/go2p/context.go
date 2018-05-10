package go2p

type Context interface {
	IO() *IO
	Peers() *Peers
}

type context struct {
	io    *IO
	peers *Peers
}

func newContext() *context {
	result := &context{
		io:    NewIO(),
		peers: newPeers(),
	}

	return result
}

func (self *context) IO() *IO {
	return self.io
}

func (self *context) Peers() *Peers {
	return self.peers
}
