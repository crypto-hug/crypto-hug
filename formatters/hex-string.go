package formatters

import (
	"encoding/hex"
	"strconv"
)

func HexString(self int64) string {
	var result = strconv.FormatInt(self, 16)
	return result
}

func HexStringFromRaw(data []byte) string {
	return hex.EncodeToString(data)
}
