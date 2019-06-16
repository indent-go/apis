package pgproxy

import (
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/crunchydata/crunchy-proxy/connect"
	"github.com/crunchydata/crunchy-proxy/protocol"
)

// SetupInput is the configuration for the Postgres Encryption proxy
type SetupInput struct {
	configFilePath string
}

// Server proxies Postgres Wire Protocol traffic while executing any Envelopes discovered.
type Server struct {
	Config
	*log.Logger

	serverTLSConfig tls.Config
	originPool      Pool
	lis             net.Listener
}

// Setup establishes connection to the origin Postgres server.
func (s *Server) Setup() error {
	if s.Logger == nil {
		s.Logger = log.New(os.Stderr, "", log.LstdFlags)
	}

	s.Printf("Connecting to origin at %s...", s.OriginHost)
	origin, err := connect.Connect(s.OriginHost)
	if err != nil {
		s.Fatalf("Failed to connect to origin '%s': %v", s.OriginHost, err)
	}

	opts := map[string]string{}
	startupMsg := protocol.CreateStartupMessage(s.OriginUsername, s.OriginDatabase, opts)
	origin.Write(startupMsg)

	response := make([]byte, 4096)
	if _, err = origin.Read(response); err != nil {
		return fmt.Errorf("error connecting to '%s': %v", s.OriginHost, err)
	}

	if authd := connect.HandleAuthenticationRequest(origin, response); !authd {
		return fmt.Errorf("Origin authentication failed")
	}

	s.Printf("Successfully connected to origin '%s'", s.OriginHost)
	s.originPool = make(Pool, 20)
	s.originPool.Add(origin)
	return nil
}

// Start begins listening for client connections and proxies them to the Origin.
func (s *Server) Start(stopCh <-chan struct{}) (err error) {
	s.Print("Starting server...")
	s.lis, err = net.Listen("tcp", s.ServerHost)
	if err != nil {
		return err
	}
	s.Printf("Listening on %s", s.ServerHost)

	for {
		select {
		case <-stopCh:
			break
		default:
			conn, err := s.lis.Accept()
			if err != nil {
				return err
			}

			if conn != nil {
				go s.HandleConnection(conn)
			}
		}
	}
}
