package commands

import (
	"../../../core"
	"../../../prompt"
	"../../utils"
	ishell "gopkg.in/abiosoft/ishell.v2"
)

func NewWalletCmd() *ishell.Cmd {
	result := &ishell.Cmd{
		Name: "wallet",
		Help: "create a priv/pub key pair",
		Func: func(c *ishell.Context) {
			prompt.Shared().Info("generating key pair ...")

			var wallet, err = core.NewWallet()
			if err != nil {
				utils.PanicExit(err)
				return
			}

			prompt.Shared().Success("keys generated:")
			prompt.Shared().Say("prv key	%v", wallet.PrivAsString())
			prompt.Shared().Say("pub key	%v", wallet.PubAsString())
			prompt.Shared().Say("address	%v", wallet.Address.Address)

		},
	}

	return result
}
