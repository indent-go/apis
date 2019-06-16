package pgproxy

import (
	"crypto/tls"
	"fmt"
	"io"
	"net"

	"github.com/crunchydata/crunchy-proxy/connect"
	"github.com/crunchydata/crunchy-proxy/protocol"
)

func (s *Server) setupTLS() {
	// load certificates
	s.serverTLSConfig = tls.Config{}
	if crt, err := tls.LoadX509KeyPair(
		s.ServerTLSCert,
		s.ServerTLSKey,
	); err != nil {
		s.Fatalf("Failed to load TLS certificate or key: %v", err)
	} else {
		s.serverTLSConfig.Certificates = []tls.Certificate{crt}
	}
}

func (s *Server) upgradeConnToTLS(conn net.Conn) (net.Conn, error) {
	// confirm TLS availability
	tlsOk := []byte{protocol.SSLAllowed}
	if _, err := connect.Send(conn, tlsOk); err != nil {
		return conn, fmt.Errorf("failed sending SSLAllowed message: %v", err)
	}

	// upgrade connection
	conn = tls.Server(conn, &s.serverTLSConfig)

	// send startup message on newly upgraded connection
	if _, _, err := connect.Receive(conn); err == io.EOF {
		return conn, fmt.Errorf("Client closed connection during TLS setup")
	} else if err != nil {
		return conn, err
	}
	return conn, nil
}

func isTLSRequested(msg []byte) bool {
	protoVersion := protocol.GetVersion(msg)
	return protoVersion == protocol.SSLRequestCode
}
