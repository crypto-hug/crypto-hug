package ctx

import (
	"../../../core"
	"../../../persistance"
	"../../utils"
)

type CHugContext struct {
	app        *utils.App
	db         *hugdb.BoltDb
	blockchain *core.Blockchain
}

var root *CHugContext

func setup() {
	var app = utils.NewApp("crypto ℍug")
	app.InteractiveWelcome = "welcome to the crypto hug interactive shell"
	app.Shell.SetPrompt("ℍ ")

	utils.SetupPrompt(app.Shell)

	db := utils.SetupDb()
	blockchain := utils.SetupBlockchain(db)

	root = &CHugContext{app: app, db: db, blockchain: blockchain}
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
func (self *CHugContext) DB() *hugdb.BoltDb {
	utils.AssertExists(self.db, "db")
	return self.db
}
func (self *CHugContext) Blockchain() *core.Blockchain {
	utils.AssertExists(self.blockchain, "blockchain")
	return self.blockchain
}
