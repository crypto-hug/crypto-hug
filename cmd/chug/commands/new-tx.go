package commands

import (
	"../../../core"
	"../../../prompt"
	"../../utils"
	"../ctx"
	ishell "gopkg.in/abiosoft/ishell.v2"
)

func NewTransactionCmd() *ishell.Cmd {
	var root = ctx.Root()

	result := &ishell.Cmd{
		Name:    "transaction",
		Help:    "create new transaction",
		Aliases: []string{"tx"},
		Func: func(c *ishell.Context) {
			addr, err := core.NewAddressFromString("VB6QzPAL7P83N48MhoFdLXuroxPmUiphp")
			if err != nil {
				utils.PanicExit(err)
			}

			tx, err := core.NewSpawnHugTransaction(addr)
			if err != nil {
				utils.PanicExit(err)
			}

			err = root.Blockchain().AddTransaction(tx)
			if err != nil {
				utils.PanicExit(err)
			}

			txTypes := []string{string(core.SpawnHugTxType), string(core.DonateHugTxType)}
			choosed := c.MultiChoice(txTypes, "tx type:")

			c.Print("your priv key; ")
			priv := c.ReadLine()

			c.Print("your address: ")
			addrStr := c.ReadLine()

			c.Print("destination address: ")
			dest := c.ReadLine()

			prompt.Shared().Info("priv: %s | addr: %s | dest: %s | choosed: %s", priv, addrStr, dest, choosed)

		},
	}

	return result
}