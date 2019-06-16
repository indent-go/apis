package pgproxy

import (
	"log"
	"net"

	"github.com/indent-go/apis/pkg/database/postgres/pgproto"
)

func proxyOriginToClient(origin, client net.Conn, logger *log.Logger) (copyActive bool, err error) {
	msg := make([]byte, 4096)
	var length int
	var msgLength int32

	for {
		log.Printf("Proxying from origin '%s' to client '%s'.", origin.RemoteAddr(), client.RemoteAddr())
		if length, err = origin.Read(msg); err != nil {
			return
		}

		if _, err = client.Write(msg[:length]); err != nil {
			logger.Printf("Error sending origin (%s) message to client (%s): %v",
				origin.RemoteAddr(), client.RemoteAddr(), err)
		}

		for i := 0; i < length; i += int(msgLength) + 1 {
			msgType := pgproto.Type(msg[i:])
			msgLength = pgproto.Len(msg[i:])

			switch msgType {
			case pgproto.ReadyForQueryMsg:
				return false, nil
			case pgproto.CopyInResponseMsg:
				return true, nil
			}
		}
	}
}

// Pool of origins shared between clients.
type Pool chan net.Conn

// Add a new origin to the Pool.
func (p Pool) Add(conn net.Conn) {
	p <- conn
}

// Pop an origin out of the Pool to be used.
func (p Pool) Pop() net.Conn {
	return <-p
}

// Available origins in Pool.
func (p Pool) Available() int {
	return len(p)
}
