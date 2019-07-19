package chug

import (
	"github.com/buger/jsonparser"
	"github.com/crypto-hug/crypto-hug/fs"
	"github.com/pkg/errors"
)

type StateStore struct {
	fs   *fs.FileSystem
	conf *Config
}

func newStateStore(fs *fs.FileSystem, conf *Config) *StateStore {
	result := new(StateStore)
	result.fs = fs
	result.conf = conf

	return result
}

func (sm *StateStore) HugExists(address string) bool {
	return sm.fs.FileExists(sm.conf.Paths.HugsDir + address + "/hug.json")
}

func (sm *StateStore) HugGetEtag(address string) (string, error) {
	path := sm.conf.Paths.HugsDir + address + "/hug.json"
	if !sm.fs.FileExists(path) {
		return "", nil
	}

	data, err := sm.fs.ReadFile(path)
	if err != nil {
		return "", err
	}

	val, _, _, err := jsonparser.Get(data, "etag")
	if err != nil {
		return "", err
	}

	return string(val), nil
}

func (sm *StateStore) HugSetEtag(address string, etag string) error {
	path := sm.conf.Paths.HugsDir + address + "/hug.json"
	data, err := sm.fs.ReadFile(path)
	if err != nil {
		return err
	}

	data, err = jsonparser.Set(data, []byte(etag), "etag")
	if err != nil {
		return err
	}

	err = sm.fs.WriteFile(path, data)

	return err
}

func (sm *StateStore) HugCreateIfNotExists(address string) error {
	var err error

	path := sm.conf.Paths.HugsDir + address + "/links.json"
	if err = sm.fs.WriteIfNotExists(path, []byte("{}")); err != nil {
		return err
	}

	path = sm.conf.Paths.HugsDir + address + "/hug.json"
	if err = sm.fs.WriteIfNotExists(path, []byte(`{
	"ver": "`+string(HugVersion)+`",
	"addr": "`+address+`",
	"etag": ""
}`)); err != nil {
		return err
	}

	return err
}

func (sm *StateStore) HugAddLinks(ctx *txProcessCtx) error {
	from := ctx.issuerAddress
	to := ctx.validatorAddress

	fromPath := sm.conf.Paths.HugsDir + from + "/links.json"
	toPath := sm.conf.Paths.HugsDir + to + "/links.json"

	fromFile, err := sm.fs.ReadFile(fromPath)
	if err != nil {
		return errors.Wrapf(err, "could not read hug file [%s]", fromPath)
	}

	toFile, err := sm.fs.ReadFile(toPath)
	if err != nil {
		return errors.Wrapf(err, "could not read hug file [%s]", toPath)
	}

	_, res, _, err := jsonparser.Get(fromFile, "links", to)
	if res != jsonparser.NotExist && err != nil {
		return errors.Wrapf(err, "could not parse json file [%s]", fromFile)
	}

	if res != jsonparser.NotExist {
		return errors.Errorf("hug %s already linked to %s", from, to)
	}

	_, res, _, err = jsonparser.Get(toFile, "links", to)
	if res != jsonparser.NotExist && err != nil {
		return errors.Wrapf(err, "could not parse json file [%s]", toFile)
	}

	if res != jsonparser.NotExist {
		return errors.Errorf("hug %s already linked to %s", to, from)
	}

	jsonparser.Set(fromFile, []byte(ctx.address), "links", to, "tx")
	jsonparser.Set(fromFile, []byte(string(ctx.tx.Timestamp)), "links", to, "date")
	if err := sm.fs.WriteFile(fromPath, fromFile); err != nil {
		return errors.Wrapf(err, "could not write hug file [%s]", from)
	}

	jsonparser.Set(toFile, []byte(ctx.address), "links", from, "tx")
	jsonparser.Set(toFile, []byte(string(ctx.tx.Timestamp)), "links", from, "date")
	if err := sm.fs.WriteFile(toPath, toFile); err != nil {
		return errors.Wrapf(err, "could not write hug file [%s]", to)
	}

	return nil
}
