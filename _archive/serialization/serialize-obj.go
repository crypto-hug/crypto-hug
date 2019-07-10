package serialization

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"github.com/pkg/errors"
)

func ObjToJsonStr(obj interface{}) (string, error) {
	var raw, err = ObjToJsonRaw(obj)
	if err != nil {
		return "", err
	}
	var result = string(raw)
	return result, errors.Wrap(err, "ObjToJsonStr")
}

func JsonParse(jsonStr string, obj interface{}) error {
	raw := []byte(jsonStr)
	err := JsonParseRaw(raw, obj)
	return errors.Wrap(err, "JsonParse")
}

func JsonParseRaw(raw []byte, obj interface{}) error {
	err := json.Unmarshal(raw, obj)
	return errors.Wrap(err, "JsonParseRaw")
}

func ObjToJsonRaw(obj interface{}) ([]byte, error) {
	var result, err = json.Marshal(obj)
	return result, errors.Wrap(err, "ObjToJsonRaw")
}

func ObjDecode(data []byte, result interface{}) error {
	var decoder = gob.NewDecoder(bytes.NewReader(data))
	var err = decoder.Decode(result)

	return errors.Wrap(err, "ObjDecode")
}

func ObjEncode(data interface{}) ([]byte, error) {
	var buffer bytes.Buffer
	var encoder = gob.NewEncoder(&buffer)
	var err = encoder.Encode(data)
	if err != nil {
		return nil, errors.Wrap(err, "ObjEncode")
	}

	result := buffer.Bytes()

	return result, nil
}
