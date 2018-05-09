package go2p

type IO struct {
	In    chan *Message
	Out   chan *Message
	Error chan error
	Done  chan struct{}
}
