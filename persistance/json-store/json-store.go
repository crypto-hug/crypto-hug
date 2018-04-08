package json_store

import (
	"path/filepath"

	"github.com/Jeffail/gabs"
	"github.com/crypto-hug/crypto-hug/errors"
	"github.com/spf13/afero"
)

type FsJsonStore struct {
	fs   afero.Fs
	root string
}

func NewFsJsonStore(root string) *FsJsonStore {
	result := FsJsonStore{}

	if len(root) == 0 {
		result.fs = afero.NewMemMapFs()
	} else {
		result.fs = afero.NewOsFs()

	}

	result.root = root
	// todo: check why 0600 or 0666 cause permission denied for CreateFile
	result.fs.MkdirAll(result.root, 0700)

	return &result
}

func (self *FsJsonStore) Fs() afero.Fs {
	return self.fs
}

func (self *FsJsonStore) ReadJsonFile(path string) (*gabs.Container, error) {
	exist, err := afero.Exists(self.fs, path)
	if err != nil {
		return nil, errors.Wrap(err, "could not check file")
	}
	if exist == false {
		return nil, nil
	}

	handle, err := self.fs.Open(path)
	if err != nil {
		return nil, errors.Wrap(err, "could not open json file")
	}
	defer handle.Close()

	bin, err := afero.ReadAll(handle)
	if err != nil {
		return nil, errors.Wrap(err, "could not read json file")
	}
	if len(bin) <= 0 {
		return nil, nil
	}

	jsonFile, err := gabs.ParseJSON(bin)
	if err != nil {
		return nil, errors.Wrap(err, "could not parse json file")
	}

	return jsonFile, nil
}

func (self *FsJsonStore) WriteJsonFile(path string, file *gabs.Container) error {
	self.fs.MkdirAll(filepath.Dir(path), 0700)

	handle, err := self.fs.Create(path)

	if err != nil {
		return errors.Wrap(err, "FsBlockSink:writeJsonFile")
	}
	defer handle.Close()

	_, err = handle.Write(file.BytesIndent("", " "))
	return err
}

func (self *FsJsonStore) ModifyJsonFile(path string, mod func(file *gabs.Container) (*gabs.Container, error)) error {
	f, err := self.ReadJsonFile(path)
	if err != nil {
		return errors.Wrap(err, "ModifyJsonFile; could not read")
	}

	f, err = mod(f)
	if err != nil {
		return errors.Wrap(err, "ModifyJsonFile: could not modify")
	}

	if f != nil {
		err = self.WriteJsonFile(path, f)
	}

	return err
}
