package main

import (
	_ "github.com/crypto-hug/crypto-hug/cmd/chug/commands"
	"github.com/crypto-hug/crypto-hug/cmd/chug/ctx"
)

func main() {
	ctx.Root().App().Run()
}
