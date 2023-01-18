package sockrus

import (
	"net"
	"time"

	logrus_logstash "github.com/bshuster-repo/logrus-logstash-hook"
	"github.com/sirupsen/logrus"
)

// Hook represents a connection to a socket
type Hook struct {
	formatter  logrus_logstash.LogstashFormatter
	protocol   string // Protocol to the socket.
	address    string // Address to the socket.
	addNewline bool   // Toggle to send newline after every message to socket.
}

// NewHook establish a socket connection.
// Protocols allowed are: "udp", "tcp", "unix" (corresponds to SOCK_STREAM),
// "unixdomain" (corresponds to SOCK_DGRAM) or "unixpacket" (corresponds to SOCK_SEQPACKET).
//
// For TCP and UDP, address must have the form `host:port`.
//
// For Unix networks, the address must be a file system path.
func NewHook(protocol, address string, addNewline bool) (*Hook, error) {
	logstashFormatter := logrus_logstash.LogstashFormatter{
		TimestampFormat: time.RFC3339Nano,
	}
	return &Hook{
		protocol:   protocol,
		address:    address,
		formatter:  logstashFormatter,
		addNewline: addNewline,
	}, nil
}

// Fire send log to the defined socket
func (h *Hook) Fire(entry *logrus.Entry) error {
	var err error
	dataBytes, err := h.formatter.Format(entry)
	if err != nil {
		return err
	}

	conn, err := net.Dial(h.protocol, h.address)
	if err != nil {
		return nil
	}
	defer conn.Close()
	// Add new line to every message if desired.
	if h.addNewline {
		dataBytes = append(dataBytes, "\n"...)
	}
	_, _ = conn.Write(dataBytes) // #nosec
	return nil
}

// Levels return an array of handled logging levels
func (h *Hook) Levels() []logrus.Level {
	return logrus.AllLevels
}
