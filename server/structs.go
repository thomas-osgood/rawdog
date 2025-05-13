package server

import (
	"net"
	"rawdog/comms"
)

// structure defining the TeamServer object.
type TeamServer struct {
	// function that will fire off when an error occurs
	// during the handling/dispatching of a request.
	internalErrorFunc comms.TcpTransmissionFunc
	// handler that will be used when an invalid endpoint
	// gets requested.
	invalidEndpointHandler TcpEndpointHandler
	// address the server will listen for connections on.
	listenAddress string
	// TCP listener for the server. this will accept the
	// incoming connections.
	listener net.Listener
	// channel designed to block until a signal is sent
	// to it, indicating the server can shutdown.
	quitChan chan struct{}
	// map holding all endpoints the server can handle.
	endpoints EndpointMap
}

// structure defining the various configuration
// options that can be set for a new TeamServer.
type TeamServerConfig struct {
	// function that will transmit error messages
	// to the client.
	InternalErrorFunc comms.TcpTransmissionFunc
	// function that will handle when an invalid
	// endpoint has been requested.
	InvalidEndpointHandler TcpEndpointHandler
	// address the server should listen on.
	//
	// should be in the form "<ip>:<port>"
	//
	// ex: 0.0.0.0:1234
	ListenAddress string
	// channel that will block to allow the server
	// to run forever.
	QuitChan chan struct{}
}
