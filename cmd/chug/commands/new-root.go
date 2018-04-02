package commands

import (
	"../ctx"
	ishell "gopkg.in/abiosoft/ishell.v2"
)

func init() {
	var root = ctx.Root()
	var shell = root.App().Shell
	var cmd = &ishell.Cmd{
		Name: "new",
		Help: "creates a new wallet or a transaction",
	}
	cmd.AddCmd(NewWalletCmd())
	cmd.AddCmd(NewTransactionCmd())
	shell.AddCmd(cmd)

}
