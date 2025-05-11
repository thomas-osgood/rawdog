package comms

// structure of the metadata block that will
// be expected in a transmission.
type TcpHeader struct {
	Agentname string `json:"agentname" xml:"agentname"`
	Endpoint  int    `json:"endpoint" xml:"endpoint"`
	Addldata  string `json:"addldata" xml:"addldata"`
}

// structure designed to convey a status.
type TcpStatusMessage struct {
	Code    int    `json:"code" xml:"code"`
	Message string `json:"message,omitempty" xml:"message"`
}

// structure of a TCP transmission.
type TcpTransmission struct {
	MdSize   uint16
	DatSize  uint64
	Data     []byte
	Metadata []byte
}
