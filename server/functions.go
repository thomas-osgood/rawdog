package server

import (
	"fmt"
	"strings"
	"time"

	"github.com/thomas-osgood/rawdog/server/internal/defaults"
	"github.com/thomas-osgood/rawdog/server/internal/messages"
)

// function designed to create and initialize a
// new teamserver based on the configuration options
// passed in.
//
// Example:
//
//	var err error
//	var newsrv *TeamServer
//
//	newsrv, err = NewTeamServer(WithListenAddress("0.0.0.0:9999"))
//	if err != nil {
//		log.Fatal(err)
//	}
func NewTeamServer(opts ...TeamServerConfigFunc) (ts *TeamServer, err error) {
	var config *TeamServerConfig = &TeamServerConfig{
		InternalErrorFunc:      nil,
		InvalidEndpointHandler: nil,
		ListenAddress:          defaults.DEFAULT_ADDRESS,
		QuitChan:               make(chan struct{}),
		RecvTimeout:            0,
		SendTimeout:            0,
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

	// if no custom endpoint map has been specified, create
	// a blank map.
	if config.Endpoints == nil {
		config.Endpoints = make(EndpointMap)
	}

	// if no InternalErrorFunc has been set, use the
	// default function.
	if config.InternalErrorFunc == nil {
		config.InternalErrorFunc = defaults.InternalErrorSender
	}

	// if no recv timeout has been set, set the recv
	// timeout to the default value.
	if config.SendTimeout == 0 {
		config.RecvTimeout = defaults.DEFAULT_READ_TIMEOUT
	}

	// if no send timeout has been set, set the send
	// timeout to the default value.
	if config.SendTimeout == 0 {
		config.SendTimeout = defaults.DEFAULT_SEND_TIMEOUT
	}

	// assign values to the teamserver that will
	// be returned by this function.
	ts = &TeamServer{
		endpoints:              config.Endpoints,
		internalErrorFunc:      config.InternalErrorFunc,
		invalidEndpointHandler: config.InvalidEndpointHandler,
		listenAddress:          config.ListenAddress,
		quitChan:               config.QuitChan,
		timeoutRecv:            config.RecvTimeout,
		timeoutSend:            config.SendTimeout,
	}

	return ts, nil
}

// function designed to set the endpoints map the server
// will use when handling requests.
func WithEndpoints(endpoints EndpointMap) TeamServerConfigFunc {
	return func(tsc *TeamServerConfig) error {

		if tsc.Endpoints != nil {
			return fmt.Errorf(messages.ERR_ENDPOINT_MAP_SET)
		}

		tsc.Endpoints = endpoints

		return nil
	}
}

// function designed to set the InvalidEndpointHandler
// that will be used with the server.
func WithInvalidEndpointHandler(handler TcpEndpointHandler) TeamServerConfigFunc {
	return func(tsc *TeamServerConfig) error {

		if tsc.InvalidEndpointHandler != nil {
			return fmt.Errorf(messages.ERR_IEH_SET)
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
			return fmt.Errorf(messages.ERR_LISTEN_EMPTY)
		}

		tsc.ListenAddress = listenaddress

		return nil
	}
}

// function designed to set the receive timeout for the server.
func WithRecvTimeout(timeout time.Duration) TeamServerConfigFunc {
	return func(tsc *TeamServerConfig) error {
		if tsc.RecvTimeout > 0 {
			return fmt.Errorf(messages.ERR_TIMEOUT_RECV_SET)
		}

		tsc.RecvTimeout = timeout

		return nil
	}
}

// function designed to set the send timeout for the server.
func WithSendTimeout(timeout time.Duration) TeamServerConfigFunc {
	return func(tsc *TeamServerConfig) error {
		if tsc.RecvTimeout > 0 {
			return fmt.Errorf(messages.ERR_TIMEOUT_SEND_SET)
		}

		tsc.SendTimeout = timeout

		return nil
	}
}
