package server

import (
	"encoding/json"
	"io"
	"log"
	"net"
	"net/http"
	"rawdog/comms"
	"rawdog/server/internal/messages"
)

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

func (ts *TeamServer) AddEndpoint(endpoint int, handler TcpEndpointHandler) {
	ts.endpoints[endpoint] = handler
}

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

func (ts *TeamServer) handleConn(conn net.Conn) {
	defer conn.Close()

	var err error
	var routeHandler TcpEndpointHandler
	var md comms.TcpHeader = comms.TcpHeader{}
	var messageBuff []byte
	var ok bool
	var response comms.TcpStatusMessage = comms.TcpStatusMessage{Code: http.StatusOK}
	var transmission *comms.TcpTransmission

	transmission, err = comms.ReadTransmission(conn)
	if err != nil {
		if err == io.EOF {
			return
		}

		log.Printf(messages.ERR_DATA_READ, err.Error())
		return
	}

	// if no metadata was received, return an error.
	if transmission.MdSize < 1 {
		comms.SendTransmission(conn, []byte("no metadata found"), "")
		return
	}

	// attempt to unmarshal the header information
	// so the data can be dispatched correctly.
	err = json.Unmarshal(transmission.Metadata, &md)
	if err != nil {
		comms.SendTransmission(conn, []byte(err.Error()), "")
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

	err = comms.SendTransmission(conn, messageBuff, "")
	if err != nil {
		log.Printf(messages.ERR_SEND_RESPONSE, err.Error())
	}
}
