package comms

import (
	"bytes"
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"net"
	"rawdog/comms/internal/constants"
	"rawdog/comms/internal/messages"
)

// function designed to read data from a TCP transmission.
//
// this will unpack the transmission and handle
// it as needed.
//
// Expected Transmission Format:
// ------------------------------
//
// First 2 bytes: uint16 representing how large (in bytes) the metadata is.
//
// Next 8 bytes: uint64 representing how large (in bytes) the payload is.
//
// Next N bytes: metadata.
//
// Remaining N bytes: payload.
func ReadTransmission(conn net.Conn) (transmission *TcpTransmission, err error) {
	var dataBuff []byte = make([]byte, constants.SZ_DATABUFF)
	var decoded []byte
	var i uint64
	var iterations uint64
	var mdBuff []byte
	var n int
	var payloadBuff []byte = make([]byte, 0)
	var sizeBuffD [constants.SZ_SIZEBLOCK_DAT]byte
	var sizeBuffM [constants.SZ_SIZEBLOCK_MD]byte

	// initialize the object that will be returned to
	// avoid NIL dereference errors.
	transmission = new(TcpTransmission)

	// first 2 bytes will hold the size of the metadata.
	n, err = conn.Read(sizeBuffM[:])
	if err != nil {
		return nil, err
	}

	// convert the transmitted metadata size to a uint16
	// so it can be used later on in this function.
	transmission.MdSize = binary.BigEndian.Uint16(sizeBuffM[:n])

	// next 8 bytes will hold the size of the payload.
	n, err = conn.Read(sizeBuffD[:])
	if err != nil {
		return nil, err
	}

	// convert the transmitted payload size to a uint64
	// so it can be used later on in this function.
	transmission.DatSize = binary.BigEndian.Uint64(sizeBuffD[:n])

	if transmission.MdSize > 0 {
		// initialize a byte for the metadata and allocate
		// the exact size of the metadata to it.
		mdBuff = make([]byte, transmission.MdSize)

		// next 1024 bytes will contain metadata.
		n, err = conn.Read(mdBuff)
		if err != nil {
			return nil, err
		}

		// because the metadata block is guarenteed to be
		// 1024 bytes, any space not used will be NULL bytes (\x00)
		// and need to be trimmed manually or there will be
		// JSON Unmarshal errors.
		transmission.Metadata = bytes.Trim(mdBuff[:n], constants.NULL_BYTE)
	} else {
		transmission.Metadata = make([]byte, 0)
	}

	// initialize the new byte slice that will hold
	// the data from the agent based on the size read
	// from the first 8 bytes of the received transmission.
	transmission.Data = new(bytes.Buffer)

	// if the payload size is 0, do not attempt to continue.
	//
	// return an empty string.
	if transmission.DatSize == 0 {
		return transmission, nil
	}

	// determine how many blocks it will take to read
	// all the data from the client.
	iterations = transmission.DatSize / constants.SZ_DATABUFF
	if (transmission.DatSize % constants.SZ_DATABUFF) > 0 {
		iterations++
	}

	// read the payload in 2048 byte chunks.
	for i = 0; i < iterations; i++ {
		n, err = conn.Read(dataBuff)
		if err != nil {

			// no more data can be read because
			// the END OF FILE has been reached,
			// so break the loop and do not attmept
			// to read anymore.
			if err == io.EOF {
				break
			}

			continue
		}

		// build the payload slice.
		payloadBuff = append(payloadBuff, dataBuff[:n]...)
	}

	// base64-decode the payload that was sent and save
	// the result in the Data field of the transmission object.
	decoded, err = base64.StdEncoding.DecodeString(string(payloadBuff))
	if err != nil {
		return nil, err
	}

	// read the decoded data into the buffer.
	transmission.Data.Read(decoded)

	return transmission, nil
}

// function designed to tranmsit data to the client
// that is connected to the teamserver via the conn
// that is passed in.
//
// Transmission Format:
// ------------------------------
//
// First 2 bytes: uint16 representing how large (in bytes) the metadata is.
//
// Next 8 bytes: uint64 representing how large (in bytes) the payload is.
//
// Next N bytes: metadata.
//
// Remaining N bytes: payload.
//
// references:
//
// https://stackoverflow.com/questions/16888357/convert-an-integer-to-a-byte-array
func SendTransmission(conn net.Conn, data *bytes.Buffer, metadata string) (err error) {
	var dataBuff *bytes.Buffer
	var dataEnc string
	var lenBuffD [constants.SZ_SIZEBLOCK_DAT]byte
	var lenBuffM [constants.SZ_SIZEBLOCK_MD]byte
	var lenData int
	var lenMd int = len(metadata)
	var mdBuff []byte = make([]byte, constants.SZ_METADATA_MAX)
	var n int64

	// make sure the metadata length does not exceeed
	// the max allowed buffer size.
	if lenMd > int(constants.SZ_METADATA_MAX) {
		return fmt.Errorf(messages.ERR_MD_LARGE, lenMd, constants.SZ_METADATA_MAX)
	}

	// base64-encode the passed in data and save the
	// length of it.
	dataEnc = base64.StdEncoding.EncodeToString(data.Bytes())
	lenData = len(dataEnc)

	// initialize byte buffer that will hold the
	// data to transmit to the client.
	dataBuff = new(bytes.Buffer)

	// set the first 2 bytes to be the length of the metadata.
	binary.BigEndian.PutUint16(lenBuffM[:], uint16(lenMd))

	// write the metadata info to the metdata buffer.
	mdBuff = make([]byte, lenMd)
	copy(mdBuff, []byte(metadata))

	// set the next 8 bytes to be the length of the payload.
	binary.BigEndian.PutUint64(lenBuffD[:], uint64(lenData))

	// set the md length packet.
	dataBuff.Write(lenBuffM[:])
	// set the length packet.
	dataBuff.Write(lenBuffD[:])
	// set the metadata.
	dataBuff.Write(mdBuff)
	// set the data packet.
	dataBuff.Write([]byte(dataEnc))

	// transmit the block to the client.
	n, err = io.Copy(conn, dataBuff)
	if err != nil {
		return err
	}

	log.Printf("sent %d response bytes\n", n)

	return nil
}
