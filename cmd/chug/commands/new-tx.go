package commands

import (
	"github.com/crypto-hug/crypto-hug/cmd/chug/ctx"
	"github.com/crypto-hug/crypto-hug/cmd/utils"
	"github.com/crypto-hug/crypto-hug/core"
	"github.com/crypto-hug/crypto-hug/core/chug"
	"github.com/crypto-hug/crypto-hug/formatters"
	ishell "gopkg.in/abiosoft/ishell.v2"
)

func spawnHug() {
	var root = ctx.Root()
	addr := core.NewAddressFromStringStrict("VB6QzPAL7P83N48MhoFdLXuroxPmUiphp")
	myPubKey, err := formatters.Base58FromString("aY3JXGjbhvc8gpyepFpGgEhoKmjLYL8piWeKP6cYNGP9U5zs9HqHJASBb7WbD5FevKTJWvhcctZd5w3Em62bquoVB6QzPAL7P83N48MhoFdLXuroxPmUiphp")

	if err != nil {
		utils.FatalExit(err)
	}

	tx, err := chug.NewSpawnHugTransaction(addr, myPubKey)
	if err != nil {
		utils.FatalExit(err)
	}

	err = root.Blockchain().AddTransaction(tx)
	if err != nil {
		utils.FatalExit(err)
	}
}

func spendHug() {
	var root = ctx.Root()
	senderAddr := core.NewAddressFromStringStrict("VB6QzPAL7P83N48MhoFdLXuroxPmUiphp")
	senderPubKey, err := formatters.Base58FromString("aY3JXGjbhvc8gpyepFpGgEhoKmjLYL8piWeKP6cYNGP9U5zs9HqHJASBb7WbD5FevKTJWvhcctZd5w3Em62bquoVB6QzPAL7P83N48MhoFdLXuroxPmUiphp")
	hugAddr := "Xk2ieh6R23dsJPQeR2np8mM85m2NYENsC"
	receipientAddr := "VB6QzPAL7P83N48MhoFdLXuroxPmUiph2"
	// receipientAddr := "VB6QzPAL7P83N48MhoFdLXuroxPmUiphp"

	if err != nil {
		utils.FatalExit(err)
	}

	tx, err := chug.NewSpendHugTransaction(senderAddr.Address, senderPubKey, hugAddr, receipientAddr)
	if err != nil {
		utils.FatalExit(err)
	}

	err = root.Blockchain().AddTransaction(tx)
	if err != nil {
		utils.FatalExit(err)
	}
}

func NewTransactionCmd() *ishell.Cmd {

	result := &ishell.Cmd{
		Name:    "transaction",
		Help:    "create new transaction",
		Aliases: []string{"tx"},
		Func: func(c *ishell.Context) {
			// spawnHug()

			spendHug()

		},
	}

	return result
}
