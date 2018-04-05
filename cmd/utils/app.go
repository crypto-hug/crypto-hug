package utils

import (
	"os"

	"github.com/crypto-hug/crypto-hug/log"
	ishell "gopkg.in/abiosoft/ishell.v2"
)

type App struct {
	Shell              *ishell.Shell
	printer            *Printer
	InteractiveWelcome string
	log                *log.Logger
}

func NewApp(name string) *App {
	var shell = ishell.New()
	var printer = Printer{shell: shell}
	var logger = log.NewLog("App")
	var result = App{Shell: shell, log: logger, printer: &printer}

	return &result
}

func (self *App) Printer() *Printer {
	return self.printer
}

func (self *App) Run() {
	if len(os.Args) >= 2 && os.Args[1] == "-i" {
		self.log.Info(self.InteractiveWelcome, nil)
		self.Shell.Run()
	} else if len(os.Args) >= 2 {
		var e = self.Shell.Process(os.Args[1:]...)
		if e != nil {
			self.Shell.Println(e.Error())
			os.Exit(1)
		}
	} else {
		self.Shell.Process("help")
	}

}
