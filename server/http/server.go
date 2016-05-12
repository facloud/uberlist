package http

import (
	"errors"
	"fmt"
	"net"
	"net/http"
	"sync"

	"github.com/Sirupsen/logrus"
	"github.com/glestaris/uberlist-server"
	"github.com/glestaris/uberlist-server/backend"
	"github.com/gorilla/websocket"
)

type Backend interface {
	SendMessage(client backend.Client, msg backend.Message) error
	Tasks() ([]uberlist.Task, error)
	Subscribe(client backend.Client) (chan backend.Message, error)
	Unsubscribe(client backend.Client) error
}

type Server struct {
	backend        Backend
	messageEncoder MessageEncoder

	listenerMutex sync.Mutex
	listener      net.Listener

	server   *http.Server
	upgrader websocket.Upgrader

	connectionsMutex sync.Mutex
	connections      []*conn
}

func NewServer(logger *logrus.Logger, backend Backend, messageEncoder MessageEncoder) *Server {
	s := &Server{
		backend:        backend,
		messageEncoder: messageEncoder,
		upgrader:       websocket.Upgrader{},
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", s.handleConnection)

	s.server = &http.Server{
		Handler: mux,
	}

	return s
}

func (s *Server) Serve(logger *logrus.Logger, ip net.IP, port uint16) error {
	if s.IsListening(logger) {
		return errors.New("already serving")
	}

	var err error
	s.listenerMutex.Lock()
	s.listener, err = net.Listen("tcp", fmt.Sprintf("%s:%d", ip, port))
	if err != nil {
		s.listenerMutex.Unlock()
		return err
	}
	s.listenerMutex.Unlock()

	if err := s.server.Serve(s.listener); err != nil {
		logger.Errorf("serving stopped: %s", err)
	}

	return nil
}

func (s *Server) IsListening(logger *logrus.Logger) bool {
	s.listenerMutex.Lock()
	defer s.listenerMutex.Unlock()

	return s.listener != nil
}

func (s *Server) Close(logger *logrus.Logger) error {
	s.listenerMutex.Lock()
	defer s.listenerMutex.Unlock()

	if s.listener == nil {
		return nil
	}

	if err := s.listener.Close(); err != nil {
		logger.Errorf("closing the listener: %s", err)
	}
	s.listener = nil

	s.closeConnections(logger)

	return nil
}

func (s *Server) handleConnection(w http.ResponseWriter, r *http.Request) {
	// open connection by upgrading the protocol
	wc, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	defer wc.Close()

	c := wrapConn(wc, s.messageEncoder)
	s.registerConnection(c)

	for {
		msg, err := c.readMessage()
		if err != nil {
			// untested - connection dropped or corrupt data
			s.deregisterConnection(c)
			break
		}

		if err := s.backend.SendMessage(backend.Client{}, msg); err != nil {
			if err := c.respondError(err); err != nil {
				logrus.Error("failed to respond: %s", err)
				break // untested
			}
		} else {
			if err := c.respondOK(); err != nil {
				logrus.Error("failed to respond: %s", err)
				break // untested
			}
		}
	}
}

func (s *Server) registerConnection(c *conn) {
	s.connectionsMutex.Lock()
	defer s.connectionsMutex.Unlock()

	s.connections = append(s.connections, c)
}

func (s *Server) deregisterConnection(c *conn) {
	s.connectionsMutex.Lock()
	defer s.connectionsMutex.Unlock()

	for i, ac := range s.connections {
		if ac == c {
			// remove element
			s.connections = append(s.connections[:i], s.connections[i+1:]...)
			return
		}
	}
}

func (s *Server) closeConnections(logger *logrus.Logger) {
	s.connectionsMutex.Lock()
	defer s.connectionsMutex.Unlock()

	for _, conn := range s.connections {
		if err := conn.close(); err != nil {
			logger.Error("closing connection: %s", err)
		}
	}

	// closeConnections does not remove the connections from s.connections. The
	// connection handler of each of them will notice that its connection
	// dropped and it will deregister it.
}
