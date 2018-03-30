package serialization

import (
	"bytes"
	"encoding/gob"
)

func SimpleDecode(data []byte, result interface{}) error {
	var decoder = gob.NewDecoder(bytes.NewReader(data))
	var err = decoder.Decode(result)

	return err
}
