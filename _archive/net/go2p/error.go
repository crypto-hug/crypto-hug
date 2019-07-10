package go2p

type NetError interface {
	error
	IsTemp() bool
}

type netError struct {
	err    error
	isTemp bool
}

func (self *netError) IsTemp() bool {
	return self.isTemp
}

func (self *netError) Error() string {
	return self.err.Error()
}
