package chug_test

import (
	"os"
	"testing"

	"github.com/crypto-hug/crypto-hug"
	"github.com/crypto-hug/crypto-hug/fs"
)

func TestNewBlockchain(t *testing.T) {
	d, _ := os.Getwd()
	fs := fs.NewFileFs(d + "/testdata/")
	chug.NewBlockchain(fs)
}
