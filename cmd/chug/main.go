package main

import (
	_ "./commands"
	"./ctx"
)

func main() {
	ctx.Root().App().Run()
}
