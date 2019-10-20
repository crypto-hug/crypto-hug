package main

import (
	"log"

	chug "github.com/crypto-hug/crypto-hug"
	"github.com/crypto-hug/crypto-hug/fs"
	"github.com/v-braun/go-must"
)

func main() {
	fs := fs.NewCWDFileFs("blockchain_data")
	cfg, err := chug.NewConfigFromFileOrDefault(fs)
	must.NoError(err, "could not get config")
	host := chug.NewNodeHost(fs, cfg)

	api := newApi(host)
	log.Fatal(api.Run())
}
