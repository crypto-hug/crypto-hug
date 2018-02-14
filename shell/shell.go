package shell

import (
	"errors"
	"log"
	"os"

	"../core"
	"../core/storage"
	ishell "gopkg.in/abiosoft/ishell.v2"
)

var shell = new(HugShell)

func init() {
	var s = ishell.New()
	shell.Console = s
	shell.Console.SetPrompt("ℍ ")
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
	cmdCtx.Println("no blockchain loaded!")
	cmdCtx.Println("use 'new' to create one")
	return false
}

func (self *HugShell) Blockchain() *core.Blockchain {
	if self.blockchain == nil {
		PanicExit(errors.New("blockchain not created or loaded"))
		return nil
	}

	return self.blockchain
}

func (self *HugShell) CreateBlockchain() {
	if self.blockchain != nil {
		PanicExit(errors.New("blockchain already created or loaded"))
		return
	}

	var sink, err = storage.NewBoltBlockStore("crypto-hug.db")
	if err != nil {
		PanicExit(err)
		return
	}

	var store core.BlockStore = sink
	chain, err := core.NewBlockchain(&store)
	if err != nil {
		PanicExit(err)
		return
	}

	self.blockchain = chain
}

func PanicExit(err error) {
	log.Fatal(err)
	os.Exit(1)
}
