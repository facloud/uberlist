package http

import (
	"errors"

	"github.com/glestaris/uberlist-server/backend"
	"github.com/gorilla/websocket"
)

type conn struct {
	conn           *websocket.Conn
	messageEncoder MessageEncoder
}

func wrapConn(wc *websocket.Conn, messageEncoder MessageEncoder) *conn {
	return &conn{
		conn:           wc,
		messageEncoder: messageEncoder,
	}
}

func (c *conn) writeMessage(msg backend.Message) error {
	data, err := c.messageEncoder.Encode(msg)
	if err != nil {
		return err
	}

	if err := c.conn.WriteMessage(websocket.TextMessage, data); err != nil {
		return ConnectionDroppedError(err.Error())
	}

	return nil
}

func (c *conn) readMessage() (backend.Message, error) {
	_, data, err := c.conn.ReadMessage()
	if err != nil {
		return nil, ConnectionDroppedError(err.Error())
	}

	return c.messageEncoder.Decode(data)
}

func (c *conn) respondOK() error {
	connErr := c.conn.WriteMessage(websocket.TextMessage, []byte("ok"))
	if connErr != nil {
		return ConnectionDroppedError(connErr.Error())
	}

	return nil
}

func (c *conn) respondError(err error) error {
	connErr := c.conn.WriteMessage(websocket.TextMessage, []byte(err.Error()))
	if connErr != nil {
		return ConnectionDroppedError(connErr.Error())
	}

	return nil
}

func (c *conn) readResponse() error {
	_, data, err := c.conn.ReadMessage()
	if err != nil {
		return ConnectionDroppedError(err.Error())
	}

	dataStr := string(data)
	if dataStr == "ok" {
		return nil
	}

	return errors.New(dataStr)
}

func (c *conn) close() error {
	return c.conn.Close()
}
