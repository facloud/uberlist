package fakes

import (
	"encoding/json"

	"github.com/glestaris/uberlist-server/backend"
)

type FakeMessageEncoder struct {
	EncodeError error
	DecodeError error
}

func (f *FakeMessageEncoder) Decode(data []byte) (backend.Message, error) {
	if f.DecodeError != nil {
		return nil, f.DecodeError
	}

	msg := new(FakeMessage)
	json.Unmarshal(data, msg)

	return msg, nil
}

func (f *FakeMessageEncoder) Encode(msg backend.Message) ([]byte, error) {
	if f.EncodeError != nil {
		return nil, f.EncodeError
	}

	data, _ := json.Marshal(msg)
	return data, nil
}
