package json

import (
	"encoding/json"
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
