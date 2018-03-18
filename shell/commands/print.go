package commands

import (
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
	if !root.AssertGenesisBlockExists(c) {
		return
	}

	var chain = root.Blockchain()
	var proof = chain.ProofAlgorithm()
	var cursor, err = chain.Cursor()
	var hasNext = true

	for hasNext {
		if err != nil {
			shell.PanicExit(err)
			return
		}
		var block = (*cursor).Current()
		root.Console.Printf(block.PrettyPrint(proof))
		root.Console.Println()
		hasNext, err = (*cursor).Next()

		if hasNext {

			root.Console.Printf("more [ENTER] | exit [e]: ")
			var s = root.Console.ReadLine()
			if s == "e" || s == "E" {
				root.Console.Printf("\n\n")
				break
			}
		}
	}

}

// var blocks = self.blockchain.GetBlocks()
// for _, block := range blocks {
// 	fmt.Printf(block.PrettyPrint())
// 	fmt.Println()
// }
