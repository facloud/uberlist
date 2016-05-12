package http

import "github.com/glestaris/uberlist-server/backend"

type MessageEncoder interface {
	Encode(msg backend.Message) ([]byte, error)
	Decode(data []byte) (backend.Message, error)
}

type ConnectionDroppedError string

func (e ConnectionDroppedError) Error() string {
	return string(e)
}
