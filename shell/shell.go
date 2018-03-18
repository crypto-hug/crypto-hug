package shell

import (
	"errors"
	"os"


	"../core"
	"../core/storage"
	ishell "gopkg.in/abiosoft/ishell.v2"
	color "github.com/fatih/color"
	"../common/prompt"
)

var shell = new(HugShell)

func init() {
	var s = ishell.New()
	shell.Console = s
	shell.Console.SetPrompt("ℍ ")


	prompt.SetupDelegatePrinter(
		func(msg string){ // debzg
			color.Set(color.FgHiBlack, color.Italic)
			defer color.Unset()
			shell.Console.Printf("🔘: %s\n", msg)
		},
		func(msg string){ // info
			color.Set(color.FgWhite)
			defer color.Unset()
			shell.Console.Printf("ℹ️️: %s\n", msg)
		},
		func(msg string){ // say
			color.Set(color.FgGreen, color.Bold)
			defer color.Unset()
			shell.Console.Printf("💬: %s\n", msg)
		},
		func(msg string){ // success
			color.Set(color.FgGreen, color.Bold)
			defer color.Unset()
			shell.Console.Printf("✅: %s\n", msg)
		},
		func(msg string){ // warning
			color.Set(color.FgYellow, color.Bold)
			defer color.Unset()
			shell.Console.Printf("⚠️: %s\n", msg)
		},
		func(msg string){ // panic
			color.Set(color.FgRed, color.Bold)
			defer color.Unset()
			shell.Console.Printf("🚨️: %s\n", msg)

		})
}

type HugShell struct {
	Console    *ishell.Shell
	blockchain *core.Blockchain
}

func Get() *HugShell {
	return shell
}

func (self *HugShell) AssertChainExist(cmdCtx *ishell.Context) bool {
	if self.blockchain != nil {
		return true
	}

	prompt.Shared().Warn("no blockchain loaded!")
	prompt.Shared().Warn("use 'new' to create one")
	return false
}

func (self *HugShell) Blockchain() *core.Blockchain {
	if self.blockchain == nil {
		PanicExit(errors.New("blockchain not created or loaded"))
		return nil
	}

	return self.blockchain
}

func (self *HugShell) CreateBlockchain(rewardAddress string) {
	if self.blockchain != nil {
		PanicExit(errors.New("blockchain already created or loaded"))
		return
	}


	var filePath = "crypto-hug.db"
	prompt.Shared().Debug("use blockchain file: %s", filePath)
	var sink, err = storage.NewBoltBlockStore(filePath)
	if err != nil {
		PanicExit(err)
		return
	}

	chain, err := core.NewBlockchain(sink, rewardAddress)
	if err != nil {
		sink.Close()
		os.Remove(filePath)
		PanicExit(err)
		return
	}

	self.blockchain = chain
}

func PanicExit(err error) {
	prompt.Shared().Panic(err.Error())
	//log.Fatal(err)
	os.Exit(1)
}
