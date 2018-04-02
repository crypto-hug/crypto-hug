package prompt

import "fmt"

type Printer interface {
	Debug(msg string,  a ...interface{})
	Info(msg string,  a ...interface{})
	Say(msg string,  a ...interface{})
	Success(msg string,  a ...interface{})
	Warn(msg string,  a ...interface{})
	Panic(msg string,  a ...interface{})
}

type PrintCB func(string)

type DelgatePrinter struct{
	DebugCB PrintCB
	InfoCB PrintCB
	SayCB PrintCB
	SuccessCB PrintCB
	WarnCB PrintCB
	PanicCB PrintCB
}

func format4Print(msg string,  a ...interface{}) string{
	if len(a) > 0 {
		return fmt.Sprintf(msg, a...)
	} else{
		return msg
	}
}


func (self *DelgatePrinter) Debug(msg string,  a ...interface{}) {
	msg = format4Print(msg, a...)
	self.DebugCB(msg)
}


func (self *DelgatePrinter) Info(msg string,  a ...interface{}) {
	msg = format4Print(msg, a...)
	self.InfoCB(msg)
}
func (self *DelgatePrinter) Say(msg string,  a ...interface{}) {
	msg = format4Print(msg, a...)
	self.SayCB(msg)
}
func (self *DelgatePrinter) Success(msg string,  a ...interface{}) {
	msg = format4Print(msg, a...)
	self.SuccessCB(msg)
}
func (self *DelgatePrinter) Warn(msg string,  a ...interface{}) {
	msg = format4Print(msg, a...)
	self.WarnCB(msg)
}

func (self *DelgatePrinter) Panic(msg string,  a ...interface{}) {
	msg = format4Print(msg, a...)
	self.PanicCB(msg)
}

var __shared Printer

func SetupDelegatePrinter(debug PrintCB, info PrintCB, say PrintCB, success PrintCB, warn PrintCB, panic PrintCB){
	var printer = &DelgatePrinter{}
	printer.DebugCB = debug
	printer.InfoCB = info
	printer.SayCB = say
	printer.SuccessCB = success
	printer.WarnCB = warn
	printer.PanicCB = panic
	__shared = printer
}

func Shared() Printer{
	if __shared == nil{
		panic("FATAL: NO printer created")
		return nil
	}

	return __shared
}