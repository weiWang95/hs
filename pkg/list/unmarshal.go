package list

import (
	"encoding/json"
	"reflect"
)

var _valueMap = map[string]any{}

func RegisterValue(v any) {
	t := reflect.TypeOf(v)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	_valueMap[t.Name()] = v
}

type valueWithName struct {
	Name string          `json:"name"`
	Data json.RawMessage `json:"data"`
}

func marshal(v any) ([]byte, error) {
	t := reflect.TypeOf(v)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	data, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}

	d := valueWithName{
		Name: t.Name(),
		Data: data,
	}
	return json.Marshal(d)
}

func unmarshal(data []byte) (any, error) {
	var d valueWithName
	if err := json.Unmarshal(data, &d); err != nil {
		return nil, err
	}

	v, ok := _valueMap[d.Name]
	if !ok {
		var r map[string]any
		err := json.Unmarshal(data, &r)
		return r, err
	}

	r := reflect.New(reflect.TypeOf(v)).Interface()

	err := json.Unmarshal(d.Data, r)
	return r, err
}
