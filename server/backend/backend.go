package backend

import (
	"fmt"
	"net"

	"github.com/glestaris/uberlist-server"
)

type ClientID string

type Client struct {
	ID ClientID
	IP net.IP
}

type backendClientEntry struct {
	client Client
	ch     chan Message
}

type Message interface {
	Apply(store Store) error
}

type Store interface {
	AddTask(task uberlist.Task) (uberlist.TaskID, error)
	UpdateTask(task uberlist.Task) error
	TaskByID(id uberlist.TaskID) (uberlist.Task, error)
	OrderedTasks() ([]uberlist.Task, error)
}

type Backend struct {
	store   Store
	clients []backendClientEntry
}

func NewBackend(store Store) *Backend {
	return &Backend{
		store:   store,
		clients: []backendClientEntry{},
	}
}

func (b *Backend) SendMessage(client Client, msg Message) error {
	if err := msg.Apply(b.store); err != nil {
		return err
	}

	for _, ce := range b.clients {
		if ce.client.ID == client.ID {
			continue
		}

		go func() {
			ce.ch <- msg
		}()
	}

	return nil
}

func (b *Backend) Tasks() ([]uberlist.Task, error) {
	return b.store.OrderedTasks()
}

func (b *Backend) Subscribe(client Client) (chan Message, error) {
	ch := make(chan Message)
	b.clients = append(b.clients, backendClientEntry{
		client: client,
		ch:     ch,
	})

	return ch, nil
}

func (b *Backend) Unsubscribe(client Client) error {
	for idx, ce := range b.clients {
		if ce.client.ID == client.ID {
			close(ce.ch)
			// remove from list of clients
			b.clients = append(b.clients[:idx], b.clients[idx+1:]...)

			return nil
		}
	}

	return fmt.Errorf("Client with ID '%s' was not found!", client.ID)
}
