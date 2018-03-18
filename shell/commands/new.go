package commands

import (
	".."
	ishell "gopkg.in/abiosoft/ishell.v2"
	"strings"
)

func init() {
	var root = shell.Get()
	root.Console.AddCmd(&ishell.Cmd{
		Name: "new",
		Help: "creates a new blockchain",
		Func: func(c *ishell.Context) {
			var address = strings.Join(c.Args, " ")
			root.CreateBlockchain(address)
		},
	})

}
