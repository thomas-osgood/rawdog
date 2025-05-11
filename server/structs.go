package server

import "net"

// structure defining the TeamServer object.
type TeamServer struct {
	invalidEndpointHandler TcpEndpointHandler
	listenAddress          string
	listener               net.Listener
	quitChan               chan struct{}
	endpoints              EndpointMap
}

// structure defining the various configuration
// options that can be set for a new TeamServer.
type TeamServerConfig struct {
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
