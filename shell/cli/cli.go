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
		loadBlockchainIfExists(s)
		s.Console.Run()
	} else if len(os.Args) >= 2 {
		loadBlockchainIfExists(s)
		var e = s.Console.Process(os.Args[1:]...)
		if e != nil {
			s.Console.Println(e.Error())
			os.Exit(1)
		}
	} else {
		loadBlockchainIfExists(s)
		s.Console.Process("help")
	}
}

func loadBlockchainIfExists(hugShell *shell.HugShell){
	if _, err := os.Stat("./crypto-hug.db"); os.IsNotExist(err) {
		prompt.Shared().Warn("no blockchain found, please create one with 'new' address")
	} else {
		hugShell.CreateBlockchain("")
	}
}
