package commands

import (
	//	"../../../core"
	"../../utils"
	"../ctx"
	color "github.com/fatih/color"
	ishell "gopkg.in/abiosoft/ishell.v2"
)

func NewPrintBlockchainCmd() *ishell.Cmd {
	var root = ctx.Root()
	result := &ishell.Cmd{
		Name: "blockchain",
		Help: "prints the entire blockahin",
		Func: func(c *ishell.Context) {
			blockchain := root.Blockchain()
			var cursor, err = blockchain.Cursor()
			if cursor.Current() == nil {
				root.App().Printer().Warn("blockchain is empty")
				return
			}
			var hasNext = true

			for hasNext {
				if err != nil {
					utils.FatalExit(err)
					return
				}
				var block = (*cursor).Current()

				color.Set(color.FgHiGreen, color.Bold)
				c.Print("BLOCK:")
				color.Unset()

				color.Set(color.FgHiCyan)
				c.Printf(block.PrettyPrint())
				color.Unset()
				c.Println()

				color.Set(color.FgHiGreen, color.Bold)
				c.Printf("\nTRANSACTIONS(%d):", len(block.Transactions))
				color.Unset()

				color.Set(color.FgHiMagenta)
				for _, tx := range block.Transactions {
					c.Println(tx.PrettyPrint())
				}
				color.Unset()

				//c.Println("\n\n")
				c.Println()
				c.Println()

				hasNext, err = (*cursor).Next()

				// if hasNext {
				// 	var s = c.ReadLine()
				// 	if s == "e" || s == "E" {
				// 		c.Printf("\n\n")
				// 		break
				// 	}
				// }
			}

		},
	}

	return result
}
