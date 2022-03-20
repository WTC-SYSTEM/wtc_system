package utils

import (
	"bytes"
	"encoding/json"
)

func CreateResponse(res map[string]any) ([]byte, error) {
	var b = new(bytes.Buffer)
	m := map[string]map[string]any{}
	m["payload"] = res
	err := json.NewEncoder(b).Encode(m)
	if err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}
