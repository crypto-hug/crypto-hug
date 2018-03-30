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
			if !root.AssertGenesisBlockExists(c) {
				return
			}

			var data = strings.Join(c.Args, " ")
			var err = root.Blockchain().AddNewBlock(data)
			if err != nil {
				shell.PanicExit(err)
			}

		},
	})

}
