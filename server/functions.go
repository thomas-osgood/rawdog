package server

import (
	"fmt"
	"rawdog/server/internal/defaults"
	"strings"
)

// function designed to create and initialize a
// new teamserver based on the configuration options
// passed in.
func NewTeamServer(opts ...TeamServerConfigFunc) (ts *TeamServer, err error) {
	var config *TeamServerConfig = &TeamServerConfig{
		InvalidEndpointHandler: nil,
		ListenAddress:          defaults.DEFAULT_ADDRESS,
		QuitChan:               make(chan struct{}),
	}
	var configFunc TeamServerConfigFunc

	// go through the custom options passed in and
	// set the values.
	for _, configFunc = range opts {
		err = configFunc(config)
		if err != nil {
			return nil, err
		}
	}

	// if no InvalidEndpointHandler has been set, use
	// the default handler function.
	if config.InvalidEndpointHandler == nil {
		config.InvalidEndpointHandler = defaults.InvalidEndpointHandler
	}

	// assign values to the teamserver that will
	// be returned by this function.
	ts = &TeamServer{
		endpoints:              make(map[int]TcpEndpointHandler),
		invalidEndpointHandler: config.InvalidEndpointHandler,
		listenAddress:          config.ListenAddress,
		quitChan:               config.QuitChan,
	}

	return ts, nil
}

// function designed to set the InvalidEndpointHandler
// that will be used with the server.
func WithInvalidEndpointHandler(handler TcpEndpointHandler) TeamServerConfigFunc {
	return func(tsc *TeamServerConfig) error {

		if tsc.InvalidEndpointHandler != nil {
			return fmt.Errorf("can only set one InvalidEndpointHandler")
		}

		tsc.InvalidEndpointHandler = handler

		return nil
	}
}

// function designed to set the listen address for
// the TeamServer.
func WithListenAddress(listenaddress string) TeamServerConfigFunc {
	return func(tsc *TeamServerConfig) (err error) {
		listenaddress = strings.TrimSpace(listenaddress)

		if len(listenaddress) < 1 {
			return fmt.Errorf("listenaddress cannot be an empty string")
		}

		tsc.ListenAddress = listenaddress

		return nil
	}
}
