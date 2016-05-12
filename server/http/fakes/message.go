package fakes

import "github.com/glestaris/uberlist-server/backend"

type FakeMessage struct {
	Id string
}

func (f *FakeMessage) Apply(_ backend.Store) error {
	return nil
}
