package commands

import (
	"strings"

	"../../common/prompt"

	".."
	ishell "gopkg.in/abiosoft/ishell.v2"
)

func init() {
	var root = shell.Get()
	root.Console.AddCmd(&ishell.Cmd{
		Name: "new",
		Help: "creates a new blockchain",
		Func: func(c *ishell.Context) {
			var address = strings.Join(c.Args, " ")
			var exists, err = root.Blockchain().HasGenesisBlock()
			if err != nil {
				prompt.Shared().Warn("command failed. error:")
				prompt.Shared().Warn(err.Error())
				return
			}

			if exists {
				prompt.Shared().Warn("blockain already created! Command ignored")
				return
			}

			err = root.Blockchain().CreateGenesisBlockIfNotExists(address)
			if err != nil {
				prompt.Shared().Warn("command failed. error:")
				prompt.Shared().Warn(err.Error())
				return
			}
		},
	})

}
