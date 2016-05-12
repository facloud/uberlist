package messages

import (
	"github.com/glestaris/uberlist-server"
	"github.com/glestaris/uberlist-server/backend"
)

type UpdateTaskMessage struct {
	Task uberlist.Task
}

func (msg UpdateTaskMessage) Apply(store backend.Store) error {
	return store.UpdateTask(msg.Task)
}
