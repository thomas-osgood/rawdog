package defaults

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"rawdog/comms"
	"rawdog/server/internal/messages"
)

func InternalErrorSender(conn net.Conn, message []byte, md string) (err error) {
	var payload []byte
	var transmission comms.TcpStatusMessage = comms.TcpStatusMessage{
		Code:    http.StatusInternalServerError,
		Message: string(message),
	}

	payload, err = json.Marshal(&transmission)
	if err != nil {
		return err
	}

	return comms.SendTransmission(conn, payload, md)
}

// function designed to handle when a request comes in
// to an endpoint that does not exist.
func InvalidEndpointHandler(conn net.Conn, md comms.TcpHeader, data *bytes.Buffer) (string, error) {

	// print out remote address and invalid endpoint
	// requested.
	log.Printf("\"%s\" requested invalid endpoint \"%d\"\n", conn.RemoteAddr(), md.Endpoint)

	return "", fmt.Errorf(messages.ERR_ENDPOINT_UNKNOWN, md.Endpoint)
}
