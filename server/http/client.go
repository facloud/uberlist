package http

import (
	"net"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/glestaris/uberlist-server/backend"
	"github.com/gorilla/websocket"
)

var DialTimeout time.Duration = 30 * time.Second

type Client struct {
	conn *conn
}

func NewClient(
	logger *logrus.Logger, endpoint string, messageEncoder MessageEncoder,
) (*Client, error) {
	// assemble dialer
	netDiealer := &net.Dialer{
		Timeout: DialTimeout,
	}
	wsDialer := &websocket.Dialer{
		NetDial: netDiealer.Dial,
	}

	// ignore redirects
	wc, _, err := wsDialer.Dial(endpoint, nil)
	if err != nil {
		return nil, err
	}

	return &Client{wrapConn(wc, messageEncoder)}, nil
}

func (c *Client) SendMessage(logger *logrus.Logger, msg backend.Message) error {
	if err := c.conn.writeMessage(msg); err != nil {
		return err
	}

	return c.conn.readResponse()
}
