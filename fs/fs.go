package fs

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	"github.com/v-braun/go-must"

	"github.com/spf13/afero"
)

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

func NewCWDFileFs(subdir string) *FileSystem {
	d, _ := os.Getwd()
	fs := NewFileFs(d + "/" + subdir)
	return fs
}

func (fs *FileSystem) ListDir(path string) ([]os.FileInfo, error) {
	result, err := afero.ReadDir(fs.Fs, path)
	return result, err
}
func (fs *FileSystem) ReadFileMust(path string) []byte {
	data, err := fs.ReadFile(path)
	must.NoError(err, "could not read file %s", path)
	return data
}
func (fs *FileSystem) ReadFile(path string) ([]byte, error) {
	f, err := fs.Open(path)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	data, err := ioutil.ReadAll(f)
	return data, errors.WithStack(err)
}

func (fs *FileSystem) FileName(path string) string {
	return filepath.Base(path)
}

func (fs *FileSystem) FileNameWithoutExt(path string) string {
	filename := fs.FileName(path)
	return strings.TrimSuffix(filename, filepath.Ext(filename))
}

func (fs *FileSystem) WriteFile(path string, data []byte) error {
	var file afero.File
	var err error

	if !fs.FileExists(path) {
		dir := fs.GetDirPath(path)
		fs.MkdirAll(dir, 0755)

		file, err = fs.Create(path)
		if err != nil {
			return err
		}
	} else {
		file, err = fs.OpenFile(path, os.O_WRONLY, os.ModeAppend)
		if err != nil {
			return err
		}
	}

	_, err = file.Write(data)

	return errors.WithStack(err)
}

func (fs *FileSystem) WriteIfNotExists(path string, data []byte) error {
	if !fs.FileExists(path) {
		err := fs.WriteFile(path, data)
		return err
	}

	return nil
}

func (fs *FileSystem) WriteIfNotExistsMust(path string, data []byte) {
	err := fs.WriteIfNotExists(path, data)
	must.NoError(err, "could not write file %s", path)
}

func (fs *FileSystem) FileExists(filePath string) bool {
	_, err := fs.Stat(filePath)
	if os.IsNotExist(err) {
		return false
	}

	must.NoError(err, "could not get file stat")
	return true
}

func (fs *FileSystem) GetDirPath(path string) string {
	if filepath.Ext(path) == "" {
		return path
	}

	return filepath.Dir(path)
}
