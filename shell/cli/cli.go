package cli

import (
	"os"

	".."
	_ "../commands"
)

func Run() {
	var s = shell.Get()

	if _, err := os.Stat("./crypto-hug.db"); os.IsNotExist(err) {
		s.Console.Println("no blockchain found, please create one with 'new' address")
	} else {
		s.CreateBlockchain()
	}

	if len(os.Args) >= 2 && os.Args[1] == "-i" {
		s.Console.Println("crypto hug interactive shell")
		s.Console.Run()
	} else if len(os.Args) >= 2 {
		var e = s.Console.Process(os.Args[1:]...)
		if e != nil {
			s.Console.Println(e.Error())
			os.Exit(1)
		}
	} else {
		s.Console.Process("help")
	}

}
