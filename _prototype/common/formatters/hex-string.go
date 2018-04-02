package formatters

import "strconv"

func HexString(self int64) string {
	var result = strconv.FormatInt(self, 16)
	return result
}
