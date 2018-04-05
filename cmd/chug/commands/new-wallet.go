package commands

import (
	"github.com/crypto-hug/crypto-hug/cmd/chug/ctx"
	"github.com/crypto-hug/crypto-hug/cmd/utils"
	"github.com/crypto-hug/crypto-hug/core"
	ishell "gopkg.in/abiosoft/ishell.v2"
)

func NewWalletCmd() *ishell.Cmd {
	var root = ctx.Root()
	result := &ishell.Cmd{
		Name: "wallet",
		Help: "create a priv/pub key pair",
		Func: func(c *ishell.Context) {
			root.App().Printer().Info("generating key pair ...")

			var wallet, err = core.NewWallet()
			if err != nil {
				utils.FatalExit(err)
				return
			}

			root.App().Printer().Success("keys generated:")
			root.App().Printer().Say("prv key	%v", wallet.PrivAsString())
			root.App().Printer().Say("pub key	%v", wallet.PubAsString())
			root.App().Printer().Say("address	%v", wallet.Address.Address)

		},
	}

	return result
}
