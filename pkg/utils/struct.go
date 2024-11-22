package utils

import "encoding/json"

func ConvertTo(source any, target any) error {
	bs, err := json.Marshal(source)
	if err != nil {
		return err
	}
	return json.Unmarshal(bs, target)
}
