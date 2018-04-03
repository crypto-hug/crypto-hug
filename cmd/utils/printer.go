package utils

import (
	"fmt"
	color "github.com/fatih/color"
	ishell "gopkg.in/abiosoft/ishell.v2"
)

type Printer struct {
	shell *ishell.Shell
}

func format4Print(msg string, a ...interface{}) string {
	if len(a) > 0 {
		return fmt.Sprintf(msg, a...)
	} else {
		return msg
	}
}

func (self *Printer) Info(msg string, a ...interface{}) {
	msg = format4Print(msg, a...)
	color.Set(color.FgWhite)
	defer color.Unset()
	self.shell.Printf("ℹ️️: %s\n", msg)
}
func (self *Printer) Say(msg string, a ...interface{}) {
	msg = format4Print(msg, a...)
	color.Set(color.FgGreen, color.Bold)
	defer color.Unset()
	self.shell.Printf("💬: %s\n", msg)
}
func (self *Printer) Success(msg string, a ...interface{}) {
	msg = format4Print(msg, a...)
	color.Set(color.FgGreen, color.Bold)
	defer color.Unset()
	self.shell.Printf("✅: %s\n", msg)
}
func (self *Printer) Warn(msg string, a ...interface{}) {
	msg = format4Print(msg, a...)
	color.Set(color.FgYellow, color.Bold)
	defer color.Unset()
	self.shell.Printf("⚠️: %s\n", msg)
}

func (self *Printer) Panic(msg string, a ...interface{}) {
	msg = format4Print(msg, a...)
	color.Set(color.FgRed, color.Bold)
	defer color.Unset()
	self.shell.Printf("🚨️: %s\n", msg)
}
