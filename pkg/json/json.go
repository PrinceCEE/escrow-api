package json

import (
	"encoding/json"
	"fmt"
	"io"
)

func ReadJSON(b io.ReadCloser, dst any) error {
	dec := json.NewDecoder(b)
	err := dec.Decode(dst)

	if err != nil {
		return err
	}

	return nil
}

func WriteJSON(data any) ([]byte, error) {
	return json.MarshalIndent(data, "", "    ")
}

func MarshalIndent(v any) (string, error) {
	js, err := json.MarshalIndent(v, "", "    ")
	if err != nil {
		return "", err
	}

	return string(js), nil
}

func Marshal(v any) ([]byte, error) {
	return json.Marshal(v)
}

func ReadTypedJSON[T any](b io.ReadCloser) (*T, error) {
	out := new(T)
	fmt.Println(out)
	dec := json.NewDecoder(b)
	err := dec.Decode(out)

	if err != nil {
		return nil, err
	}

	return out, nil
}

func Unmarshal(data []byte, dst any) error {
	return json.Unmarshal(data, dst)
}
