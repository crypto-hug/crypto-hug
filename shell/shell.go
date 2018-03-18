package shell

import (
	"errors"
	"os"

	"../common/prompt"
	"../core"
	"../core/storage"
	color "github.com/fatih/color"
	ishell "gopkg.in/abiosoft/ishell.v2"
)

var shell = new(HugShell)

func init() {
	var s = ishell.New()
	shell.Console = s
	shell.Console.SetPrompt("ℍ ")

	prompt.SetupDelegatePrinter(
		func(msg string) { // debzg
			color.Set(color.FgHiBlack, color.Italic)
			defer color.Unset()
			shell.Console.Printf("🔘: %s\n", msg)
		},
		func(msg string) { // info
			color.Set(color.FgWhite)
			defer color.Unset()
			shell.Console.Printf("ℹ️️: %s\n", msg)
		},
		func(msg string) { // say
			color.Set(color.FgGreen, color.Bold)
			defer color.Unset()
			shell.Console.Printf("💬: %s\n", msg)
		},
		func(msg string) { // success
			color.Set(color.FgGreen, color.Bold)
			defer color.Unset()
			shell.Console.Printf("✅: %s\n", msg)
		},
		func(msg string) { // warning
			color.Set(color.FgYellow, color.Bold)
			defer color.Unset()
			shell.Console.Printf("⚠️: %s\n", msg)
		},
		func(msg string) { // panic
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

func (self *HugShell) AssertGenesisBlockExists(c *ishell.Context) bool {
	var hasGenBlock, err = self.Blockchain().HasGenesisBlock()
	if err != nil {
		PanicExit(err)
		return false
	}

	if !hasGenBlock {
		prompt.Shared().Warn("no genesis block found. command ignored")
		return false
	}

	return hasGenBlock
}

func (self *HugShell) Blockchain() *core.Blockchain {
	if self.blockchain == nil {
		self.createBlockchain()
	}

	return self.blockchain
}

func (self *HugShell) createBlockchain() {
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

	chain := core.NewBlockchain(sink)
	self.blockchain = chain
}

func PanicExit(err error) {
	prompt.Shared().Panic(err.Error())
	//log.Fatal(err)
	os.Exit(1)
}
