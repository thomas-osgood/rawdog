package server

import (
	"bytes"
	"net"
	"rawdog/comms"
)

// type defining a function that can be used to
// set a value during a TeamServer initialization.
type TeamServerConfigFunc func(*TeamServerConfig) error

// alias for a map that can be used to hold endpoint
// handlers for a TeamServer.
type EndpointMap map[int]TcpEndpointHandler

// type defining the shape of a function that can be
// used as a "TCP Endpoint".
type TcpEndpointHandler func(net.Conn, comms.TcpHeader, *bytes.Buffer) (string, error)
