package wait

type EndSignal interface {
	HasSignal() bool
	SendSignal()
	C() <-chan struct{}
}

type endSignal struct {
	c     chan struct{}
	ended bool
}

func NewEndSignal() EndSignal {
	result := endSignal{}
	result.ended = false
	result.c = make(chan struct{})

	return &result
}

func (self *endSignal) SendSignal() {
	self.ended = true
	close(self.c)
}

func (self *endSignal) HasSignal() bool {
	return self.ended
}

func (self *endSignal) C() <-chan struct{} {
	return self.c
}
