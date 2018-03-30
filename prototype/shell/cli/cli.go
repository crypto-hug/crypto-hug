package cli

import (
	"os"

	".."
	"../../common/prompt"
	_ "../commands"
)

func Run() {
	var s = shell.Get()

	if len(os.Args) >= 2 && os.Args[1] == "-i" {
		prompt.Shared().Say("welcome to the crypto hug interactive shell")
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
