package protocol

import "encoding/json"

var Default Coder = JsonCoder{}

type Coder interface {
	Encode(v any) ([]byte, error)
	Decode(data []byte, v any) error
}

type JsonCoder struct{}

func (JsonCoder) Encode(v any) ([]byte, error) {
	return json.Marshal(v)
}

func (JsonCoder) Decode(data []byte, v any) error {
	return json.Unmarshal(data, v)
}
