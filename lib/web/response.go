package web

import "encoding/json"

type JSONResponse map[string]any

func (r JSONResponse) ByteArray() (res []byte) {
	json, err := json.Marshal(r)
	if err != nil {
		res = nil
		return
	}
	res = json
	return
}
