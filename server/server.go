package server

import (
	"bytes"
	"encoding/json"
	"log"
	"net"
	"net/http"

	"github.com/thomas-osgood/rawdog/comms"

	"github.com/thomas-osgood/rawdog/server/internal/messages"
)

// function designed to start the server listening
// for incoming connections.
func (ts *TeamServer) Start() (err error) {
	ts.listener, err = net.Listen("tcp", ts.listenAddress)
	if err != nil {
		return err
	}
	defer ts.listener.Close()
	log.Printf("listening on \"%s\"\n", ts.listenAddress)

	go ts.acceptConnections()

	<-ts.quitChan
	return nil
}

// function designed to add a new endpoint handler
// to the server's endpoint map.
func (ts *TeamServer) AddEndpoint(endpoint int, handler TcpEndpointHandler) {
	ts.endpoints[endpoint] = handler
}

// function designed to act as the main loop for
// the server. this will accept connections and
// spawn the handleConn function in a new go routine.
func (ts *TeamServer) acceptConnections() {
	var conn net.Conn
	var err error

	for {
		conn, err = ts.listener.Accept()
		if err != nil {
			log.Printf("ERROR Accepting Conn: %s\n", err.Error())
			continue
		}
		log.Printf("new connection from \"%s\" ...\n", conn.RemoteAddr())

		go ts.handleConn(conn)
	}
}

// function designed to handle the connection from
// the client. this will dispatch the request to the
// correct endpoint and respond to the client accordingly.
func (ts *TeamServer) handleConn(conn net.Conn) {
	defer conn.Close()

	var err error
	var routeHandler TcpEndpointHandler
	var md comms.TcpHeader = comms.TcpHeader{}
	var messageBuff []byte
	var ok bool
	var response comms.TcpStatusMessage = comms.TcpStatusMessage{Code: http.StatusOK}
	var transmission *comms.TcpTransmission

	// read client request.
	transmission, err = comms.ReadTransmission(conn)
	if err != nil {
		log.Printf(messages.ERR_DATA_READ, err.Error())
		return
	}

	// if no metadata was received, return an error.
	if transmission.MdSize < 1 {
		ts.internalErrorFunc(conn, []byte("no metadata found"), "")
		return
	}

	// attempt to unmarshal the header information
	// so the data can be dispatched correctly.
	err = json.Unmarshal(transmission.Metadata, &md)
	if err != nil {
		ts.internalErrorFunc(conn, []byte(err.Error()), "")
		return
	}

	log.Printf("ENDPOINT: %d\n", md.Endpoint)

	// determine the route handler based on the endpoint.
	//
	// if the route is not found in the map, the handler
	// will be set to the invalid endpoint handler.
	routeHandler, ok = ts.endpoints[md.Endpoint]
	if !ok {
		routeHandler = ts.invalidEndpointHandler
	}

	// execute the correct handler function and
	// process the request.
	response.Message, err = routeHandler(conn, md, transmission.Data)

	// if there was an error handling the transmission,
	// set the message to the error.
	if err != nil {
		response.Code = http.StatusInternalServerError
		response.Message = err.Error()
	}

	// JSON encode the response that will be written to
	// the message buffer.
	messageBuff, err = json.Marshal(&response)
	if err != nil {
		log.Printf(messages.ERR_MARSHAL_RESPONSE, err.Error())
		return
	}

	// send response to the client.
	err = comms.SendTransmission(conn, bytes.NewBuffer(messageBuff), "")
	if err != nil {
		log.Printf(messages.ERR_SEND_RESPONSE, err.Error())
	}
}
