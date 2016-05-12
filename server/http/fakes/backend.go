package fakes

import (
	"github.com/glestaris/uberlist-server"
	"github.com/glestaris/uberlist-server/backend"
)

type SendMessageEntry struct {
	Client backend.Client
	Msg    backend.Message
}

type FakeBackend struct {
	Messages           []SendMessageEntry
	SendMessageReturns error
}

func (f *FakeBackend) SendMessage(client backend.Client, msg backend.Message) error {
	if f.Messages == nil {
		f.Messages = []SendMessageEntry{}
	}
	f.Messages = append(f.Messages, SendMessageEntry{client, msg})

	return f.SendMessageReturns
}

func (f *FakeBackend) Tasks() ([]uberlist.Task, error) {
	return nil, nil
}

func (f *FakeBackend) Subscribe(client backend.Client) (chan backend.Message, error) {
	return nil, nil
}

func (f *FakeBackend) Unsubscribe(client backend.Client) error {
	return nil
}
