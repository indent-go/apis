package pgproxy

import (
	"encoding/json"
	"io"
	"net"

	"github.com/crunchydata/crunchy-proxy/connect"
	access "go.indent.com/apis/pkg/access/v1"
	"go.indent.com/apis/pkg/database/postgres/pgproto"
)

// HandleConnection is the handler for a single TCP connection
func (s *Server) HandleConnection(client net.Conn) {
	var (
		startupRcvd = false

		// tracks bytes remaining in copy operation (if 0 copy is done)
		copyLeft   int
		copyActive bool

		origin net.Conn

		msg    []byte
		length int
		err    error
	)

	for {
		s.Printf("Ready to receive messages from client (%s)...", client.RemoteAddr())
		if msg, length, err = connect.Receive(client); err != nil {
			if err == io.EOF {
				s.Printf("Client '%s' closed connection to server.", client.RemoteAddr())
			} else {
				s.Printf("Error receiving from client '%s': %v", client.RemoteAddr(), err)
			}
			break
		}

		msgType := pgproto.Type(msg)
		s.Printf("Received %s message.", msgType.String(true))

		if msgType == pgproto.QueryMsg {
			msgBody := string(pgproto.QueryBody(msg))
			req := access.Request{
				Actor:   access.Actor{ID: client.RemoteAddr().String()},
				Actions: access.Actions{"indent:actions::sql.select"},
				Resources: access.Resources{
					"indent:resources::tables:users",
					"indent:resources::columns:users.*",
				},
				Metadata: &access.Metadata{
					Labels: access.Labels{
						"originalQuery": msgBody,
					},
				},
			}

			s.Printf("Query = %s", string(msgBody))

			jm, _ := json.Marshal(req)

			s.Printf("Audit = %s", string(jm))
		}

		// process first message as startup message
		if !startupRcvd {
			if err = s.handleStartup(client, msg, length); err != nil {
				s.Printf("Startup error with client '%s': %v", client.RemoteAddr().String(), err)
				return
			}

			startupRcvd = true
			continue
		}

		if msgType == pgproto.TerminateMsg {
			s.Printf("Client '%s' sent terminate message.", client.RemoteAddr())
			return
		}

		if origin == nil {
			s.Print("Setting up origin...")
			origin = s.originPool.Pop()
		}

		s.Printf("Proxying message from client '%s' to origin '%s'", client.RemoteAddr(), origin.RemoteAddr())
		if _, err = connect.Send(origin, msg); err != nil {
			s.Printf("Proxy error from client '%s' to origin '%s': %v", client.RemoteAddr(), origin.RemoteAddr(), err)
		}

		// handle copy operations if active
		if copyActive {
			copyActive, copyLeft = processCopy(msg, copyLeft, length, s.Logger)
			s.Printf("Copied data, %d left", copyLeft)

			if copyActive {
				continue
			}
			s.Printf("Finished data transfer (got %d left) msg: %d.", copyLeft, length)
		}

		if copyActive, err = proxyOriginToClient(origin, client, s.Logger); err != nil {
			s.Printf("Error proxying from origin (%s) to client (%s): %v",
				origin.RemoteAddr(), client.RemoteAddr(), err)
			return
		} else if copyActive {
			continue
		}

		s.Printf("Closing connection to origin '%s'", origin.RemoteAddr())
		s.originPool.Add(origin)
		origin = nil
		s.Printf("Switching to receive for client '%s'", client.RemoteAddr())
	}
}

func (s *Server) handleStartup(client net.Conn, msg []byte, length int) (err error) {
	// upgrade connection if TLS is being used
	if isTLSRequested(msg) {
		if client, err = s.upgradeConnToTLS(client); err != nil {
			s.Printf("Failed to upgrade client connection to TLS: %v", err)
			return
		}
	}

	s.Logger.Printf("Authenticating client '%s'...", client.RemoteAddr())
	if authd, err := connect.AuthenticateClient(client, msg, length); err != nil {
		s.Logger.Printf("Error authing client '%s': %v", client.RemoteAddr(), err)
		return err
	} else if !authd {
		s.Logger.Printf("Auth failed for client '%s'", client.RemoteAddr())
		return err
	}
	return nil
}
