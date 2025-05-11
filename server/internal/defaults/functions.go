package defaults

import (
	"fmt"
	"log"
	"net"
	"rawdog/comms"
	"rawdog/server/internal/messages"
)

// function designed to handle when a request comes in
// to an endpoint that does not exist.
func InvalidEndpointHandler(conn net.Conn, md comms.TcpHeader, data []byte) (string, error) {

	// print out remote address and invalid endpoint
	// requested.
	log.Printf("\"%s\" requested invalid endpoint \"%d\"\n", conn.RemoteAddr(), md.Endpoint)

	return "", fmt.Errorf(messages.ERR_ENDPOINT_UNKNOWN, md.Endpoint)
}
