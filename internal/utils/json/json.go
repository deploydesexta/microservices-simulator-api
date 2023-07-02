package json

import (
	"github.com/goccy/go-json"
)

func ToJson(i interface{}) ([]byte, error) {
	b, err := json.Marshal(i)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func FromJson(data []byte, i interface{}) error {
	return json.Unmarshal(data, i)
}
