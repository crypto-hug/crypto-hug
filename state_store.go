package chug

import (
	"github.com/bitly/go-simplejson"
	"github.com/buger/jsonparser"
	"github.com/crypto-hug/crypto-hug/fs"
	"github.com/pkg/errors"
	must "github.com/v-braun/go-must"
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

	data, err = jsonparser.Set(data, jsonparserSetWrapStr(etag), "etag")
	if err != nil {
		return err
	}

	err = sm.fs.WriteFile(path, data)

	return err
}

func jsonparserSetWrapStr(val string) []byte {
	return []byte("\"" + val + "\"")
}

func (sm *StateStore) HugCreateIfNotExists(address string) error {
	var err error

	path := sm.conf.Paths.HugsDir + address + "/links.json"
	sm.fs.WriteIfNotExistsMust(path, []byte(`{
		"links": {}
	}`))

	path = sm.conf.Paths.HugsDir + address + "/hug.json"
	sm.fs.WriteIfNotExistsMust(path, []byte(`{
	"ver": "`+string(HugVersion)+`",
	"addr": "`+address+`",
	"etag": ""
}`))

	return err
}

func (sm *StateStore) HugAddLinks(ctx *txProcessCtx) error {
	from := ctx.issuerAddress
	to := ctx.validatorAddress

	fromPath := sm.conf.Paths.HugsDir + from + "/links.json"
	toPath := sm.conf.Paths.HugsDir + to + "/links.json"

	fromFile := readJSONFromFsMust(sm.fs, fromPath)
	toFile := readJSONFromFsMust(sm.fs, toPath)

	res := fromFile.GetPath("links", to)
	if res.Interface() != nil {
		return errors.Errorf("hug %s already linked to %s", from, to)
	}

	res = toFile.GetPath("links", to)
	if res.Interface() != nil {
		return errors.Errorf("hug %s already linked to %s", to, from)
	}

	fromFile.SetPath([]string{"links", to, "tx"}, ctx.address)
	fromFile.SetPath([]string{"links", to, "date"}, ctx.tx.Timestamp)
	storeJSONToFsMust(fromFile, sm.fs, fromPath)

	toFile.SetPath([]string{"links", from, "tx"}, ctx.address)
	toFile.SetPath([]string{"links", from, "date"}, ctx.tx.Timestamp)
	storeJSONToFsMust(toFile, sm.fs, toPath)

	return nil
}

func readJSONFromFsMust(fs *fs.FileSystem, path string) *simplejson.Json {
	data, err := fs.ReadFile(path)
	must.NoError(err, "could not read json file [%s]", path)

	file, err := simplejson.NewJson(data)
	must.NoError(err, "could parse json file [%s]", path)

	return file
}

func storeJSONToFsMust(json *simplejson.Json, fs *fs.FileSystem, path string) {
	data, err := json.Encode()
	must.NoError(err, "could not gen raw data from json to store in file [%s]", path)

	err = fs.WriteFile(path, data)
	must.NoError(err, "could not write file [%s]", path)
}
