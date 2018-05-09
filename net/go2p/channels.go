package go2p

type IO struct {
	In    chan *Message
	Out   chan *Message
	Error chan error
	Done  chan struct{}
}

func NewIO() *IO {
	result := &IO{}
	result.In = make(chan *Message, 250)
	result.Out = make(chan *Message, 250)
	result.Error = make(chan error)
	result.Done = make(chan struct{})

	return result
}
