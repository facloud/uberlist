package messages

import (
	"github.com/glestaris/uberlist-server"
	"github.com/glestaris/uberlist-server/backend"
)

type NewTaskMessage struct {
	NewTask uberlist.Task
}

func (msg NewTaskMessage) Apply(store backend.Store) error {
	_, err := store.AddTask(msg.NewTask)
	return err
}
