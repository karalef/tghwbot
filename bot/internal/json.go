package internal

import (
	"encoding/json"
	"io"
)

// DecodeJSON decodes reader into object or
// returns raw json data if error occured.
func DecodeJSON[T any](r io.Reader) (*T, []byte, error) {
	var v T
	dec := json.NewDecoder(r)
	err := dec.Decode(&v)
	if err == nil {
		return &v, nil, nil
	}
	b, _ := io.ReadAll(io.MultiReader(dec.Buffered(), r))
	return nil, b, err
}

// Empty type is used to avoid spending resources on unmarshaling.
type Empty struct{}

// UnmarshalJSON is json.Unmarshaler implementation.
func (e *Empty) UnmarshalJSON([]byte) error {
	return nil
}
