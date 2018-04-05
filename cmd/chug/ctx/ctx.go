package ctx

import (
	"github.com/crypto-hug/crypto-hug/cmd/utils"
	"github.com/crypto-hug/crypto-hug/core"
)

type CHugContext struct {
	app        *utils.App
	blockchain *core.Blockchain
}

var root *CHugContext

func setup() {
	var app = utils.NewApp("crypto ℍug")
	app.InteractiveWelcome = "welcome to the crypto hug interactive shell"
	app.Shell.SetPrompt("ℍ ")

	//utils.SetupPrompt(app.Shell)

	blockchain, err := utils.SetupBlockchain()
	if err != nil {
		utils.FatalExit(err)
	}

	root = &CHugContext{app: app, blockchain: blockchain}
}

func Root() *CHugContext {
	if root == nil {
		setup()
	}

	return root
}

func (self *CHugContext) App() *utils.App {
	utils.AssertExists(self.app, "app")
	return self.app
}

func (self *CHugContext) Blockchain() *core.Blockchain {
	utils.AssertExists(self.blockchain, "blockchain")
	return self.blockchain
}
