package redis

import (
	"bytes"
	"encoding/gob"
)

func EncodeGob(data any) ([]byte, error) {
	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(data); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func DecodeGob(data []byte, to any) error {
	buf := bytes.NewBuffer(data)
	if err := gob.NewDecoder(buf).Decode(to); err != nil {
		return err
	}
	return nil
}
