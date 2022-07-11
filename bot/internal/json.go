package internal

import (
	"encoding/json"
	"io"
)

// DecodeJSON decodes reader into object or
// returns raw json data if error occured.
func DecodeJSON(r io.Reader, res interface{}) ([]byte, error) {
	if res == nil {
		return io.ReadAll(r)
	}
	dec := json.NewDecoder(r)
	err := dec.Decode(res)
	if err == nil {
		return nil, nil
	}
	return io.ReadAll(io.MultiReader(dec.Buffered(), r))
}
