package utils

import (
	"bytes"
	"encoding/json"
	"io"
)

func JsonNopCloser(s interface{}) io.ReadCloser {
	j := Must(json.Marshal(s)).([]byte)

	return io.NopCloser(bytes.NewBuffer(j))
}
