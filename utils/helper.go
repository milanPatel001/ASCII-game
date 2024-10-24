package utils

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"net"
)

func ConvBytesToIpv4(ip [4]byte) net.IP {
	return net.IPv4(ip[0], ip[1], ip[2], ip[3])
}

func ConvIpv4ToBytes(ip net.IP) [4]byte {
	return [4]byte(ip.To4())
}

func ConvSimplePayloadToBytes(payload any) ([]byte, error) {
	var payloadBuf bytes.Buffer

	err := binary.Write(&payloadBuf, binary.BigEndian, payload)

	if err != nil {
		return nil, err
	}

	return payloadBuf.Bytes(), nil

}

func ConvComplexPayloadToBytes(payload any) ([]byte, error) {
	var buf bytes.Buffer

	enc := gob.NewEncoder(&buf)

	err := enc.Encode(payload)

	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
