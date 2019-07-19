package chug_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/crypto-hug/crypto-hug"
	"github.com/crypto-hug/crypto-hug/fs"
)

func TestNewBlockchain(t *testing.T) {
	d, _ := os.Getwd()
	fs := fs.NewFileFs(d + "/testdata/")

	cfg, err := chug.NewConfigFromFileOrDefault(fs)
	assert.NoError(t, err)

	bc := chug.NewBlockchain(fs, cfg)
	bc.CreateGenesisBlockIfNotExists()
}

// func TestCreateGenesisTx(t *testing.T) {
// 	d, _ := os.Getwd()
// 	fs := fs.NewFileFs(d + "/testdata/")
// 	bc := chug.NewBlockchain(fs)
// 	bc.CreateGenesisTxIfNotExists()
// }
