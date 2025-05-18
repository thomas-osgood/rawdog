package defaults

import "time"

// default address and port the TeamServer will
// listen for incoming connections on.
const DEFAULT_ADDRESS string = "0.0.0.0:8080"

// default read timeout for server.
const DEFAULT_READ_TIMEOUT time.Duration = 15 * time.Second

// default send timeout for server.
const DEFAULT_SEND_TIMEOUT time.Duration = 5 * time.Second
