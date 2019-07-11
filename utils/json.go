package utils

import (
	"encoding/json"

	"github.com/pkg/errors"
)

func JsonSerializeStr(obj interface{}) (string, error) {
	var raw, err = JsonSerializeRaw(obj)
	if err != nil {
		return "", err
	}
	var result = string(raw)
	return result, errors.WithStack(err)
}

func JsonSerializeRaw(obj interface{}) ([]byte, error) {
	var result, err = json.Marshal(obj)
	return result, errors.WithStack(err)
}

func JsonParseStr(jsonStr string, obj interface{}) error {
	raw := []byte(jsonStr)
	err := JsonParseRaw(raw, obj)
	return errors.WithStack(err)
}

func JsonParseRaw(raw []byte, obj interface{}) error {
	err := json.Unmarshal(raw, obj)
	return errors.WithStack(err)
}
