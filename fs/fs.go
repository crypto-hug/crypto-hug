package fs

import "github.com/spf13/afero"

type FileSystem struct {
	afero.Fs
}

func NewFs4Tests() *FileSystem {
	fs := new(FileSystem)
	fs.Fs = afero.NewMemMapFs()

	return fs
}

func NewFileFs(base string) *FileSystem {
	fs := new(FileSystem)
	fs.Fs = afero.NewBasePathFs(afero.NewOsFs(), base)
	return fs
}
