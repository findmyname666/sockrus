package sockrus

import (
	"net"

	"github.com/Sirupsen/logrus"
	"github.com/Sirupsen/logrus/formatters/logstash"
)

// Hook represents a connection to a socket
type Hook struct {
	conn net.Conn
}

// NewHook establish a socket connection.
// Protocols allowed are: "udp", "tcp" or "unix".
// For TCP and UDP, address must have the form `host:port`.
// For Unix networks, the address must be a file system path.
func NewHook(protocol, address string) (*Hook, error) {
	conn, err := net.Dial(protocol, address)
	if err != nil {
		return nil, err
	}
	return &Hook{conn: conn}, nil
}

// Fire send log to the defined socket
func (h *Hook) Fire(entry *logrus.Entry) error {
	formatter := logstash.LogstashFormatter{}
	dataBytes, err := formatter.Format(entry)
	if err != nil {
		return err
	}
	if _, err = h.conn.Write(dataBytes); err != nil {
		return err
	}
	return nil
}

// Levels return an array of handled logging levels
func (h *Hook) Levels() []logrus.Level {
	return logrus.AllLevels
}
