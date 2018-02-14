package commands

import (
	"fmt"

	".."
	ishell "gopkg.in/abiosoft/ishell.v2"
)

func init() {
	var root = shell.Get()
	root.Console.AddCmd(&ishell.Cmd{
		Name: "print",
		Help: "print the blockchain",
		Func: func(c *ishell.Context) {
			print(root, c)
		},
	})
}

func print(root *shell.HugShell, c *ishell.Context) {
	if !root.AssertChainExist(c) {
		return
	}

	var cursor, err = root.Blockchain().Cursor()
	var hasNext = true
	for hasNext {
		if err != nil {
			shell.PanicExit(err)
			return
		}
		var block = (*cursor).Current()
		fmt.Printf(block.PrettyPrint())
		fmt.Println()
		hasNext, err = (*cursor).Next()

		root.Console.Printf("[more]")
		root.Console.ReadLine()

	}
}

// var blocks = self.blockchain.GetBlocks()
// for _, block := range blocks {
// 	fmt.Printf(block.PrettyPrint())
// 	fmt.Println()
// }
