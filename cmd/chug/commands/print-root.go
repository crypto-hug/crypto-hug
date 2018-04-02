package commands

import (
	"../ctx"
	ishell "gopkg.in/abiosoft/ishell.v2"
)

func init() {
	var root = ctx.Root()
	var shell = root.App().Shell
	var cmd = &ishell.Cmd{
		Name: "print",
		Help: "print blockchain",
	}
	cmd.AddCmd(NewPrintBlockchainCmd())

	shell.AddCmd(cmd)

}
