package utils

import (
	"os"
	"../../prompt"
	ishell "gopkg.in/abiosoft/ishell.v2"
)

type App struct{
	Shell *ishell.Shell
	InteractiveWelcome string
}

func NewApp(name string) *App{
	var shell = ishell.New()
	var result = App{Shell: shell}

	return &result
}

func (self *App) Run(){
	if len(os.Args) >= 2 && os.Args[1] == "-i" {
		prompt.Shared().Say(self.InteractiveWelcome)
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