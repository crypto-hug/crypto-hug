package go2p

import (
	"fmt"
	"net"
)

type MiddlewareContext interface {
	Next()
	Fail(err error)
	Done()
	Message() *Message
	Connection() *net.Conn
}

type Middleware interface {
	Name() string
}
type InMiddleware interface {
	Middleware
	Handle(message *Message)
}

type OutMiddleware interface {
	Middleware
	Handle(message *Message)
}
type AcceptMiddleware interface {
	Middleware
	Handle()
}

type middlewareResult string

const (
	middlewareResultFail   middlewareResult = "FAIL"
	middlewareResultDone   middlewareResult = "DONE"
	middlewareResultNext   middlewareResult = "NEXT"
	middlewareResultNotSet middlewareResult = "NOTSET"
)

type middlewareContext struct {
	name    string
	message *Message
	result  middlewareResult
	err     error
	conn    *net.Conn
}

func (self *middlewareContext) Next() {
	self.result = middlewareResultNext
}
func (self *middlewareContext) Done() {
	self.result = middlewareResultDone
}
func (self *middlewareContext) Fail(err error) {
	self.result = middlewareResultFail
	if err == nil {
		panic(fmt.Sprintf("middleware %s failed withot to provide an error", self.name))
	}
	self.err = err
}
func (self *middlewareContext) Connection() *net.Conn {
	return self.conn
}

func (self *middlewareContext) mustHaveResult() {
	if self.result == middlewareResultNotSet {
		panic(fmt.Sprintf("middleware %s should call next/done or fail", self.name))
	}
}
