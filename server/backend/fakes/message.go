package fakes

import "github.com/glestaris/uberlist-server/backend"

type FakeMessage struct {
	Applied        bool
	AppliedToStore backend.Store
}

func (f *FakeMessage) Apply(store backend.Store) error {
	f.Applied = true
	f.AppliedToStore = store
	return nil
}
