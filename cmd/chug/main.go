package main

import (
	"./ctx"
	_ "./commands"
)


func main() {
	
	ctx.Root().App().Run()

}
