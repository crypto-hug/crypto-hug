package go2p

type IO struct {
	In   <-chan *Message
	Out  chan<- *Message
	Err  <-chan NetError
	Done chan<- struct{}

	in     chan<- *Message
	out    <-chan *Message
	err    chan<- NetError
	done   <-chan struct{}
	newCon chan *peer
	endCon chan *peer
}

func NewIO() *IO {
	result := &IO{}
	in := make(chan *Message, 250)
	result.in = in
	result.In = in

	out := make(chan *Message, 250)
	result.out = out
	result.Out = out

	err := make(chan NetError)
	result.err = err
	result.Err = err

	done := make(chan struct{})
	result.done = done
	result.Done = done

	result.newCon = make(chan *peer)

	return result
}
