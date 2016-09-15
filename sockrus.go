package sockrus

import (
	"net"

	"github.com/Sirupsen/logrus"
)

// Hook represents a connection to a socket
type Hook struct {
	conn     net.Conn
	protocol string
	address  string
}

// NewHook establish a socket connection.
// Protocols allowed are: "udp", "tcp", "unix" (corresponds to SOCK_STREAM),
// "unixdomain" (corresponds to SOCK_DGRAM) or "unixpacket" (corresponds to SOCK_SEQPACKET).
//
// For TCP and UDP, address must have the form `host:port`.
//
// For Unix networks, the address must be a file system path.
func NewHook(protocol, address string) (*Hook, error) {
	return &Hook{conn: nil, protocol: protocol, address: address}, nil
}

// Fire send log to the defined socket
func (h *Hook) Fire(entry *logrus.Entry) error {
	var err error
	if h.conn == nil {
		err = h.dialSock()
		if err != nil {
			return err
		}
	}
	formatter := logrus.JSONFormatter{}
	dataBytes, err := formatter.Format(entry)
	if err != nil {
		return err
	}
	if _, err = h.conn.Write(dataBytes); err != nil {
		_ = h.closeSock() // #nosec
		return err
	}
	return nil
}

// Levels return an array of handled logging levels
func (h *Hook) Levels() []logrus.Level {
	return logrus.AllLevels
}

// closeSock tries to close connection to Unix socket
func (h *Hook) closeSock() error {
	if h.conn == nil {
		return nil
	}
	err := h.conn.Close()
	h.conn = nil
	return err
}

// dialSock tries to connect to Unix socket
func (h *Hook) dialSock() error {
	conn, err := net.Dial(h.protocol, h.address)
	if err != nil {
		h.conn = nil
		return err
	}
	h.conn = conn
	return nil
}
