package commands

import (
	"strings"

	".."
	ishell "gopkg.in/abiosoft/ishell.v2"
)

func init() {
	var root = shell.Get()
	root.Console.AddCmd(&ishell.Cmd{
		Name: "add",
		Help: "add a new block to the blockchain",
		Func: func(c *ishell.Context) {
			if !root.AssertChainExist(c) {
				return
			}

			var data = strings.Join(c.Args, " ")
			root.Blockchain().AddNewBlock(data)

		},
	})

}
